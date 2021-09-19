package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"time"

	"github.com/Pauloo27/gmail-notifier/gapi"
	"github.com/Pauloo27/gmail-notifier/utils"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

const clientCount = 3

func fetchMessages(srv *gmail.Service) []*gmail.Message {
	r, err := srv.Users.Messages.List("me").LabelIds("UNREAD").IncludeSpamTrash(true).MaxResults(10).Do()
	if err != nil {
		utils.HandleFatal("Unable to retrieve message", err)
	}

	return r.Messages
}

func runDaemon(askLogin bool) {
	usr, err := user.Current()
	if err != nil {
		utils.HandleFatal("Cannot get user", err)
	}

	secretFolder := usr.HomeDir + "/.cache/gmail-notifier/secret/gmail"
	credentialsFile := secretFolder + "credentials.json"

	gmailServices := []*gmail.Service{}

	for i := 0; i < clientCount; i++ {
		fmt.Println("Loading client", i)
		tokFile := fmt.Sprintf("%stoken-%d.json", secretFolder, i)

		b, err := ioutil.ReadFile(credentialsFile)
		if err != nil {
			utils.HandleFatal("Unable to read client secret file", err)
		}

		config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
		if err != nil {
			utils.HandleFatal("Unable to parse client secret file to config", err)
		}

		client := gapi.GetClient(config, tokFile, askLogin)
		fmt.Printf("Client %d loaded\n", i)

		srv, err := gmail.New(client)
		if err != nil {
			utils.HandleFatal("Unable to retrieve Gmail client", err)
		}

		gmailServices = append(gmailServices, srv)
	}

	fmt.Printf("Started. Logging status to %s\n", utils.StatusFile)

	for {
		status := []int{}
		fmt.Println("Fetching...")
		for _, srv := range gmailServices {
			messages := fetchMessages(srv)
			messageCount := len(messages)
			status = append(status, messageCount)
		}
		utils.LogStatus(status)
		fmt.Printf("Found %d\n", status)
		time.Sleep(3 * time.Minute)
	}
}

type PolybarActionButton struct {
	Index            uint
	Display, Command string
}

func (a PolybarActionButton) String() string {
	return fmt.Sprintf("%%{A%d:%s:}%s%%{A}", a.Index, a.Command, a.Display)
}

func main() {
	if len(os.Args) == 2 {
		switch os.Args[1] {
		case "start":
			runDaemon(false)
		case "stop":
			out, err := exec.Command("pgrep", "-f", "gmail-notifier start").Output()
			if err != nil {
				utils.HandleFatal("Cannot find daemon process", err)
			}
			pid := strings.TrimSuffix(string(out), "\n")
			err = exec.Command("kill", pid).Run()
			if err != nil {
				utils.HandleFatal("Cannot kill daemon process", err)
			}
			return
		case "login":
			runDaemon(true)
			return
		case "polybar":
			buttons := []string{}
			buffer, err := ioutil.ReadFile(utils.StatusFile)
			color := "#50fa7b"

			if err == nil {
				statusLine := strings.Split(string(buffer), "\n")[0]
				if statusLine == "-" {
					color = "#ff5555"
					btn := PolybarActionButton{
						1,
						"Error",
						"systemctl restart gmail-notifier",
					}
					buttons = append(buttons, btn.String())
				} else {
					statusList := strings.Split(statusLine, " ")
					for i, unread := range statusList {
						if unread != "0" {
							color = "#ffb86c"
						}
						btn := PolybarActionButton{
							1,
							fmt.Sprintf(": %s", unread),
							fmt.Sprintf("brave https\\://mail.google.com/mail/u/%d &", i),
						}
						buttons = append(buttons, btn.String())
					}
				}
			}
			fmt.Printf("%%{u%s}%s%%{u-}", color, strings.Join(buttons, " "))
			return
		}
	}
	fmt.Println("Invalid operation. Operations: start, stop, login, polybar")
}
