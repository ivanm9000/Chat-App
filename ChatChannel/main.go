package main

import (
	"bufio"
	//"fmt"
	"strconv"
	//"log"
	"net"
	//	"time"
)

// Client --- username and message
type Client struct {
	Socket   net.Conn
	Username string
	Message  string
}

func main() {

	newConnection := make(chan *Client, 128)
	removeConnection := make(chan *Client, 128)
	broadcastMessage := make(chan *Client, 256)
	connectionsList := make(map[*Client]struct{})

	listen, err := net.Listen("tcp", ":8181")
	if err != nil {
		panic(err)
	}
	defer listen.Close()
	go func() {

		// listen non-stop new connection on port
		id := 0

		for {

			conn, err := listen.Accept()
			if err != nil {
				panic(err)
			}
			id++
			client := &Client{Socket: conn, Username: "User_" + strconv.Itoa(id)}

			newConnection <- client

		}
	}()
	//go func() {

	// wait non-stop for read, write or delete

	for {
		select {
		case conn := <-newConnection:

			// read messages from all clients using go routine, when new connection appears read go routine is called

			connectionsList[conn] = struct{}{}

			go func(connect *Client) {
				for {
					message, err := bufio.NewReader(connect.Socket).ReadString('\n')
					if err != nil {
						removeConnection <- connect
						break

					} else {

						connect.Message = message

						broadcastMessage <- connect

					}

				}

			}(conn)

		case conn := <-removeConnection:

			// delete client from clients list

			conn.Socket.Close()
			delete(connectionsList, conn)

		case mess := <-broadcastMessage:

			// write message to all clients using go routine

			if string(mess.Message[0]) == "-" {

				oldName := mess.Username
				mess.Username = string(mess.Message[1 : len(mess.Message)-1]) //last char is \n so client cannot read from chan cause it signals end

				//write to all clients

				for conn := range connectionsList {
					go func(connect *Client) {

						_, err := connect.Socket.Write([]byte(oldName + " CHANGED IT'S NAME INTO " + mess.Username + "\n"))
						if err != nil {
							removeConnection <- connect
						}

					}(conn)
				}

			} else {

				//write to all clients

				for conn := range connectionsList {
					go func(connect *Client) {

						_, err := connect.Socket.Write([]byte(mess.Username + ": " + mess.Message + "\n"))
						if err != nil {
							removeConnection <- connect
						}

					}(conn)
				}
			}
		}

	}
	//	}()
	//fmt.Println("123")
	//time.Sleep(3 * time.Minute)
}
