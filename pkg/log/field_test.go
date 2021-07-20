package log

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"math"
	"net"
	"testing"
	"time"
)

type username string

func (n username) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("username", string(n))
	return nil
}

func TestFieldConstructors(t *testing.T) {
	// Interface types.
	addr := net.ParseIP("1.2.3.4")
	name := username("phil")
	ints := []int{5, 6}

	// Helpful values for use in constructing pointers to primitives below.
	var (
		boolVal       bool          = true
		complex128Val complex128    = complex(0, 0)
		complex64Val  complex64     = complex(0, 0)
		durationVal   time.Duration = time.Second
		float64Val    float64       = 1.0
		float32Val    float32       = 1.0
		intVal        int           = 1
		int64Val      int64         = 1
		int32Val      int32         = 1
		int16Val      int16         = 1
		int8Val       int8          = 1
		stringVal     string        = "hello"
		timeVal       time.Time     = time.Unix(100000, 0)
		uintVal       uint          = 1
		uint64Val     uint64        = 1
		uint32Val     uint32        = 1
		uint16Val     uint16        = 1
		uint8Val      uint8         = 1
		uintptrVal    uintptr       = 1
	)

	tests := []struct {
		name   string
		field  Field
		expect Field
	}{
		{"Binary", Field{zapField: zap.Field{Key: "k", Type: zapcore.BinaryType, Interface: []byte("ab12")}}, Binary("k", []byte("ab12"))},
		{"Bool", Field{zapField: zap.Field{Key: "k", Type: zapcore.BoolType, Integer: 1}}, Bool("k", true)},
		{"Bool", Field{zapField: zap.Field{Key: "k", Type: zapcore.BoolType, Integer: 1}}, Bool("k", true)},
		{"ByteString", Field{zapField: zap.Field{Key: "k", Type: zapcore.ByteStringType, Interface: []byte("ab12")}}, ByteString("k", []byte("ab12"))},
		{"Complex128", Field{zapField: zap.Field{Key: "k", Type: zapcore.Complex128Type, Interface: 1 + 2i}}, Complex128("k", 1+2i)},
		{"Complex64", Field{zapField: zap.Field{Key: "k", Type: zapcore.Complex64Type, Interface: complex64(1 + 2i)}}, Complex64("k", 1+2i)},
		{"Duration", Field{zapField: zap.Field{Key: "k", Type: zapcore.DurationType, Integer: 1}}, Duration("k", 1)},
		{"Int", Field{zapField: zap.Field{Key: "k", Type: zapcore.Int64Type, Integer: 1}}, Int("k", 1)},
		{"Int64", Field{zapField: zap.Field{Key: "k", Type: zapcore.Int64Type, Integer: 1}}, Int64("k", 1)},
		{"Int32", Field{zapField: zap.Field{Key: "k", Type: zapcore.Int32Type, Integer: 1}}, Int32("k", 1)},
		{"Int16", Field{zapField: zap.Field{Key: "k", Type: zapcore.Int16Type, Integer: 1}}, Int16("k", 1)},
		{"Int8", Field{zapField: zap.Field{Key: "k", Type: zapcore.Int8Type, Integer: 1}}, Int8("k", 1)},
		{"String", Field{zapField: zap.Field{Key: "k", Type: zapcore.StringType, String: "foo"}}, String("k", "foo")},
		{"Time", Field{zapField: zap.Field{Key: "k", Type: zapcore.TimeType, Integer: 0, Interface: time.UTC}}, Time("k", time.Unix(0, 0).In(time.UTC))},
		{"Time", Field{zapField: zap.Field{Key: "k", Type: zapcore.TimeType, Integer: 1000, Interface: time.UTC}}, Time("k", time.Unix(0, 1000).In(time.UTC))},
		{"Time", Field{zapField: zap.Field{Key: "k", Type: zapcore.TimeType, Integer: math.MinInt64, Interface: time.UTC}}, Time("k", time.Unix(0, math.MinInt64).In(time.UTC))},
		{"Time", Field{zapField: zap.Field{Key: "k", Type: zapcore.TimeType, Integer: math.MaxInt64, Interface: time.UTC}}, Time("k", time.Unix(0, math.MaxInt64).In(time.UTC))},
		{"Time", Field{zapField: zap.Field{Key: "k", Type: zapcore.TimeFullType, Interface: time.Time{}}}, Time("k", time.Time{})},
		{"Time", Field{zapField: zap.Field{Key: "k", Type: zapcore.TimeFullType, Interface: time.Unix(math.MaxInt64, 0)}}, Time("k", time.Unix(math.MaxInt64, 0))},
		{"Uint", Field{zapField: zap.Field{Key: "k", Type: zapcore.Uint64Type, Integer: 1}}, Uint("k", 1)},
		{"Uint64", Field{zapField: zap.Field{Key: "k", Type: zapcore.Uint64Type, Integer: 1}}, Uint64("k", 1)},
		{"Uint32", Field{zapField: zap.Field{Key: "k", Type: zapcore.Uint32Type, Integer: 1}}, Uint32("k", 1)},
		{"Uint16", Field{zapField: zap.Field{Key: "k", Type: zapcore.Uint16Type, Integer: 1}}, Uint16("k", 1)},
		{"Uint8", Field{zapField: zap.Field{Key: "k", Type: zapcore.Uint8Type, Integer: 1}}, Uint8("k", 1)},
		{"Uintptr", Field{zapField: zap.Field{Key: "k", Type: zapcore.UintptrType, Integer: 10}}, Uintptr("k", 0xa)},
		{"Reflect", Field{zapField: zap.Field{Key: "k", Type: zapcore.ReflectType, Interface: ints}}, Reflect("k", ints)},
		{"Reflect", Field{zapField: zap.Field{Key: "k", Type: zapcore.ReflectType}}, Reflect("k", nil)},
		{"Stringer", Field{zapField: zap.Field{Key: "k", Type: zapcore.StringerType, Interface: addr}}, Stringer("k", addr)},
		{"Object", Field{zapField: zap.Field{Key: "k", Type: zapcore.ObjectMarshalerType, Interface: name}}, Object("k", name)},
		{"Any:ObjectMarshaler", Any("k", name), Object("k", name)},
		{"Any:Stringer", Any("k", addr), Stringer("k", addr)},
		{"Any:Bool", Any("k", true), Bool("k", true)},
		{"Any:Byte", Any("k", byte(1)), Uint8("k", 1)},
		{"Any:Bytes", Any("k", []byte{1}), Binary("k", []byte{1})},
		{"Any:Complex128", Any("k", 1+2i), Complex128("k", 1+2i)},
		{"Any:Complex64", Any("k", complex64(1+2i)), Complex64("k", 1+2i)},
		{"Any:Float64", Any("k", 3.14), Float64("k", 3.14)},
		{"Any:Float32", Any("k", float32(3.14)), Float32("k", 3.14)},
		{"Any:Int", Any("k", 1), Int("k", 1)},
		{"Any:Int64", Any("k", int64(1)), Int64("k", 1)},
		{"Any:Int32", Any("k", int32(1)), Int32("k", 1)},
		{"Any:Int16", Any("k", int16(1)), Int16("k", 1)},
		{"Any:Int8", Any("k", int8(1)), Int8("k", 1)},
		{"Any:Rune", Any("k", rune(1)), Int32("k", 1)},
		{"Any:String", Any("k", "v"), String("k", "v")},
		{"Any:Uint", Any("k", uint(1)), Uint("k", 1)},
		{"Any:Uint64", Any("k", uint64(1)), Uint64("k", 1)},
		{"Any:Uint32", Any("k", uint32(1)), Uint32("k", 1)},
		{"Any:Uint16", Any("k", uint16(1)), Uint16("k", 1)},
		{"Any:Uint8", Any("k", uint8(1)), Uint8("k", 1)},
		{"Any:Uint8s", Any("k", []uint8{1}), Binary("k", []uint8{1})},
		{"Any:Uintptr", Any("k", uintptr(1)), Uintptr("k", 1)},
		{"Any:Time", Any("k", time.Unix(0, 0)), Time("k", time.Unix(0, 0))},
		{"Any:Duration", Any("k", time.Second), Duration("k", time.Second)},
		{"Any:Fallback", Any("k", struct{}{}), Reflect("k", struct{}{})},
		{"Ptr:Bool", Boolp("k", &boolVal), Bool("k", boolVal)},
		{"Any:PtrBool", Any("k", &boolVal), Bool("k", boolVal)},
		{"Ptr:Complex128", Complex128p("k", &complex128Val), Complex128("k", complex128Val)},
		{"Any:PtrComplex128", Any("k", &complex128Val), Complex128("k", complex128Val)},
		{"Ptr:Complex64", Complex64p("k", &complex64Val), Complex64("k", complex64Val)},
		{"Any:PtrComplex64", Any("k", &complex64Val), Complex64("k", complex64Val)},
		{"Ptr:Duration", Durationp("k", &durationVal), Duration("k", durationVal)},
		{"Any:PtrDuration", Any("k", &durationVal), Duration("k", durationVal)},
		{"Ptr:Float64", Float64p("k", &float64Val), Float64("k", float64Val)},
		{"Any:PtrFloat64", Any("k", &float64Val), Float64("k", float64Val)},
		{"Ptr:Float32", Float32p("k", &float32Val), Float32("k", float32Val)},
		{"Any:PtrFloat32", Any("k", &float32Val), Float32("k", float32Val)},
		{"Ptr:Int", Intp("k", &intVal), Int("k", intVal)},
		{"Any:PtrInt", Any("k", &intVal), Int("k", intVal)},
		{"Ptr:Int64", Int64p("k", &int64Val), Int64("k", int64Val)},
		{"Any:PtrInt64", Any("k", &int64Val), Int64("k", int64Val)},
		{"Ptr:Int32", Int32p("k", &int32Val), Int32("k", int32Val)},
		{"Any:PtrInt32", Any("k", &int32Val), Int32("k", int32Val)},
		{"Ptr:Int16", Int16p("k", &int16Val), Int16("k", int16Val)},
		{"Any:PtrInt16", Any("k", &int16Val), Int16("k", int16Val)},
		{"Ptr:Int8", Int8p("k", &int8Val), Int8("k", int8Val)},
		{"Any:PtrInt8", Any("k", &int8Val), Int8("k", int8Val)},
		{"Ptr:String", Stringp("k", &stringVal), String("k", stringVal)},
		{"Any:PtrString", Any("k", &stringVal), String("k", stringVal)},
		{"Ptr:Time", Timep("k", &timeVal), Time("k", timeVal)},
		{"Any:PtrTime", Any("k", &timeVal), Time("k", timeVal)},
		{"Ptr:Uint", Uintp("k", &uintVal), Uint("k", uintVal)},
		{"Any:PtrUint", Any("k", &uintVal), Uint("k", uintVal)},
		{"Ptr:Uint64", Uint64p("k", &uint64Val), Uint64("k", uint64Val)},
		{"Any:PtrUint64", Any("k", &uint64Val), Uint64("k", uint64Val)},
		{"Ptr:Uint32", Uint32p("k", &uint32Val), Uint32("k", uint32Val)},
		{"Any:PtrUint32", Any("k", &uint32Val), Uint32("k", uint32Val)},
		{"Ptr:Uint16", Uint16p("k", &uint16Val), Uint16("k", uint16Val)},
		{"Any:PtrUint16", Any("k", &uint16Val), Uint16("k", uint16Val)},
		{"Ptr:Uint8", Uint8p("k", &uint8Val), Uint8("k", uint8Val)},
		{"Any:PtrUint8", Any("k", &uint8Val), Uint8("k", uint8Val)},
		{"Ptr:Uintptr", Uintptrp("k", &uintptrVal), Uintptr("k", uintptrVal)},
		{"Any:PtrUintptr", Any("k", &uintptrVal), Uintptr("k", uintptrVal)},
	}

	for _, tt := range tests {
		if !assert.Equal(t, tt.expect, tt.field, "Unexpected output from convenience field constructor %s.", tt.name) {
			t.Logf("type expected: %T\nGot: %T", tt.expect.zapField.Interface, tt.field.zapField.Interface)
		}
	}
}
