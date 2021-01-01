package decode

type Func func([]byte, interface{}) error

type Decoder interface {
	DecodeFunc() Func
}
