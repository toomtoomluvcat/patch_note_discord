package service

import (
	"dockertest/connect"
	"dockertest/format"
	"dockertest/schema"

	"github.com/gin-gonic/gin"
)

func Cronnotified30minute(c *gin.Context) {
	var req []schema.Calendar
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON payload" + err.Error()})
		return
	}
	for _, r := range req {
		if r.AccountID == "" {
			c.JSON(400, gin.H{"error": "AccountID is empty"})
			return
		}
		channel, err := connect.DG.UserChannelCreate(r.AccountID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		if channel == nil || channel.ID == "" {
			c.JSON(500, gin.H{"error": "Failed to create DM channel or channel ID empty"})
			return
		}

		_, err = connect.DG.ChannelMessageSend(channel.ID, format.NotifyTaskDueSoon(r.Title, r.DueDateTime, r.ID))
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}
}
