package main

import (
	"fmt"
	"strings"
)

type Question struct {
	Text    string
	Options []string
	Answer  string
	Type    string
}

func main() {
	questionsText := `What colour is the sky?
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
	
	true` // Your input text here

	lines := strings.Split(questionsText, "\n")
	var questions []Question
	var q Question
	var inQuestion bool

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			if inQuestion {
				questions = append(questions, q)
				q = Question{}
				inQuestion = false
			}
			continue
		}

		if strings.HasSuffix(line, "?") {
			q.Text = line
			inQuestion = true
		} else if inQuestion {
			q.Options = append(q.Options, line)
		}
	}

	// Output the questions for testing
	for _, q := range questions {
		fmt.Println("Question:", q.Text)
		fmt.Println("Options:", q.Options)
	}
}
