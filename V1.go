package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	VERSION      = 1
	INPUT_FOLDER = "./data"
	OUTPUT_FILE  = "./output.txt"
)

func main() {
	if len(os.Args) > 1 {
		INPUT_FOLDER = os.Args[len(os.Args)-1]
	}

	fmt.Printf("Starting program V%d!\n", VERSION)
	processDataFolder(INPUT_FOLDER)
}

/*
	Structure:
		list() -filePaths-> read() -rawData-> parse() -parsedData-> aggregate() -processedData-> write()
*/
func processDataFolder(folder string) {
	timer := time.Now()

	var filePaths = make(chan string)     // Carries paths of input files
	var rawData = make(chan string)       // Carries raw input lines
	var parsedData = make(chan *Visit)    // Carries parsed Visit structs
	var processedData = make(chan *Visit) // Carries aggregated Visit structs

	go list(folder, filePaths)

	go read(filePaths, rawData)

	go parse(rawData, parsedData)

	go aggregate(parsedData, processedData)

	write(processedData)

	fmt.Printf("Processing folder %s took %dms (%v)\n", folder, MillisecondsSince(timer), time.Since(timer))
}

/*
	========== READ STEP ===========
*/
func list(folder string, filePaths chan string) {
	defer close(filePaths)

	var fileList []os.FileInfo
	fileList, err := ioutil.ReadDir(folder)
	fail(err)

	for _, fileInfo := range fileList {
		if !fileInfo.IsDir() {
			filePaths <- folder + "/" + fileInfo.Name()
		}
	}
}

func read(filePaths chan string, out chan string) {
	defer close(out)
	counter, timer := 0, time.Now()

	for filePath := range filePaths {
		file, err := os.Open(filePath)
		fail(err)

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			out <- scanner.Text() // Send data on the output channel
		}
		fail(scanner.Err())

		file.Close()
		counter++
	}
	fmt.Printf("\tRead %d files in %dms\n", counter, MillisecondsSince(timer))
}

/*
	========== PARSE STEP ==========
	Expected input format (tsv):
		Date    UserID    IP    OS    Browser
*/
const INPUT_SEP = "\t"
const INPUT_FIELDS = 5

func parse(in chan string, out chan *Visit) {
	defer close(out)
	counter, timer := 0, time.Now()

	for line := range in {
		parts := strings.Split(line, INPUT_SEP)
		if len(parts) != INPUT_FIELDS {
			fail(fmt.Errorf("Wrong number of fields in line: %s", line))
		}
		userID, err := strconv.ParseUint(parts[1], 10, 32)
		fail(err)

		doBusyWork() // Real parsing would take more computing than this

		out <- NewVisit(uint32(userID), parts[2], parts[3], parts[4])
		counter++
	}
	fmt.Printf("\tParsed %d lines in %dms\n", counter, MillisecondsSince(timer))
}

/*
	======== AGGREGATE STEP ========
*/
func aggregate(in chan *Visit, out chan *Visit) {
	defer close(out)
	var hash = make(map[string]*Visit)

	for visit := range in {
		var key string = visit.GetKey()

		cachedVisit, exists := hash[key]
		if exists {
			cachedVisit.Count += visit.Count
		} else {
			hash[key] = visit
		}
	}

	counter, timer := 0, time.Now()
	for _, visit := range hash {
		out <- visit
		counter++
	}
	fmt.Printf("\tWrote %d lines in %dms\n", counter, MillisecondsSince(timer))
}

/*
	========== WRITE STEP ==========
*/
func write(in chan *Visit) {
	file, err := os.OpenFile(OUTPUT_FILE, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	fail(err)

	for visit := range in {
		fmt.Fprintf(file, "%s\n", visit)
	}
}

/*
	================================
	======== DATA STRUCTURE ========
	================================
*/
type Visit struct {
	UserID  uint32
	IP      string
	OS      string
	Browser string
	Count   uint32
}

func NewVisit(userID uint32, IP, OS, Browser string) *Visit {
	if userID > 0 {
		IP, OS, Browser = "", "", ""
	}
	return &Visit{
		UserID:  userID,
		IP:      IP,
		OS:      OS,
		Browser: Browser,
		Count:   1,
	}
}

func (v *Visit) String() string {
	return fmt.Sprintf("%d\t%d\t%s\t%s\t%s", v.Count, v.UserID, v.IP, v.OS, v.Browser)
}

func (v *Visit) GetKey() string {
	return fmt.Sprintf("%d\t%s\t%s\t%s", v.UserID, v.IP, v.OS, v.Browser)
}

/*
	================================
	====== UTILITY FUNCTIONS =======
	================================
*/
// Utility function for measuring time
func MillisecondsSince(t time.Time) int64 {
	return int64(time.Since(t) / time.Millisecond)
}

// Utility function for error checking
func fail(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func doBusyWork() {
	for i := 2; i < 50; i++ {
		isPrime := true
		for j := 2; j < i; j++ {
			if i%j == 0 {
				isPrime = false
			}
		}
		_ = isPrime
	}
}
