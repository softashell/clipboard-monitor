package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
	"os/signal"
)

var (
	inChan   = make(chan string)
	outChan  = make(chan string)
	doneChan = make(chan bool)
)

func main() {
	app := cli.NewApp()

	app.Description = "Automatically translates text in clipboard"

	app.HideVersion = true

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose, v",
			Usage: "Spam debug messages",
		},
		cli.BoolFlag{
			Name:  "vnr",
			Usage: "Load shared VNR dictionary",
		},
	}

	app.Action = func(c *cli.Context) error {
		startMonitoring(c)

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func startMonitoring(c *cli.Context) {
	if c.Bool("verbose") {
		log.SetLevel(log.DebugLevel)
	}

	if c.Bool("vnr") {
		loadSharedDictionary()
	}

	log.Info("Starting system clipboard monitoring~")
	go monitorCliboard()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	for {
		select {
		case sig := <-sigChan:
			if sig.String() == "interrupt" {
				log.Debug("Closing done channel.")
				close(doneChan)
			}
		case text := <-outChan:
			processInput(text)

		case <-doneChan:
			log.Info("Cleaning up~")
			os.Exit(1)
		}
	}
}

func processInput(text string) {
	log.Debugf("Input: %q", text)
	text = cleanInput(text)

	if len(text) < 1 || !isJapanese(text) {
		return
	} else if len(text) > 300 {
		log.Debugf("Input too long (%d)", len(text))
		return
	}

	out, err := translateString(text)
	if err != nil {
		log.Errorf("Failed to translate")
		return
	}

	out = cleanOutput(out)

	fmt.Printf("\n%s\n", out)
}
