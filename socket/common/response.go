package common

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Response struct {
	Error error       `json:"error"`
	Data  interface{} `json:"data"`
	To    string      `json:"to"`
}

func (r *Response) MarshalJSON() ([]byte, error) {
	var err interface{} = r.Error
	if err != nil {
		err = r.Error.Error()
	}
	return json.Marshal(map[string]interface{}{
		"data":  r.Data,
		"to":    r.To,
		"error": err,
	})
}

func getStringField(res map[string]interface{}, fieldName string) (string, error) {
	rawValue, found := res[fieldName]
	if !found {
		return "", fmt.Errorf("field %s not found", fieldName)
	}
	strValue, ok := rawValue.(string)
	if !ok {
		return "", fmt.Errorf("field %s is not a string", fieldName)
	}
	return strValue, nil
}

func (r *Response) UnmarshalJSON(b []byte) error {
	var res map[string]interface{}
	err := json.Unmarshal(b, &res)
	if err != nil {
		return err
	}

	r.Data = res["data"]

	to, err := getStringField(res, "to")
	if err != nil {
		return err
	}
	r.To = to

	resErr, err := getStringField(res, "error")
	if err != nil {
		r.Error = nil
	} else {
		r.Error = errors.New(resErr)
	}

	return nil
}
