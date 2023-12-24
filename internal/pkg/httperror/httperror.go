package httperror

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type body struct {
	Message string `json:"message"`
}

func Write(w http.ResponseWriter, status int, message ...any) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)

	raw, err := json.Marshal(body{
		Message: fmt.Sprint(message...),
	})

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(raw)
}

func Writef(w http.ResponseWriter, status int, format string, message ...any) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)

	raw, err := json.Marshal(body{
		Message: fmt.Sprintf(format, message...),
	})

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(raw)
}
