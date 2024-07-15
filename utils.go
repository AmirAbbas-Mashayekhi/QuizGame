package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"log"
	"os"
	"strings"
	"time"
)

// Types
type quizStats struct {
	score int
	total int
}

type problem struct {
	question string
	answer   string
}

func createProblemSet(csvRecords [][]string) []problem {
	var problemSet []problem
	for _, record := range csvRecords {
		problemSet = append(
			problemSet,
			problem{
				question: record[0],
				answer:   record[1],
			},
		)
	}
	return problemSet
}

func parseFlags() (*string, *int64) {
	csvPtr := flag.String(
		"csv",
		"problems.csv",
		"a csv file containing the problems fmt -> 'question,answer'",
	)
	timeLimitPtr := flag.Int64(
		"limit",
		30,
		"the time limit for the quiz in seconds",
	)
	flag.Parse()
	return csvPtr, timeLimitPtr
}

func readCSV(filename string) [][]string {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal("Unable to read input file "+filename, err)
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(f)

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filename, err)
	}

	return records
}

func Quiz(problems []problem, qs *quizStats, c chan string) {
	for _, prob := range problems {
		var ans string
		fmt.Printf("%v: ", prob.question)
		_, err := fmt.Scanln(&ans)
		if err != nil {
			log.Fatal(err)
		}
		scoreQuestion(prob, ans, qs)
	}
	c <- color.BlueString("completed")
}

func scoreQuestion(p problem, ans string, qs *quizStats) {
	if cleanString(p.answer) == cleanString(ans) {
		qs.score++
	}
}

func cleanString(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}

func showStats(qs *quizStats) {
	fmt.Print(color.GreenString("%v / ", qs.score))
	fmt.Print(color.WhiteString("%v\n", qs.total))
}

func timer(l int64, c chan string) {
	time.Sleep(time.Duration(l) * time.Second)
	c <- color.RedString("Time Limit Exceeded")
}
