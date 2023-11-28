package main

import (
	"bufio"
	"fmt"
	"strings"
)

func main() {
	input := `[INPUT QUESTIONS AND ANSWERS HERE]`

	scanner := bufio.NewScanner(strings.NewReader(input))
	var question string
	var answers []string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		// Split the line by tab
		splitLine := strings.Split(line, "\t")
		if len(splitLine) > 1 {
			// If there is a tab in the line, treat the first part as the question and the second part as the answer
			question = splitLine[0]
			answers = append(answers, splitLine[1])
		} else if question == "" {
			question = line
		} else {
			answers = append(answers, line)
		}

		if isTFQuestion(answers) {
			printTFQuestion(question, answers)
			question = ""
			answers = nil
		} else if len(answers) > 0 {
			printQuestion(question, answers)
			question = ""
			answers = nil
		}

	}

	if question != "" && len(answers) > 0 {
		if isTFQuestion(answers) {
			printTFQuestion(question, answers)
		} else {
			printQuestion(question, answers)
		}
	}
}

func printQuestion(question string, answers []string) {
	fmt.Println()
	fmt.Println("NewQuestion,SA")
	fmt.Println("ID,")
	fmt.Println("Title,")
	fmt.Printf("QuestionText,%s\n", question)
	fmt.Println("Points,")
	fmt.Println("Difficulty,")
	fmt.Println("Image,")
	fmt.Println("InputBox,")

	for _, answer := range answers {
		answer = strings.TrimSpace(answer)
		// Split the answer by "or"
		splitAnswers := strings.Split(answer, "or")
		for _, splitAnswer := range splitAnswers {
			// Trim the spaces and print each answer
			fmt.Printf("Answer,,%s\n", strings.TrimSpace(splitAnswer))
		}
	}

	fmt.Println("Hint,")
	fmt.Println("Feedback,")
}

func printTFQuestion(question string, answers []string) {
	fmt.Println()
	fmt.Println("NewQuestion,TF")
	fmt.Println("ID,")
	fmt.Println("Title,")
	fmt.Printf("QuestionText,%s\n", question)
	fmt.Println("Points,")
	fmt.Println("Difficulty,")
	fmt.Println("Image,")

	for _, answer := range answers {
		answer = strings.TrimSpace(answer)
		if strings.HasPrefix(answer, "Answer: ") || strings.HasPrefix(answer, "Answer - ") || strings.HasPrefix(strings.ToLower(answer), "answer") {
			answer = strings.TrimPrefix(answer, "Answer: ")
			answer = strings.TrimPrefix(answer, "Answer - ")
			answer = strings.TrimPrefix(strings.ToLower(answer), "answer")
		}
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

func isTFQuestion(answers []string) bool {
	for _, answer := range answers {
		answer = strings.TrimSpace(answer)
		if strings.HasPrefix(answer, "Answer: ") || strings.HasPrefix(answer, "Answer - ") || strings.HasPrefix(strings.ToLower(answer), "answer") {
			answer = strings.TrimPrefix(answer, "Answer: ")
			answer = strings.TrimPrefix(answer, "Answer - ")
			answer = strings.TrimPrefix(strings.ToLower(answer), "answer")
		}
		if answer == "True" || answer == "False" || answer == "true" || answer == "false" {
			return true
		}
	}
	return false
}
