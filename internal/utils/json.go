package utils

import (
	"bytes"
	"encoding/json"
)

func MarshalJSON(data interface{}) ([]byte, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func UnmarshalJSON(data []byte, writeTo interface{}) error {
	err := json.NewDecoder(bytes.NewReader(data)).Decode(writeTo)
	return err
}
