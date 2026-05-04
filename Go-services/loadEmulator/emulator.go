package main

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
	"log"
)

// client to send requests
func main() {

	if len(os.Args) > 1 {
		count, _ := strconv.Atoi(os.Args[1])
		for i := 0; i < count; i++ {
			resp, err := http.Get("http://10.0.0.100:8080/hello")

			time.Sleep(time.Millisecond * 100)

			if err != nil {
				log.Println(err)
			}
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Println(err)
			}
			log.Println(string(body))
		}
	}
}
