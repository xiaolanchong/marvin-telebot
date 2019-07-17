package bot

import (
	"log"
	"time"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var dialog *Dialog

func StartTeleBot(botToken string) (*tgbotapi.BotAPI, tgbotapi.UpdatesChannel, error) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	return bot, updates, err
}

func ProcessTeleBotUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update, dataRootDir string) {
	if update.Message == nil {
		return
	}
	
	message := update.Message

	log.Printf("[%s] %s", message.From.UserName, message.Text)
	
	if(dialog == nil) {
		sender := func (msgText string) {
						msg := tgbotapi.NewMessage(message.Chat.ID, "")
						msg.ParseMode = tgbotapi.ModeMarkdown
						msg.Text = msgText
						/*a := "1"
						msg.ReplyMarkup = &tgbotapi.InlineKeyboardMarkup{
											InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{ 
																[]tgbotapi.InlineKeyboardButton{ 
																	tgbotapi.InlineKeyboardButton{Text: "1", CallbackData: &a, },
																	tgbotapi.InlineKeyboardButton{Text: "2", CallbackData: &a, },
																},
											},
										  }
						*/
						bot.Send(msg)
					}
		dialog = NewDialog( sender,
							time.Second * 30,
							message.Chat.UserName,
							dataRootDir)
	}

	if message.IsCommand() {
		args := strings.Split(message.CommandArguments(), " ")
		dialog.OnCommand(message.Command(), args)
	} else {
		dialog.OnMessage(message.Text)
	}
}