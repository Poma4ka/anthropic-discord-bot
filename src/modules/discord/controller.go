package discord

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Controller struct {
	Service *Service
}

func (c *Controller) messageCreate(client *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == client.State.User.ID {
		return
	}

	if m.GuildID == "" {
		return
	}

	if !strings.Contains(m.Content, client.State.User.Mention()) && (m.ReferencedMessage == nil || m.ReferencedMessage.Author.ID != client.State.User.ID) {
		return
	}

	reply, err := c.Service.MessageCreate(client, m.Message)
	if err != nil {
		c.Service.logger.Error("Error create message", err)
		_, err := editReplyOrReply(client, reply, m.Message, "Что-то я затупил, может быть пора отдохнуть 😞")
		if err != nil {
			c.Service.logger.Error("Error send error message", err)
			return
		}
	}
}
