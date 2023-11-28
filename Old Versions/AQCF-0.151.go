package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

type Question struct {
	Type    string
	Text    string
	Options map[string]string
	Answer  string
}

var removeListPrefixes = true // Change this to false if you want to keep the prefixes
var prefixes = []string{"answer ", "answer: ", "answer- ", "answers ", "answers: ", "answers- ", "correct answer ", "correct answer: ", "correct answer- ", "correct answers: ", "correct answers- "}

func main() {
	inputQuestionsText := `
	What is the primary purpose of using weight machines in resistance training?
•	A) To provide a restricted and controlled plane of motion
•	B) To increase cardiovascular endurance
•	C) To improve balance and coordination 
•	D) To facilitate high-speed movements
•	Answer: A) To provide a restricted and controlled plane of motion
+++
What is the capital of Australia?

a. Sydney
b. Melbourne
*c. Canberra
d. Adelaide
---

What is the capital of Australia?

a. Sydney
b. Melbourne
c. Canberra
d. Adelaide

Answer c
---
What colour is the sky?

Azure or blue; orange
---
Write the value of the 5 in the number 8526.	500
---
Write the value of the 7 in the number 97 450.	7000 or 7 000 or 7,000
---
Write this as a number. 7 tens of thousands, 4 hundreds and 2 ones	70402 or 70 402 or 70,402
---
9 + 6 =	15
---
16 + 7 =	23
---
∫_a^x▒∫_a^s▒〖f(y)\,dy\,ds〗= ∫_a^x▒〖f(y)(x-y)\,dy〗
An answer
---
24 + 5 =	29
---
628 - 284 =	344
---
The Earth is the only planet in our solar system with liquid water on its surface. 	False
---
What is the colour of grass?
answer green or brown
---
Short answer no question mark
a short answer
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
	questions := []Question{}
	q := Question{Options: make(map[string]string)}
	optKey := ""
	lastOptionRune := rune('0')

	// Loop through the lines
	lineCount := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		line = strings.TrimPrefix(line, "•\t")
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "*") {
			line = strings.TrimPrefix(line, "*") // this line is the correct answer...
			answerKey := string(line[0]) // Get the first character of the line as the answer key
			q.Answer = answerKey
		}

		// Remove "mml:" prefix
		line = strings.ReplaceAll(line, "mml:", "")

		// Split the line by tab
		splitLine := strings.Split(line, "\t")
		if len(splitLine) > 1 {
			// If there is a tab in the line, treat the first part as the question and the second part as the answer
			q.Text = splitLine[0]
			q.Options["1"] = splitLine[1]
			questions = append(questions, q)
			q = Question{Options: make(map[string]string)}
		}

		for strings.Contains(line, "  ") {
			line = strings.ReplaceAll(line, "  ", " ")
		}
		if line == "" {
			continue
		}

		if line == "---" || line == "+++" {
			if q.Text != "" && len(q.Options) > 0 {
				questions = append(questions, q)
				q = Question{Options: make(map[string]string)}
			}
			lineCount = 0
			continue
		}

			// Check if line contains "---"
			if strings.Contains(line, "---") {		
			parts := strings.Split(line, "---")
			line = parts[0] // Only consider the part before the "---" delimiter
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
		if len(line) > 1 && isValidCharacter(strings.ToLower(string(line[0]))) {
			if isValidListItemDelimiter(string(line[1])) {
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
		}
		


		// If line is not a question or an option, it must be an answer
		if len(q.Options) == 0 {
			q.Options["1"] = line
		} else {
			// Check if line is an option
			if len(line) > 1 && isValidCharacter(strings.ToLower(string(line[0:1]))) && (line[1] == ')' || line[1] == '.') {
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

	printQuestions(questions, prefixes)

}

func processConvertQuestions(questions *[]Question) {
	for i, q := range *questions {
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
	}
}

func isValidCharacter(s string) bool {
	for _, c := range s {
		if !unicode.IsLetter(c) && !unicode.IsNumber(c) {
			return false
		}
	}
	return true
}

func isValidListItemDelimiter(s string) bool {
	for _, c := range s {
		if c != ')' && c != '.' && c != '-' {
			return false
		}
	}
	return true
}

func printQuestions(questions []Question, prefixes []string) {
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
				if removeListPrefixes {
					fmt.Printf("Option,%s,%s\n", score, option)
				} else {
					fmt.Printf("Option,%s,%s %s\n", score, key, option)
				}
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
			// Get the first (and only) option
			var firstOption string
			for _, option := range q.Options {
				firstOption = option
				break
			}
			r := regexp.MustCompile(`\s+or\s+|\s*;\s*|\s*\t\s*`)
			answers := r.Split(firstOption, -1)
			for i, answer := range answers {
				answer = strings.TrimSpace(answer)
				if i == 0 {
					lowerAnswer := strings.ToLower(answer)
					for _, prefix := range prefixes {
						if strings.HasPrefix(lowerAnswer, strings.ToLower(prefix)) {
							answer = strings.TrimPrefix(answer, prefix)
							break
						}
					}
				}
				fmt.Printf("Answer,100,%s\n", answer)
			}
		}

		fmt.Println("Hint,")
		fmt.Println("Feedback,")
	}
}
