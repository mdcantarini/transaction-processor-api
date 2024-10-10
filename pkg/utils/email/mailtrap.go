package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/mdcantarini/transaction-processor-api/pkg/utils"
)

type Mailtrap struct {
	FromEmail string
	Host      string
	Token     string
}

type emailData struct {
	From        map[string]string   `json:"from"`
	To          []map[string]string `json:"to"`
	Subject     string              `json:"subject"`
	Text        string              `json:"text"`
	HTML        string              `json:"html"`
	Attachments []map[string]string `json:"attachments"`
}

// SendEmail sends an email using the Mailtrap API: https://api-docs.mailtrap.io/docs/mailtrap-api-docs
func (mt Mailtrap) SendEmail(to, subject, body string, attachments ...string) error {
	emailData, err := buildEmailData(mt.FromEmail, to, subject, body, attachments...)
	if err != nil {
		return err
	}

	jsonPayload, err := json.Marshal(emailData)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", mt.Host, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", mt.Token)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// check if message was sent
	if !strings.Contains(string(resBody), `"success":true`) {
		return fmt.Errorf("unable to send email to %s; error: %s", to, string(resBody))
	}

	return nil
}

func buildEmailData(fromEmail, toEmail, subject, body string, attachments ...string) (*emailData, error) {
	data := &emailData{}

	data.From = map[string]string{"email": fromEmail}
	data.To = []map[string]string{
		{"email": toEmail},
	}
	data.Subject = subject
	data.HTML = body

	for _, attattachment := range attachments {
		attachmentContentEncoded, err := utils.EncodeFileContent(attattachment)
		if err != nil {
			return nil, err
		}

		data.Attachments = []map[string]string{
			{
				"content":  attachmentContentEncoded,
				"filename": attattachment,
			},
		}
	}

	return data, nil
}
