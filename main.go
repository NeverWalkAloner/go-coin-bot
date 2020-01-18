package main

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	tb "gopkg.in/tucnak/telebot.v2"
)

const readyStage int = 1

var tgToken string

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error when loading env variables")
	}
	tgToken = os.Getenv("TGTOKEN")
}

func main() {
	b, err := tb.NewBot(tb.Settings{
		Token:  tgToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	inlineBtn1 := tb.InlineButton{
		Unique: tails,
		Text:   tails,
	}

	inlineBtn2 := tb.InlineButton{
		Unique: heads,
		Text:   heads,
	}

	inlineKeys := [][]tb.InlineButton{
		[]tb.InlineButton{inlineBtn1, inlineBtn2},
	}

	b.Handle(&inlineBtn1, func(c *tb.Callback) {
		b.Respond(c, &tb.CallbackResponse{
			ShowAlert: false,
		})
		currentState := getFlipStage(c.Sender.ID)
		if currentState.stage != 2 {
			b.Send(c.Sender, "Start protocol by sending /flip command")
		} else {
			b.Send(c.Sender, getResultMessage(tails, currentState))
			b.Send(c.Sender, "Check this server seed to validate hash: "+currentState.serverSeed)
		}
	})

	b.Handle(&inlineBtn2, func(c *tb.Callback) {
		b.Respond(c, &tb.CallbackResponse{
			ShowAlert: false,
		})
		currentState := getFlipStage(c.Sender.ID)
		if currentState.stage != 2 {
			b.Send(c.Sender, "Start protocol by sending /flip command")
		} else {
			b.Send(c.Sender, getResultMessage(heads, currentState))
			b.Send(c.Sender, "Check this server seed to validate hash: "+currentState.serverSeed)
		}
	})

	b.Handle("/flip", func(m *tb.Message) {
		flip := Flip{
			userID:     m.Sender.ID,
			stage:      readyStage,
			userSeed:   "",
			serverSeed: "",
			coinSide:   "",
		}
		createUpdateFlip(flip)
		b.Send(m.Sender, "Ok, now send me please a random string")
	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		currentState := getFlipStage(m.Sender.ID)
		switch currentStage := currentState.stage; currentStage {
		case 0:
			// protocol has not been started yet
			b.Send(m.Sender, "Start protocol by sending /flip command")
		case 1:
			// first step, user should send random seed to the server to initiate coin flip
			flipResult := flipCoin()
			serverSeed := generateSeed()
			flip := Flip{
				userID:     m.Sender.ID,
				stage:      currentStage + 1,
				userSeed:   m.Text,
				serverSeed: serverSeed,
				coinSide:   flipResult,
			}
			createUpdateFlip(flip)
			commitString := fmt.Sprintf("%s %s %s", flipResult, serverSeed, m.Text)
			b.Send(m.Sender, "Done! Hash of the result: "+getFlipHash(commitString))
			b.Send(m.Sender, "You can check it later by computing SHA256(<flip result> <server seed> <your random string>)")
			b.Send(
				m.Sender,
				"Tails or heads?",
				&tb.ReplyMarkup{InlineKeyboard: inlineKeys})
		}
	})

	b.Start()
}

// get final response text based on user's choice and coin flip result
func getResultMessage(userChoice string, currentState Flip) string {
	if userChoice == currentState.coinSide {
		return "You are right! Result is " + currentState.coinSide
	} else if userChoice != currentState.coinSide {
		return "You are wrong! Result is " + currentState.coinSide
	}
	return "Something went wrong. Try again."
}
