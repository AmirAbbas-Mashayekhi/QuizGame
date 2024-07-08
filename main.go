package main

import "fmt"

func main() {
	// initialization
	c := make(chan string)
	csvPtr, limitPtr := parseFlags()
	records := readCSV(*csvPtr)
	qs := quizStats{score: 0, total: len(records)}

	// flow
	go Quiz(records, &qs, c)
	go exitWhenReachesLimit(*limitPtr, c)

	// final
	fmt.Println(<-c)
	showStats(&qs)
}
