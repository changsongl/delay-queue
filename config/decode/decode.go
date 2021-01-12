package decode

// Func is common decode function for files
type Func func([]byte, interface{}) error

// Decoder interface to return a decode function
// for files. different type of file has different
// decoder
type Decoder interface {
	DecodeFunc() Func
}
