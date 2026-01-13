package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

type Result struct {
	keyword string
	count   int
}

func worker(
	lines <-chan string,
	results chan<- Result,
	keywords []string,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for line := range lines {
		text := strings.ToLower(line)
		for _, kw := range keywords {
			c := strings.Count(text, kw)
			if c > 0 {
				results <- Result{keyword: kw, count: c}
			}
		}
	}
}

func main() {
	keywords := []string{"go", "concurrency", "goroutines", "channels"}

	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	lines := make(chan string)
	results := make(chan Result)

	var wg sync.WaitGroup
	workerCount := 4

	// Start worker goroutines
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go worker(lines, results, keywords, &wg)
	}

	// Read file and send lines
	go func() {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines <- scanner.Text()
		}
		close(lines)
	}()

	// Close results when workers finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	total := make(map[string]int)
	for r := range results {
		total[r.keyword] += r.count
	}

	fmt.Println("Keyword counts:")
	for _, kw := range keywords {
		fmt.Printf("%-12s : %d\n", kw, total[kw])
	}
}
