package svc

// Slim views of system monitor blobs (full objects are large and vary;
// only audit-relevant fields are kept, the rest is ignored on decode).

type SupportInfo struct {
	Status string `json:"status,omitempty"`
	Level  string `json:"support_level,omitempty"`
	Expiry int64  `json:"expires,omitempty"`
}

type Forticare struct {
	Status  string `json:"status,omitempty"`
	Account string `json:"account,omitempty"`
	Support struct {
		Hardware SupportInfo `json:"hardware"`
		Enhanced SupportInfo `json:"enhanced"`
	} `json:"support"`
}

type LicenseFG struct {
	Supported bool   `json:"supported,omitempty"`
	Connected bool   `json:"connected,omitempty"`
	Server    string `json:"server_address,omitempty"`
}

type License struct {
	Forticare  Forticare `json:"forticare"`
	Fortiguard LicenseFG `json:"fortiguard"`
}

type StatWindow struct {
	Hour []int64 `json:"1_hour,omitempty"`
	Day  []int64 `json:"24_hour,omitempty"`
	Week []int64 `json:"1_week,omitempty"`
}

type Fortiguard map[string]StatWindow

type NtpServer struct {
	Server    string  `json:"server,omitempty"`
	IP        string  `json:"ip,omitempty"`
	Reachable bool    `json:"reachable"`
	Stratum   int     `json:"stratum"`
	Offset    float64 `json:"offset"`
}

type DnsEntry struct {
	Interface string `json:"interface,omitempty"`
	Server    string `json:"server,omitempty"`
}

type DhcpEntry struct {
	Interface string `json:"interface,omitempty"`
	Status    string `json:"status,omitempty"`
}

type Snmp struct {
	Status      string `json:"status"`
	Description string `json:"description,omitempty"`
	Contact     string `json:"contact-info,omitempty"`
	Location    string `json:"location,omitempty"`
}
