package input

import (
	zmq "github.com/alecthomas/gozmq"
)

type ZeroMQ struct {
	sock zmq.Socket
}

func NewZeroMQ() (*ZeroMQ) {
    context, _ := zmq.NewContext()
    socket, _ := context.NewSocket(zmq.PULL)
    socket.Connect("tcp://127.0.0.1:5000")
    socket.Connect("tcp://127.0.0.1:6000")
	return &ZeroMQ{socket}
}

func (z *ZeroMQ) Receive(output chan []byte) {
    for {
      	msg, _ := z.sock.Recv(0)
		output <- msg
    }
}