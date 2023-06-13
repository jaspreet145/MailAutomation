package main

import (
	"bufio"
	"crypto/tls"
	"encoding/csv"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	gomail "gopkg.in/mail.v2"
)

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

func sendMail(to string, body string) {

	const email = "EMAIL"
	const password = "PASSWORD"
	const header = "HEADER"

	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", os.Getenv(email))

	// Set E-Mail receivers
	m.SetHeader("To", to)

	// Set E-Mail subject
	m.SetHeader("Subject", os.Getenv(header))

	// Attaach resume or file
	m.Attach("Resume.pdf")
	m.Attach("CoverLetter.docx")

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", body)

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv(email), os.Getenv(password))

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		log.Fatal(err)
	}
}

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal(err)
	}

	records := readCsvFile("./contacts.csv")

	for _, record := range records[1:] {

		file, err := os.Open("./body.txt")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		emailBody := ""
		for scanner.Scan() {
			body := scanner.Text()

			body = strings.Replace(body, "{@name}", record[1], 1)
			body = strings.Replace(body, "{@company}", record[0], 1)
			emailBody = emailBody + "\n" + body
		}
		sendMail(record[2], emailBody)
	}

	return
}
