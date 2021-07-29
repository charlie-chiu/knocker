package main

import (
	"log"

	"knocker"
)

func main() {
	//default
	//knocker https://daniu.cool@45.60.64.140
	//knocker https://daniu.cool@8.212.8.138:7788

	rawURL := "http://daniu.cool/admin/login/index.html"
	//modifiedIP := "45.60.64.140"
	//modifiedIP := "52.229.224.88"
	//port := 443

	//statusCode, err := knocker.Knock2(rawURL, modifiedIP, 0, true)
	statusCode, err := knocker.Knock2(rawURL, "", 0, true)
	failOnError(err, "failed to knock")

	log.Printf("got status %d\n", statusCode)

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
