package main

import "github.com/mostafah/mandrill"

// TODO: msg.FromEmail must be kkowalczyk@gmail.com. Look into having it
// come from @databaseworkbench.com
// kkowalczyk@gmail is fine for e-mails that go to us but not fine for things
// like e-mail confirmation e-mails, marketing e-mails etc.
const (
	mandrillAPIKey = "_cL1X8-B1GTqTu2IPuxsMg"
)

func sendTestEmail() {
	mandrill.Key = mandrillAPIKey
	err := mandrill.Ping()
	if err != nil {
		LogErrorf("mandrill.Ping() failed with '%s'\n", err)
		return
	}
	msg := mandrill.NewMessageTo("kkowalczyk@gmail.com", "")
	msg.Text = "This is a test message"
	msg.Subject = "Database Workbench: test e-mail"
	msg.FromEmail = "kkowalczyk@gmail.com"
	msg.FromName = "Database Workbench"
	_, err = msg.Send(false)
	if err != nil {
		LogErrorf("msg.Send() failed with %s\n", err)
	}
}
