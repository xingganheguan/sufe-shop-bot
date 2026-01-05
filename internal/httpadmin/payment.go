package httpadmin

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	logger "shop-bot/internal/log"
	"shop-bot/internal/metrics"
	"shop-bot/internal/notification"
	payment "shop-bot/internal/payment/epay"
	"shop-bot/internal/store"
)

// handlePaymentReturn handles the payment return page
func (s *Server) handlePaymentReturn(c *gin.Context) {
	// Check if this is a payment result with parameters
	tradeStatus := c.Query("trade_status")
	outTradeNo := c.Query("out_trade_no")
	productName := c.Query("name")
	amount := c.Query("money")
	payType := c.Query("type")

	// Get bot username for the link
	var botUsername string
	if s.bot != nil {
		botUsername = s.bot.Self.UserName
	}

	// Prepare order info for display
	orderInfo := gin.H{}
	if outTradeNo != "" {
		orderInfo["OrderNo"] = outTradeNo
	}
	if productName != "" && productName != "product" {
		orderInfo["ProductName"] = productName
	}
	if amount != "" {
		orderInfo["Amount"] = amount
	}
	if payType != "" {
		orderInfo["PayType"] = payType
	}

	// Prepare bot link
	botLink := fmt.Sprintf("https://t.me/%s", botUsername)
	if botUsername == "" {
		botLink = "#"
	}

	if tradeStatus == "TRADE_SUCCESS" && outTradeNo != "" {
		// This looks like a payment notification via GET
		// Convert query params to form values for compatibility
		params := make(url.Values)
		for k, v := range c.Request.URL.Query() {
			params[k] = v
		}

		logger.Info("Processing payment return as notification", "out_trade_no", outTradeNo, "params", params)

		// Process as payment notification
		s.processPaymentNotification(c, params)

		// Show beautiful success page
		c.HTML(http.StatusOK, "payment_success.html", gin.H{
			"orderInfo": orderInfo,
			"botLink":   botLink,
		})
		return
	}

	// Simple return page
	c.HTML(http.StatusOK, "payment_success.html", gin.H{
		"botLink": botLink,
	})
}

// handleEpayNotify handles payment callbacks from EPay
func (s *Server) handleEpayNotify(c *gin.Context) {
	metrics.PaymentCallbacksReceived.Inc()

	// Parse form data
	if err := c.Request.ParseForm(); err != nil {
		logger.Error("Failed to parse form", "error", err)
		c.String(http.StatusBadRequest, "fail")
		return
	}

	params := c.Request.Form
	s.processPaymentNotification(c, params)

	logger.Info("Payment processed successfully")
	c.String(http.StatusOK, "success")
}

