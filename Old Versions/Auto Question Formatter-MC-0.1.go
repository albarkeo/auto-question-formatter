package main

import (
	"bufio"
	"fmt"
	"strings"
)

func isTFQuestion(answers []string) bool {
	for _, answer := range answers {
		answer = strings.TrimSpace(answer)
		if strings.HasPrefix(answer, "Answer: ") || strings.HasPrefix(answer, "Answer - ") || strings.HasPrefix(strings.ToLower(answer), "answer") {
			answer = strings.TrimPrefix(answer[8:], "Answer: ")
			answer = strings.TrimPrefix(answer[8:], "Answer - ")
			answer = strings.TrimPrefix(strings.ToLower(answer[6:]), "answer")
		}
		if answer == "True" || answer == "False" || answer == "true" || answer == "false" {
			return true
		}
	}
	return false
}

func printTFQuestion(question string, answers []string) {
	fmt.Println()
	fmt.Println("NewQuestion,TF")
	fmt.Println("ID,")
	fmt.Println("Title,This is a True/False question")
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


func isMCQuestion(options []string) bool {
    // Check if there are at least two options
    if len(options) < 2 {
        return false
    }

    // Check if each option line has the expected format
    for _, option := range options {
        option = strings.TrimSpace(option)
        if len(option) < 3 || !(option[1] == '.' || option[1] == ' ' || option[1] == ')' || option[1] == '-') ||
            !((option[0] >= 'a' && option[0] <= 'z') ||
                (option[0] >= 'A' && option[0] <= 'Z') ||
                (option[0] >= '1' && option[0] <= '9')) {
            return false
        }
    }

    return true
}

func isCorrectAnswer(line string) (bool, string) {
    line = strings.TrimSpace(line)

    if strings.HasPrefix(strings.ToLower(line), "correct answer") ||
       strings.HasPrefix(strings.ToLower(line), "answer") {
        // Remove the "Correct answer" or "Answer" prefix and return
        line = strings.TrimPrefix(line, "Correct answer: ")
        line = strings.TrimPrefix(line, "Answer: ")
        line = strings.TrimPrefix(strings.ToLower(line), "answer")
        return true, line
    }

    // Check if the line has only one character (A-Z, a-z, 0-9)
    if len(line) == 1 && ((line[0] >= 'a' && line[0] <= 'z') ||
        (line[0] >= 'A' && line[0] <= 'Z') ||
        (line[0] >= '1' && line[0] <= '9')) {
        return true, line
    }

    // Check if the line has two characters and the second character is ')' or '-'
    if len(line) == 2 && (line[1] == ')' || line[1] == '-') &&
        ((line[0] >= 'a' && line[0] <= 'z') ||
            (line[0] >= 'A' && line[0] <= 'Z') ||
            (line[0] >= '1' && line[0] <= '9')) {
        return true, string(line[0])
    }

    return false, ""
}

func printMCQuestion(question string, options []string, correctAnswer string) {
    fmt.Println()
    fmt.Println("NewQuestion,MC")
    fmt.Println("ID,")
    fmt.Println("Title,")
    fmt.Printf("QuestionText,%s\n", question)
    fmt.Println("Points,")
    fmt.Println("Difficulty,")
    fmt.Println("Image,")

    for _, option := range options {
        if strings.TrimSpace(option[2:]) == correctAnswer {
            fmt.Printf("Option,,100,%s\n", option[2:])
        } else {
            fmt.Printf("Option,,0,%s\n", option[2:])
        }
    }

    fmt.Println("Hint,")
    fmt.Println("Feedback,\n")
}


func printSAQuestion(question string, answers []string) {
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
		splitAnswers := strings.Split(answer, "or")
		for _, splitAnswer := range splitAnswers {
			fmt.Printf("Answer,,%s\n", strings.TrimSpace(splitAnswer))
		}
	}

	fmt.Println("Hint,")
	fmt.Println("Feedback,")
}

func main() {
	input := `This is the question text for MC1

a. This is the correct answer
b. This is incorrect answer 1
c. This is incorrect answer 2
d. This is incorrect answer 3

Correct answer a.

This is the question text for MC2

1 This is an incorrect answer
2 This is an incorrect answer
3 This is the correct answer
4 This is an incorrect answer

3

This is the question text for MC3

A This is the correct answer
B This is incorrect answer 1
C This is incorrect answer 2
D This is incorrect answer 3

A`

scanner := bufio.NewScanner(strings.NewReader(input))
var question string
var answers []string
var options []string
var correctAnswer string

for scanner.Scan() {
    line := strings.TrimSpace(scanner.Text())
    fmt.Println("Read line:", line)  // Print the line that's being read
    if line == "" {
        continue
    }
	if ok, ans := isCorrectAnswer(line); ok {
        correctAnswer = ans
        continue
    }
    if len(line) == 1 || isMCQuestion([]string{line}) {
        if question != "" && len(options) > 0 && isMCQuestion(options) {
            fmt.Println("Calling printMCQuestion with question:", question, "options:", options, "correctAnswer:", correctAnswer)  // Print the variables before calling printMCQuestion
            printMCQuestion(question, options, correctAnswer)
            question = ""
            options = nil
            correctAnswer = ""
        } else if question != "" && len(answers) > 0 {
            if isTFQuestion(answers) {
                fmt.Println("Calling printTFQuestion with question:", question, "answers:", answers)  // Print the variables before calling printTFQuestion
                printTFQuestion(question, answers)
            } else {
                fmt.Println("Calling printSAQuestion with question:", question, "answers:", answers)  // Print the variables before calling printSAQuestion
                printSAQuestion(question, answers)
            }
            question = ""
            answers = nil
        }
        question = line
    } else {
        options = append(options, line)
        answers = append(answers, line)
    }
}

if err := scanner.Err(); err != nil {
    log.Fatal(err)
}

if question != "" && len(options) > 0 && isMCQuestion(options) {
    fmt.Println("Calling printMCQuestion with question:", question, "options:", options, "correctAnswer:", correctAnswer)  // Print the variables before calling printMCQuestion
    printMCQuestion(question, options, correctAnswer)
} else if question != "" && len(answers) > 0 {
    if isTFQuestion(answers) {
        fmt.Println("Calling printTFQuestion with question:", question, "answers:", answers)  // Print the variables before calling printTFQuestion
        printTFQuestion(question, answers)
    } else {
        fmt.Println("Calling printSAQuestion with question:", question, "answers:", answers)  // Print the variables before calling printSAQuestion
        printSAQuestion(question, answers)
    }
}



}