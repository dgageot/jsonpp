package buffer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type Json interface {
	io.Reader
	fmt.Stringer
	Skip(skip bool)
	StartArray() error
	EndArray() error
	StartObject() error
	EndObject() error
	WriteComma() error
	WriteKey(key string) error
	WriteValue(value interface{}) error
}

type jsonBuffer struct {
	buffer *bytes.Buffer
	skip   bool
}

func NewJson() Json {
	return &jsonBuffer{
		buffer: new(bytes.Buffer),
	}
}

func (j *jsonBuffer) Skip(skip bool) {
	j.skip = skip
}

func (j *jsonBuffer) writeString(value string) error {
	if j.skip {
		return nil
	}
	_, err := j.buffer.WriteString(value)
	return err
}

func (j *jsonBuffer) StartArray() error {
	return j.writeString("[")
}

func (j *jsonBuffer) EndArray() error {
	return j.writeString("]")
}

func (j *jsonBuffer) StartObject() error {
	return j.writeString("{")
}

func (j *jsonBuffer) EndObject() error {
	return j.writeString("}")
}

func (j *jsonBuffer) WriteComma() error {
	return j.writeString(",")
}

func (j *jsonBuffer) WriteKey(key string) error {
	return j.writeString(`"` + key + `":`)
}

func (j *jsonBuffer) WriteValue(value interface{}) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return j.writeString(string(bytes))
}

func (j *jsonBuffer) String() string {
	return j.buffer.String()
}

func (j *jsonBuffer) Read(p []byte) (n int, err error) {
	return j.buffer.Read(p)
}
