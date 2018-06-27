package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/alecthomas/kingpin"
)

type LogEntry struct {
	URL       string
	CreatedAt int
}

var debug = kingpin.Flag("concurrency", "Describe concurrency").Int()

func main() {
	kingpin.Parse()

	filename := "log.csv"

	byt, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	buffer := []LogEntry{}

	lines := strings.Split(string(byt), "\n")
	for _, line := range lines {
		data := strings.Split(line, ",")
		if len(data) != 2 {
			break
		}

		timeVal, _ := strconv.Atoi(data[0])
		url := data[1]

		buffer = append(buffer, LogEntry{
			URL:       url,
			CreatedAt: timeVal,
		})
	}

	fmt.Printf("%+v\n", buffer)

	// var g errgroup.Group

	// for _, url := range urls {
	// 	// Launch a goroutine to fetch the URL.
	// 	url := url // https://golang.org/doc/faq#closures_and_goroutines
	// 	g.Go(func() error {
	// 		// Fetch the URL.
	// 		resp, err := http.Get(url)
	// 		if err == nil {
	// 			resp.Body.Close()
	// 		}
	// 		return err
	// 	})
	// }
	// // Wait for all HTTP fetches to complete.
	// if err := g.Wait(); err == nil {
	// 	fmt.Println("Successfully fetched all URLs.")
	// }

}
