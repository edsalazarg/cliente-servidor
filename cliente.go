package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"strconv"
	"time"
)

var canal chan string

type Proceso struct {
	Id       uint64
	Contador uint64
}

func impresora() {
	for {
		fmt.Println(<-canal)
	}
}

func procesar(p *Proceso) {
	for {
		p.Contador++
		canal <- strconv.FormatUint(p.Id, 10) + " : " + strconv.FormatUint(p.Contador, 10)
		time.Sleep(time.Millisecond * 500)
	}
}

func main() {
	var input string
	canal = make(chan string)
	proceso := Proceso{0, 0}
	conexion, err := net.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = gob.NewEncoder(conexion).Encode(proceso)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = gob.NewDecoder(conexion).Decode(&proceso)
	if err != nil {
		fmt.Println(err)
		return
	}
	go impresora()
	go procesar(&proceso)
	fmt.Scanln(&input)
	err = gob.NewEncoder(conexion).Encode(proceso)
	if err != nil {
		fmt.Println(err)
	}
	conexion.Close()
}
