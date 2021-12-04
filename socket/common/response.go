package common

import "encoding/json"

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
