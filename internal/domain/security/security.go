package security

import "context"

type IpsAnomaly struct {
	ID        int64  `json:"id"`
	Name      string `json:"name,omitempty"`
	Threshold int64  `json:"threshold,omitempty"`
	Action    int    `json:"action,omitempty"`
}

type WafClass struct {
	ID      int64  `json:"id"`
	Name    string `json:"name,omitempty"`
	Comment string `json:"comment,omitempty"`
}

type DlpEntry struct {
	Pattern  string `json:"pattern,omitempty"`
	FileType string `json:"file-type,omitempty"`
}

type DlpPattern struct {
	ID      int64      `json:"id"`
	Name    string     `json:"name,omitempty"`
	Comment string     `json:"comment,omitempty"`
	Entries []DlpEntry `json:"entries,omitempty"`
}

type Repository interface {
	Ips(ctx context.Context, vdom string) ([]IpsAnomaly, error)
	ListWaf(ctx context.Context, vdom string) ([]WafClass, error)
	GetWaf(ctx context.Context, id int64, vdom string) (WafClass, error)
	ListDlp(ctx context.Context, vdom string) ([]DlpPattern, error)
	GetDlp(ctx context.Context, id int64, vdom string) (DlpPattern, error)
	DownloadPac(ctx context.Context, vdom string) ([]byte, error)
}
