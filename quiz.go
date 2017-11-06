package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/google/go-cmp/cmp"
)

// Improvements:
//		1. Use csv.Reader() as described below to read the input CSV file
//		2. Use fmt.Scanf() as described below to get input from stdin
//		3. Could have used a "problems" struct to package the question and answer
//		4. Could have written an "exit" function. Arguably, this isn't much of, if any, improvement

func main() {
	fname := flag.String("fname", "problems.csv", "Provide the name of the  input file")
	timeout := flag.Int("timeout", 30, "Max time (in seconds) to complete the quiz")
	flag.Parse()

	f, err := os.Open(*fname)
	if err != nil {
		fmt.Printf("Error opening file %s: %s\n", *fname, err)
		os.Exit(-1)
	}
	defer f.Close()

	// Could have used csv.Reader and would have avoided all the parsing in my solution
	// E.g.,
	// 		r := csv.NewReader(f)
	// 		lines := r.ReadAll()
	var bytes []byte
	if bytes, err = ioutil.ReadAll(f); err != nil {
		fmt.Printf("Error reading bytes from file %s: Error %s", *fname, err)
		os.Exit(-2)
	}

	// Get set of questions delimited by newlines
	quests := strings.Split(string(bytes), "\n")
	total := float32(len(quests))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(*timeout))
	defer cancel()

	timer := time.NewTimer(time.Second * time.Duration(*timeout))
	c := make(chan string)
	var correct float32

Loop:
	for _, qna := range quests {
		q := qna[:strings.Index(qna, ",")]
		a := qna[strings.Index(qna, ",")+1:]

		go getUserInput(ctx, q, c)

		select {
		case <-timer.C:
			fmt.Println("\n\nTime's up!")
			break Loop
		case ans := <-c:
			ans = strings.TrimRight(ans, "\n")
			if cmp.Equal(ans, a) {
				correct++
			}
		}
	}

	close(c)
	fmt.Printf("\nTotal number of questions: %.f, correct answers: %.f, score:%2.f percent\n",
		total, correct, float32((correct/total)*100))
}

func getUserInput(ctx context.Context, q string, c chan<- string) {
	fmt.Print(q, " = ")

	// Could have used fmt.Scanf to avoid creating a new Reader and then calling ReadString()
	// E.g.,
	//		fmt.Scanf("%s\n", &answer)
	// will read a string ending in "\n" from stdin and scan it into &answer
	in := bufio.NewReader(os.Stdin)
	ans, _ := in.ReadString('\n')
	select {
	case c <- ans:
		return
	case <-ctx.Done():
		return
	}
}
