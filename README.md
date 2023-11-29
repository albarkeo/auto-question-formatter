# Question Formatter
A program written in Go to format minimally altered input text directly from Word to a Brightspace import-ready CSV File.

It currently supports:
- Multiple Choice (MC)
- Short Answer (SA)
- True or False (TF)
- Written Response (WR)

It trades extra formatting in Word for a few features available in other programs such as Respondus and Brightspace's own quiz converter: https://community.d2l.com/brightspace/kb/articles/4161-quiz-question-converter

It is limited to focusing on the main bulk question types we receive.

See a text-only version (v1.51) here https://go.dev/play/p/FTMU7afwqd-, (Copy in Questions, press run, then you can copy the output and use Text to Columns in Excel using commas as the delimiter)

## Preformatting Requirements
Add a "---" or "+++" between each question

## Example of Accepted Inputs
Use the following sample input as a test in the text converter (https://go.dev/play/p/FTMU7afwqd-) or with the latest version of the program (download the .exe, install and run, follow the text prompts)

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
9 + 6 =
15
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

### Prefixes for Answers
The following will be removed from an answer if written in the Word document (not case sensitive):
"answer ", "answer: ", "answer- ", "answers ", "answers: ", "answers- ", "correct answer ", "correct answer: ", "correct answer- ", "correct answers: ", "correct answers- " "*"

### Prefixes for Options
By default removeListPrefixes = true

This removes numbered list prefixes such as a,b,c,d or 1), 2), 3), 4), or A-, B-, C-, D- or w., x., y., z. etc when printing the options to allow for Brightspaces own ordering to be used.

The second character can be ' ' ')' '.' '-' or ' -'

### Currently Accepted Variations
It should work with most enters and line breaks.
Tabs are considered as a potential new line and therefore option or answer for the question.

### Short Answer Question Type Answers
These can be separated by the word 'or' or a semicolon ';'

e.g. 
```
What colour is the sky?
Azure or blue; orange; light blue or sky blue
```


