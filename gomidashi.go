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

type PubSubMessage struct {
	Data string `json:"data"`
}

type Calendar struct {
	Date string `json:"date"`
	Type string `json:"type"`
}

func Gomidashi(ctx context.Context, m PubSubMessage) error {
	log.Printf("message: %v", m)
	var webhookUrl string = os.Getenv("SLACK_WEBHOOK_URL")
	name := "gomidashi"
	channel := "gomi"

	text, err := createText()
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = postMessage(name, text, channel, webhookUrl)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func postMessage(name string, msg string, channel string, webhookUrl string) error {
	jsonStr := `{"channel":"` + channel + `","username":"` + name + `","text":"` + msg + `"}`

	req, err := http.NewRequest(
		"POST",
		webhookUrl,
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

	// log.Println(resp)
	defer resp.Body.Close()
	return nil
}

func createText() (string, error) {
	bytes, err := ioutil.ReadFile("garbage.json")
	if err != nil {
		return "err", err
	}
	var cal []Calendar
	if err := json.Unmarshal(bytes, &cal); err != nil {
		return "err", err
	}

	const format = "2006-01-02 00:00:00 +0900"
	today := time.Now()
	begin, _ := time.Parse(format, cal[0].Date)
	diff := today.Sub(begin)
	idx := (int)(diff.Hours() / 24)
	log.Printf("today : %v\n", today)
	log.Printf("cal[0]: %v\n", cal[0].Date)
	log.Printf("begin : %v\n", begin)
	log.Printf("idx   : %v\n", idx)

	res := "<!channel> " + today.Format(format) + "\n" + cal[idx].Type

	return res, nil
}
