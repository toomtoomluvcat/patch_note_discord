package service

import (
	"dockertest/connect"
	"dockertest/schema"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

func CreateAccount(i *discordgo.InteractionCreate, inputPassKey string) string {

	err := connect.DB.Transaction(func(tx *gorm.DB) error {

		var passKey schema.PassKey
		if err := tx.
			Where("key = ? AND is_used = false", inputPassKey).
			First(&passKey).Error; err != nil {
			return err
		}

		passKey.DiscordID = i.Member.User.ID
		passKey.IsUsed = true
		passKey.Message = i.Member.User.ID

		if err := tx.Save(&passKey).Error; err != nil {
			return err
		}

		if err := tx.Model(&schema.Account{}).Create(map[string]interface{}{"DiscordID": i.Member.User.ID, "Username": i.Member.User.Username}).Error; err != nil {
			return err //"ดูเหมือนคุณสมัครไว้แล้วนะ  \nลองตรวจสอบอีกทีได้ไหม? 🔍"
		}
		return nil
	})

	if err != nil {
		return " ดูเหมือนจะมีข้อผิดพลาดเกิดขึ้น📋 \nลองติดต่อแอดมินได้ไหม😵‍💫"
	}

	return "ลงทะเบียนเรียบร้อยแล้วนะ 📝 \nพร้อมลุยรึยัง? 🔥"

}

func CheckRegister(userID string) error {
	var Account schema.Account
	if err := connect.DB.Where("discord_id = ?", userID).First(&Account).Error; err != nil {
		return err
	}
	return nil
}

func GetUserID(i *discordgo.InteractionCreate) string {
	if i.Member != nil && i.Member.User != nil {
		return i.Member.User.ID
	}
	return i.User.ID
}

func GetUsername(i *discordgo.InteractionCreate) string {
	if i.Member != nil && i.Member.User != nil {
		return i.Member.User.Username
	}
	return i.User.Username
}
