package main

import (
	"bufio"
	"fmt"
	"log"
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
    fmt.Println("isMCQuestion received options:", options)  // Print the received options
    // Check if there are at least two options
    if len(options) < 2 {
        fmt.Println("isMCQuestion returns: false")  // Print the return value
        return false
    }

    // Check if each option line has the expected format
    var lastOptionPrefix rune
    for _, option := range options {
        option = strings.TrimSpace(option)
        if len(option) < 3 {
            fmt.Println("isMCQuestion returns: false")  // Print the return value
            return false
        }
        optionPrefix := rune(option[0])
        if lastOptionPrefix != 0 && optionPrefix != lastOptionPrefix+1 {
            fmt.Println("isMCQuestion returns: false")  // Print the return value
            return false
        }
        lastOptionPrefix = optionPrefix
    }

    fmt.Println("isMCQuestion returns: true")  // Print the return value
    return true
}


func isCorrectAnswer(line string) (bool, string) {
    fmt.Println("isCorrectAnswer received line:", line)  // Print the received line
    line = strings.TrimSpace(line)

    if strings.HasPrefix(strings.ToLower(line), "correct answer ") ||
        strings.HasPrefix(strings.ToLower(line), "answer ") {
        // Remove the "Correct answer " or "Answer " prefix and return
        line = strings.TrimPrefix(line, "Correct answer ")
        line = strings.TrimPrefix(line, "Answer ")
        line = strings.TrimPrefix(strings.ToLower(line), "answer ")
        fmt.Println("isCorrectAnswer returns:", true, line)  // Print the return values
        return true, line
    }

    // Check if the line has only one character (A-Z, a-z, 0-9)
    if len(line) == 1 && ((line[0] >= 'a' && line[0] <= 'z') ||
        (line[0] >= 'A' && line[0] <= 'Z') ||
        (line[0] >= '1' && line[0] <= '9')) {
 //       fmt.Println("isCorrectAnswer returns:", true, line)  // Print the return values
        return true, line
    }

    // Check if the line has two characters and the second character is ')' or '-'
    if len(line) == 2 && (line[1] == ')' || line[1] == '-') &&
        ((line[0] >= 'a' && line[0] <= 'z') ||
            (line[0] >= 'A' && line[0] <= 'Z') ||
            (line[0] >= '1' && line[0] <= '9')) {
        fmt.Println("isCorrectAnswer returns:", true, string(line[0]))  // Print the return values
        return true, string(line[0])
    }

    fmt.Println("isCorrectAnswer returns:", false, "")  // Print the return values
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
	correctAnswerFound := false

	var lastOptionPrefix rune
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if ok, ans := isCorrectAnswer(line); ok {
			if question != "" && len(options) > 0 && isMCQuestion(options) {
				fmt.Println("Calling printMCQuestion with question:", question, "options:", options, "correctAnswer:", ans)
				printMCQuestion(question, options, ans)
				question = ""
				options = nil
			}
			correctAnswerFound = true
			continue
		}
		if correctAnswerFound {
			// Reset options and correctAnswerFound if a correct answer was found in the last iteration
			options = nil
			correctAnswerFound = false
		}
		// Check if the line is a new question or an option for the current question
		if len(line) == 1 || (len(line) > 2 && (line[1] == ' ')) {
			if len(line) == 1 {
				question = line
			} else {
				optionPrefix := rune(line[0])
				if lastOptionPrefix != 0 && optionPrefix != lastOptionPrefix+1 {
					continue
				}
				options = append(options, line)
				lastOptionPrefix = optionPrefix
				if isMCQuestion(options) {
					fmt.Println("isMCQuestion received options:", options)
				}
			}
		} else {
			answers = append(answers, line)
		}
	}
	
	if question != "" && len(options) > 0 && isMCQuestion(options) {
		fmt.Println("Calling printMCQuestion with question:", question, "options:", options, "correctAnswer:", correctAnswer)
		printMCQuestion(question, options, correctAnswer)
	} else if question != "" && len(answers) > 0 {
		if isTFQuestion(answers) {
			fmt.Println("Calling printTFQuestion with question:", question, "answers:", answers)
			printTFQuestion(question, answers)
		} else {
			fmt.Println("Calling printSAQuestion with question:", question, "answers:", answers)
			printSAQuestion(question, answers)
		}
	}

}
