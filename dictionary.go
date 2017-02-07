package main

import (
	"compress/gzip"
	"encoding/xml"
	log "github.com/Sirupsen/logrus"
	"io"
	"os"
	"time"
)

type dictionaryTerm struct {
	ID             int    `xml:"id,attr"`
	Type           string `xml:"type,attr"`
	Disabled       bool   `xml:"disabled,attr"`
	GameID         int    `xml:"gameId"`
	Regex          bool   `xml:"regex"`
	SourceLanguage string `xml:"sourceLanguage"`
	Language       string `xml:"language"`
	Pattern        string `xml:"pattern"`
	Text           string `xml:"text"`
}

var (
	inputReplacement  = make(map[string]string)
	transReplacement  = make(map[string]string)
	outputReplacement = make(map[string]string)
)

func loadSharedDictionary() {
	log.Infof("Loading shared VNR dictionary")

	start := time.Now()

	file, err := os.Open("./gamedic.xml.gz")
	if err != nil {
		log.Errorln("Error opening file:", err)
		return
	}
	defer file.Close()

	gz, err := gzip.NewReader(file)
	if err != nil {
		log.Fatal(err)
	}
	defer gz.Close()

	decoder := xml.NewDecoder(gz)

	total := 0
	added := 0

	// Read tokens from the XML document in a stream.
	for {
		t, err := decoder.Token()
		if err == io.EOF {
			break
		}

		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "term" {
				var t dictionaryTerm
				err := decoder.DecodeElement(&t, &se)
				if err != nil {
					log.Warning("xml element decoder:", err)
					continue
				}

				if parseTerm(&t) {
					added++
				}
				total++
			}
		}

		//		if added >= 10 {
		//			break
		//		}
	}

	log.WithFields(log.Fields{
		"time": time.Since(start),
	}).Infof("Added %d out of %d entries~", added, total)

	log.Debugf("Input: %d, Trans: %d, Output: %d", len(inputReplacement), len(transReplacement), len(outputReplacement))
}

func parseTerm(t *dictionaryTerm) bool {
	if t.Disabled || t.GameID > 0 ||
		// t.SourceLanguage != "ja" || // Always seems to be japanese
		(t.Language != "ja" && t.Language != "en") {
		return false
	}

	if t.Regex {
		// Need to implement macros first and only few generic ones actually use those
		return false
	}

	switch t.Type {
	case "input":
		//		log.Debugf("%+v", t)
		inputReplacement[t.Pattern] = t.Text
	case "output":
		//		log.Debugf("%+v", t)
		outputReplacement[t.Pattern] = t.Text
	case "trans":
		//		log.Debugf("%+v", t)
		transReplacement[t.Pattern] = t.Text
	default:
		// game, ocr, yomi, proxy, tts, name, prefix, suffix, name, macro
		return false
	}

	return true
}
