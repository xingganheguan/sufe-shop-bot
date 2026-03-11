package store

import (
	"gorm.io/gorm"
)

// Settings keys
const (
	SettingOrderExpireMinutes = "order_expire_minutes"
	SettingOrderCleanupDays   = "order_cleanup_days"
	SettingEnableAutoExpire   = "enable_auto_expire"
	SettingEnableAutoCleanup  = "enable_auto_cleanup"
)

// GetSetting retrieves a setting by key
func GetSetting(db *gorm.DB, key string) (string, error) {
	var setting SystemSetting
	err := db.Where("key = ?", key).First(&setting).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Backward compatibility: read old key if new key not found
			if key == SettingOrderExpireMinutes {
				var oldSetting SystemSetting
				if oldErr := db.Where("key = ?", "order_expire_hours").First(&oldSetting).Error; oldErr == nil {
					return oldSetting.Value, nil
				}
			}

			// Return default values
			switch key {
			case SettingOrderExpireMinutes:
				return "10", nil
			case SettingOrderCleanupDays:
				return "7", nil
			case SettingEnableAutoExpire:
				return "true", nil
			case SettingEnableAutoCleanup:
				return "true", nil
			default:
				return "", nil
			}
		}
		return "", err
	}
	return setting.Value, nil
}

// SetSetting updates or creates a setting
func SetSetting(db *gorm.DB, key, value, description, settingType string) error {
	var setting SystemSetting
	err := db.Where("key = ?", key).First(&setting).Error

	if err == gorm.ErrRecordNotFound {
		// Create new setting
		setting = SystemSetting{
			Key:         key,
			Value:       value,
			Description: description,
			Type:        settingType,
		}
		return db.Create(&setting).Error
	}

	// Update existing setting
	return db.Model(&setting).Updates(map[string]interface{}{
		"value":       value,
		"description": description,
		"type":        settingType,
	}).Error
}

// GetAllSettings retrieves all settings
func GetAllSettings(db *gorm.DB) ([]SystemSetting, error) {
	var settings []SystemSetting
	err := db.Order("key").Find(&settings).Error
	return settings, err
}

// InitializeSettings creates default settings if they don't exist
func InitializeSettings(db *gorm.DB) error {
	// Migrate old key to new key if needed
	var oldSetting SystemSetting
	if err := db.Where("key = ?", "order_expire_hours").First(&oldSetting).Error; err == nil {
		var newSetting SystemSetting
		if err := db.Where("key = ?", SettingOrderExpireMinutes).First(&newSetting).Error; err == gorm.ErrRecordNotFound {
			if err := db.Create(&SystemSetting{
				Key:         SettingOrderExpireMinutes,
				Value:       oldSetting.Value,
				Description: "订单过期时间（分钟）",
				Type:        "int",
			}).Error; err != nil {
				return err
			}
		}
	}

	defaultSettings := []struct {
		Key         string
		Value       string
		Description string
		Type        string
	}{
		{
			Key:         SettingOrderExpireMinutes,
			Value:       "10",
			Description: "订单过期时间（分钟）",
			Type:        "int",
		},
		{
			Key:         SettingOrderCleanupDays,
			Value:       "7",
			Description: "清理过期订单的天数",
			Type:        "int",
		},
		{
			Key:         SettingEnableAutoExpire,
			Value:       "true",
			Description: "启用订单自动过期",
			Type:        "bool",
		},
		{
			Key:         SettingEnableAutoCleanup,
			Value:       "true",
			Description: "启用过期订单自动清理",
			Type:        "bool",
		},
	}

	for _, s := range defaultSettings {
		var existing SystemSetting
		err := db.Where("key = ?", s.Key).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&SystemSetting{
				Key:         s.Key,
				Value:       s.Value,
				Description: s.Description,
				Type:        s.Type,
			}).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// SettingsMap is a helper to get all settings as a map
func GetSettingsMap(db *gorm.DB) (map[string]string, error) {
	settings, err := GetAllSettings(db)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, s := range settings {
		result[s.Key] = s.Value
	}

	// Add defaults for missing settings
	if _, ok := result[SettingOrderExpireMinutes]; !ok {
		result[SettingOrderExpireMinutes] = "10"
	}
	if _, ok := result[SettingOrderCleanupDays]; !ok {
		result[SettingOrderCleanupDays] = "7"
	}
	if _, ok := result[SettingEnableAutoExpire]; !ok {
		result[SettingEnableAutoExpire] = "true"
	}
	if _, ok := result[SettingEnableAutoCleanup]; !ok {
		result[SettingEnableAutoCleanup] = "true"
	}

	return result, nil
}
