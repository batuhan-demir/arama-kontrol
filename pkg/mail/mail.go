package mail

import (
	"bytes"
	"fmt"
	"net/http"
)

func SendMail(to string, subject string, body string) {

	requestJson := fmt.Sprintf(`{
		"to": "%s",
		"subject": "%s",
		"body": "%s"
		}`, to, subject, body)

	requestBody := []byte(requestJson)

	http.Post(
		"https://mail-service-weld.vercel.app/",
		"application/json",
		bytes.NewReader(requestBody))
}
