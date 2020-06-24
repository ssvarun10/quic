package main

import (
	//"bufio"
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
	//"strings"
	quic "github.com/lucas-clemente/quic-go"
)

var bufferone string

func tcpclient() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide host:port.")
		return
	}

	CONNECT := arguments[1]
	c, err := net.Dial("tcp", CONNECT)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("asking for unique id from server ")
	buffertwo := make([]byte, 1024)

	c.Read(buffertwo)
	bufferone = string(buffertwo)
	fmt.Println("the id from tcpserver is ", bufferone)
	time.Sleep(5 * time.Second)

}
func udpclient() {
	//arguments := os.Args
	//if len(arguments) == 1 {
	//fmt.Println("Please provide host:port.")
	//return

	CONNECT := "127.0.0.1:1235"

	s, err := net.ResolveUDPAddr("udp4", CONNECT)
	c, err := net.DialUDP("udp4", nil, s)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("The UDP server is %s\n", c.RemoteAddr().String())
	defer c.Close()
	fmt.Println("asking for unique id from udpserver ")
	data := []byte("please send me the connection id")
	_, err1 := c.Write(data)
	if err1 != nil {
		fmt.Println(err)
		return
	}

	if bufferone == "" {
		buffertwo := make([]byte, 1024)

		c.Read(buffertwo)
		bufferone = string(buffertwo)

		fmt.Println("the id from server udpserveris ", bufferone)
	} else {

		fmt.Println("the connection id is already generated ", bufferone)
	}
	time.Sleep(5 * time.Second)

}
func qclient() {
	timeoutDuration := flag.Int("rtt", 50, "timeout duration (in ms)")

	flag.Parse()
	addr := "127.0.0.1:1236"
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-sendreceive"},
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
	message := "send the id quicserver"
	fmt.Printf("Client: Sending '%s'\n", message)
	_, err = stream.Write([]byte(message))
	if err != nil {
		panic(err)
	}
	resp := make(chan string)
	go func() {
		buff := make([]byte, 1024)
		_, err = io.ReadAtLeast(stream, buff, 3)
		if err != nil {
			panic(err)
		}

		resp <- string(buff)
	}()

	select {
	case reply := <-resp:
		log.Printf("Client: Got '%s'\n", reply)
	case <-time.After(timeout):
		log.Printf("Client: Timed out\n")
	}

}
func main() {

	tcpclient()
	fmt.Println("bufferone value after exiting tcpclient is ", bufferone)
	udpclient()
	fmt.Println("bufferone value after exiting udpclient is ", bufferone)
	qclient()
}
