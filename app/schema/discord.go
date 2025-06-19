package schema

import "time"

type Account struct {
	DiscordID string `gorm:"primaryKey"`
	Username  string
	Calendars []Calendar
}

type Calendar struct {
	ID          uint `gorm:"primaryKey"`
	Title       string
	DueDateTime time.Time
	Notify_1    bool `gorm:"default:false"`
	Notify_2    bool `gorm:"default:false"`
	AccountID   string
	Account     Account `gorm:"foreignKey:AccountID" json:"-"`
}

type PassKey struct {
	Message   string
	DiscordID string
	Key       string `gorm:"primaryKey"`
	IsUsed    bool
}
