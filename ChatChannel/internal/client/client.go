package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {

	conn, _ := net.Dial("tcp", ":8181")

	go func() {
		for {
			message, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				fmt.Println("Error occured, try reconnect")
				os.Exit(1)
			}
			fmt.Println(message)
		}
	}()

	for {
		text, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		fmt.Fprintf(conn, text+"\n")
	}

}
