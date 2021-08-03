package main

import (
	"fmt"
	"log"

	"knocker"
)

func main() {
	//default

	//rawURL := "http://daniu.cool/admin/login/index.html"
	rawURL := "http://daniu.cool"
	//rawURL := "https://image06.fenhao24.com/"
	//modifiedIP := ""
	//modifiedIP := "45.60.64.140"
	modifiedIP := "52.229.224.88"
	//modifiedIP := "1.1.1.1"
	//port := 443

	door := knocker.Door{
		URL:       rawURL,
		IPAddress: modifiedIP,
		Port:      0,
		WithTrace: false,
		IgnoreSSL: true,
	}

	results, err := knocker.Knock(door)
	failOnError(err, "failed to knock")

	fmt.Printf("%+v", results)

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
