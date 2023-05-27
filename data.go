package main

import (
	"log"
)

// Define a struct to represent the XML data
type Author struct {
	Name string `xml:"name"`
	URI  string `xml:"uri"`
}

type Data struct {
	VideoId   string `xml:"id"`
	YtVideoID string `xml:"http://www.youtube.com/xml/schemas/2015 videoId"`
	YtChnID   string `xml:"http://www.youtube.com/xml/schemas/2015 channelId"`
	Title     string `xml:"title"`
	Link      struct {
		Href string `xml:"href,attr"`
	} `xml:"link"`
	Author    Author `xml:"author"`
	Published string `xml:"published"`
	Updated   string `xml:"updated"`
}

type Subscribe struct {
	Sub int `json:"sub"`
}

type Feed struct {
	Data Data `xml:"entry"`
}

type application struct {
	errorLog *log.Logger
	infoLog *log.Logger 
}