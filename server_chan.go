package main

import (
	//"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"net"

	//"strings"
	"context"
	"math/big"
	//"time"

	"github.com/lithammer/shortuuid"
	quic "github.com/lucas-clemente/quic-go"
)

var a string

func main() {

	//fmt.Println("the value of a ", a)
	tcpserver()
	udpserver()
	quicserver()

}
func tcpserver() {

	PORT := ":1234"
	fmt.Println("tcp server listening on ", PORT)
	l, err := net.Listen("tcp", PORT)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()
	c, err1 := l.Accept()
	if err1 != nil {
		fmt.Println(err)
		return
	}
	if a == "" {
		a = genshortUUID()
		fmt.Println("unique id generated from tcpserver is ", a)
	}

	c.Write([]byte(a))
	handling := "tcp"
	handleClient(c, handling)
	fmt.Println("Done Handling")
	c.Write([]byte("File is read"))

	c.Close()
}

func udpserver() {

	PORT := ":1235"
	s, err := net.ResolveUDPAddr("udp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}

	connection, err1 := net.ListenUDP("udp4", s)
	fmt.Println("conenction ", connection)
	if err1 != nil {
		fmt.Println(err)
		return
	}
	//	defer connection.Close()
	//	bufferlocal := make([]byte, 1024)
	//	n, _, err := connection.ReadFromUDP(bufferlocal)
	//	fmt.Print("msg from udpclient ", string(bufferlocal[0:n-1]))
	buffer := make([]byte, 1024)
	n, addr, err := connection.ReadFromUDP(buffer)
	fmt.Print("-> ", string(buffer[0:n]))
	if err != nil {
		fmt.Println(err)
		return
	}
	if a == "" {
		a = genshortUUID()
		fmt.Println("\n unique id generated from udpserver is ", a)
		data := []byte(a)
		_, err = connection.WriteToUDP(data, addr)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		fmt.Println("the connection id already generated  in tcp is ", a)
	}

	//_, err = connection.WriteToUDP(data, addr)
	newbuffer := make([]byte, 1024)

	for {
		n, _, err2 := connection.ReadFromUDP(newbuffer)

		//		if string(buffer[0:n]) != "" {

		//			break

		//		}
		//fmt.Println("byteshere  read ", n)
		if n < 1024 {
			fmt.Println("condition for exiting")
			break
		}
		if err2 != nil {
			fmt.Println(err2)
			return
		}
	}

	fmt.Println("Done Handling")
	data := []byte("over from udp server")
	connection.WriteToUDP(data, addr)

}
func quicserver() {
	PORT := "localhost:1236"
	fmt.Println("Server running @", PORT)
	listener, err := quic.ListenAddr(PORT, generateTLSConfig(), nil)
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
	normalbuff := make([]byte, 100)
	io.ReadAtLeast(stream, normalbuff, 1)
	if a == "" {
		a = genshortUUID()
		fmt.Println("\n unique id generated from quicpserver is ", a)
	} else {
		fmt.Println("the connection id already generated  in tcp is ", a)
	}

	counter := 0

	for {

		if counter == 1 {
			break
		}

		io.WriteString(stream, a)

		counter++

	}

	for {
		buf := make([]byte, 1024)
		n, err := io.ReadAtLeast(stream, buf, 1)
		if err != nil {
			fmt.Println(err)
		}
		//fmt.Println("bytes read: ", n)
		if n < 1024 {
			break
		}

	}
	fmt.Println("done handling")
	io.WriteString(stream, "last message from server")
	for {

	}

}

func genshortUUID() string {

	id := shortuuid.New()

	return id
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
		NextProtos:   []string{"quic-sendreceive"},
	}
}
func handleClient(c net.Conn, method string) {
	buf := make([]byte, 1024)

	for {
		n, err := c.Read(buf)
		if err != nil {
			panic(err)
		}

		//fmt.Println("bytes read: ", string(buf[0:n]))
		if n < 1024 {
			fmt.Println("condition for exiting")
			break
		}

	}
}
