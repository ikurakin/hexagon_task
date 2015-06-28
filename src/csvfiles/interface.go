package csvfiles

type HexWriter interface {
	GetData() interface{}
	Write(p []byte) (n int, err error)
}
