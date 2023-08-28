package gotube

import (
	"net/http"
	"fmt"
)

func (app *application) webhook(res http.ResponseWriter, req *http.Request) {
	fmt.Println("webhook")
	challenge := req.URL.Query().Get("hub.challenge")
	if challenge == "" {
		contentType := req.Header.Get("Content-Type")
		if contentType == "application/json" {
			renew := app.renewSub()
			app.infoLog.Println("Subscribtion renewed")
			app.infoLog.Println(renew)
			return
		}

		data, err := app.paseXML(req)
		if err != nil {
			app.errorLog.Println(err.Error())
			http.Error(res, "Internal Server Error", 500)	
		}

		app.infoLog.Printf("Titles: %s\n", data.Data.Title)
		app.infoLog.Printf("Author: %s\n", data.Data.Author.Name)
		app.infoLog.Printf("URL: %s\n", data.Data.Link.Href)

		vidToSend := app.checkDataHistory(data.Data.YtVideoID)
		if vidToSend == nil {
			return
		} else {
			app.infoLog.Printf("Pushing %s to Slack\n", data.Data.Link.Href)
			app.pushtoSlack(data.Data.Link.Href)
		}
	}
	// Renew subscription if there is a challenge
	res.Write([]byte(challenge))
}