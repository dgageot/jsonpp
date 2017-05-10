package pull

import (
	"encoding/json"
	"fmt"
	"io"
	"github.com/dgageot/jsonpp/buffer"
	"github.com/dgageot/jsonpp/paths"
	"github.com/pkg/errors"
)

type parser struct {
	*json.Decoder
	buffer.Json

	skipMatching bool
	path         []string
	setValue     interface{}
}

func Replace(reader io.Reader, buf buffer.Json, path []string, value interface{}) error {
	decoder := json.NewDecoder(reader)

	parser := &parser{
		Decoder:      decoder,
		Json:         buf,
		skipMatching: true,
		path:         path,
		setValue:     value,
	}

	buf.Skip(false)

	return parser.parseJson(nil)
}

func Extract(reader io.Reader, buf buffer.Json, path []string) error {
	decoder := json.NewDecoder(reader)

	parser := &parser{
		Decoder:      decoder,
		Json:         buf,
		skipMatching: false,
		path:         path,
		setValue:     nil,
	}

	buf.Skip(true)

	return parser.parseJson(nil)
}

func (p *parser) parseJson(currentPath []string) error {
	if p.path != nil && paths.Matches(currentPath, p.path) {
		p.Skip(p.skipMatching)
	}

	token, err := p.Token()
	if err == io.EOF {
		return err
	}

	switch token := token.(type) {
	case json.Delim:
		switch token.String() {
		case "{":
			if err := p.parseObject(currentPath); err != nil {
				return err
			}
		case "[":
			if err := p.parseArray(currentPath); err != nil {
				return err
			}
		default:
			return errors.New("Invalid delimiter")
		}
	default:
		if err := p.WriteValue(token); err != nil {
			return err
		}
	}

	if p.path != nil && paths.Matches(currentPath, p.path) {
		p.Skip(!p.skipMatching)
		if p.skipMatching {
			return p.WriteValue(p.setValue)
		}
	}

	return nil
}

func (p *parser) parseArray(currentPath []string) error {
	if err := p.StartArray(); err != nil {
		return err
	}

	needsComma := false
	for i := 0; p.More(); i++ {
		if needsComma {
			if err := p.WriteComma(); err != nil {
				return err
			}
		}

		newPath := append(currentPath, fmt.Sprintf("[%d]", i))
		if err := p.parseJson(newPath); err != nil {
			return err
		}

		needsComma = true
	}

	token, err := p.Token()
	if err == io.EOF {
		return err
	}

	delimiter, ok := token.(json.Delim)
	if !ok || delimiter.String() != "]" {
		return errors.New("Invalid end of array: " + delimiter.String())
	}

	return p.EndArray()
}

func (p *parser) parseObject(currentPath []string) error {
	if err := p.StartObject(); err != nil {
		return err
	}

	needsComma := false
	for p.More() {
		token, err := p.Token()
		if err == io.EOF {
			return err
		}

		key, ok := token.(string)
		if !ok {
			return errors.New("Invalid key")
		}

		if needsComma {
			if err := p.WriteComma(); err != nil {
				return err
			}
		}

		if err := p.WriteKey(key); err != nil {
			return err
		}

		newPath := append(currentPath, key)
		if err = p.parseJson(newPath); err != nil {
			return err
		}

		needsComma = true
	}

	token, err := p.Token()
	if err == io.EOF {
		return err
	}

	delimiter, ok := token.(json.Delim)
	if !ok || delimiter.String() != "}" {
		return errors.New("Invalid end of object: " + delimiter.String())
	}

	return p.EndObject()
}
