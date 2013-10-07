package main

import (
	"os"
	"bufio"
	"log"
)

func logger(c chan string) {
	file, err := os.OpenFile("/tmp/out", os.O_RDWR | os.O_APPEND | os.O_CREATE, 0700)
	if err != nil {
		log.Fatal(err)
	}
	for {
		msg := <- c
		file.WriteString(msg)
	}
}

func main() {

	var c chan string = make(chan string)

	go logger(c)

	for {
		reader := bufio.NewReader(os.Stdin)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			c <- line
		}
	}
}
