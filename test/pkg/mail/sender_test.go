package mail

import (
	"github.com/NeptuneYeh/simplerecommend/init/config"
	"github.com/NeptuneYeh/simplerecommend/pkg/mail"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSendEmailWithGmail(t *testing.T) {
	configModule := config.NewModule("../../..")

	sender := mail.NewGmailSender(configModule.EmailSenderName, configModule.EmailSenderAddress, configModule.EmailSenderPassword)
	subject := "A test Verify email"
	content := `<h1>your code: 12345678</h1>`
	to := []string{"alex554833@gmail.com"}

	err := sender.SendEmail(subject, content, to, nil, nil, nil)
	require.NoError(t, err)
}
