package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/mailgun/mailgun-go.v1"
)

func main() {
	// Create an instance of the Mailgun Client
	mg := mailgun.NewMailgun(os.Getenv("MAILGUN_DOMAIN"), os.Getenv("MAILGUN_PRIVATE"), os.Getenv("MAILGUN_PUBLIC"))

	body := "Hello from Mailgun Go!"

	sendMessage(mg, body)
}

func sendMessage(mg mailgun.Mailgun, body string) {
	message := mg.NewMessage("danielsussa@gmail.com", "Stock", body, "danielsussa@gmail.com")
	resp, id, err := mg.Send(message)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID: %s Resp: %s\n", id, resp)
}
