package bot

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"shop-bot/internal/bot/messages"
	logger "shop-bot/internal/log"
	"shop-bot/internal/store"
)

// handleGroupMessage handles messages from groups
func (b *Bot) handleGroupMessage(message *tgbotapi.Message) {
	// Check if bot was mentioned or if it's a command
	if !message.IsCommand() && !strings.Contains(message.Text, "@"+b.api.Self.UserName) {
		return
	}

	// Handle group commands
	if message.IsCommand() {
		switch message.Command() {
		case "register":
			b.handleGroupRegister(message)
		case "unregister":
			b.handleGroupUnregister(message)
		case "settings":
			b.handleGroupSettings(message)
		case "help":
			b.handleGroupHelp(message)
		}
	}
}

// handleGroupRegister registers a group for notifications
func (b *Bot) handleGroupRegister(message *tgbotapi.Message) {
	// Check if user is group admin
	chatConfig := tgbotapi.ChatConfigWithUser{
		ChatID: message.Chat.ID,
		UserID: message.From.ID,
	}

	member, err := b.api.GetChatMember(tgbotapi.GetChatMemberConfig{ChatConfigWithUser: chatConfig})
	if err != nil {
		logger.Error("Failed to get chat member", "error", err)
		return
	}

	// Only administrators can register groups
	if member.Status != "administrator" && member.Status != "creator" {
		msg := tgbotapi.NewMessage(message.Chat.ID, "只有群组管理员可以注册此群组 / Only group administrators can register this group")
		b.api.Send(msg)
		return
	}

	// Get or create user
	user, err := store.GetOrCreateUser(b.db, message.From.ID, message.From.UserName)
	if err != nil {
		logger.Error("Failed to get/create user", "error", err)
		return
	}

	// Register group
	groupType := "group"
	if message.Chat.IsSuperGroup() {
		groupType = "supergroup"
	} else if message.Chat.IsChannel() {
		groupType = "channel"
	}

	group, err := store.RegisterGroup(b.db, message.Chat.ID, message.Chat.Title, groupType, user.ID)
	if err != nil {
		if err == store.ErrGroupExists {
			msg := tgbotapi.NewMessage(message.Chat.ID, "✅ 群组已经注册 / Group already registered")
			b.api.Send(msg)
		} else {
			logger.Error("Failed to register group", "error", err)
			msg := tgbotapi.NewMessage(message.Chat.ID, "❌ 注册失败，请稍后再试 / Registration failed, please try again later")
			b.api.Send(msg)
		}
		return
	}

	// Send success message
	lang := messages.GetUserLanguage(group.Language, "")
	successMsg := b.msg.Format(lang, "group_registered", map[string]interface{}{
		"GroupName": group.GroupName,
	})

	msg := tgbotapi.NewMessage(message.Chat.ID, successMsg)
	msg.ParseMode = "Markdown"
	b.api.Send(msg)

	logger.Info("Group registered", "group_id", group.ID, "tg_group_id", group.TgGroupID)
}

// handleGroupUnregister unregisters a group
func (b *Bot) handleGroupUnregister(message *tgbotapi.Message) {
	// Check if user is group admin
	chatConfig := tgbotapi.ChatConfigWithUser{
		ChatID: message.Chat.ID,
		UserID: message.From.ID,
	}

	member, err := b.api.GetChatMember(tgbotapi.GetChatMemberConfig{ChatConfigWithUser: chatConfig})
	if err != nil {
		logger.Error("Failed to get chat member", "error", err)
		return
	}

	// Only administrators can unregister groups
	if member.Status != "administrator" && member.Status != "creator" {
		msg := tgbotapi.NewMessage(message.Chat.ID, "只有群组管理员可以取消注册 / Only group administrators can unregister")
		b.api.Send(msg)
		return
	}

	// Get group
	group, err := store.GetGroup(b.db, message.Chat.ID)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ 群组未注册 / Group not registered")
		b.api.Send(msg)
		return
	}

	// Deactivate group
	if err := store.DeactivateGroup(b.db, group.ID); err != nil {
		logger.Error("Failed to deactivate group", "error", err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ 操作失败，请稍后再试 / Operation failed, please try again later")
		b.api.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "✅ 群组已取消注册，不再接收通知 / Group unregistered, will no longer receive notifications")
	b.api.Send(msg)

	logger.Info("Group unregistered", "group_id", group.ID)
}

