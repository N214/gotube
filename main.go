package main

import (
	"log"
	 gotube "github.com/N214/gotube/cmd"
)


func main() {
	if err := gotube.Run(); err != nil {
		log.Fatal(err)
	}
}