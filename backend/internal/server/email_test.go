package server

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T){

	sender := NewGmailSender("Vas deal tracker",
	 "eebenezer949@gmail.com", 
	"hfqiiyhadcguwqmh")

	subject := "A test email"
	content := `
	<h1>Hello world<h1>
	<p> This is a test message
	`
	to := []string{"ekefan4@gmail.com"}
	err := sender.sendEmail(subject, content, to)
	require.NoError(t, err)
}