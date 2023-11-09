package main

import (
	"fmt"
	"sort"
	"strings"
)

type Question struct {
	Type    string
	Text    string
	Options map[string]string
	Answer  string
}

func main() {
	inputQuestionsText := `What colour is the sky?

answer Azure or blue or orange
---
What is the colour of grass?
green or brown
---
Short answer no question mark
short answer answer
---
What is the capital of Australia?

a. Sydney
b. Melbourne
c. Canberra
d. Adelaide

Answer c
---
Which planet is known as the Red Planet?

1) Venus
2) Mars
3) Jupiter
4) Saturn

2
---
Mars is known as the Red Planet
TRUE
---
Venus is known as the Red Planet
F
---
Jupiter is known as the Red Planet
False
---
Who wrote the novel "Pride and Prejudice"?

A) Charles Dickens
B) Jane Austen
C) Mark Twain
D) George Orwell

Correct answer B.
---
What is the square root of 81?

a - 8
b - 9
c - 10
d - 11
e - 999

b
---
True or false, this has been a difficult but rewarding process

true
---
`

lines := strings.Split(inputQuestionsText, "\n")
var questions []Question
var q Question = Question{Options: make(map[string]string)}
var optKey string
var lastOptionRune rune = '0'

// Loop through the lines
lineCount := 0
for _, line := range lines {
    line = strings.ReplaceAll(line, "\t", "")
    for strings.Contains(line, "  ") {
        line = strings.ReplaceAll(line, "  ", " ")
    }
    if line == "" {
        continue
    }

    //fmt.Println("Debug: Processing line - ", line)

    if line == "---" {
        if q.Text != "" && len(q.Options) > 0 {
            questions = append(questions, q)
            q = Question{Options: make(map[string]string)}
        }
        lineCount = 0
        continue
    }

    lineCount++

    // Check if line is a question
    if lineCount == 1 {
        q.Text = line
        lastOptionRune = '0'
        continue
    }

// Check if line is an option
if len(line) > 1 && strings.ContainsAny(strings.ToLower(string(line[0])), "abcdefghijklmnopqrstuvwxyz1234567890") && (line[1] == ')' || line[1] == '.' || (len(line) > 2 && line[1] == ' ' && line[2] == '-')) {
    optKey = string(line[0])
    if line[1] == ' ' && line[2] == '-' {
        q.Options[optKey] = strings.TrimSpace(line[3:])
    } else {
        q.Options[optKey] = strings.TrimSpace(line[2:])
    }
    if lastOptionRune < rune(optKey[0]) {
        lastOptionRune = rune(optKey[0])
    }
    continue
}


// If line is not a question or an option, it must be an answer
if len(q.Options) == 0 {
    q.Options["1"] = line
} else {
    // Check if line is an option
    if len(line) > 1 && strings.ContainsAny(strings.ToLower(string(line[0])), "abcdefghijklmnopqrstuvwxyz1234567890") && (line[1] == ')' || line[1] == '.') {
        optKey = string(line[0])
        q.Options[optKey] = strings.TrimSpace(line[2:])
        if lastOptionRune < rune(optKey[0]) {
            lastOptionRune = rune(optKey[0])
        }
    } else {
        // Check if line is an answer
        if len(line) == 1 { // Check if line is only one character long
            for key := range q.Options {
                if strings.ToLower(key[0:1]) == strings.ToLower(line[0:1]) {
                    q.Answer = key
                    break
                }
            }
        } else {
            prefixes := []string{"correct answer", "answer", "correct answers", "answers", "correct answers -", "answers -"}
            lowerLine := strings.ToLower(line)
            for _, prefix := range prefixes {
                if strings.HasPrefix(lowerLine, prefix) {
                    answerKey := strings.Fields(strings.TrimPrefix(lowerLine, prefix))[0]
                    for key := range q.Options {
                        if strings.ToLower(key[0:1]) == strings.ToLower(answerKey[0:1]) {
                            q.Answer = key
                            break
                        }
                    }
                    break
                }
            }
        }
        // If there is only one option, treat the answer line as both an answer and an option
        if len(q.Options) == 1 {
			q.Type = "TF"
            answerKey := string(lastOptionRune + 1)
            q.Options[answerKey] = line
            q.Answer = answerKey
        } else {
        q.Type = "MC"
		}
    }
}

}

// Add the last question
if q.Text != "" && len(q.Options) > 0 {
    questions = append(questions, q)
}

processConvertQuestions(&questions)

printQuestions(questions)

}
func processConvertQuestions(questions *[]Question) {
    for i, q := range *questions {
        fmt.Println("Before processing:", q)
        if len(q.Options) == 1 {
            for key, value := range q.Options {
                lowerOption := strings.ToLower(value)
                if lowerOption == "t" || lowerOption == "f" || lowerOption == "true" || lowerOption == "false" {
                    (*questions)[i].Type = "TF"
                } else {
                    (*questions)[i].Type = "SA"
                }
                // Check if the key is "1" and remove it
                if key == "1" {
                    (*questions)[i].Answer = value
                } else {
                    (*questions)[i].Answer = key + " " + value
                }
                prefixes := []string{"correct answer ", "answer ", "answer: ", "correct answers ", "answers ", "answers: ", "answers -"}
                lowerAnswer := strings.ToLower((*questions)[i].Answer)
                for _, prefix := range prefixes {
                    if strings.HasPrefix(lowerAnswer, strings.ToLower(prefix)) {
                        (*questions)[i].Answer = strings.TrimPrefix(lowerAnswer, strings.ToLower(prefix))
                        break
                    }
                }
            }
        } else {
            (*questions)[i].Type = "MC"
        }
        fmt.Println("After processing:", (*questions)[i])
    }
}




func printQuestions(questions []Question) {
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
		} else   if q.Type == "SA" {
            // Get the first (and only) option
            var firstOption string
            for _, option := range q.Options {
                firstOption = option
                break
            }
            answers := strings.Split(firstOption, " or ")
            for _, answer := range answers {
                answer = strings.TrimSpace(answer)
                fmt.Printf("Answer,100,%s\n", answer)
            }
        }

        fmt.Println("Hint,")
        fmt.Println("Feedback,")
    }
}