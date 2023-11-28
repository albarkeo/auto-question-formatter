# Question Formatter
Go program to format minimally altered input text directly from Word to a Brightspace import ready CSV File
It currently supports:
- Multiple Choice (MC)
- Short Answer (SA)
- True or False (TF)

With the next planned implementation being Written Response (WR)

See a text only version here https://go.dev/play/p/FTMU7afwqd-, (Press run, then you can copy out the output and use Text to Columns in Excel using a comma as the delimiter)

## Preformatting Requirements
Add a "---" or "+++" between each question

## Example of Accepted Inputs
```
What is the capital of Australia?

a. Sydney
b. Melbourne
*c. Canberra
d. Adelaide
+++

What is the capital of Australia?

a. Sydney
b. Melbourne
c. Canberra
d. Adelaide

Answer c
+++
What colour is the sky?
Azure or blue; orange
+++
Write the value of the 5 in the number 8526.	500
+++
Write the value of the 7 in the number 97 450.	7000 or 7 000 or 7,000
+++
Write this as a number. 7 tens of thousands, 4 hundreds and 2 ones	70402 or 70 402 or 70,402
+++
9 + 6 =		15
+++
16 + 7 =	23
---
24 + 5 =	29
---
628 - 284 =  344
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
```

## Accepted Input Types
The below table shows an example of accepted input types for questions and answers, it's a bit of a mess but I'm working on it:

| Multiple Choice | Multiple Choice | Short Answer | Short Answer | True or False | True or False | Written Response | Written Response |
|---|---|---|---|---|---|---|---|
| *Accepted Inputs* | *Accepted Answer Inputs* | *Accepted Inputs* | *Accepted Answer Inputs* | *Accepted Inputs* | *Accepted Answer Inputs* | *Accepted Inputs* | *Accepted Answer Inputs* |
|  Question text<br>a<br>*b<br>c<br>d | *b | Question text | Single answer | Question text | TRUE | Question text | *None required* |
|  Question text<br>1<br>2<br>3<br>4<br><br>  Correct answer 2 | Correct answer 2 | Question text ending in a question mark? | answer 1 or answer 2 or answer 3 | Question text ending in a question mark? | T | Question text ending in a question mark? |  |
| Question text<br>w<br>x<br>y<br>z<br><br>  Answer x | Answer x | Question text | answer 1; answer 2; answer 3 | | true  |  |  |
| Question text<br>a<br>b<br>c<br>d<br>e<br>f<br>... | correct answer: b |  | answer 1 or answer 2; answer 3 |  |  |  |  |

### Prefixes for Answers
"answer ", "answer: ", "answer- ", "answers ", "answers: ", "answers- ", "correct answer ", "correct answer: ", "correct answer- ", "correct answers: ", "correct answers- " "*"

### Allowed Variations
It should work with most enters and line breaks.
Tabs are consider as a potential new line and therefore option or answer for the question. 
