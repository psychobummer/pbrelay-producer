package producer

type Producer interface {
	Stream() <-chan []byte
}