// processPaymentNotification processes payment notification
func (s *Server) processPaymentNotification(c *gin.Context, params url.Values) {
	metrics.PaymentCallbacksReceived.Inc()

	traceID := c.GetString("trace_id")
	logger.Info("Processing payment notification", "params", params, "trace_id", traceID)

	// Verify signature
	if s.epay == nil || !s.epay.VerifyNotify(params) {
		logger.Error("Invalid callback signature", "params", params)
		return
	}

	// Parse notification
	notify := payment.ParseNotify(params)

	// Check trade status
	if notify.TradeStatus != "TRADE_SUCCESS" {
		logger.Info("Trade not successful", "status", notify.TradeStatus)
		return
	}

	// Find order by out_trade_no
	var order store.Order
	if err := s.db.Preload("User").Preload("Product").Where("epay_out_trade_no = ?", notify.OutTradeNo).First(&order).Error; err != nil {
		// Try parsing order ID from out_trade_no (format: orderID-timestamp)
		parts := strings.Split(notify.OutTradeNo, "-")
		if len(parts) > 0 {
			if orderID, err := strconv.ParseUint(parts[0], 10, 32); err == nil {
				err = s.db.Preload("User").Preload("Product").First(&order, orderID).Error
			}
		}

		if err != nil {
			logger.Error("Order not found", "out_trade_no", notify.OutTradeNo, "error", err)
			metrics.PaymentCallbacksFailed.Inc()
			return
		}
	}

	// Check if already processed
	if order.Status != "pending" {
		logger.Info("Order already processed", "order_id", order.ID, "status", order.Status, "trace_id", traceID)
		return
	}

	// Use transaction to process payment
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Update order status
		updates := map[string]interface{}{
			"status":        "paid",
			"epay_trade_no": notify.TradeNo,
		}
		if err := tx.Model(&order).Updates(updates).Error; err != nil {
			return err
		}

		// Track metric
		metrics.OrdersPaid.Inc()

		// Handle product delivery or balance recharge
		if order.ProductID != nil {
			// Product order - try to claim code
			code, err := store.ClaimOneCodeTx(context.Background(), tx, *order.ProductID, order.ID)
			if err != nil {
				if err == store.ErrNoStock {
					// Update status to paid_no_stock
					if err := tx.Model(&order).Update("status", "paid_no_stock").Error; err != nil {
						return err
					}

					// Track no stock metric
					metrics.OrdersNoStock.Inc()

					// Will notify user and admin after transaction
					return nil
				}
				return err
			}

			// Update order status to delivered
			if err := tx.Model(&order).Update("status", "delivered").Error; err != nil {
				return err
			}

			// Track delivered metric
			metrics.OrdersDelivered.Inc()

			// Send code to user after transaction
			go s.sendCodeToUser(&order, code)
		} else {
			// Balance recharge
			if err := store.AddBalance(tx, order.UserID, order.AmountCents, "recharge",
				fmt.Sprintf("充值订单 #%d", order.ID), nil, &order.ID); err != nil {
				return err
			}

			// Update order status to delivered
			if err := tx.Model(&order).Update("status", "delivered").Error; err != nil {
				return err
			}

			// Send notification after transaction
			go s.sendRechargeSuccessMessage(&order)
		}

		return nil
	})

	if err != nil {
		logger.Error("Failed to process payment", "order_id", order.ID, "error", err, "trace_id", traceID)
		metrics.PaymentCallbacksFailed.Inc()
		return
	}

	logger.Info("Order payment confirmed", "order_id", order.ID, "trade_no", notify.TradeNo, "trace_id", traceID)

	// Send notification to admins
	if s.notification != nil {
		productName := "余额充值"
		if order.Product != nil {
			productName = order.Product.Name
		}
		s.notification.NotifyAdmins(notification.EventOrderPaid, map[string]interface{}{
			"order_id":       order.ID,
			"user_id":        order.UserID,
			"product_name":   productName,
			"amount":         order.AmountCents,
			"payment_method": "Epay",
		})
	}

	// Handle no stock notification
	if order.Status == "paid_no_stock" {
		go s.notifyNoStock(&order)
	}
}

// sendCodeToUser sends the purchased code to the user
func (s *Server) sendCodeToUser(order *store.Order, code string) {
	if s.bot == nil {
		return
	}

	message := fmt.Sprintf(
		"🎉 购买成功！\n\n"+
			"订单号: #%d\n"+
			"商品: %s\n"+
			"金额: ¥%.2f\n\n"+
			"📦 您的卡密信息：\n"+
			"<code>%s</code>\n\n"+
			"感谢您的购买！如有问题请联系客服。",
		order.ID,
		order.Product.Name,
		float64(order.AmountCents)/100,
		code,
	)
	msg := tgbotapi.NewMessage(order.User.TgUserID, message)
	msg.ParseMode = "HTML"
	s.bot.Send(msg)
}

// sendRechargeSuccessMessage sends recharge success message to user
func (s *Server) sendRechargeSuccessMessage(order *store.Order) {
	if s.bot == nil {
		return
	}

	newBalance, _ := store.GetUserBalance(s.db, order.UserID)
	message := fmt.Sprintf(
		"✅ 充值成功！\n\n"+
			"订单号: #%d\n"+
			"充值金额: ¥%.2f\n"+
			"当前余额: ¥%.2f\n\n"+
			"感谢您的充值！",
		order.ID,
		float64(order.AmountCents)/100,
		float64(newBalance)/100,
	)
	msg := tgbotapi.NewMessage(order.User.TgUserID, message)
	s.bot.Send(msg)
}

// notifyNoStock notifies user and admin about no stock
func (s *Server) notifyNoStock(order *store.Order) {
	// Notify user
	if s.bot != nil {
		message := fmt.Sprintf(
			"⚠️ 抱歉，商品 %s 暂时缺货\n\n"+
				"您的订单 #%d 已支付成功，但商品暂时无货。\n"+
				"请联系客服处理退款或等待补货。\n\n"+
				"给您带来的不便深感抱歉！",
			order.Product.Name,
			order.ID,
		)
		msg := tgbotapi.NewMessage(order.User.TgUserID, message)
		s.bot.Send(msg)
	}

	// Notify admins
	if s.notification != nil {
		productName := "Unknown"
		if order.Product != nil {
			productName = order.Product.Name
		}
		s.notification.NotifyAdmins(notification.EventNoStock, map[string]interface{}{
			"order_id":     order.ID,
			"product_name": productName,
			"user_id":      order.UserID,
			"amount":       order.AmountCents,
		})
	}
}
