package main

import (
	log "github.com/Sirupsen/logrus"
	"golang.org/x/text/width"
	"regexp"
	"strings"
	"time"
	"unicode"
)

func constantRepeatFilter(str string, rep int) string {
	if rep <= 1 {
		return str
	}

	log.Debugf("Cleaning %d repeated characters", rep)

	var out string
	count := 0

	// TODO: Maybe double check if it's actually repeating since automatic detection isn't too good
	for _, c := range str {
		if count < rep {
			count++
			continue
		}

		count = 0
		out += string(c)
	}

	return out
}

func autoRepeatFilter(str string) string {
	// TODO: Should make it handle small variations, since it's possible to have a bit of garbage which isn't repeated at start or end

	//Euclidean algorithm
	gcd := func(x, y int) int {
		for y != 0 {
			x, y = y, x%y
		}

		return x
	}

	var last rune

	count := 1
	repeat := 0

	for _, c := range str {
		if c == last {
			count++
			continue
		}

		last = c
		count = 1
		repeat = gcd(count, repeat)
	}

	repeat = gcd(count, repeat)

	// TODO: Handle blocks of text as well
	return constantRepeatFilter(str, repeat)
}

func cleanInput(text string) string {
	start := time.Now()

	// Characters that can't be printed are replaced with space
	isValid := func(r rune) rune {
		if !unicode.IsPrint(r) && !unicode.IsSpace(r) {
			return ' '
		}

		return r
	}

	text = strings.Map(isValid, text)
	text = strings.TrimSpace(text)

	// Convert full width to half width
	text = width.Narrow.String(text)

	// Money
	replaceRegex(text, `^(((\d+:)+)?\d+[デﾙ]{2,4})+`, "")

	text = autoRepeatFilter(text)

	// Compress repeated whitespace
	text = replaceRegex(text, `\s{2,}`, " ")

	// TODO: Probably should be handled externally as well
	splitter := regexp.MustCompile(`(?:^(.{0,15}?)?(?::?\s+?)?)?(?:(?:[｢『(])([\s\S]*?).(?:[」』）]?)$)`) // Just end me
	split := splitter.FindStringSubmatch(text)
	if len(split) == 3 {
		if len(split[1]) > 0 && len(split[2]) > 0 {
			// TODO: Maybe attempt to translate name as well and return slice
			text = split[2]

			log.WithFields(log.Fields{
				"name": split[1],
				"text": split[2],
			}).Debug("Input was split!")
		}
	}

	// TODO: Move these to external files

	// Save screen
	text = replaceRegex(text, `ﾌｲﾙ\d\d?`, "")

	// Cistina
	text = replaceRegex(text, `^[｢『(]\d+日 \(.*?[｣』）](?:\s(.*))?`, "$1")

	// HP,MP,TP
	text = replaceRegex(text, `([HMTP]+[\d]+/?([\d]+)?)+`, "")

	// Remove extra quotes
	text = replaceRegex(text, `^[｢『(]+(.*?)[｣』）]+$`, "$1")

	for old, new := range inputReplacement {
		text = strings.Replace(text, old, new, -1)
	}

	log.WithFields(log.Fields{
		"time": time.Since(start),
	}).Debugf("Cleaned: %q", text)

	return text
}

func cleanOutput(text string) string {
	start := time.Now()

	// TODO: External config
	text = width.Narrow.String(text)

	for old, new := range transReplacement {
		text = strings.Replace(text, old, new, -1)
	}

	for old, new := range outputReplacement {
		text = strings.Replace(text, old, new, -1)
	}

	// Repeated whitespace
	text = replaceRegex(text, `\s{2,}`, " ")

	// ... ... ......
	text = replaceRegex(text, `\s+?\.{3,}(\s\.+)?`, "...")
	text = replaceRegex(text, `\s+[\.·]`, ".")

	// TODO: Probably should fix other characters being spammed so much as well
	log.WithFields(log.Fields{
		"time": time.Since(start),
	}).Debugf("Cleaned: %q", text)

	return text
}

func replaceRegex(text string, expr string, repl string) string {
	// TODO: Cache compiled expressions in case this gets too slow
	regex := regexp.MustCompile(expr)

	return regex.ReplaceAllString(text, repl)
}

func isJapanese(text string) bool {
	regex := regexp.MustCompile(`(\p{Hiragana}|\p{Katakana}|\p{Han})`)
	matches := regex.FindAllString(text, 1)

	if len(matches) >= 1 {
		return true // Probably
	}

	return false
}
