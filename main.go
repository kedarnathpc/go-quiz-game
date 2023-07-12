package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type problem struct {
	que string
	ans string
}

func main() {
	// create a flag variable to take input from the user in the terminal
	// the file will be in the format of .csv
	// The return value is the address of a string variable that stores the value of the flag.
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question answer'")
	flag.Parse()

	// open the file
	file, err := os.Open(*csvFilename)

	// check for error
	if err != nil {
		fmt.Printf("Failed to open the CSV file: %s.\n", *csvFilename)
		os.Exit(1)
	}

	// create a reader to read the contents of the file
	r := csv.NewReader(file)
	// ReadAll reads all the remaining records from r. Each record is a slice of fields.
	lines, err := r.ReadAll()
	if err != nil {
		fmt.Println("Failed to parse the provided CSV file.")
	}

	problems := parseLines(lines)
	askQuestions(problems)
}

// parseLines function takes a 2D slice of strings,
// where each row represents a question and answer pair,
// and converts it into a slice of problem structs,
// where each struct holds the question and answer extracted from the input.
func parseLines(lines [][]string) []problem {

	// create a new slice of type struct with length = no. of rows in the lines slice
	ret := make([]problem, len(lines))

	for i, line := range lines {

		// assign a new struct to the ith element of ret
		ret[i] = problem{

			// assign the question to the q field of struct
			que: line[0],

			// avoid errors if there are any spaces in the answers
			// assign the answer to the a field of struct
			ans: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

func askQuestions(problems []problem) {

	// create a time limit flag to pass the custom time limit for the quiz
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")

	// convert into seconds
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	correct := 0

	// range over the slice of problems
problemloop:
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.que)

		// make a channel of string to recieve the user answer input
		answerCh := make(chan string)

		// fire a go routine to take the user input
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		// one of two cases :
		// 1. timer channel will recieve default time and end the quiz
		// 2. answer will recieve value from answer channel
		select {
		case <-timer.C:
			fmt.Println()
			break problemloop
		case answer := <-answerCh:

			// check if the answer is correct
			if answer == p.ans {
				correct++
			}
		}
	}

	// print the score
	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}
