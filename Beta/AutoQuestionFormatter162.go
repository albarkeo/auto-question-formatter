package main

import (
	"bufio"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"syscall/js"
	"unicode"
)

type Question struct {
	Type           string
	Text           string
	Options        map[string]string
	Answer         string
	ID             int
	Image          string
	OptionKey      string
	AnswerKey      string
	Feedback       string
	OptionFeedback map[string]string
}

var removeListPrefixes = true // Change this to false if you want to keep the prefixes
var prefixes = []string{"answer ", "answer: ", "answer- ", "answers ", "answers: ", "answers- ", "correct answer ", "correct answer: ", "correct answer- ", "correct answers: ", "correct answers- "}

func printWelcomeMessage() {
	welcomeMessage := `
+++

Welcome to the Auto Question Formatter
This tool is used to convert text copied from Word or equivalent
It will generate a .csv file formatted for Brightspace
v1.6

+++

Paste all question text below, type 'END' at the end of the question block, then press 'Enter': 
`
	fmt.Println(welcomeMessage)
}

func main() {
	c := make(chan struct{}, 0)

	js.Global().Set("RunMain", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		csvData := RunMain(args[0].String())
		return csvData
	}))
	<-c
}

func RunMain(input string) string {
	printWelcomeMessage()

	// Convert the input string to a strings.Reader
	strReader := strings.NewReader(input)

	// Wrap the strings.Reader with a bufio.Reader
	reader := bufio.NewReader(strReader)

	inputQuestionsText := readInput(reader)
	questionNumber := 1
	questions := []Question{}
	lastOptionRune := rune('0')
	lineCount := 0

	preTabSplitLines := strings.Split(inputQuestionsText, "\n")
	lines := []string{}

	q := Question{
		Options:        make(map[string]string),
		OptionFeedback: make(map[string]string),
	}

	// Preprocess the lines to handle tab-separated lines
	for _, line := range preTabSplitLines {
		tabSplitLines := strings.Split(line, "\t")
		for _, tabSplitLine := range tabSplitLines {
			lines = append(lines, tabSplitLine)
		}
	}

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

		handleQuestionTextLine(line, &q, &lineCount)

		if handleImages(line, &q) {
			continue
		}

		q, line, _ := handleCorrectAnswerLine(line, &q)

		ok := handleFirstOption(line, q, &lastOptionRune)
		if ok {
			continue
		}
		if handleAdditionalOptions(line, q, &lastOptionRune) {
			continue
		}
		if handleOptionFeedback(line, q) {
			continue
		}

		handleTabSeparatedLine(line, q, &questions)

		handleAnswerLine(line, q)

		handleFeedback(line, q)

	}

	// Add the last question
	if q.Text != "" {
		questions = append(questions, q)
	}

	processQuestionType(&questions, &lastOptionRune)

	printQuestions(questions, prefixes)

	writeQuestionsToCSV(questions, prefixes)

	fmt.Println()
	fmt.Println("Success!")

	csvData := writeQuestionsToCSV(questions, prefixes)
	return csvData
}

func readInput(reader *bufio.Reader) string {
	var inputQuestionsText string
	for {
		line, err := reader.ReadString('\n')
		// If the error is io.EOF, we've reached the end of the input
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
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
		*q = Question{
			Options:        make(map[string]string),
			OptionFeedback: make(map[string]string),
		}
	}
}

func handleCorrectAnswerLine(line string, q *Question) (*Question, string, bool) {
	if strings.HasPrefix(line, "*") || (len(line) > 2 && line[1:2] == "*") {
		line = strings.Replace(line, "*", "", 1) // Remove the asterisk from the line
		q.AnswerKey = string(line[0])            // Get the first character of the line as the answer key
		q.Answer = q.AnswerKey
		return q, line, true
	}
	return q, line, false
}


