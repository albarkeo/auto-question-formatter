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
	questionsText := `What is the capital of Australia?

a. Sydney
b. Melbourne
c. Canberra
d. Adelaide

answer c

Which planet is known as the Red Planet?

1) Venus
2) Mars
3) Jupiter
4) Saturn

2

Mars is known as the Red Planet
TRUE

Venus is known as the Red Planet
F

Jupiter is known as the Red Planet
False

Who wrote the novel "Pride and Prejudice"?

A) Charles Dickens
B) Jane Austen
C) Mark Twain
D) George Orwell

Correct answer B.

What is the square root of 81?

a - 8
b - 9
c - 10
d - 11
e - 999

b

True or false, this has been a difficult but rewarding process

true
`

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
					answerKey := strings.Fields(strings.TrimPrefix(lowerLine, strings.ToLower(prefix)))[0]
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
			if strings.HasPrefix(strings.ToLower(line), "true or false") {
				q.Type = "TF"
			} else {
				q.Type = "MC"
			}
			continue
		}

		optKey = strings.Fields(line)[0]
		if len(optKey)+1 <= len(line) {
			q.Options[optKey] = line[len(optKey)+1:]
		} else {
			q.Options[optKey] = ""
		}
		lastOptionRune = rune(optKey[0])

		if strings.ToLower(optKey) == "true" || strings.ToLower(optKey) == "false" || strings.ToLower(optKey) == "t" || strings.ToLower(optKey) == "f" {
			q.Type = "TF"
			q.Answer = optKey
			if len(q.Options) > 0 {
				questions = append(questions, q)
				q = Question{Options: make(map[string]string)}
			}
			lastOptionRune = '0'
			afterAnswer = true
		}

	}

	if q.Text != "" && len(q.Options) > 0 {
		questions = append(questions, q)
	}

	for _, q := range questions {
		fmt.Println()
		fmt.Printf("NewQuestion,%s\n", q.Type)
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

		if q.Type == "MC" {
			for _, key := range keys {
				option := q.Options[key]
				score := "0"
				if key == q.Answer {
					score = "100"
				}
				fmt.Printf("Option,%s,%s %s\n", score, key, option)
			}
		} else if q.Type == "TF" {
			if strings.ToLower(q.Answer) == "true" || strings.ToLower(q.Answer) == "t" {
				fmt.Println("TRUE,100")
				fmt.Println("FALSE,0")
			} else {
				fmt.Println("TRUE,0")
				fmt.Println("FALSE,100")
			}
		}

		fmt.Println("Hint,")
		fmt.Println("Feedback,")
	}
}