package network

import "context"

type DnsFilter struct {
	ID      int64  `json:"id"`
	Name    string `json:"name,omitempty"`
	Comment string `json:"comment,omitempty"`
}

type DnsFilterRepository interface {
	List(ctx context.Context, vdom string) ([]DnsFilter, error)
	Get(ctx context.Context, id int64, vdom string) (DnsFilter, error)
}
