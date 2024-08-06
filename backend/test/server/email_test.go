package test

import (
	"testing"

	"github.com/ekefan/deal-tracker/internal/server"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	config, err := server.LoadConfig("./../../.env")
	require.NoError(t, err)

	sender := server.NewGmailSender(config.EmailSenderName,
	 config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "A test email"
	content := `
	<h1>Hello world<h1>
	<p> This is a test message
	`
	to := []string{"ekefan4@gmail.com"}
	err = sender.SendEmail(subject, content, to)
	require.NoError(t, err)
}
