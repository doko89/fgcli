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

var Pretty bool

func Print(data any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	if Pretty {
		enc.SetIndent("", "  ")
	}
	_ = enc.Encode(Envelope{OK: true, Data: data})
}

func Fail(code string, err error) {
	msg := ""
	if err != nil {
		msg = err.Error()
	}
	env := Envelope{OK: false, Error: &Err{Code: code, Message: msg}}
	if Pretty {
		b, _ := json.MarshalIndent(env, "", "  ")
		fmt.Fprintln(os.Stdout, string(b))
		return
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(env)
}

func Logf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}
