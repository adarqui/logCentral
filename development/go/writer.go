package main

import (
	"fmt"
	"os"
//	"io/ioutil"
	"bufio"
	"log"
)

func main() {
	fmt.Println("writer.go")
	file, err := os.OpenFile("/tmp/out", os.O_RDWR | os.O_APPEND | os.O_CREATE, 0700)
	if err != nil {
		log.Fatal(err)
	}
	for {
/*
		bytes, err := ioutil.ReadAll(os.Stdin)
		log.Println(err, string(bytes))
		n_written, err := file.Write(bytes)
		if err != nil {
			log.Println("err",err)
		} else {
			log.Println("n_written",n_written)
		}
*/
		reader := bufio.NewReader(os.Stdin)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			log.Println("line", line)
			file.WriteString(line)
		}
	}
}