// handleGroupSettings shows/updates group settings
func (b *Bot) handleGroupSettings(message *tgbotapi.Message) {
	// Get group
	group, err := store.GetGroup(b.db, message.Chat.ID)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ 群组未注册，请先使用 /register 注册 / Group not registered, please use /register first")
		b.api.Send(msg)
		return
	}

	// Parse command arguments
	args := strings.Fields(message.CommandArguments())

	// If no arguments, show current settings
	if len(args) == 0 {
		lang := messages.GetUserLanguage(group.Language, "")
		settingsMsg := b.msg.Format(lang, "group_settings", map[string]interface{}{
			"GroupName":   group.GroupName,
			"NotifyStock": formatBool(group.NotifyStock, lang),
			"NotifyPromo": formatBool(group.NotifyPromo, lang),
			"Language":    group.Language,
		})

		// Add inline keyboard for settings
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf("库存通知 Stock: %s", formatBool(group.NotifyStock, "en")),
					fmt.Sprintf("group_toggle_stock:%d", group.ID),
				),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf("促销通知 Promo: %s", formatBool(group.NotifyPromo, "en")),
					fmt.Sprintf("group_toggle_promo:%d", group.ID),
				),
			),
		)

		msg := tgbotapi.NewMessage(message.Chat.ID, settingsMsg)
		msg.ReplyMarkup = keyboard
		msg.ParseMode = "Markdown"
		b.api.Send(msg)
	}
}

// handleGroupHelp shows help for group commands
func (b *Bot) handleGroupHelp(message *tgbotapi.Message) {
	helpText := `🤖 *群组命令 / Group Commands*

/register - 注册群组接收通知 / Register group for notifications
/unregister - 取消群组注册 / Unregister group
/settings - 查看和修改群组设置 / View and modify group settings
/help - 显示此帮助信息 / Show this help message

*管理员命令 / Admin Commands*
只有群组管理员可以使用这些命令 / Only group administrators can use these commands`

	msg := tgbotapi.NewMessage(message.Chat.ID, helpText)
	msg.ParseMode = "Markdown"
	b.api.Send(msg)
}

// formatBool formats boolean value based on language
func formatBool(value bool, lang string) string {
	if value {
		if lang == "zh" {
			return "✅ 开启"
		}
		return "✅ ON"
	}
	if lang == "zh" {
		return "❌ 关闭"
	}
	return "❌ OFF"
}

// handleGroupToggle handles group setting toggles
func (b *Bot) handleGroupToggle(callback *tgbotapi.CallbackQuery) {
	// Parse callback data
	parts := strings.Split(callback.Data, ":")
	if len(parts) != 2 {
		return
	}

	var groupID uint
	fmt.Sscanf(parts[1], "%d", &groupID)

	// Get group
	var group store.Group
	if err := b.db.First(&group, groupID).Error; err != nil {
		logger.Error("Failed to get group", "error", err, "group_id", groupID)
		return
	}

	// Check if user is admin
	chatConfig := tgbotapi.ChatConfigWithUser{
		ChatID: group.TgGroupID,
		UserID: int64(callback.From.ID),
	}

	member, err := b.api.GetChatMember(tgbotapi.GetChatMemberConfig{ChatConfigWithUser: chatConfig})
	if err != nil {
		logger.Error("Failed to get chat member", "error", err)
		return
	}

	if member.Status != "administrator" && member.Status != "creator" {
		b.api.Request(tgbotapi.NewCallback(callback.ID, "只有管理员可以修改设置 / Only admins can modify settings"))
		return
	}

	// Toggle setting
	if strings.Contains(callback.Data, "stock") {
		group.NotifyStock = !group.NotifyStock
		store.UpdateGroupSettings(b.db, group.ID, group.NotifyStock, group.NotifyPromo)
	} else if strings.Contains(callback.Data, "promo") {
		group.NotifyPromo = !group.NotifyPromo
		store.UpdateGroupSettings(b.db, group.ID, group.NotifyStock, group.NotifyPromo)
	}

	// Update keyboard
	lang := messages.GetUserLanguage(group.Language, "")
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("库存通知 Stock: %s", formatBool(group.NotifyStock, lang)),
				fmt.Sprintf("group_toggle_stock:%d", group.ID),
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("促销通知 Promo: %s", formatBool(group.NotifyPromo, lang)),
				fmt.Sprintf("group_toggle_promo:%d", group.ID),
			),
		),
	)

	edit := tgbotapi.NewEditMessageReplyMarkup(callback.Message.Chat.ID, callback.Message.MessageID, keyboard)
	b.api.Send(edit)

	// Send callback response
	b.api.Request(tgbotapi.NewCallback(callback.ID, "✅ 设置已更新 / Settings updated"))
}