func handleEndOfQuestion(line string, q *Question, questions *[]Question) bool {
	if strings.HasSuffix(line, "+++") || strings.HasSuffix(line, "---") {
		if q.Text != "" {
			*questions = append(*questions, *q)
		}
		// Reset the question
		*q = Question{
			Options:        make(map[string]string),
			OptionFeedback: make(map[string]string),
		}
		return true
	}
	return false
}

func handleNewQuestion(line string, q *Question, questions *[]Question, questionNumber *int) bool {
	if strings.HasPrefix(line, strconv.Itoa(*questionNumber)+".") {
		if q.Text != "" && len(q.Options) > 0 {
			*questions = append(*questions, *q)
		}
		*q = Question{
			Options:        make(map[string]string),
			OptionFeedback: make(map[string]string),
		}

		// Remove the question number, '.' and any leading white space from the line
		questionLineNoNumber := strings.TrimSpace(strings.TrimPrefix(line, strconv.Itoa(*questionNumber)+"."))

		// Assign the formatted line to q.Text
		q.Text = questionLineNoNumber
		q.ID = *questionNumber
		(*questionNumber)++
		return true
	}
	return false
}

func handleQuestionTextLine(line string, q *Question, lineCount *int) {
	if strings.HasPrefix(line, "@") {
		return
	}

	if *lineCount == 1 {
		if q.Text == "" {
			q.Text = line
		}
		(*lineCount)++
	}
}

func handleImages(line string, q *Question) bool {
	if strings.HasPrefix(line, "[[") && strings.HasSuffix(line, "]]") {
		line = strings.Trim(line, "[]")
		q.Image = url.QueryEscape(line)
		return true
	}
	return false
}

func handleFirstOption(line string, q *Question, lastOptionRune *rune) bool {

	for _, prefix := range prefixes {
		lowerLine := strings.ToLower(line)
		if strings.HasPrefix(lowerLine, prefix) {
			line = strings.TrimPrefix(lowerLine, prefix)
		}
	}

	if strings.ToLower(line) == "true" || strings.ToLower(line) == "false" || strings.ToLower(line) == "t" || strings.ToLower(line) == "f" {
		q.Options[line] = line
		q.AnswerKey = line
		q.Answer = line
		return true
	}

	if (q.Type == "" || q.Type == "WR") && !strings.HasPrefix(line, "@") && line != q.Text {

		if len(line) > 0 && isValidCharacter(string(line[0])) {
			delimiter, length := isValidListItemDelimiter(line[1:])
			if delimiter != "" {
				q.Type = "MC"
				q.OptionKey = string(line[0])
				q.Options[q.OptionKey] = strings.TrimSpace(line[1+length:]) // Trim spaces from the option text
				*lastOptionRune = rune(line[0])
				return true
			}
		}

		// Check if the line is an option and answer for a Short Answer question
		hasPrefix := false
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(line), prefix) {
				line = strings.TrimPrefix(line, prefix)
				hasPrefix = true
				break
			}
		}
		if !hasPrefix || len(line) == 1 {
			q.Options[line] = line
			return true
		}
	}
	return false
}

func handleAdditionalOptions(line string, q *Question, lastOptionRune *rune) bool {
	if lastOptionRune != nil && len(line) > 0 && isValidCharacter(string(line[0])) {
		delimiter, length := isValidListItemDelimiter(line[1:])
		if delimiter != "" && rune(line[0]) == *lastOptionRune+1 {
			q.OptionKey = string(line[0])
			q.Options[q.OptionKey] = strings.TrimSpace(line[1+length:]) // Trim spaces from the option text
			if *lastOptionRune < rune(q.OptionKey[0]) {
				*lastOptionRune = rune(q.OptionKey[0])
			}
			return true
		}
	}
	return false
}

func handleAnswerLine(line string, q *Question) {
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
}

func processQuestionType(questions *[]Question, lastOptionRune *rune) {
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
	validDelimiters := []string{")", ") ", ". ", "- ", " - ", "). ", ")."}
	for _, delimiter := range validDelimiters {
		if strings.HasPrefix(s, delimiter) {
			return delimiter, len(delimiter)
		}
	}
	return "", -1
}

