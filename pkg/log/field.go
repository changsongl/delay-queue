package log

import (
	"go.uber.org/zap"
	"time"
)

type FieldType int

const (
	FieldTypeInt FieldType = iota + 1
	FieldTypeString
	FieldTypeDuration
	FieldTypeInterface
)

type Field struct {
	Key       string
	Type      FieldType
	Integer   int64
	String    string
	Interface interface{}
}

func Int8(key string, value int8) Field {
	return Field{Key: key, Integer: int64(value), Type: FieldTypeInt}
}

func Int16(key string, value int16) Field {
	return Field{Key: key, Integer: int64(value), Type: FieldTypeInt}
}

func Int32(key string, value int32) Field {
	return Field{Key: key, Integer: int64(value), Type: FieldTypeInt}
}

func Int64(key string, value int64) Field {
	return Field{Key: key, Integer: value, Type: FieldTypeInt}
}

func Int(key string, value int) Field {
	return Field{Key: key, Integer: int64(value), Type: FieldTypeInt}
}

func String(key string, value string) Field {
	return Field{Key: key, String: value, Type: FieldTypeString}
}

func Duration(key string, value time.Duration) Field {
	return Field{Key: key, Interface: value, Type: FieldTypeDuration}
}

func Reflect(key string, value interface{}) Field {
	return Field{Key: key, Interface: value, Type: FieldTypeInterface}
}

func getZapFields(fields ...Field) []zap.Field {
	var fs []zap.Field
	for _, f := range fields {
		if f.Type == FieldTypeInt {
			fs = append(fs, zap.Int64(f.Key, f.Integer))
		} else if f.Type == FieldTypeDuration {
			t, ok := f.Interface.(time.Duration)
			if ok {
				fs = append(fs, zap.Duration(f.Key, t))
			}
		} else if f.Type == FieldTypeInterface {
			fs = append(fs, zap.Reflect(f.Key, f.String))
		} else {
			fs = append(fs, zap.String(f.Key, f.String))
		}
	}

	return fs
}
