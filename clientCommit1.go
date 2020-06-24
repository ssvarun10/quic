package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	quic "github.com/lucas-clemente/quic-go"
)

func qclient(addr string) {
	//hostName := flag.String("hostname", "localhost", "hostname/ip of the server")
	//portNum := flag.String("port", "4242", "port number of the server")
	//numEcho := flag.Int("necho", 100, "number of echos")
	timeoutDuration := flag.Int("rtt", 10000, "timeout duration (in ms)")

	flag.Parse()

	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo"},
	}

	session, err := quic.DialAddr(addr, tlsConf, nil)
	if err != nil {
		panic(err)
	}

	stream, err := session.OpenStreamSync(context.Background())
	if err != nil {
		panic(err)
	}

	timeout := time.Duration(*timeoutDuration) * time.Millisecond

	resp := make(chan string)

	message := "hello mydear server"

	log.Printf("Client: Sending '%s'\n", message)
	_, err = stream.Write([]byte(message))
	if err != nil {
		panic(err)
	}

	go func() {
		buff := make([]byte, len(message))
		_, err = io.ReadFull(stream, buff)
		if err != nil {
			panic(err)
		}

		resp <- string(buff)
	}()

	select {
	case reply := <-resp:
		log.Printf("client: Got '%s'\n", reply)
	case <-time.After(timeout):
		log.Printf("cient: Timed out\n")
	}

}
func tcpclient(addr string) {
	//timeoutDuration := flag.Int("rtt", 3000, "timeout duration (in ms)")
	//flag.Parse()
	fmt.Println("address toconnect is ", addr)
	c, err := net.Dial("tcp", "127.0.0.1:4242")
	if err != nil {
		fmt.Println(err)
		return
	}
	//timeout := time.Duration(*timeoutDuration) * time.Millisecond

	//	resp := make(chan string)
	for {
		// read in input from stdin
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Text to send: ")
		text, _ := reader.ReadString('\n')
		// send to socket
		fmt.Fprintf(c, text+"\n")
		// listen for reply
		message, _ := bufio.NewReader(c).ReadString('\n')
		fmt.Print("Message from server: " + message)
	}
	//	go func() {
	//		buff := make([]byte, len(message))
	//		_, err := io.ReadFull(c, buff)
	//		if err != nil {
	//			panic(err)
	//		}
	//		resp <- string(buff)
	//	}()
	//	select {
	//	case reply := <-resp:
	//		log.Printf("client got reply '\n", reply)
	//	case <-time.After(timeout):
	//		log.Printf("client timed out \n")

	//	}
}

func main() {

	hostName := flag.String("hostname", "localhost", "hostname/ip of the server")
	portNum := flag.String("port", "4242", "port number of the server")

	flag.Parse()

	addr := *hostName + ":" + *portNum
	tcpclient(addr)
}
