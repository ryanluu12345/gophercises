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
	number   int
	question string
	answer   string
}

func main() {
	csvFileName := flag.String("csv", "problems.csv", "a csv file in the format of 'question, answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	// File operations and serializing
	lines := openAndReadFile(*csvFileName)
	problems := parseCsvToProblem(lines)

	// Variables for statistics
	var incorrectQuestions []problem
	totalQuestions := len(problems)
	totalCorrect := 0

	// Driver code for presenting questions and getting answers
	promptQuestionAndGetAnswer(&problems, timeLimit, &totalCorrect, &incorrectQuestions)

	// Prints statistics which represents how well the user did
	printEndingStatistics(totalCorrect, totalQuestions, incorrectQuestions)
}

func openAndReadFile(csvFileName string) [][]string {
	file, err := os.Open(csvFileName)

	if err != nil {
		exit(fmt.Sprintf("Failed to open the csv file: %s\n", csvFileName))
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()

	if err != nil {
		exit("Trouble reading file!")
	}

	return lines
}

func parseCsvToProblem(content [][]string) []problem {
	problemArr := make([]problem, len(content))

	for idx, items := range content {
		problemArr[idx] = problem{
			number:   idx + 1,
			question: strings.TrimSpace(items[0]),
			answer:   strings.TrimSpace(items[1])}
	}

	return problemArr
}

func promptQuestionAndGetAnswer(problems *[]problem, timeLimit, totalCorrect *int, incorrectQuestions *[]problem) {
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	for _, item := range *problems {
		fmt.Printf("Problem %d: %s\n", item.number, item.question)
		answerCh := make(chan string)

		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			return
		case answer := <-answerCh:
			if answer == item.answer {
				*totalCorrect++
			} else {
				*incorrectQuestions = append(*incorrectQuestions, item)
			}
		}
	}
}

func printEndingStatistics(totalCorrect, totalQuestions int, wrongAnswers []problem) {
	fmt.Printf("You got %d out of %d correct!\n", totalCorrect, totalQuestions)

	for _, item := range wrongAnswers {
		fmt.Printf("Problem %d: %s = %s\n", item.number, item.question, item.answer)
	}
}

func exit(message string) {
	fmt.Printf(message)
	os.Exit(1)
}
