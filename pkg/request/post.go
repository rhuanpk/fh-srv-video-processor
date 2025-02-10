package request

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func Put(url, contentType string, body any) error {
	marshaled, err := json.Marshal(body)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(marshaled))
	if err != nil {
		return err
	}
	if _, err := new(http.Client).Do(request); err != nil {
		return err
	}

	return nil
}