func handleFeedback(line string, q *Question) {
	if strings.HasPrefix(line, "@") {
		feedback := strings.TrimSpace(strings.TrimPrefix(line, "@"))
		q.Feedback = feedback
		return
	}
}

func handleOptionFeedback(line string, q *Question) bool {
	if strings.HasPrefix(line, "@@") {
		// Remove "@@" from the start of the line and trim spaces
		feedback := strings.TrimSpace(strings.TrimPrefix(line, "@@"))
		// Store the feedback for the last option
		q.OptionFeedback[q.OptionKey] = feedback
		return true
	}
	return false
}

func printQuestions(questions []Question, prefixes []string) {
	for _, q := range questions {
		fmt.Println()
		fmt.Printf("QuestionType,%s\n", q.Type)
		if q.ID != 0 {
			fmt.Printf("ID,%d\n", q.ID)
		} else {
			fmt.Println("ID,")
		}
		fmt.Println("Title,")
		fmt.Printf("QuestionText,%s\n", q.Text)
		fmt.Println("Points,")
		fmt.Println("Difficulty,")
		if q.Image != "" {
			fmt.Printf("Image,images/%s\n", q.Image)
		}

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
				feedback := q.OptionFeedback[key]
				if removeListPrefixes {
					fmt.Printf("Option,%s,%s,%s\n", score, option, feedback)
				} else {
					fmt.Printf("Option,%s,%s %s,%s\n", score, key, option, feedback)
				}
			}
		} else if q.Type == "TF" {
			q.Answer = strings.TrimSpace(q.Answer)
			if strings.ToLower(q.Answer) == "true" || strings.ToLower(q.Answer) == "t" || strings.ToLower(q.Answer) == "true true" {
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
		fmt.Printf("Feedback,%s\n", q.Feedback)

	}
}

func writeQuestionsToCSV(questions []Question, prefixes []string) string {
	// Create a strings.Builder to hold the CSV data
	var b strings.Builder

	// Loop through the questions
	for _, q := range questions {
		fmt.Fprintf(&b, "QuestionType,%s\n", q.Type)
		if q.ID != 0 {
			fmt.Fprintf(&b, "ID,%d\n", q.ID)
		} else {
			b.WriteString("ID,\n")
		}

		b.WriteString("Title,\n")
		fmt.Fprintf(&b, "QuestionText,%s\n", q.Text)
		b.WriteString("Points,\n")
		b.WriteString("Difficulty,\n")
		if q.Image != "" {
			fmt.Fprintf(&b, "Image,images/%s\n", q.Image)
		}
		if q.Type == "WR" {
			b.WriteString("InitialText,\n")
			b.WriteString("AnswerKey,\n")
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
					feedback := q.OptionFeedback[key] // Get the feedback for the option
					if removeListPrefixes {
						fmt.Fprintf(&b, "Option,%s,%s,%s\n", score, option, feedback)
					} else {
						fmt.Fprintf(&b, "Option,%s,%s %s,%s\n", score, key, option, feedback)
					}
				}

			} else if q.Type == "TF" {
				q.Answer = strings.TrimSpace(q.Answer)
				if strings.ToLower(q.Answer) == "true" || strings.ToLower(q.Answer) == "t" || strings.ToLower(q.Answer) == "true true" {
					b.WriteString("TRUE,100\n")
					b.WriteString("FALSE,0\n")
				} else {
					b.WriteString("TRUE,0\n")
					b.WriteString("FALSE,100\n")
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
					fmt.Fprintf(&b, "Answer,100,%s\n", answer)
				}
			}

			b.WriteString("Hint,\n")
			fmt.Fprintf(&b, "Feedback,%s\n", q.Feedback)

			// Add an empty line after each question
			b.WriteString("///\n")
		}
	}

	// Return the CSV data as a string
	return b.String()
}
