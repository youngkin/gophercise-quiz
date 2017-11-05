package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	fname := flag.String("fname", "problems.csv", "Provide the name of the  input file")
	timeout := flag.Int("timeout", 30, "Max time (in seconds) to complete the quiz")
	flag.Parse()

	f, err := os.Open(*fname)
	if err != nil {
		fmt.Printf("Error opening file %s: %s\n", *fname, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	in := bufio.NewReader(os.Stdin)
	var total, correct int

	timer := time.NewTimer(time.Second * time.Duration(*timeout))

Loop:
	for scanner.Scan() {
		select {
		case <-timer.C:
			fmt.Println("Time's up!")
			break Loop
		default:
			total++
			qna := strings.Split(scanner.Text(), ",")
			fmt.Print(qna[0], " = ")
			ans, _ := in.ReadString('\n')
			ans = strings.TrimRight(ans, "\n")
			if ans == qna[1] {
				correct++
			}
		}
	}

	fmt.Printf("\nTotal number of questions: %d, correct answers: %d, score:%d\n",
		total, correct, correct/total*100)
}
