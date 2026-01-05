package store

import (
	"bytes"
	"html/template"

	"gorm.io/gorm"
)

// GetMessageTemplate retrieves a message template by code and language
func GetMessageTemplate(db *gorm.DB, code, language string) (*MessageTemplate, error) {
	var tmpl MessageTemplate

	// Try to get template for specific language
	err := db.Where("code = ? AND language = ? AND is_active = ?", code, language, true).
		First(&tmpl).Error

	if err == gorm.ErrRecordNotFound {
		// Fallback to English
		err = db.Where("code = ? AND language = ? AND is_active = ?", code, "en", true).
			First(&tmpl).Error
	}

	if err != nil {
		return nil, err
	}

	return &tmpl, nil
}

// RenderTemplate renders a message template with variables
func RenderTemplate(tmplContent string, data interface{}) (string, error) {
	// Parse template
	tmpl, err := template.New("message").Parse(tmplContent)
	if err != nil {
		return "", err
	}

	// Render template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// CreateDefaultTemplates creates default message templates
func CreateDefaultTemplates(db *gorm.DB) error {
	templates := []MessageTemplate{
		{
			Code:      "order_paid",
			Language:  "en",
			Name:      "Order Paid Message",
			Content:   "🎉 Payment successful!\n\nOrder ID: {{.OrderID}}\nProduct: {{.ProductName}}\nCode: `{{.Code}}`\n\nThank you for your purchase!",
			Variables: `["OrderID", "ProductName", "Code"]`,
			IsActive:  true,
		},
		{
			Code:      "order_paid",
			Language:  "zh",
			Name:      "订单支付成功消息",
			Content:   "🎉 支付成功！\n\n订单号：{{.OrderID}}\n商品：{{.ProductName}}\n卡密：`{{.Code}}`\n\n感谢您的购买！",
			Variables: `["OrderID", "ProductName", "Code"]`,
			IsActive:  true,
		},
		{
			Code:      "no_stock",
			Language:  "en",
			Name:      "No Stock Message",
			Content:   "⚠️ Payment received but product is out of stock\n\nOrder ID: {{.OrderID}}\nProduct: {{.ProductName}}\n\nPlease contact support for refund or wait for restock.\nWe apologize for the inconvenience.",
			Variables: `["OrderID", "ProductName"]`,
			IsActive:  true,
		},
		{
			Code:      "no_stock",
			Language:  "zh",
			Name:      "库存不足消息",
			Content:   "⚠️ 已收到付款但商品缺货\n\n订单号：{{.OrderID}}\n商品：{{.ProductName}}\n\n请联系客服退款或等待补货。\n给您带来的不便深表歉意。",
			Variables: `["OrderID", "ProductName"]`,
			IsActive:  true,
		},
		{
			Code:      "balance_recharged",
			Language:  "en",
			Name:      "Balance Recharged Message",
			Content:   "💰 Balance recharged successfully!\n\nAmount: {{.Currency}}{{.Amount}}\nNew Balance: {{.Currency}}{{.NewBalance}}\nCard: {{.CardCode}}",
			Variables: `["Currency", "Amount", "NewBalance", "CardCode"]`,
			IsActive:  true,
		},
		{
			Code:      "balance_recharged",
			Language:  "zh",
			Name:      "余额充值成功消息",
			Content:   "💰 余额充值成功！\n\n充值金额：{{.Currency}}{{.Amount}}\n当前余额：{{.Currency}}{{.NewBalance}}\n充值卡：{{.CardCode}}",
			Variables: `["Currency", "Amount", "NewBalance", "CardCode"]`,
			IsActive:  true,
		},
	}

	for _, tmpl := range templates {
		// Check if template already exists
		var existing MessageTemplate
		err := db.Where("code = ? AND language = ?", tmpl.Code, tmpl.Language).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			// Create new template
			if err := db.Create(&tmpl).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// UpdateMessageTemplate updates a message template
func UpdateMessageTemplate(db *gorm.DB, id uint, content string, isActive bool) error {
	return db.Model(&MessageTemplate{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"content":   content,
			"is_active": isActive,
		}).Error
}

// GetAllTemplates returns all message templates
func GetAllTemplates(db *gorm.DB) ([]MessageTemplate, error) {
	var templates []MessageTemplate
	err := db.Order("code, language").Find(&templates).Error
	return templates, err
}

// ValidateTemplateVariables validates template variables
func ValidateTemplateVariables(content string, allowedVars []string) error {
	// Parse template to check variables
	tmpl, err := template.New("validate").Parse(content)
	if err != nil {
		return err
	}

	// Create test data with all allowed variables
	testData := make(map[string]interface{})
	for _, v := range allowedVars {
		testData[v] = "test"
	}

	// Try to execute template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, testData); err != nil {
		return err
	}

	return nil
}

// GetTemplateVariables returns the variables for a template code
func GetTemplateVariables(code string) []string {
	varMap := map[string][]string{
		"order_paid":        {"OrderID", "ProductName", "Code"},
		"no_stock":          {"OrderID", "ProductName"},
		"balance_recharged": {"Currency", "Amount", "NewBalance", "CardCode"},
		"order_created":     {"ProductName", "Price", "OrderID"},
		"profile_info":      {"UserID", "Username", "Language", "JoinedDate", "Balance"},
	}

	if vars, ok := varMap[code]; ok {
		return vars
	}

	return []string{}
}
