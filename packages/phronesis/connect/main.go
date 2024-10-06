package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
	"github.com/jmoiron/sqlx"

	"github.com/prdsrm/std/messages"
	"github.com/prdsrm/std/session/postgres"
)

var db *sqlx.DB
var data = make(map[int64]string)

func extractCodeFromText(text string) (string, error) {
	var code string
	var codeRegex = regexp.MustCompile(`(?m).*: (\d{5})`)
	for _, match := range codeRegex.FindAllStringSubmatch(text, -1) {
		code = match[1]
		_, err := strconv.Atoi(code)
		if err != nil {
			return "", fmt.Errorf("couldn't extract OTP code from text: %s", text)
		}
		return code, nil
	}
	if code == "" {
		var webCodeRegex = regexp.MustCompile(`(?m):\n([a-zA-Z0-9-_-]{11})`)
		for _, match := range webCodeRegex.FindAllStringSubmatch(text, -1) {
			code = match[1]
			return code, nil
		}
	}
	return "", errors.New("no code found.")
}

func echo(ctx messages.MonitoringContext) error {
	msg := ctx.GetMessage()
	log.Println("received: ", msg.ID, msg.Message)
	user, err := ctx.GetClient().Self(ctx.Ctx)
	if err != nil {
		return err
	}
	code, err := extractCodeFromText(msg.Message)
	if err != nil {
		return err
	}
	data[user.ID] = code
	return messages.EndConversation
}

func listen(ctx context.Context, client *telegram.Client, dispatcher tg.UpdateDispatcher, options telegram.Options) error {
	user, err := client.Self(ctx)
	if err != nil {
		return err
	}
	errChan := make(chan error)
	go func() {
		log.Println("Connected to user account: ", user.ID, user.FirstName)
		monitoring := messages.NewMonitoring(dispatcher, 777000, false)
		monitoring.Handle(regexp.MustCompile(".*"), echo)
		err = monitoring.Listen(ctx, client)
		if err != nil {
			errChan <- err
		}
		errChan <- nil
	}()
	select {
	case err := <-errChan:
		if err != nil {
			return err
		}
		log.Println("successfully received Telegram message with bot: ", user.ID)
	case <-time.After(120 * time.Second):
		return errors.New("connected, but 120 seconds timeout while waiting for a Telegram message")
	}
	return nil
}

type ConnectToBotRequest struct {
	UserID int64 `json:"user_id"`
}

type ConnectToBotOutput struct {
	Message string `json:"message"`
}

// Request takes in the user ID.
type Request struct {
	// User ID of the telegram account you want to connect to.
	// It must be present in the database.
	UserID int64 `json:"user_id"`
}

// Response returns back the http code, type of data, and the received code.
type Response struct {
	// StatusCode is the http code that will be returned back to the user.
	StatusCode int `json:"statusCode,omitempty"`
	// Headers is the information about the type of data being returned back.
	Headers map[string]string `json:"headers,omitempty"`
	// Body will contain the received code.
	Body string `json:"body,omitempty"`
}

func Main(in Request) (*Response, error) {
	conn, err := postgres.OpenDBConnection()
	if err != nil {
		log.Fatalln("can't connect to database: ", err)
	}
	db = conn

	bot, err := postgres.GetBotByUserID(db, in.UserID)
	if err != nil {
		return &Response{StatusCode: http.StatusNotFound}, err
	}
	log.Println("Bot from the db: ", bot.UserID)
	err = postgres.ConnectToBotFromDatabase(db, bot, listen)
	if err != nil {
		return &Response{StatusCode: http.StatusInternalServerError}, err
	}
	text := data[bot.UserID]
	return &Response{StatusCode: 200, Body: text}, nil
}
