package main

import (
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

	var response translateResponse

	request := translateRequest{
		From: "ja",
		To:   "en",
		Text: text,
	}

	// TODO: Short timeout on connection failure, give up after few and cancel request if there's a new string to translate
	resp, reply, errs := gorequest.New().Post("http://127.0.0.1:3000/api/translate").
		Type("json").SendStruct(&request).EndStruct(&response)
	for _, err := range errs {
		log.WithFields(log.Fields{
			"response": resp,
			"reply":    reply,
		}).Error(err)

		return "", err
	}

	out := response.TranslationText

	log.WithFields(log.Fields{
		"time": time.Since(start),
	}).Debugf("Translated: %q", out)

	return out, nil
}
