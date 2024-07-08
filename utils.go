package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type quizStats struct {
	score int
	total int
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

func Quiz(csvRecords [][]string, qs *quizStats, c chan string) {
	for _, record := range csvRecords {
		var ans string
		fmt.Printf("%v: ", record[0])
		_, err := fmt.Scanln(&ans)
		if err != nil {
			log.Fatal(err)
		}
		scoreQuestion(record, ans, qs)
	}
	c <- "completed"
}

func scoreQuestion(q []string, ans string, qs *quizStats) {
	if cleanString(q[1]) == cleanString(ans) {
		qs.score++
	}
}

func cleanString(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}

func showStats(qs *quizStats) {
	fmt.Printf("%v / %v", qs.score, qs.total)
}

func exitWhenReachesLimit(l int64, c chan string) {
	time.Sleep(time.Duration(l) * time.Second)
	c <- "Time Limit Exceeded"
}
