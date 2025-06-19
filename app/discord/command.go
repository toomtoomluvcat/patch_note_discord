package discord_service

import (
	"dockertest/process"
	"dockertest/service"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func ReadMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if s.State.User.ID == m.Author.ID {
		return
	}
	fmt.Printf("Message from %s:%s\n", m.Author.Username, m.Content)
	if strings.Contains(m.Content, "โง่") {
		ResponseMessage(s, m.ChannelID, "Why You So Rude (；￣Д￣)")
		s.MessageReactionAdd(m.ChannelID, m.ID, "😠")
	} else if strings.Contains(m.Content, "เทพ") {
		s.MessageReactionAdd(m.ChannelID, m.ID, "🔥")
	} else if strings.Contains(m.Content, "เต้") {
		ResponseMessage(s, m.ChannelID, "ไม่ไหว แล้วพี่จ๋า จะแตกแล้ววววว")
	} else if strings.Contains(m.Content, "ตูม") {
		s.ChannelMessageDelete(m.ChannelID, m.ID)
		ResponseMessage(s, m.ChannelID, "อย่าให้มีครั้งที่สอง")
	}
}

func ResponseMessage(s *discordgo.Session, ChannelID string, message string) {
	_, err := s.ChannelMessageSend(ChannelID, message)
	if err != nil {
		fmt.Println("Fail to send message")
	}
}
func HandleComands(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		data := i.ApplicationCommandData()
		switch data.Name {

		case "hi":
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "PactchNote 📝✨เป็น AI ผู้ช่วยบันทึกโน้ตของเธอเอง! แค่เล่าให้เราฟัง ไม่ว่าเรื่องอะไร เราก็จัดให้ได้หมดเลย~ \nถ้าเธออยากทำอะไรก็ /help ได้เลยนะ🎤🧠",
				},
			})
		case "note":

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "กำลังเขียนโน๊ตใหม่...",
				},
			})
			var content string
			for _, opt := range data.Options {
				if opt.Name == "content" {
					content = opt.StringValue()
				}
			}
			response := process.GenAI(content, i)

			s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
				Content: response,
			})
		case "register":
			var passkey string
			for _, p := range data.Options {
				if p.Name == "key" {
					passkey = p.StringValue()
				}
			}

			Message := service.CreateAccount(i, passkey)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: Message,
				},
			})
		case "task":
			var period int64
			for _, p := range data.Options {
				if p.Name == "period" {
					period = p.IntValue()
				}
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			})
			var msg string
			response := service.GetTaskByPeriod(service.GetUserID(i), period)
			if len(response) > 2000 {
				msg = "Discord ไม่รองรับข้อมูลขนาดใหญ่ กรุณาติดต่อแอดมิน"
			} else {
				msg = response
			}
			s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
				Content: msg,
			})
		case "deltask":
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			})
			var taskID int64
			for _, t := range data.Options {
				if t.Name == "taskid" {
					taskID = t.IntValue()
				}
			}
			s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
				Content: service.DeleteTaskByID(service.GetUserID(i), uint(taskID)),
			})

		case "find":
			var periodTime string
			for _, p := range data.Options {
				if p.Name == "period" {
					periodTime = p.StringValue()
				}
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			})
			s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
				Content: process.GenAISearchTimeRange(periodTime, i),
			})
		case "help":
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "🌸 **PactchNote** เราเป็น AI ที่จะช่วยเธอจดโน๊ตให้ง่ายที่สุดเลยล่ะ~ (๑˃̵ᴗ˂̵)و 💬\n━━━━━━━━━━━━━━\n📚 **ลิสต์สิ่งที่เราทำได้ตามนี้เลย~**\n\n🔑 **/register**\nใส่ Passkey ที่เธอได้มา แล้วเราจะได้เป็นเพื่อนช่วยงานกันน้า~ (｡•̀ᴗ-)✧\n\n📝 **/note**\nบอกแค่งานกับวันเวลาก็พอ~ เดี๋ยวเราจัดให้อย่างเป๊ะ! (≧◡≦) ♡\n\n🔍 **/find**\nตามด้วยช่วงเวลา~ เราจะช่วยหางานของเธอได้ทันใจเลย! (╹◡╹)ノ💕\n\n🗑️ **/deltask**\nใส่ ID ของงาน แล้วเราจะลบให้อย่างเรียบร้อย~ 🧼(｡•̀ᴗ-)✧\n━━━━━━━━━━━━━━\n⏰ และแน่นอน~\nเราจะแจ้งเตือนก่อนงานเริ่ม 30 นาทีด้วยนะ! (•̀ᴗ•́)و ̑̑ 🔔",
				},
			})
		}
	}
}
