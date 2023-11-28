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
	Type    string
}

func main() {
	questionsText := `What is the square root of 81?
 
a - 8
b - 9
c - 10
d - 11
 
b
 
Who painted the Mona Lisa?
 
A) Vincent van Gogh
B) Pablo Picasso
C) Leonardo da Vinci
D) Claude Monet
Correct answer C.

True or false, this has been a difficult but rewarding process

True`

	lines := strings.Split(questionsText, "\n")
	var questions []Question
	var q Question = Question{Options: make(map[string]string)}
	var optKey string
	var lastProcessedOptionChar rune = '0'
	var afterAnswer bool = false

	for _, line := range lines {
		if line == "" {
			if len(q.Options) > 0 {
				questions = append(questions, q)
				q = Question{Options: make(map[string]string)}
			}
			lastProcessedOptionChar = '0'
			afterAnswer = false
			continue
		}

		if afterAnswer {
			q.Text = line
			afterAnswer = false
			continue
		}

		fields := strings.Fields(line)
		if len(fields) > 0 {
			optKey = fields[0]
			if len(optKey) > 0 {
				if len(optKey)+1 <= len(line) {
					q.Options[optKey] = line[len(optKey)+1:]
				}
				lastProcessedOptionChar = rune(optKey[0])
			}
		}

		if lastProcessedOptionChar != '0' && (rune(line[0]) != lastProcessedOptionChar+1) {
			processAnswerLine(&q, line)
			afterAnswer = true
			continue
		}

	}

	if q.Text != "" && len(q.Options) > 0 {
		questions = append(questions, q)
	}

	for _, q := range questions {
		printQuestion(q)
	}
}

func processAnswerLine(q *Question, line string) {
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
			if len(key) > 0 && len(line) > 0 && strings.ToLower(key[0:1]) == strings.ToLower(line[0:1]) {
				q.Answer = key
				break
			}
		}

	}

	if len(q.Options) == 2 {
		var options []string
		for _, option := range q.Options {
			options = append(options, strings.ToLower(strings.TrimSpace(option)))
		}
		sort.Strings(options)
		if options[0] == "false" && options[1] == "true" {
			q.Type = "TF"
		} else {
			q.Type = "MC"
		}
	}

}

func printQuestion(q Question) {
	if q.Type == "TF" {
		printTFQuestion(q)
	} else {
		printMCQuestion(q)
	}
}

func printMCQuestion(q Question) {
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

func printTFQuestion(q Question) {
	fmt.Println()
	fmt.Println("NewQuestion,TF")
	fmt.Println("ID,")
	fmt.Println("Title,")
	fmt.Printf("QuestionText,%s\n", q.Text)
	fmt.Println("Points,")
	fmt.Println("Difficulty,")
	fmt.Println("Image,")

	for _, answer := range q.Options {
		answer = strings.TrimSpace(answer)
		if strings.EqualFold(answer, "True") {
			fmt.Println("TRUE,100,")
			fmt.Println("FALSE,0,")
		} else if strings.EqualFold(answer, "False") {
			fmt.Println("TRUE,0,")
			fmt.Println("FALSE,100,")
		}
	}

	fmt.Println("Hint,")
	fmt.Println("Feedback,\n")
}

func isTFQuestion(q Question) bool {
	for _, answer := range q.Options {
		answer = strings.TrimSpace(answer)
		if answer == "True" || answer == "False" || answer == "true" || answer == "false" {
			return true
		}
	}
	return false
}