## Example Output
|NewQuestion|MC| |
|:----|:----|:----|
|ID| | |
|Title| | |
|QuestionText|What is the capital of Australia?| |
|Points| | |
|Difficulty| | |
|Image| | |
|Option|0|Sydney|
|Option|0|Melbourne|
|Option|100|Canberra|
|Option|0|Adelaide|
|Hint| | |
|Feedback| | |
| | | |
|NewQuestion|MC| |
|ID| | |
|Title| | |
|QuestionText|What is the capital of Australia?| |
|Points| | |
|Difficulty| | |
|Image| | |
|Option|0|Sydney|
|Option|0|Melbourne|
|Option|100|Canberra|
|Option|0|Adelaide|
|Hint| | |
|Feedback| | |
| | | |
|NewQuestion|SA| |
|ID| | |
|Title| | |
|QuestionText|Write the value of the 5 in the number 8526.| |
|Points| | |
|Difficulty| | |
|Image| | |
|Answer|100|500|
|Hint| | |
|Feedback| | |
| | | |
|NewQuestion|SA| |
|ID| | |
|Title| | |
|QuestionText|Write the value of the 7 in the number 97 450.| |
|Points| | |
|Difficulty| | |
|Image| | |
|Answer|100|7000|
|Answer|100|7 000|
|Answer|100|7,000|
|Hint| | |
|Feedback| | |
| | | |
|NewQuestion|SA| |
|ID| | |
|Title| | |
|QuestionText|Write this as a number. 7 tens of thousands, 4 hundreds and 2 ones| |
|Points| | |
|Difficulty| | |
|Image| | |
|Answer|100|70402|
|Answer|100|70 402|
|Answer|100|70,402|
|Hint| | |
|Feedback| | |
| | | |
|NewQuestion|SA| |
|ID| | |
|Title| | |
|QuestionText|9 + 6 =| |
|Points| | |
|Difficulty| | |
|Image| | |
|Answer|100|15|
|Hint| | |
|Feedback| | |
| | | |
|NewQuestion|SA| |
|ID| | |
|Title| | |
|QuestionText|16 + 7 =| |
|Points| | |
|Difficulty| | |
|Image| | |
|Answer|100|23|
|Hint| | |
|Feedback| | |
| | | |
|NewQuestion|SA| |
|ID| | |
|Title| | |
|QuestionText|24 + 5 =| |
|Points| | |
|Difficulty| | |
|Image| | |
|Answer|100|29|
|Hint| | |
|Feedback| | |
| | | |
|NewQuestion|TF| |
|ID| | |
|Title| | |
|QuestionText|The Earth is the only planet in our solar system with liquid water on its surface. | |
|Points| | |
|Difficulty| | |
|Image| | |
|TRUE|0| |
|FALSE|100| |
|Hint| | |
|Feedback| | |
| | | |
|NewQuestion|SA| |
|ID| | |
|Title| | |
|QuestionText|What is the colour of grass?| |
|Points| | |
|Difficulty| | |
|Image| | |
|Answer|100|green|
|Answer|100|brown|
|Hint| | |
|Feedback| | |
| | | |
|NewQuestion|SA| |
|ID| | |
|Title| | |
|QuestionText|Short answer no question mark| |
|Points| | |
|Difficulty| | |
|Image| | |
|Answer|100|a short answer|
|Hint| | |
|Feedback| | |
| | | |
|NewQuestion|MC| |
|ID| | |
|Title| | |
|QuestionText|Which planet is known as the Red Planet?| |
|Points| | |
|Difficulty| | |
|Image| | |
|Option|0|Venus|
|Option|100|Mars|
|Option|0|Jupiter|
|Option|0|Saturn|
|Hint| | |
|Feedback| | |
| | | |
|NewQuestion|TF| |
|ID| | |
|Title| | |
|QuestionText|Mars is known as the Red Planet| |
|Points| | |
|Difficulty| | |
|Image| | |
|TRUE|100| |
|FALSE|0| |
|Hint| | |
|Feedback| | |
| | | |
|NewQuestion|TF| |
|ID| | |
|Title| | |
|QuestionText|Venus is known as the Red Planet| |
|Points| | |
|Difficulty| | |
|Image| | |
|TRUE|0| |
|FALSE|100| |
|Hint| | |
|Feedback| | |
| | | |
|NewQuestion|TF| |
|ID| | |
|Title| | |
|QuestionText|Jupiter is known as the Red Planet| |
|Points| | |
|Difficulty| | |
|Image| | |
|TRUE|0| |
|FALSE|100| |
|Hint| | |
|Feedback| | |
| | | |
|NewQuestion|MC| |
|ID| | |
|Title| | |
|QuestionText|Who wrote the novel "Pride and Prejudice"?| |
|Points| | |
|Difficulty| | |
|Image| | |
|Option|0|Charles Dickens|
|Option|100|Jane Austen|
|Option|0|Mark Twain|
|Option|0|George Orwell|
|Hint| | |
|Feedback| | |
| | | |
|NewQuestion|MC| |
|ID| | |
|Title| | |
|QuestionText|What is the square root of 81?| |
|Points| | |
|Difficulty| | |
|Image| | |
|Option|0|8|
|Option|100|9|
|Option|0|10|
|Option|0|11|
|Option|0|999|
|Hint| | |
|Feedback| | |
| | | |
|NewQuestion|TF| |
|ID| | |
|Title| | |
|QuestionText|True or false, this has been a difficult but rewarding process| |
|Points| | |
|Difficulty| | |
|Image| | |
|TRUE|100| |
|FALSE|0| |
|Hint| | |
|Feedback| | |
| | | |
|NewQuestion|WR| |
|ID| | |
|Title| | |
|QuestionText|What colour is the sky?| |
|Points| | |
|Difficulty| | |
|Image| | |
|InitialText| | |
|AnswerKey| | |


## Accepted Input Types
The below table shows an example of accepted input types for questions and answers, it's a bit of a mess, the previous examples might be clearer:

| Multiple Choice | Multiple Choice | Short Answer | Short Answer | True or False | True or False | Written Response | Written Response |
|---|---|---|---|---|---|---|---|
| *Accepted Inputs* | *Accepted Answer Inputs* | *Accepted Inputs* | *Accepted Answer Inputs* | *Accepted Inputs* | *Accepted Answer Inputs* | *Accepted Inputs* | *Accepted Answer Inputs* |
|  Question text<br><br>a<br>*b<br>c<br>d | *b | Question text<br><br>Single answer | Single answer | Question text<br><br>TRUE | TRUE | Question text | *None required* |
|  Question text<br><br>1<br>2<br>3<br>4<br><br>  Correct answer 2 | Correct answer 2 | Question text ending in a question mark?<br><br>answer 1 or answer 2 or answer 3 | answer 1 or answer 2 or answer 3 | Question text ending in a question mark?<br><br>T | T | Question text ending in a question mark? |  |
| Question text<br><br>w<br>x<br>y<br>z<br><br>  Answer x | Answer x | Question text<br><br>answer 1; answer 2; answer 3 | answer 1; answer 2; answer 3 | | true  |  |  |
| Question text<br><br>a<br>b<br>c<br>d<br>e<br>f<br>...<br><br>correct answer: b | correct answer: b |  | answer 1 or answer 2; answer 3 |  | FaLsE  |  |  |


*Developed by Alex Barnes-Keoghan*
