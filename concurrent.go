package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type HttpResponse struct {
	url      string
	response *http.Response
	err      error
}

func concurrentRequests(urls []string, ch chan *HttpResponse) {
	for _, url := range urls {
		go func(url string) {
			resp, err := http.Get(url)
			if err != nil && resp != nil && resp.StatusCode == http.StatusOK {
				resp.Body.Close()
			}
			ch <- &HttpResponse{url, resp, err}
		}(url)
	}

	return
}

func main() {
	count := 0
	var URL *url.URL
	URL, err := url.Parse("http://localhost:12345/")
	if err != nil {
		panic(err)
	}

	ch := make(chan *HttpResponse)
	urls := make([]string, 10)

	// compose urls for the first 10 factorial numbers
	for n := range urls {
		parameters := url.Values{}
		parameters.Add("n", strconv.Itoa(n))
		URL.RawQuery = parameters.Encode()
		urls[n] = URL.String()
	}

	concurrentRequests(urls, ch)

	for count < 10 {
		r, ok := <-ch
		if !ok {
			break
		}

		if r.err != nil {
			panic(r.err)
		}
		body, err := ioutil.ReadAll(r.response.Body)
		if err != nil {
			panic(err)
		}
		var n = r.url[len(r.url)-1:] // pop the number n off the url
		fmt.Printf("N: %v, A: %v\n", n, string(body))
		count++
	}
}
