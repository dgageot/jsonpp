package jsonpp

import (
	"encoding/json"
	"io"

	"github.com/dgageot/jsonpp/buffer"
	"github.com/dgageot/jsonpp/pull"
)

func Read(reader io.Reader, path []string, v interface{}) (error) {
	buf := buffer.NewJson()

	if err := pull.Extract(reader, buf, path); err != nil {
		return err
	}

	return json.NewDecoder(buf).Decode(v)
}

func Rewrite(reader io.Reader, path []string, value interface{}) (string, error) {
	buf := buffer.NewJson()

	if err := pull.Replace(reader, buf, path, value); err != nil {
		return "", err
	}

	return buf.String(), nil
}
