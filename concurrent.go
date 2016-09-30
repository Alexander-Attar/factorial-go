package main

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
)

func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func main() {

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	var c = make(chan *big.Int)
	go series(c)

	active := true
	go func() {
		<-sigs
		signal.Stop(sigs)
		active = false
	}()

	for active {
		var URL *url.URL
		URL, err := url.Parse("http://localhost:12345/")
		if err != nil {
			panic(err)
		}

		n := <-c
		parameters := url.Values{}
		parameters.Add("n", n.String())
		URL.RawQuery = parameters.Encode()

		resp, err := http.Get(URL.String())
		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		fmt.Printf("N: %v, A: %v\n", n.String(), string(body))
	}

	fmt.Println("See you later")
}
