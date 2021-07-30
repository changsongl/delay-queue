package encode

import (
	"bytes"
	"encoding/binary"
	"github.com/changsongl/delay-queue/job"
)

// |tag|[length]|data

// tag s
const (
	TagID uint64 = iota
	TagTopic
	TagDelay
	TagTTR
	TagBody
	TagVersion
)

const (
	// TagLength tag length
	TagLength = 1

	// MaxUInt64Length variant max uint64 length
	MaxUInt64Length = 10
)

type compress struct {
}

// NewCompress create a json encoder
func NewCompress() Encoder {
	return &compress{}
}

// Encode compress encode, not using reflect
func (c *compress) Encode(j *job.Job) ([]byte, error) {
	buf := make([]byte, c.bufLength(j))
	written := 0
	if !j.Delay.IsEmpty() {
		written += c.PutUInt64(TagDelay, uint64(j.Delay), buf[written:])
	}
	if !j.TTR.IsEmpty() {
		written += c.PutUInt64(TagTTR, uint64(j.TTR), buf[written:])
	}
	if !j.ID.IsEmpty() {
		written += c.PutString(TagID, string(j.ID), buf[written:])
	}
	if !j.Body.IsEmpty() {
		written += c.PutString(TagBody, string(j.Body), buf[written:])
	}
	if !j.Topic.IsEmpty() {
		written += c.PutString(TagTopic, string(j.Topic), buf[written:])
	}
	written += c.PutUInt64(TagVersion, j.Version.UInt64(), buf[written:])

	return buf[:written], nil
}

// Decode compress decode
func (c *compress) Decode(b []byte, j *job.Job) error {
	index := 0
	for index < len(b) {
		tag, err := binary.ReadUvarint(bytes.NewBuffer(b[index:]))
		if err != nil {
			return err
		}
		index++

		switch tag {
		case TagID:
			id, indexInc := c.ReadString(b[index:])
			j.ID = job.ID(id)
			index += indexInc
		case TagTopic:
			topic, indexInc := c.ReadString(b[index:])
			j.Topic = job.Topic(topic)
			index += indexInc
		case TagBody:
			body, indexInc := c.ReadString(b[index:])
			j.Body = job.Body(body)
			index += indexInc
		case TagTTR:
			ttr, indexInc := c.ReadUint64(b[index:])
			j.TTR = job.TTR(ttr)
			index += indexInc
		case TagDelay:
			delay, indexInc := c.ReadUint64(b[index:])
			j.Delay = job.Delay(delay)
			index += indexInc
		case TagVersion:
			ts, indexInc := c.ReadUint64(b[index:])
			j.SetVersion(int64(ts))
			index += indexInc
		}
	}
	return nil
}

func (c *compress) bufLength(j *job.Job) int {
	l := (TagLength+MaxUInt64Length)*5 + len(j.ID) + len(j.Topic)
	if j.Body != "" {
		l += TagLength + MaxUInt64Length + len(j.Body)
	}
	return l
}

func (c *compress) ReadUint64(buf []byte) (uint64, int) {
	return binary.Uvarint(buf)
}

func (c *compress) PutUInt64(tag uint64, num uint64, buf []byte) int {
	written := binary.PutUvarint(buf, tag)
	written += binary.PutUvarint(buf[written:], num)
	return written
}

func (c *compress) ReadString(buf []byte) (string, int) {
	l, inc := binary.Uvarint(buf)
	end := inc + int(l)
	return string(buf[inc:end]), end
}

func (c *compress) PutString(tag uint64, str string, buf []byte) int {
	l := len(str)
	written := binary.PutUvarint(buf, tag)
	written += binary.PutUvarint(buf[written:], uint64(l))
	chs := make([]uint8, 0, l)
	for _, ch := range str {
		chs = append(chs, uint8(ch))
	}

	for _, ch := range []byte(str) {
		buf[written] = ch
		written++
	}

	return written
}
