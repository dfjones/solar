package main

import (
	"encoding/binary"
	"github.com/tarm/serial"
	"log"
	"time"
)

func main() {

	c := &serial.Config{Name: "/dev/cu.usbserial-A60205DM", Baud: 9600}

	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	defer s.Close()

	time.Sleep(3 * time.Second)

	binary.Write(s, binary.LittleEndian, []byte("255,0,0,"))

	s.Flush()

	in := make([]byte, 256)
	n, err := s.Read(in)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%q", in[:n])

	for {
		time.Sleep(10 * time.Second)
		log.Println("Waiting...")
	}
}
