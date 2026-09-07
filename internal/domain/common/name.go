package common

import "encoding/json"

// Name is a reference to another object. FortiOS returns these as
// {"name": "..."} objects (sometimes bare strings); it marshals back
// as a plain string for AI-friendly output.
type Name string

func (n *Name) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		*n = Name(s)
		return nil
	}
	var o struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(b, &o); err != nil {
		return err
	}
	*n = Name(o.Name)
	return nil
}
