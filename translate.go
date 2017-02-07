package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/parnurzeal/gorequest"
	"time"
)

type translateRequest struct {
	Text string `json:"text"`
	From string `json:"from"`
	To   string `json:"to"`
}

type translateResponse struct {
	Text            string `json:"text"`
	From            string `json:"from"`
	To              string `json:"to"`
	TranslationText string `json:"translationText"`
}

func translateString(text string) (string, error) {
	start := time.Now()
	jsonStruct := translateRequest{
		From: "ja",
		To:   "en",
		Text: text,
	}

	// Convert json object to string
	jsonString, err := json.Marshal(jsonStruct)
	if err != nil {
		log.Error("Failed to marshal JSON API request", err.Error())
	}

	// Post the request
	resp, reply, errs := gorequest.New().Post("http://127.0.0.1:3000/api/translate").Send(string(jsonString)).EndBytes()
	// TODO: Short timeout on connection failure, give up after few and cancel request if there's a new string to translate
	for _, err := range errs {
		log.WithFields(log.Fields{
			"response": resp,
			"reply":    reply,
		}).Error(err.Error())
		return "", err
	}

	var response translateResponse

	if err := json.Unmarshal(reply, &response); err != nil {
		log.Error("Failed to unmarshal JSON API response", err.Error())
		return "", err
	}

	var out string
	out = response.TranslationText

	log.WithFields(log.Fields{
		"time": time.Since(start),
	}).Debugf("Translated: %q", out)

	return out, nil
}
