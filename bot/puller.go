
package bot

import (
	"fmt"
	"time"
	"net/http"
	"io/ioutil"
)

type UpdatesChannel chan *Message

func (ch UpdatesChannel) pull(bot_token string, offset *int) {
	//fmt.Printf("URL ID %+v\n", *offset)
	url := GetUpdateMsgUrl(bot_token, *offset)
	if resp, err := http.Get(url); err == nil {
		defer resp.Body.Close()
		if body_bytes, err_body := ioutil.ReadAll(resp.Body); err_body == nil {
			//fmt.Printf("Received a new message, %v bytes", len(body_bytes))
			update_result, err_dec := DecodeUpdate(body_bytes)
			if err_dec != nil {
				fmt.Printf("Error on message parsing: %+v\n", string(body_bytes))
			} else {
				for _, update := range(update_result) {
					if(update.UpdateID < *offset) {  // kludge, why Telegram sends old messages?
						continue
					}
					//fmt.Printf("ID %+v\n", *offset)
					*offset = update.UpdateID + 1
					//handler(&update)
					ch <- update.Message
				}
			}
		}
	}
}

type Puller struct {
	Ticker		*time.Ticker
	UpdatesChannel 	UpdatesChannel
}

func New(milliseconds int, bot_token string) (*Puller, error) {
	puller := Puller{
		Ticker:     time.NewTicker(time.Duration(milliseconds) * time.Millisecond),
		UpdatesChannel: make(UpdatesChannel, 100),
	}
	var offset *int = new(int)
	*offset = 0
	go func() {
		for range puller.Ticker.C {
			puller.UpdatesChannel.pull(bot_token, offset)
		}
	}()
	return &puller, nil
}

func (puller *Puller) Close() {
	puller.Ticker.Stop()
}

