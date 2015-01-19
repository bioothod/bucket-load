package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/bioothod/bucket-load/reader"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	path := flag.String("path", "", "Path to log file to parse")
	max_requests := flag.Int("requests", 1024, "Maximum number of pending reqeusts to remote server")
	addr := flag.String("addr", "", "Remote backrunner proxy address (format: address:port)")

	flag.Parse()

	if *path == "" {
		log.Fatalf("You must specify log file to parse")
	}

	if *addr == "" {
		log.Fatalf("You must specify remote backrunner proxy address to test")
	}

	f, err := os.Open(*path)
	if err != nil {
		log.Fatalf("Could not open log file '%s' to parse: %v\n", *path, err)
	}

	br := bufio.NewReaderSize(f, 10240)

	re, err := reader.NewReader()
	if err != nil {
		log.Fatalf("Could not initialize regex: %v\n", err)
	}

	ch := make(chan *reader.Entry, *max_requests)

	for i := 0; i < *max_requests; i++ {
		go func(i int) {
			client := &http.Client {}

			for entry := range ch {
				if strings.HasPrefix(entry.Path, "/nobucket_upload/") {
					continue
				}
				if strings.HasPrefix(entry.Path, "/upload/") {
					entry.Path = strings.Replace(entry.Path, "/upload/", "/get/", 1)
					continue
				}

				url := fmt.Sprintf("http://%s%s", *addr, entry.Path)

				res, err := client.Get(url)
				if err != nil {
					log.Fatalf("%d: url: %s: get error: %v\n", i, url, err)
				}

				data, err := ioutil.ReadAll(res.Body)
				res.Body.Close()
				if err != nil {
					log.Fatalf("%d: url: %s: could not read response: %v\n", i, url, err)
				}

				fmt.Printf("%d: url: %s, data-size: %d, log-data-size: %d\n", i, url, len(data), entry.Size)
			}
		}(i)
	}

	for {
		line, _, err := br.ReadLine()
		if err != nil {
			log.Fatalf("Completed reading log file '%s': %v", *path, err)
		}

		entry, err := re.Run(string(line))
		if err != nil {
			continue
		}

		ch <- entry
	}
}
