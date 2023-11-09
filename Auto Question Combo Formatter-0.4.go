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
	questionsText := `
	What colour is the sky?
answer Azure or blue or orange
	
	What colour is grass?
correct answers green or brown

How many billion people are there in the world?
7.88billion or 7.8 or 7.88 or 7.8 billion or 7.8 Billion
	
	What is the capital of Australia?

a. Sydney
b. Melbourne
c. Canberra
d. Adelaide

Answer c

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
	var findNextQuestion bool = true
	var prevLine string

	for _, line := range lines {
		line = strings.ReplaceAll(line, "\t", "") // Replace tabs with a single space
		for strings.Contains(line, "  ") {
			line = strings.ReplaceAll(line, "  ", " ") // Replace double spaces with a single space
		}
		if line == "" {
			continue
		}

		
		fields := strings.Fields(line)
		if len(fields) > 0 {
			optKey = fields[0]
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
					findNextQuestion = true
				}
				lastOptionRune = '0'
			}
		} else {
			findNextQuestion = true
			continue
		}
			if lastOptionRune != '0' && (rune(line[0]) != lastOptionRune+1) {
				// Check for MC and TF type questions
				if len(q.Options) > 0 {
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
					questions = append(questions, q)
					q = Question{Options: make(map[string]string)}
				}
				lastOptionRune = '0'
				findNextQuestion = true
				continue
			}

			if strings.HasSuffix(line, "?") {
				// This line is assumed a new question
				q.Text = line
				fmt.Println("Reading question:", q.Text)
				findNextQuestion = false
				prevLine = line
				continue
				} else if !findNextQuestion && !strings.HasPrefix(strings.ToLower(line), "true") && !strings.HasPrefix(strings.ToLower(line), "false") && !strings.Contains(line, " or ") && len(prevLine) > 0 && strings.HasSuffix(prevLine, "?") {
					// This line is an answer to the previous question
					q.Type = "SA"
					q.Answer = line
					findNextQuestion = true
					prevLine = line
					continue
				}else if !strings.HasSuffix(line, "?") && !findNextQuestion {
					// This line is an answer
					q.Answer = line
					findNextQuestion = true
					prevLine = line
					continue
				}

		if findNextQuestion {
			q.Text = line
			fmt.Println("Reading question:", q.Text)  // This line will print the question text
			if strings.HasPrefix(strings.ToLower(line), "true or false") {
				q.Type = "TF"
				findNextQuestion = false
			} else {
				q.Type = "MC"
				findNextQuestion = false
			}
			prevLine = line
			continue
		}
		
	}
	// Process SA type question
		for i, q := range questions {
			if len(q.Options) == 1 && q.Type != "TF" {
				for key, value := range q.Options {
					questions[i].Type = "SA"
					questions[i].Answer = key + " " + value
					// Remove prefixes from the answer
					prefixes := []string{"correct answer ", "answer ", "answer: ", "correct answers ", "answers ", "answers: ", "answers - "}
					lowerAnswer := strings.ToLower(questions[i].Answer)
					for _, prefix := range prefixes {
						if strings.HasPrefix(lowerAnswer, strings.ToLower(prefix)) {
							questions[i].Answer = strings.TrimPrefix(lowerAnswer, strings.ToLower(prefix))
							break
						}
					}
					findNextQuestion = true
				}
			}
		fmt.Printf("Question %d - Type: %s, Answer: %s, Options: %v\n", i, questions[i].Type, questions[i].Answer, questions[i].Options)
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
		} else if q.Type == "SA" {
			answers := strings.Split(q.Answer, " or ")
			for _, answer := range answers {
				answer = strings.TrimSpace(answer)
				fmt.Printf("Answer,100,%s\n", answer)
			}
		}

		fmt.Println("Hint,")
		fmt.Println("Feedback,")
	}
}
