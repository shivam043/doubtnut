package handler

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"github.com/doubtnut/logging"
	"github.com/doubtnut/redis"
	"github.com/gorilla/mux"
	"github.com/jung-kurt/gofpdf"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

var Logger = logging.NewLogger()

// generate pdf and email
func Send(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		msg := "Content-Type header is not application/json"
		http.Error(w, msg, http.StatusUnsupportedMediaType)
		return
	}
	vars := mux.Vars(r)
	userId := vars["id"]
	//set timestamp value in redis
	if (!redis.SetValue("userId-"+userId, time.Now().String(), 0)) {
		Logger.Errorf("error setting timestamp for userId " + userId)
	}
	var pdfOutput map[string]interface{}
	//decode json
	err := json.NewDecoder(r.Body).Decode(&pdfOutput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.Image("static/db.jpeg", 15, 15, 267, 0, false, "", 0, "")
	pdf.SetFont("Arial", "B", 16)
	count := 0
	for _, value := range pdfOutput {
		word := value.(string)
		wo := ""
		var st []string
		for i := range word {
			if (i == 0) {
				wo = strconv.Itoa(count+1) + "."
			}
			wo = wo + string(word[i])
			if (len((wo)) > 60) {
				st = append(st, wo)
				wo = ""
			}
		}
		st = append(st, wo)
		for i := range st {
			pdf.Cell(float64(1), 10, st[i])
			pdf.Ln(6)

		}
		pdf.Ln(12)
		count++
	}

	err = pdf.OutputFileAndClose("static/pdfs/" + userId + ".pdf")
	if err != nil {
		Logger.Errorf("error generating pdf for userId " + userId)
	}
	err = SendEmail("frommail@gmail.com", "tomail@gmail.com", userId)
	if err != nil {
		Logger.Errorf("error sending mail " + err.Error())
	}
	w.Write([]byte("foo"))
	//fmt.Fprintf(w, "Category: %v\n", vars)

}

//uses sendgrid and send email
func SendEmail(fromMail string, toMail string, userId string) error {
	from := mail.NewEmail("DoubtNut Support", fromMail)
	subject := "Similar Questions Pdf"
	to := mail.NewEmail("Enter Name Here", toMail)
	plainTextContent := "and easy to do anywhere, even with Go"
	htmlContent := "<strong>DoubtNut</strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	pdf := mail.NewAttachment()
	fileHandle, err := os.Open("static/pdfs/" + userId + ".pdf")
	defer fileHandle.Close()
	if err != nil {
		return err
	}
	reader := bufio.NewReader(fileHandle)
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	// Encode as base64.
	encoded := base64.StdEncoding.EncodeToString(content)
	pdf.SetContent(encoded)
	pdf.SetType("application/pdf")
	pdf.SetFilename("questions.pdf")
	pdf.SetDisposition("attachment")

	//message.AddAttachment(fileHandle)
	message.AddAttachment(pdf)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	_, err = client.Send(message)
	if err != nil {
		return err
	}

	return nil

}
