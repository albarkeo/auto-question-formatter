package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type Question struct {
	Type      string
	Text      string
	Options   map[string]string
	Answer    string
	ID        int
	OptionKey string
	AnswerKey string
}

var removeListPrefixes = true // Change this to false if you want to keep the prefixes
var prefixes = []string{"answer ", "answer: ", "answer- ", "answers ", "answers: ", "answers- ", "correct answer ", "correct answer: ", "correct answer- ", "correct answers: ", "correct answers- "}

func printWelcomeMessage() {
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println("+++")
	fmt.Println()
	fmt.Println("Welcome to the Auto Question Formatter")
	fmt.Println("This tool is used to convert text copied from Word or equivalent")
	fmt.Println("It will generate a .csv file formatted for Brightspace")
	fmt.Println("v1.5.5")
	fmt.Println()
	fmt.Println("+++")
	fmt.Println()
	fmt.Println("Paste all question text below, type 'END' at the end of the question block, then press 'Enter': ")
	fmt.Println()
}

func main() {
	printWelcomeMessage()

	reader := bufio.NewReader(os.Stdin)
	inputQuestionsText := readInput(reader)
	lines := strings.Split(inputQuestionsText, "\n")
	questionNumber := 1
	questions := []Question{}
	q := Question{Options: make(map[string]string)}
	lastOptionRune := rune('0')
	lineCount := 0
	// Loop through the lines

	for _, line := range lines {
		line = processLine(line)
		if line == "" {
			continue
		}

		if handleNewQuestion(line, &q, &questions, &questionNumber) {
			lineCount++
			continue
		}

		lineCount++

		if handleEndOfQuestion(line, &q, &questions) {
			lineCount = 0
			continue
		}

		handleQuestionTextLine(line, &q, lineCount)

		q, line, _ := handleCorrectAnswerLine(line, &q)

		handleOptionLine(line, q, &lastOptionRune)

		handleTabSeparatedLine(line, q, &questions)

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

func readInput(reader *bufio.Reader) string {
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
	return inputQuestionsText
}

func processLine(line string) string {
	line = strings.TrimSpace(line)
	line = strings.TrimPrefix(line, "•\t")
	line = strings.ReplaceAll(line, "mml:", "")
	for strings.Contains(line, "  ") {
		line = strings.ReplaceAll(line, "  ", " ")
	}
	return line
}

func handleTabSeparatedLine(line string, q *Question, questions *[]Question) {
	splitLine := strings.Split(line, "\t")
	if len(splitLine) > 1 {
		// If there is a tab in the line, treat the first part as the question and the second part as the answer
		q.Text = splitLine[0]
		q.Options["1"] = splitLine[1]
		*questions = append(*questions, *q)
		*q = Question{Options: make(map[string]string)}
	}

}

func handleCorrectAnswerLine(line string, q *Question) (*Question, string, bool) {
	if strings.HasPrefix(line, "*") {
		line = strings.TrimPrefix(line, "*") // this line is the correct answer...
		q.AnswerKey = string(line[0])        // Get the first character of the line as the answer key
		q.Answer = q.AnswerKey
		return q, line, true
	}
	return q, line, false
}

func handleEndOfQuestion(line string, q *Question, questions *[]Question) bool {
	if strings.HasSuffix(line, "+++") || strings.HasSuffix(line, "---") {
		if q.Text != "" && len(q.Options) > 0 {
			*questions = append(*questions, *q)
			*q = Question{Options: make(map[string]string)}
		}
		// Reset the question
		*q = Question{Options: make(map[string]string)}
		return true
	}
	return false
}

func handleNewQuestion(line string, q *Question, questions *[]Question, questionNumber *int) bool {
	if strings.HasPrefix(line, strconv.Itoa(*questionNumber)) {
		if q.Text != "" && len(q.Options) > 0 {
			*questions = append(*questions, *q)
		}
		*q = Question{Options: make(map[string]string)}

		// Remove the question number, '.' and any leading white space from the line
		questionLineNoNumber := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(line, strconv.Itoa(*questionNumber)), "."))

		// Assign the formatted line to q.Text
		q.Text = questionLineNoNumber
		q.ID = *questionNumber
		(*questionNumber)++
		return true
	}
	return false
}

func handleQuestionTextLine(line string, q *Question, lineCount int) {
	if lineCount == 1 {
		if q.Text == "" {
			q.Text = line
		}
	}
}

func handleOptionLine(line string, q *Question, lastOptionRune *rune) bool {
	if len(line) > 1 && isValidCharacter(strings.ToLower(string(line[0]))) {
		delimiter, optionStartIndex := isValidListItemDelimiter(line[1:])
		if delimiter != "" {
			q.OptionKey = string(line[0])
			q.Options[q.OptionKey] = strings.TrimSpace(line[optionStartIndex+1:])
			if *lastOptionRune < rune(q.OptionKey[0]) {
				*lastOptionRune = rune(q.OptionKey[0])
			}
			return true
		}
	}
	return false
}

func handleAnswerLine(line string, q *Question, lastOptionRune *rune) {
	// If line is not a question or an option, it must be an answer
	if len(q.Options) == 0 {
		q.Options["1"] = line
		q.Type = "WR"
	} else {
		// Check if line is an option
		if len(line) > 1 && isValidCharacter(strings.ToLower(string(line[0:1]))) && (line[1] == ')' || line[1] == '.') {
			q.OptionKey = string(line[0])
			q.Options[q.OptionKey] = strings.TrimSpace(line[2:])
			if *lastOptionRune < rune(q.OptionKey[0]) {
				*lastOptionRune = rune(q.OptionKey[0])
			}
		} else {
			// Check if line is a single character answer
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
						q.AnswerKey = strings.Fields(strings.TrimPrefix(lowerLine, prefix))[0]
						for key := range q.Options {
							if strings.ToLower(key[0:1]) == strings.ToLower(q.AnswerKey[0:1]) {
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
					q.AnswerKey = string(*lastOptionRune + 1)
					q.Options[q.AnswerKey] = line
					q.Answer = q.AnswerKey
				}
			default:
				q.Type = "MC"
			}
		}
	}
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

func isValidListItemDelimiter(s string) (string, int) {
	validDelimiters := []string{") ", ". ", "- ", " - ", "). ", ")."}
	for _, delimiter := range validDelimiters {
		if strings.HasPrefix(s, delimiter) {
			return delimiter, len(delimiter)
		}
	}
	return "", -1
}

func printQuestions(questions []Question, prefixes []string) {
	for _, q := range questions {
		fmt.Println()
		fmt.Printf("NewQuestion,%s\n", q.Type)
		fmt.Printf("ID,%d\n", q.ID)
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
