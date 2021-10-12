package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"time"
)

var gBot *discordgo.Session
var ownerID = "268431730967314435" //Please change this when using my bot
var gPrefix = ">"
var commandMap = make(map[string]commandHandler)

func setupFunctionsMap() {
	commandMap["status"] = commandSendIntraStatus
}

// startBot Starts discord bot
func startBot() *discordgo.Session {
	discordBot, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	CheckError(err)
	discordBot.AddHandler(messageHandler)
	err = discordBot.Open()
	CheckError(err)
	fmt.Println("Discord bot created")
	if os.Getenv("SEGBOT_IN_PROD") == "" {
		channel, err := discordBot.UserChannelCreate(ownerID)
		if err != nil {
			return nil
		}
		hostname, _ := os.Hostname()
		_, _ = discordBot.ChannelMessageSend(channel.ID, "Bot up - "+
			time.Now().Format(time.Stamp)+" - "+hostname)
	}
	if gPrefix == "" {
		gPrefix = "!"
	}
	SetUpCloseHandler(discordBot)

	return discordBot
}

// prepFileSystem Create required directories
func prepFileSystem() error {
	err := CreateDirIfNotExist("./data")
	if err != nil {
		return err
	}
	err = CreateDirIfNotExist("./data/guilds")
	if err != nil {
		return err
	}
	err = CreateDirIfNotExist("./data/targets")
	if err != nil {
		return err
	}
	err = CreateDirIfNotExist("./data/users")
	return err
}

func main() {
	if len(os.Args) < 2 {
		_, _ = fmt.Fprintln(os.Stderr, "You must provide and env file")
		return
	}

	CheckError(prepFileSystem())
	setupFunctionsMap()
	gBot = startBot()

	for {
		time.Sleep(time.Second * 3)
	}
}
