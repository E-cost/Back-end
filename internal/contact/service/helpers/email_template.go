package helpers

import (
	"Ecost/internal/contact/dto"
	"bytes"
	"fmt"
	"html/template"
)

type TemplateData struct {
	Name      string
	VerifyURL string
}

func EmailTemplate(dto dto.CreateContactDto, verifyURL string) (subject string, msg []byte, to []string, err error) {
	t, err := template.ParseFiles("/usr/local/src/internal/contact/mail/template/index.html")
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to parse HTML template: %v", err)
	}

	var buff bytes.Buffer

	data := TemplateData{
		Name:      dto.Name,
		VerifyURL: verifyURL,
	}
	subject = "E-cost order"

	if err := t.Execute(&buff, data); err != nil {
		return "", nil, nil, fmt.Errorf("failed to execute HTML template: %v", err)
	}

	msg = []byte(buff.Bytes())

	to = []string{dto.Email}

	return subject, msg, to, nil
}
