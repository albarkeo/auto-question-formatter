package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
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
	fmt.Println("Enter questions text, type 'END' at the end of the question block, then press 'Enter': ")
	reader := bufio.NewReader(os.Stdin)

	var inputQuestionsText string
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		// Check if the line is the termination string
		if strings.TrimSpace(line) == "END" {
			break
		}

		// Append the line to the questions text
		inputQuestionsText += line
	}

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
			answerKey := string(line[0])         // Get the first character of the line as the answer key
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

		// Check if line ends with '+++' or '---'
		if strings.HasSuffix(line, "+++") || strings.HasSuffix(line, "---") {
			// Process the current question here
			if q.Text != "" && len(q.Options) > 0 {
				questions = append(questions, q)
				q = Question{Options: make(map[string]string)}
			}
			// Reset the question
			q = Question{Options: make(map[string]string)}
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
		if len(line) > 1 && isValidCharacter(strings.ToLower(string(line[0]))) {
			if isValidListItemDelimiter(string(line[1])) || (len(line) > 2 && isValidListItemDelimiter(string(line[1:3]))) {
				optKey = string(line[0])
				if len(line) > 2 && line[1] == ' ' && line[2] == '-' {
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
			q.Type = "WR"
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

				switch len(q.Options) {
				case 1:
					if q.Type != "WR" {
						q.Type = "TF"
						answerKey := string(lastOptionRune + 1)
						q.Options[answerKey] = line
						q.Answer = answerKey
					}
				default:
					q.Type = "MC"
				}
			}
		}
	}

	// Add the last question
	if q.Text != "" {
		questions = append(questions, q)
	}

	processConvertQuestions(&questions)

	printQuestions(questions, prefixes)

	writeQuestionsToCSV(questions, prefixes)

	fmt.Println()
	fmt.Println("Success! CSV file saved to location of this program")
}

func processConvertQuestions(questions *[]Question) {
	for i, q := range *questions {
		switch len(q.Options) {
		case 0:
			(*questions)[i].Type = "WR"
		case 1:
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
		default:
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
	validDelimiters := []string{")", ".", "-", " -"}
	for _, delimiter := range validDelimiters {
		if s == delimiter {
			return true
		}
	}
	return false
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

func writeQuestionsToCSV(questions []Question, prefixes []string) {
	// Get the current date and time
	now := time.Now()

	// Format the date and time as a string
	timestamp := now.Format("20060102_1504")
	// Create a CSV file with the timestamp in the name
	file, err := os.Create("Formatted_questions_" + timestamp + ".csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Loop through the questions
	for _, q := range questions {
		err := writer.Write([]string{"NewQuestion", q.Type})
		if err != nil {
			log.Fatal(err)
		}
		writer.Write([]string{"ID"})
		writer.Write([]string{"Title"})
		writer.Write([]string{"QuestionText", q.Text})
		writer.Write([]string{"Points"})
		writer.Write([]string{"Difficulty"})
		writer.Write([]string{"Image"})
		if q.Type == "WR" {
			writer.Write([]string{"InitialText"})
			writer.Write([]string{"AnswerKey"})
		} else {

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
						writer.Write([]string{"Option", score, option})
					} else {
						writer.Write([]string{"Option", score, key + " " + option})
					}
				}
			} else if q.Type == "TF" {
				if strings.ToLower(q.Answer) == "true" || strings.ToLower(q.Answer) == "t" {
					writer.Write([]string{"TRUE", "100"})
					writer.Write([]string{"FALSE", "0"})
				} else {
					writer.Write([]string{"TRUE", "0"})
					writer.Write([]string{"FALSE", "100"})
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
					writer.Write([]string{"Answer", "100", answer})
				}
			}

			writer.Write([]string{"Hint"})
			writer.Write([]string{"Feedback"})

			// Add an empty line after each question
			writer.Write([]string{""})
		}
	}
}
