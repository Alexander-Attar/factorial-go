package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

func worker(id int, jobs <-chan int, results chan<- map[string]string) {
	var URL *url.URL
	URL, err := url.Parse("http://localhost:12345/")
	if err != nil {
		panic(err)
	}

	for n := range jobs {
		parameters := url.Values{}
		parameters.Add("n", strconv.Itoa(n))
		URL.RawQuery = parameters.Encode()

		resp, err := http.Get(URL.String())
		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		results <- map[string]string{"n": strconv.Itoa(n), "a": string(body)}
	}
}

func sorter(id int, sortingJobs <-chan map[string]string, sortedResults chan<- string) {
	for j := range sortingJobs {
		byteSlice := []byte(j["a"])
		intSlice := make([]int, len(byteSlice))
		output := make([]string, len(byteSlice))

		// Built an integer slice from each byte character for sorting
		i := 0
		for _, a := range byteSlice {
			a, err := strconv.Atoi(string(a))
			if err != nil {
				panic(err)
			}

			intSlice[i] = a
			i++
		}

		sortable := sort.IntSlice(intSlice)
		sort.Sort(sortable)

		// Convert the sorted output to a string
		for number := range intSlice {
			output = append(output, strconv.Itoa(intSlice[number]))
		}

		sortedString := strings.Join(output, "")
		sortedResults <- fmt.Sprintf("N: %v\n  A: %v\n  S: %v\n", j["n"], j["a"], sortedString)
	}
}

func main() {
	jobs := make(chan int, 10)
	sortingJobs := make(chan map[string]string, 10)
	results := make(chan map[string]string)
	sortedResults := make(chan string)

	// Create 4 HTTP request workers
	for w := 1; w <= 4; w++ {
		go worker(w, jobs, results)
	}

	// Send jobs to workers for factorial numbers 30 - 40
	for n := 30; n <= 40; n++ {
		jobs <- n
	}
	close(jobs)

	// Create 2 sorting workers
	for s := 1; s <= 2; s++ {
		go sorter(s, sortingJobs, sortedResults)
	}

	// Send the results from the HTTP workers to the sorting workers
	for a := 0; a <= 10; a++ {
		sortingJobs <- <-results
	}
	close(sortingJobs)

	// Print the formatted results for the first 5 values
	for sr := 0; sr < 5; sr++ {
		fmt.Print(<-sortedResults)
	}
	close(results)
	close(sortedResults)
}
