package mail

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSendMailWithGmail(t *testing.T) {
	sender := NewGmailSender("e-cost.by", "hmurauyou.test@gmail.com", "nhpxotugfxlzbsrg")

	subject := "A test mail"
	content := `
		<h1>Thank you for confirmation your email</h1>
	`
	to := []string{"georgiemurauyou@gmail.com"}
	attachFile := []string{"./txt.txt"}

	err := sender.SendMail(subject, content, to, nil, nil, attachFile)
	require.NoError(t, err)
}
