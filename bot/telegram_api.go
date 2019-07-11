package bot

import (
	"encoding/json"
	"fmt"
)

const (
	main_bot_url = "https://api.telegram.org/bot%s/"
)

type Message struct {
	MessageID int `json:"message_id"`
	From      struct {
		ID        int    `json:"id"`
		IsBot     bool   `json:"is_bot"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Username  string `json:"username"`
	} `json:"from"`
	Chat struct {
		ID        int    `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Username  string `json:"username"`
		Type      string `json:"type"`
	} `json:"chat"`
	Date     int    `json:"date"`
	Text     string `json:"text"`
	Entities []struct {
		Offset int    `json:"offset"`
		Length int    `json:"length"`
		Type   string `json:"type"`
	} `json:"entities"`
}

type UpdateMessage struct {
	UpdateID int `json:"update_id"`
	Message  *Message `json:"message,omitempty"`
}

type UpdateResponse struct {
	Ok     bool `json:"ok"`
	Result []UpdateMessage `json:"result"`
}

func DecodeUpdate(msg_bytes []byte) ([]UpdateMessage, error) {
	res := UpdateResponse{}
	err := json.Unmarshal(msg_bytes, &res)
	return res.Result, err
}

func GetBotUrl(bot_token string) string {
	return fmt.Sprintf(main_bot_url, bot_token)
}

func GetUpdateMsgUrl(bot_token string, offset int) string {
	get_update_url := main_bot_url + "getUpdates?offset=%s"
	return fmt.Sprintf(get_update_url, bot_token, offset)
}