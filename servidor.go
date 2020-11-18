package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"strconv"
	"time"
)

//Variables globales
var canal chan string
var cantidadProcesos uint64
var procesos []Proceso
var eliminarProceso uint64
var estadoDeProcesos map[uint64]bool

// Estructura de Proceso que va a guardar la id del proceso y el contador que va a ir cambiando cada ciertos milisegundos
type Proceso struct {
	Id       uint64
	Contador uint64
}

// Funcion que se llama para generar un slice de la cantidad de procesos asignada
func generarProcesos(cantidad uint64) []Proceso {
	procesos := make([]Proceso, cantidad)
	for id := uint64(0); id < cantidad; id++ {
		procesos[id] = Proceso{id, 0}
	}
	return procesos
}

// Gorutine que va a imprimir el mensaje actual
func impresora() {
	for {
		fmt.Println(<-canal)
	}
}

// Funcion con la cual se va calculando el valor actual a menos que se pida eliminar el proceso
func proceso(p *Proceso) {
	for eliminarProceso != p.Id {
		p.Contador++
		canal <- strconv.FormatUint(p.Id, 10) + " : " + strconv.FormatUint(p.Contador, 10)
		time.Sleep(time.Millisecond * 500)
	}
}

// Funcion net que genera el servidor tcp
func server() {
	receptor, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		cliente, err := receptor.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handler(cliente)
	}
}

// gorutine que va a estar controlando las diferentes conneciones por parte de los clientes
func handler(cliente net.Conn) {
	var proceso Proceso
	for {
		err := gob.NewDecoder(cliente).Decode(&proceso)
		if err != nil {
			fmt.Println(err)
			return
		} else {
			switch proceso.Contador {
			case 0:
				enviarProceso(cliente)
			default:
				recibirProceso(proceso)
			}
		}
	}
}

// Funcion con la cual mandaremos un proceso a un cliente
func enviarProceso(cliente net.Conn) {
	procesoAEnviar := cantidadProcesos
	for i := uint64(0); i < cantidadProcesos; i++ {
		if estadoDeProcesos[i] == true {
			procesoAEnviar = i
			break
		}
	}
	if procesoAEnviar == cantidadProcesos {
		return
	}
	eliminarProceso = procesoAEnviar
	err := gob.NewEncoder(cliente).Encode(procesos[procesoAEnviar])
	if err != nil {
		fmt.Println(err)
		go proceso(&procesos[procesoAEnviar])
		return
	} else {
		estadoDeProcesos[procesoAEnviar] = false
	}
}

// Funcion con la cual vamos a recibir un proceso y crearemos una instancia nueva con los datos actuales
func recibirProceso(p Proceso) {
	eliminarProceso = cantidadProcesos
	procesos[p.Id] = p
	estadoDeProcesos[p.Id] = true
	go proceso(&procesos[p.Id])
}

func main() {
	canal = make(chan string)
	cantidadProcesos = uint64(5)
	estadoDeProcesos = make(map[uint64]bool)
	procesos = generarProcesos(cantidadProcesos)
	eliminarProceso = cantidadProcesos
	var pausa string
	go impresora()
	for i, _ := range procesos {
		go proceso(&procesos[i])
		estadoDeProcesos[uint64(i)] = true
	}
	go server()
	fmt.Scanln(&pausa)
}
