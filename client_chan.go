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
	nanosBegin := timehandle()
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
	bufferthree := make([]byte, 1024)
	nanosEnd := timehandle()
	timediff := nanosEnd - nanosBegin
	fmt.Println("delay in making a connection and receiving uniq id is ", timediff)
	nanosBegin = timehandle()
	sendfiletoserver(c)

	c.Read(bufferthree)

	fmt.Println("the final say from tcp server is  ", string(bufferthree))

	nanosEnd = timehandle()
	timediff = nanosEnd - nanosBegin
	fmt.Println("delay in sending a file and reading from the receiver side in nanoseconds is ", timediff)
	time.Sleep(5 * time.Second)

}
func udpclient() {
	nanosBegin := timehandle()
	CONNECT := "127.0.0.1:1235"

	s, err := net.ResolveUDPAddr("udp4", CONNECT)
	c, err := net.DialUDP("udp4", nil, s)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("The UDP server is %s\n", c.RemoteAddr().String())
	defer c.Close()
	data := []byte("requesting for id")
	_, err = c.Write(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	var addr1 *net.UDPAddr
	if bufferone == "" {
		buffertwo := make([]byte, 1024)
		for {
			fmt.Println("waitinf for data")
			n, addr, err := c.ReadFromUDP(buffertwo)
			addr1 = addr

			if string(buffertwo[0:n]) != "" {
				fmt.Println("id generated is ", string(buffertwo[0:n]))

				break

			}
			if err != nil {
				fmt.Println(err)
				return
			}
		}

	} else {

		fmt.Println("the connection id is already generated ", bufferone)
	}
	nanosEnd := timehandle()

	timediff := nanosEnd - nanosBegin
	fmt.Println("delay in making a connection in udp  in nanoseconds ", timediff)
	fmt.Println("value of address is ", addr1)
	nanosBegin = timehandle()
	sendfiletoserver(c)

	fmt.Println("watiing for last reply from server")
	buffertwo := make([]byte, 1024)
	for {

		n, _, err := c.ReadFromUDP(buffertwo)

		if string(buffertwo[0:n]) != "" {

			break

		}
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	nanosEnd = timehandle()
	timediff = nanosEnd - nanosBegin
	fmt.Println("delay in sending a file and reading from the receiver side through udp  in nanoseconds is ", timediff)
	time.Sleep(5 * time.Second)
}
func qclient() {
	nanosBegin := timehandle()
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
	fmt.Println("the id from tcpserver is ", bufferone)
	nanosEnd := timehandle()
	timediff := nanosEnd - nanosBegin
	fmt.Println("delay in making a connection and receving uniq id in quic is  ", timediff)
	nanosBegin = timehandle()
	file, _ := os.Open("/home/varun/phone_backup_evolx")

	defer file.Close()
	buf := make([]byte, 1024)

	for {
		n, err := file.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}

			break

		}
		stream.Write(buf[0:n])
		//fmt.Println("bytes read: ", n)

	}
	fmt.Println("exiting from loop")
	buff := make([]byte, 1024)
	_, err = io.ReadAtLeast(stream, buff, 1)
	fmt.Println("last value", string(buff))
	if err != nil {
		panic(err)
	}
	nanosEnd = timehandle()
	timediff = nanosEnd - nanosBegin
	fmt.Println("delay in sending a file  throughquic in nanoseconds is ", timediff)
	stream.Close()
}
func main() {

	tcpclient()
	//fmt.Println("bufferone value after exiting tcpclient is ", bufferone)
	udpclient()
	//	fmt.Println("bufferone value after exiting udpclient is ", bufferone)
	qclient()
}
func timehandle() int64 {

	timeNow := time.Now()
	timeInNanos := timeNow.Unix()
	return timeInNanos
}
func sendfiletoserver(c net.Conn) {

	file, _ := os.Open("/home/varun/phone_backup_evolx")

	defer file.Close()
	buf := make([]byte, 1024)

	for {
		n, err := file.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}

			break

		}
		c.Write(buf[0:n])
		//fmt.Println("bytes read: ", n)

	}

}
