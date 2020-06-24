package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"log"

	"net"
	"strings"
	"time"

	"math/big"

	quic "github.com/lucas-clemente/quic-go"
)

func tcp(PORT string) {
	fmt.Println("listening on port", PORT)
	l, err := net.Listen("tcp", ":4242")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("1")
	defer l.Close()

	c, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		// will listen for message to process ending in newline (\n)
		message, _ := bufio.NewReader(c).ReadString('\n')
		// output message received
		fmt.Print("Message Received:", string(message))
		// sample process for string received
		newmessage := strings.ToUpper(message)
		// send new string back to client
		c.Write([]byte(newmessage + "\n"))
	}
}
func udp(PORT string) {
	udpaddr, _ := net.ResolveUDPAddr("udp", PORT)
	ServerConn, _ := net.ListenUDP("udp", udpaddr)
	defer ServerConn.Close()
	buf := make([]byte, 1024)

	for {
		n, addr, _ := ServerConn.ReadFromUDP(buf)
		fmt.Println("Received ", string(buf[:n]), " from ", addr)

		cmdvar := string(buf)
		cmdvar = strings.TrimSpace(cmdvar)
		if cmdvar == "STOP" {
			fmt.Println("Exiting UDP server!")
			return
		}
	}
}
func quicserver(addr string) {
	timeoutDuration := flag.Int("rtt", 3000, "timeout duration (in ms)")
	flag.Parse()
	fmt.Println("ghai server")
	listener, err := quic.ListenAddr(addr, generateTLSConfig(), nil)
	if err != nil {
		panic(err)
	}

	sess, err := listener.Accept(context.Background())
	if err != nil {
		panic(err)
	}

	stream, err := sess.AcceptStream(context.Background())
	if err != nil {
		panic(err)
	}

	timeout := time.Duration(*timeoutDuration) * time.Millisecond

	resp := make(chan string)

	message := "hello mydear client"

	log.Printf("server: Sending '%s'\n", message)
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
		log.Printf("server: Got '%s'\n", reply)
	case <-time.After(timeout):
		log.Printf("server: Timed out\n")
	}

}

func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quic-echo"},
	}
}

func main() {
	hostName := flag.String("hostname", "localhost", "hostname/ip of the server")
	portNum := flag.String("port", "4242", "port number of the server")

	flag.Parse()

	addr := *hostName + ":" + *portNum
	tcp(addr)
	//quicserver(addr)
}
