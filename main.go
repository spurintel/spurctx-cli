package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

var ipRegex = regexp.MustCompile(`\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`)

func main() {
	ipFlag := flag.String("ip", "", "IP addresses separated by comma")
	fileFlag := flag.String("f", "", "File with IP addresses")
	garbageFlag := flag.Bool("g", false, "Garbage input")
	garbageFileFlag := flag.String("gf", "", "Garbage file input")
	parallelismFlag := flag.Int("n", runtime.NumCPU()*2, "Override the parallelism")
	flag.Parse()

	// If there are no input flags or stdin, print usage
	hasStdin := false
	stat, err := os.Stdin.Stat()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if (stat.Size()) > 0 {
		hasStdin = true
	}

	if *ipFlag == "" && *fileFlag == "" && !hasStdin && *garbageFileFlag == "" {
		flag.Usage()
		os.Exit(1)
	}

	token := os.Getenv("SPUR_TOKEN")
	if token == "" {
		fmt.Println("SPUR_TOKEN environment variable is not set")
		os.Exit(1)
	}

	ipChan := make(chan string)
	var wg sync.WaitGroup

	for i := 0; i < *parallelismFlag; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for ip := range ipChan {
				resp, err := queryAPI(ip, token)
				if err != nil {
					fmt.Println(err)
					continue
				}
				jsonResp, err := json.Marshal(resp)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println(string(jsonResp))
			}
		}()
	}

	processIPFlag(*ipFlag, ipChan)
	processFileFlag(*fileFlag, ipChan)
	processGarbageFileFlag(*garbageFileFlag, ipChan)
	processStdin(*garbageFlag, ipChan)

	close(ipChan)
	wg.Wait()
}

func queryAPI(ip, token string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", "https://api.spur.us/v2/context/"+ip, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("TOKEN", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func processIPFlag(ipFlag string, ipChan chan<- string) {
	if ipFlag != "" {
		ips := strings.Split(ipFlag, ",")
		for _, ip := range ips {
			ipChan <- ip
		}
	}
}

func processFileFlag(fileFlag string, ipChan chan<- string) {
	if fileFlag != "" {
		file, err := os.Open(fileFlag)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			ipChan <- scanner.Text()
		}
	}
}

func processGarbageFileFlag(garbageFileFlag string, ipChan chan<- string) {
	if garbageFileFlag != "" {
		content, err := os.ReadFile(garbageFileFlag)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		ips := ipRegex.FindAllString(string(content), -1)
		for _, ip := range ips {
			ipChan <- ip
		}
	}
}

func processStdin(garbageFlag bool, ipChan chan<- string) {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		if garbageFlag {
			for scanner.Scan() {
				ips := ipRegex.FindAllString(scanner.Text(), -1)
				for _, ip := range ips {
					ipChan <- ip
				}
			}
		} else {
			for scanner.Scan() {
				ipChan <- scanner.Text()
			}
		}
	}
}
