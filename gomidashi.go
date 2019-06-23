package functions

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// PubSubMessage : published message from Cloud Pub/Sub
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// SubscribedMessage : subscribed message from Cloud Pub/Sub
// Mention      :: ["channel", "here"] (empty: garbage type only)
// Channel 		:: specify the channel to send messages
// Area			:: ["north", "south", "east", "west"]
type SubscribedMessage struct {
	Mention string `json:"mention,omitempty"`
	Channel string `json:"channel"`
	Area    string `json:"area"`
}

// DecodeMessage : decoding published message from Cloud Pub/Sub
func (msg PubSubMessage) DecodeMessage() (msgData SubscribedMessage, err error) {
	if err = json.Unmarshal(msg.Data, &msgData); err != nil {
		log.Printf("Message[%v] ... Could not decode subscribing data: %v", msg, err)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d", e.Offset)
		}
		return
	}
	return
}

// Gomidashi : entry point
func Gomidashi(ctx context.Context, m PubSubMessage) error {
	msg, err := m.DecodeMessage()
	if err != nil {
		log.Fatal(err)
		return err
	}

	filename := "dist/2019" + msg.Area + ".json"
	text, err := createText(filename, msg.Mention)
	if err != nil {
		log.Fatal(err)
		return err
	}

	var webhookURL = os.Getenv("SLACK_WEBHOOK_URL")
	err = postMessage("gomidashi", text, msg.Channel, webhookURL)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func postMessage(name string, msg string, channel string, webhookURL string) error {
	jsonStr := `{"channel":"` + channel + `","username":"` + name + `","text":"` + msg + `"}`

	req, err := http.NewRequest(
		"POST",
		webhookURL,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}

func createText(filename string, mention string) (string, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return "err", err
	}

	// store date and garbage type data
	var cal map[string]string

	if err := json.Unmarshal(bytes, &cal); err != nil {
		return "err", err
	}

	const format = "2006-01-02"
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	today := time.Now().UTC().In(jst)
	date := today.Format(format)
	log.Printf("today : %v\n", today)
	log.Printf("idx   : %v\n", date)

	var res string
	if mention == "here" || mention == "channel" {
		res = "<!" + mention + "> " + date + "\n" + cal[date]
	} else {
		res = date + "\n" + cal[date]
	}

	return res, nil
}
