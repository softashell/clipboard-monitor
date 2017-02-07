package main

import (
	"crypto/sha256"
	"fmt"
	"github.com/tjgq/clipboard"
	"io"
)

func monitorCliboard() {
	clipboard.Notify(inChan)

	var old string

	for str := range inChan {
		if len(str) < 1 {
			continue
		}

		hash := sha256.New()
		io.WriteString(hash, str)
		sha := fmt.Sprintf("%x", hash.Sum(nil))

		if sha == old {
			continue
		}

		old = sha

		outChan <- str
	}
}
