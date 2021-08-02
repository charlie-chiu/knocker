package main

import (
	"log"

	"knocker"
)

func main() {
	//default

	//rawURL := "http://daniu.cool/admin/login/index.html"
	rawURL := "http://daniu.cool"
	//rawURL := "https://image06.fenhao24.com:443/"
	//modifiedIP := "45.60.64.140"
	modifiedIP := "52.229.224.88"
	//port := 443

	statusCode, err := knocker.Knock2(rawURL, modifiedIP, 0, true)
	failOnError(err, "failed to knock")

	log.Printf("got status %d\n", statusCode)

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
