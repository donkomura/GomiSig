package format

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"log"
)

type Calendar struct {
	Date string `json:"date"`
	Type string `json:"type"`
}

type Rawdata struct {
	Data string `json:data`
}

func main {
	filename := os.Args[1]
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	var data []Rawdata
	if err := json.Unmarshal(bytes, &data); err != nil {
		log.Fatal(err)
	}

	for _, d := range data {
		slice := strings.Split(d, ":")
		
	}
}
