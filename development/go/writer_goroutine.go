package main

import (
	"bufio"
	"log"
	"os"
)

const Logfile = "/tmp/out"

func logger(name string, c chan string) {
	f, err := os.OpenFile(name, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	for {
		msg := <-c
		_, err = f.WriteString(msg)
		if err != nil {
			break
		}
	}
}

func main() {
	ch := make(chan string, 16)
	go logger(Logfile, ch)

	r := bufio.NewReader(os.Stdin)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		ch <- line
	}
}
