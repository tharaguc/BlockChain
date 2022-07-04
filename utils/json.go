package utils

import "encoding/json"

//jsonに変換
func JsonStatus(message string) []byte {
	m, _ := json.Marshal(struct {
		Message string `json:"message"`
	}{
		Message: message,
	})
	return m
}
