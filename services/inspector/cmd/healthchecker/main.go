package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)

func main() {
	port := flag.Int("port", 80, "port of localhost to ping")
	flag.Parse()

	resp, err := http.Get(fmt.Sprintf("http://localhost:%v/v1/ping", *port))
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println(err)
		os.Exit(1)
	}
}
