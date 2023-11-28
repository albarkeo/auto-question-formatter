package main

import (
	"fmt"
	"sort"
	"strings"
)

type Question struct {
	Text    string
	Options map[string]string
	Answer  string
}

func main() {
	questionsText := ``

	lines := strings.Split(questionsText, "\n")
	var questions []Question
	var q Question = Question{Options: make(map[string]string)}
	var optKey string
	var lastOptionRune rune = '0'
	var afterAnswer bool = true

	for _, line := range lines {
		if line == "" {
			continue
		}

		if lastOptionRune != '0' && (rune(line[0]) != lastOptionRune+1) {
			prefixes := []string{"correct answer ", "answer ", "answer: "}
			lowerLine := strings.ToLower(line)
			for _, prefix := range prefixes {
				if strings.HasPrefix(lowerLine, strings.ToLower(prefix)) {
					answerKey := strings.Fields(strings.TrimPrefix(line, prefix))[0]
					for key := range q.Options {
						if strings.ToLower(key[0:1]) == strings.ToLower(answerKey[0:1]) {
							q.Answer = key
							break
						}
					}
					break
				}
			}
			if len(line) == 1 { // Check if line is only one character long
				for key := range q.Options {
					if strings.ToLower(key[0:1]) == strings.ToLower(line[0:1]) {
						q.Answer = key
						break
					}
				}
			}
			if len(q.Options) > 0 {
				questions = append(questions, q)
				q = Question{Options: make(map[string]string)}
			}
			lastOptionRune = '0'
			afterAnswer = true
			continue
		}

		if afterAnswer {
			q.Text = line
			afterAnswer = false
			continue
		}

		optKey = strings.Fields(line)[0]
		q.Options[optKey] = line[len(optKey)+1:]
		lastOptionRune = rune(optKey[0])

	}

	if q.Text != "" && len(q.Options) > 0 {
		questions = append(questions, q)
	}

	for _, q := range questions {
		fmt.Println()
		fmt.Println("NewQuestion,MC")
		fmt.Println("ID,")
		fmt.Println("Title,")
		fmt.Printf("QuestionText,%s\n", q.Text)
		fmt.Println("Points,")
		fmt.Println("Difficulty,")
		fmt.Println("Image,")

		var keys []string
		for k := range q.Options {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, key := range keys {
			option := q.Options[key]
			score := "0"
			if key == q.Answer {
				score = "100"
			}
			fmt.Printf("Option,%s,%s %s\n", score, key, option)
		}

		fmt.Println("Hint,")
		fmt.Println("Feedback,")
	}
}
