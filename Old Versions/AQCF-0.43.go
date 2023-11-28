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
	questionsText := `What colour is the sky?

answer Azure or blue or orange

What is the colour of grass?
green or brown

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
var findNextQuestion bool
var shortAnswerLine bool


for _, line := range lines {
	fmt.Println("Debug: Processing line - ", line)
	line = strings.ReplaceAll(line, "\t", "")
	for strings.Contains(line, "  ") {
		line = strings.ReplaceAll(line, "  ", " ") // Replace double spaces with a single space
	}
	if line == "" {
		fmt.Println("Debug: Skipping empty line")
		continue
	}
	if strings.HasSuffix(line, "?") {
		fmt.Println("Debug: Line ends with '?'")
		if q.Text == "" {
			q.Text = line
			q.Type = "MC"
			fmt.Println("Debug: Set question text and type")
		} else {
			questions = append(questions, q)
			q = Question{Options: make(map[string]string)}
			fmt.Println("Debug: Appended question to questions and reset q")
			q.Text = line
			q.Type = "MC"
			fmt.Println("Debug: Set question text and type for new question")
		}
		shortAnswerLine = true
		findNextQuestion = false
		continue
	}
	if !shortAnswerLine {
	if lastOptionRune != '0' && (rune(line[0]) != lastOptionRune+1) && (rune(line[0]) != lastOptionRune) {
		if len(q.Options) > 0 {
			fmt.Println("Debug: Options length is greater than 0")
			q.Answer = line
			prefixes := []string{"correct answer ", "answer ", "answer: "}
			lowerLine := strings.ToLower(line)
			for _, prefix := range prefixes {
				if strings.HasPrefix(lowerLine, strings.ToLower(prefix)) {
					fmt.Println("Debug: Line starts with prefix: ", prefix)
					answerKey := strings.Fields(strings.TrimPrefix(lowerLine, strings.ToLower(prefix)))[0]
					for key := range q.Options {
						if strings.ToLower(key[0:1]) == strings.ToLower(answerKey[0:1]) {
							q.Answer = key
							fmt.Println("Debug: Answer key found: ", key)
							break
						}
					}
					break
				}
			}
			if len(line) == 1 && strings.ToLower(line) != "t" && strings.ToLower(line) != "f" { // Check if line is only one character long and not "T" or "F"
				fmt.Println("Debug: Line is one character long and not 'T' or 'F'")
				for key := range q.Options {
					if strings.ToLower(key[0:1]) == strings.ToLower(line[0:1]) {
						q.Answer = key
						lastOptionRune = '0'
						fmt.Println("Debug: Answer key found: ", key)
						break
					}
				}

			} else if len(line) > 1 { 
				q.Answer = line
				fmt.Println("Debug: Answer key found: ", line)
			}
			questions = append(questions, q)
			q = Question{Options: make(map[string]string)}
			fmt.Println("Debug: Appended question to questions and reset q")
		}
		lastOptionRune = '0'
		findNextQuestion = true
		fmt.Println("Debug: Reset lastOptionRune and set findNextQuestion to true")
		continue
	}
	

	if findNextQuestion {
		q.Text = line
		q.Type = "MC"
		findNextQuestion = false
		fmt.Println("Debug: Set question text and type and reset findNextQuestion")
		continue
	}
	// Process options
	if len(strings.Fields(line)) > 0 {
		//determine Option's first character
		optKey = strings.Fields(line)[0]
		if len(optKey)+1 <= len(line) {
			q.Options[optKey] = line[len(optKey)+1:]
		} else {
			q.Options[optKey] = ""
		}
		lastOptionRune = rune(optKey[0])
		fmt.Printf("Debug: Current Question - Text: %s, Type: %s, Answer: %s, Options: %v\n", q.Text, q.Type, q.Answer, q.Options)
		fmt.Printf("Debug: lastOptionRune is now: %c\n", lastOptionRune) // Print the lastOptionRune
	}
}
if shortAnswerLine && q.Options[line] != line {
	q.Options[line] = line
	shortAnswerLine = false
	fmt.Println("Debug: Added option for question")
	continue
}
}


	processTFQuestions(&questions)
	processSAQuestions(&questions)
	printQuestions(questions)

	if q.Text != "" && len(q.Options) > 0 {
		questions = append(questions, q)
		q = Question{Options: make(map[string]string)}
	}
}

func processTFQuestions(questions *[]Question) {
    for i, q := range *questions {
        for option, _ := range q.Options {
            lowerOption := strings.ToLower(option)
            if lowerOption == "t" || lowerOption == "f" || lowerOption == "true" || lowerOption == "false" {
                (*questions)[i].Type = "TF"
                (*questions)[i].Answer = option
                break
            }
        }
    }
}

	
func processSAQuestions(questions *[]Question) {
    for i, q := range *questions {
        if len(q.Options) == 1 && q.Type != "TF" {
			for key, value := range q.Options {
				(*questions)[i].Type = "SA"
				(*questions)[i].Answer = key + " " + value
				// Remove prefixes from the answer
				prefixes := []string{"correct answer ", "answer ", "answer: ", "correct answers ", "answers ", "answers: ", "answers -"}
				lowerAnswer := strings.ToLower((*questions)[i].Answer)
				for _, prefix := range prefixes {
					if strings.HasPrefix(lowerAnswer, strings.ToLower(prefix)) {
						(*questions)[i].Answer = strings.TrimPrefix(lowerAnswer, strings.ToLower(prefix))
						break
					}
				}
			}
		}
	//	fmt.Printf("Question %d - Text: %s, Type: %s, Answer: %s, Options: %v\n", i, (*questions)[i].Text, (*questions)[i].Type, (*questions)[i].Answer, (*questions)[i].Options)

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
