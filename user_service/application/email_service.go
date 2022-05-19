package application

import (
	"bytes"
	"github.com/dislinkt/user_service/startup/config"
	"html/template"
	"net/smtp"
)

type EmailService struct {
	emailSender   string
	emailPassword string
	emailHost     string
}

func NewEmailService(c *config.Config) *EmailService {
	return &EmailService{
		emailSender:   c.EmailSender,
		emailPassword: c.EmailPassword,
		emailHost:     c.EmailHost,
	}
}

type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

func NewRequest(to []string, subject, body string) *Request {
	return &Request{
		to:      to,
		subject: subject,
		body:    body,
	}
}

var auth smtp.Auth

func (service *EmailService) SendActivationMail(receiver string, name string, activationId string) {
	auth = smtp.PlainAuth("Dislinkt", service.emailSender, service.emailPassword, service.emailHost)
	templateData := struct {
		Name string
		URL  string
	}{
		Name: name,
		URL:  "http//:localhost:8000/api/users/activate/" + activationId,
	}
	r := NewRequest([]string{receiver}, "Hello "+name+"!", "Hello "+name+"!")
	if err := r.parseTemplate("mailActivation.html", templateData); err == nil {
		if ok, _ := r.sendEmail(); ok {
			//logger.Logger.WithFields(logrus.Fields{"activation_id": activationId}).Info("Activation e-mail sent")
		} else {
			//logger.Logger.WithFields(logrus.Fields{"activation_id": activationId}).Error("Sending activation e-mail")
		}
	}
}

func (service *EmailService) SendResetPasswordMail(receiver string, name string, activationId string) {
	auth = smtp.PlainAuth("Dislinkt", service.emailSender, service.emailPassword, service.emailHost)
	templateData := struct {
		Name string
		URL  string
	}{
		Name: name,
		URL:  "http//:localhost:8000/api/users/reset-password/" + activationId,
	}
	r := NewRequest([]string{receiver}, "Hello "+name+"!", "Hello "+name+"!")
	if err := r.parseTemplate("mailResetPassword.html", templateData); err == nil {
		if ok, _ := r.sendEmail(); ok {
			//logger.Logger.WithFields(logrus.Fields{"reset_password_id": activationId}).Info("Reset password e-mail sent")
		} else {
			//logger.Logger.WithFields(logrus.Fields{"reset_password_id": activationId}).Error("Sending reset password e-mail")
		}
	}
}

func (r *Request) sendEmail() (bool, error) {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + r.subject + "!\n"
	msg := []byte(subject + mime + "\n" + r.body)
	conf := config.NewConfig()
	addr := conf.EmailHost + ":" + conf.EmailPort

	if err := smtp.SendMail(addr, auth, conf.EmailSender, r.to, msg); err != nil {
		return false, err
	}
	return true, nil
}

func (r *Request) parseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}
