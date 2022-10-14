package main

import (
	"log"
	"memo/services"
)

func main() {
	err := services.Start("0.0.0.0:80")
	if err != nil {
		log.Fatalln(err.Error())
	}
}
