package main

import (
	"context"
	"log"
	"regexp"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"

	"github.com/prdsrm/std/messages"
	"github.com/prdsrm/std/session/postgres"
)

func echo(ctx messages.MonitoringContext) error {
	log.Println("Message: ", ctx.GetMessage().Message)
	return nil
}

func listen(ctx context.Context, client *telegram.Client, dispatcher tg.UpdateDispatcher, options telegram.Options) error {
	user, err := client.Self(ctx)
	if err != nil {
		return err
	}
	log.Println("Connected to user account: ", user.ID, user.FirstName)
	monitoring := messages.NewMonitoring(dispatcher, 777000, false)
	monitoring.Handle(regexp.MustCompile(".*"), echo)
	err = monitoring.Listen()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	db, err := postgres.OpenDBConnection()
	if err != nil {
		log.Fatalln("can't connect to database: ", err)
	}
	bot, err := postgres.GetBotByUserID(db, 0)
	if err != nil {
		log.Fatalln("can't get bot: ", err)
	}
	log.Println("Bot from the db: ", bot.UserID)
	err = postgres.ConnectToBotFromDatabase(db, bot, listen)
	if err != nil {
		log.Fatalln("can't connect: ", err)
	}
}
