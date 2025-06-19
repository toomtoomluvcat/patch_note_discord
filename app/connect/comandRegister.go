package connect

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func SetCommands(dg *discordgo.Session) {
	log.Println("🔧 Start setting commands...")

	commands := []discordgo.ApplicationCommand{
		{
			Name:        "hi",
			Description: "Start talking with Haru!",
		},
		{
			Name:        "note",
			Description: "Let Create Your note",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "content",
					Description: "Your note content",
					Required:    true,
				},
			},
		},
		{
			Name:        "find",
			Description: "Find your task with period time you want",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "period",
					Description: "Your task period",
					Required:    true,
				},
			},
		},
		{
			Name:        "register",
			Description: "Input your password Key",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "key",
					Description: "Add your key. You can get key from admin",
					Required:    true,
				},
			},
		},
		{
			Type:        discordgo.ChatApplicationCommand,
			Name:        "task",
			Description: "Select your period time to show task",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "period",
					Description: "Select time period",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{Name: "Day", Value: 1},
						{Name: "Week", Value: 7},
						{Name: "Month", Value: 30},
					},
				},
			},
		},
		{
			Type:        discordgo.ChatApplicationCommand,
			Name:        "deltask",
			Description: "Enter your TaskID to delete permanent",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "taskid",
					Description: "Enter TaskID in your list",
					Required:    true,
				},
			},
		},
		{
			Type:        discordgo.ChatApplicationCommand,
			Name:        "help",
			Description: "Check PacthNote Comand",
		},
	}

	log.Println("🔍 Fetching bot application ID...")
	app, err := dg.Application("@me")
	if err != nil {
		log.Fatalf("❌ Failed to get bot application: %v", err)
	}

	log.Printf("✅ Got Application ID: %s\n", app.ID)

	for _, v := range commands {
		log.Printf("📦 Preparing to register command: %s\n", v.Name)

		// Optional: Remove deprecated commands
		if v.Name == "ping" || v.Name == "sum" {
			log.Printf("🗑️ Attempting to delete old command: %s\n", v.Name)
			err := dg.ApplicationCommandDelete(app.ID, "", v.ID)
			if err != nil {
				log.Println("❌ Failed to delete command:", v.Name, err)
			} else {
				log.Println("✅ Deleted command:", v.Name)
			}
		}

		log.Printf("🚀 Creating command: %s\n", v.Name)
		createdCmd, err := dg.ApplicationCommandCreate(app.ID, "", &v)
		if err != nil {
			log.Printf("❌ Cannot create command %s: %v\n", v.Name, err)
		} else {
			log.Printf("✅ Created command: %s with ID: %s\n", createdCmd.Name, createdCmd.ID)
		}
	}

	log.Println("🎉 Finished setting commands.")
}
