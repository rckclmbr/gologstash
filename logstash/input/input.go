package input

type InputType interface {
    Receive(chan []byte) error
}