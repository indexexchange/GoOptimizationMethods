package main

import (
	"fmt"
	"flag"
	"math/rand"
	"math"
	"time"
)

var (
	NUM_LINES int = 1000000
	
	DATE string = "2018-06-05 06:00:00"
	
	USER_FRACTION int = 66
	USER_RANGE int = 1000
	
	IP_PREFIX string = "192.168.0."
	IP_RANGE int = 256
	
	OS_LIST = []string{"Windows", "Mac OS X", "Android", "iOS", "Ubuntu"}
	
	BROWSER_LIST = []string{"Chrome", "Safari", "Firefox", "Opera", "IE"}
	
	RANDOM_TIME bool = false
	RANDOM_SEED int64 = 4855279955359852901
)

var gen *rand.Rand

func main() {
	flag.IntVar(&NUM_LINES, "n", NUM_LINES, "The number of lines to generate")
	flag.BoolVar(&RANDOM_TIME, "r", RANDOM_TIME, "If set, uses a unique random seed; otherwise, always uses the same random seed")
	flag.Parse()
	
	if RANDOM_TIME {
		RANDOM_SEED = time.Now().Unix()
	}
	gen = rand.New(rand.NewSource(RANDOM_SEED))
	
	for i := 0; i < NUM_LINES; i++ {
		fmt.Println(randVisit())
	}
}

/* Generator Functions */

func randVisit() string {
	return fmt.Sprintf("%s\t%d\t%s\t%s\t%s", DATE, randUserID(), randIP(), randOS(), randBrowser())
}

func randUserID() int {
	percent := gen.Intn(100)
	if percent < USER_FRACTION {
		return randIntFrontWeighted(1, USER_RANGE)
	} else {
		return 0
	}
}

func randIP() string {
	return fmt.Sprintf("%s%d", IP_PREFIX, randIntFrontWeighted(0, IP_RANGE))
}

func randOS() string {
	return OS_LIST[randIntFrontWeighted(0, len(OS_LIST))]
}

func randBrowser() string {
	return BROWSER_LIST[randIntFrontWeighted(0, len(BROWSER_LIST))]
}

/* Support Functions */

func randIntFrontWeighted(start, end int) int {
	length := end - start
	return start + int((math.Pow(101, gen.Float64()) - 1) / 100 * float64(length))
}

func randIntMiddleWeighted(start, end int) int {
	length := end - start
	return start + (gen.Intn(length) + gen.Intn(length)) / 2
}

