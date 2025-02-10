package request

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func Post(url, contentType string, body any) error {
	marshaled, err := json.Marshal(body)
	if err != nil {
		return err
	}

	if _, err := http.Post(url, contentType, bytes.NewReader(marshaled)); err != nil {
		return err
	}

	return nil
}
