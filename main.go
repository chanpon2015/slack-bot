package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(err)
		return
	}
	token, ok := os.LookupEnv("SLACK_TOKEN")
	if !ok {
		log.Println("Missing SLACK_TOKEN in environment")
	}
	channelID, ok := os.LookupEnv("MAIN_CHANNEL_ID")
	if !ok {
		log.Println("Missing MAIN_CHANNEL_ID in environment")
		return
	}
	api := slack.New(
		token,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Llongfile|log.LstdFlags)),
	)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		fmt.Print("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.ConnectedEvent:
			fmt.Println("connected", ev.Info, ev.ConnectionCount)
			rtm.SendMessage(rtm.NewOutgoingMessage("test", channelID))
		case *slack.MessageEvent:
			fmt.Println("message")
		case *slack.ReactionAddedEvent:
			fmt.Println("reaction added")
			rtm.SendMessage(rtm.NewOutgoingMessage("test", channelID))
		default:
			fmt.Println(ev)
		}
	}
}
