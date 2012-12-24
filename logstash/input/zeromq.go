package input

import (
	zmq "github.com/alecthomas/gozmq"
)

type ZeroMQ struct {
	sock zmq.Socket
    conn_string string
}

func NewZeroMQ() (*ZeroMQ) {
    context, _ := zmq.NewContext()
    socket, _ := context.NewSocket(zmq.PULL)
	return &ZeroMQ{socket, ""}
}

func (z *ZeroMQ) parseConfig(args map[string]interface{}) {
    z.conn_string = args["conn_string"].(string)
}

// Called once for every time a "grok" config is used
func (z *ZeroMQ) Register(args map[string]interface{}) (error) {
    z.parseConfig(args)
    z.sock.Connect(z.conn_string)
    return nil
}

func (z *ZeroMQ) Receive(output chan []byte) (error) {
    for {
      	msg, _ := z.sock.Recv(0)
		output <- msg
    }
    return nil
}