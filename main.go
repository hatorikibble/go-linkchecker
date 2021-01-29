package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"net/http"

	//"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var file string

func init() {
	// init random generator
	rand.Seed(time.Now().UnixNano())
	// read Flags
	flag.StringVar(&file, "file", "urls.csv", "CSV file with URLs")
  flag.Parse()
}

type link struct {
	url  string
	name string
}

func checkLink(wg *sync.WaitGroup, checkLink link, nr int) {
	defer wg.Done()

	//resp, err := http.Get(checkLink.url)
	//if err != nil {
	//	exit(fmt.Sprintf("Problem %s %s", checkLink.url, resp.Status))
	//}
	//defer resp.Body.Close()

	resp, err := http.Get(checkLink.url)
	if err != nil{
		fmt.Printf("%3d Fehler bei %s, Status %s: %s\n",nr, checkLink.url,err)
	  return
	}
		fmt.Printf("%3d Checking %s (%s): %s\n", nr, checkLink.name, checkLink.url,resp.Status)
  defer resp.Body.Close()

}

func parseLines(lines [][]string) []link {
	ret := make([]link, len(lines))
	for i, line := range lines {
		ret[i] = link{name: strings.TrimSpace(line[0]),
			url: strings.TrimSpace(line[2]),
		}
	}
	return ret
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func main() {

	// open the csv file
	csvfile, err := os.Open(file)
	if err != nil {
		exit(fmt.Sprintf("Couldn't open the csv file: %s", file))
	}

	// Parse the file
	r := csv.NewReader(csvfile)
	lines, err := r.ReadAll()

	if err != nil {
		exit("Failed to parse CSV file")
	}

	links := parseLines(lines)

	var wg sync.WaitGroup

	count := 0
	for i, l := range links {
		wg.Add(1)
		go checkLink(&wg, l, i+1)
		count++
	}

	wg.Wait()
	fmt.Printf("Checked %d Links..\n", count)

}
