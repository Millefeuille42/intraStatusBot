package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

// sendMessageWithMention Sends a discord message according to user params
func sendMessageToUser(message, channel, userID, chanParam string, agent discordAgent) {
	switch chanParam {
	case "none":
		_, _ = agent.session.ChannelMessageSend(channel, message)
	case "channel":
		_, _ = agent.session.ChannelMessageSend(channel, fmt.Sprintf("<@%s>\n%s", userID, message))
	case "dm":
		dmChan, err := agent.session.UserChannelCreate(userID)
		if err == nil {
			_, _ = agent.session.ChannelMessageSend(dmChan.ID, message)
		}
		_, _ = agent.session.ChannelMessageSend(channel, message)
	case "all":
		dmChan, err := agent.session.UserChannelCreate(userID)
		if err == nil {
			_, _ = agent.session.ChannelMessageSend(dmChan.ID, message)
		}
		_, _ = agent.session.ChannelMessageSend(channel, fmt.Sprintf("<@%s>\n%s", userID, message))
	}
}

// sendMessageWithMention Sends a discord message prepending a mention + \n to the message, if id == "", id becomes the message author
func sendMessageWithMention(message, id string, agent discordAgent) {
	if id == "" {
		id = agent.message.Author.ID
	}

	if agent.channel == "" {
		agent.channel = agent.message.ChannelID
	}
	_, err := agent.session.ChannelMessageSend(agent.channel, fmt.Sprintf("<@%s>\n%s", id, message))
	if err != nil {
		LogError(err)
	}
}

// getUser Returns associated user of provided id
func getUser(session *discordgo.Session, id string) string {
	ret, err := session.User(id)
	if err != nil {
		return ""
	}
	return ret.Username
}

// getChannelName Returns associated channel name of provided id
func getChannelName(session *discordgo.Session, id string) string {
	ret, _ := session.Channel(id)
	return ret.Name
}

// logErrorToChan Sends plain error to command channel
func logErrorToChan(agent discordAgent, err error) {
	if err == nil {
		return
	}
	LogError(err)
	_, _ = agent.session.ChannelMessageSend(agent.channel,
		fmt.Sprintf("An Error Occured, Please Try Again Later {%s}", err.Error()))
}
