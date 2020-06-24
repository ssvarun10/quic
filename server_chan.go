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
	"time"

	"github.com/lithammer/shortuuid"
	quic "github.com/lucas-clemente/quic-go"
)

var a string

func main() {

	//	channelone := make(chan string)
	//	channeltwo := make(chan string)
	//go tcpserver(channelone)
	//go udpserver(channeltwo)

	fmt.Println("the value of a ", a)
	tcpserver()
	udpserver()
	quicserver()
	//udpserver()

	//		select{

	//			case msg1:=<-channelone:

	//			case msg2:=<channeltwo:

	//		}

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

	//	netDatafromclient, err := bufio.NewReaderSize(c, 1024).ReadString('\n')
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}

}

func udpserver() {

	PORT := ":1235"
	s, err := net.ResolveUDPAddr("udp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}

	connection, err1 := net.ListenUDP("udp4", s)
	if err1 != nil {
		fmt.Println(err)
		return
	}
	defer connection.Close()
	bufferlocal := make([]byte, 1024)
	n, addr, err := connection.ReadFromUDP(bufferlocal)
	fmt.Print("msg from udpclient ", string(bufferlocal[0:n-1]))
	if a == "" {
		a = genshortUUID()
		fmt.Println("\n unique id generated from udpserver is ", a)
	} else {
		fmt.Println("the connection id already generated  in tcp is ", a)
	}

	data := []byte(a)
	_, err = connection.WriteToUDP(data, addr)

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
	if a == "" {
		a = genshortUUID()
		fmt.Println("\n unique id generated from quicpserver is ", a)
	} else {
		fmt.Println("the connection id already generated  in tcp is ", a)
	}

	counter := 0
	//buff := make([]byte, 1024)
	for {

		if counter == 1 {
			break
		}

		// Echo through the loggingWriter

		io.WriteString(stream, a)
		time.Sleep(5 * time.Second)
		counter++

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
