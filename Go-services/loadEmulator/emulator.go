package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

// client to send requests
func main() {

	if len(os.Args) > 1 {
		count, _ := strconv.Atoi(os.Args[1])
		for i := 0; i < count; i++ {
			resp, err := http.Get("http://10.0.0.100:8080/hello")

			time.Sleep(time.Millisecond * 100)

			if err != nil {
				fmt.Println(err)
			}
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(body))
		}
	}
}
