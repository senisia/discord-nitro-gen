package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

var (
	url           = "https://api.discord.gx.games/v1/direct-fulfillment"
	partnerUserID = "50b1bf177eca2a06f77680c1aa6277e1d5a44eb6d8b4a72545348e4828cf0753"
	headers       = map[string]string{
		"authority":          "api.discord.gx.games",
		"accept":             "*/*",
		"accept-language":    "en-US,en;q=0.9",
		"content-type":       "application/json",
		"origin":             "https://www.opera.com",
		"referer":            "https://www.opera.com/",
		"sec-ch-ua":          `"Opera GX";v="105", "Chromium";v="119", "Not?A_Brand";v="24"`,
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": `"Windows"`,
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "cross-site",
		"user-agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36 OPR/105.0.0.0",
	}
	wg sync.WaitGroup
)

func main() {
	var threads int

	fmt.Print("How much threads do you want to use? :  ")
	fmt.Scan(&threads)

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go nitroRequest(i)
	}
	wg.Wait()
}

func nitroRequest(threadID int) {
	defer wg.Done()
	fmt.Printf("Thread %d has started\n", threadID)
	for {
		response, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(fmt.Sprintf(`{"partnerUserId": "%s"}`, partnerUserID))))
		if err != nil {
			fmt.Printf("Request failed with error: %v\n", err)
			continue
		}
		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("Failed to read response body: %v\n", err)
			continue
		}
		if response.StatusCode == http.StatusOK {
			var result map[string]interface{}
			if err := json.Unmarshal(body, &result); err != nil {
				fmt.Printf("Failed to unmarshal JSON response: %v\n", err)
				continue
			}
			token := result["token"].(string)
			file, err := os.OpenFile("codes.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Printf("Failed to open file: %v\n", err)
				continue
			}
			defer file.Close()
			if _, err := file.WriteString(fmt.Sprintf("https://discord.com/billing/partner-promotions/1180231712274387115/%s\n", token)); err != nil {
				fmt.Printf("Failed to write token to file: %v\n", err)
				continue
			}
			fmt.Println("Token saved to codes.txt file.")
		} else {
			fmt.Printf("Request failed with status code %d\n", response.StatusCode)
			fmt.Printf("Error message: %s\n", body)
		}
	}
}
