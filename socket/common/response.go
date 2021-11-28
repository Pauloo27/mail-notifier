package common

import "encoding/json"

type Response struct {
	Error error       `json:"error"`
	Data  interface{} `json:"data"`
}

func (r *Response) MarshalJSON() ([]byte, error) {
	var err interface{} = r.Error
	if err != nil {
		err = r.Error.Error()
	}
	return json.Marshal(map[string]interface{}{
		"data":  r.Data,
		"error": err,
	})
}

/* TODO
func (r *Response) UnmarshalJSON() ([]byte, error) {
	return nil, nil
}
*/
