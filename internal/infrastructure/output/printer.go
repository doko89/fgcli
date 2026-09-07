package output

import (
	"encoding/json"
	"fmt"
	"os"
)

type Envelope struct {
	OK    bool `json:"ok"`
	Data  any  `json:"data,omitempty"`
	Error *Err `json:"error,omitempty"`
}

type Err struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Print(data any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(Envelope{OK: true, Data: data})
}

func Fail(code string, err error) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	msg := ""
	if err != nil {
		msg = err.Error()
	}
	_ = enc.Encode(Envelope{OK: false, Error: &Err{Code: code, Message: msg}})
}

func Logf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}
