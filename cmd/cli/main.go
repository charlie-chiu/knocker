package main

import (
	"fmt"
	"log"

	"knocker"
)

func main() {
	//door := knocker.Door{
	//	//URL: "https://daniu.cool",
	//	URL:       "https://daniu.cool/admin/login/index.html",
	//	Host:      "52.229.224.88",
	//	IgnoreSSL: true,
	//}

	//door := knocker.Door{
	//	URL:       "https://cdn.dev3x.club:22443",
	//	//Host:      "",
	//	IgnoreSSL: false,
	//}

	//door := knocker.Door{
	//	URL: "https://image06.fenhao24.com/",
	//	//IgnoreSSL: true,
	//}

	//door := knocker.Door{URL: "https://news.baidu.com"}
	door := knocker.Door{URL: "https://leaidc.com"}

	results, err := knocker.Knock(door)
	if err != nil {
		fmt.Printf("knock error: %v", err)
	}

	knocker.PrintResults(results, false)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
