package utils

import (
"github.com/zerobfd/mailbuilder"
"github.com/usercenter/usercenter/settings"
)

func mail(from, to mailbuilder.Address, subject, content string) {
  message := mailbuilder.NewMessage()
  message.AddTo(to)
  message.From = from
  message.Subject = subject
  body := mailbuilder.NewSimplePart()
  message.SetBody(body)
  body.AddHeader("Content-Type", "text/plain; charset=utf8")
  body.AddHeader("Content-Transfer-Encoding", "quoted-printable")
  body.Content = content
  auth := smtp.PlainAuth("", settings.SMTP_USER, settings.SMTP_PASS, settings.SMTP_SERVER)
  err := smtp.SendMail(settings.SMTP_SERVER+":"SMTP_PORT,
                auth,
                message.From.Email,
                message.Recipients(),
                message.Bytes())
  if (err != nil) {fmt.Printf("%v", err)}
}