package fortios

import (
	"context"
	"strconv"
	"strings"

	dsec "github.com/local/fgcli/internal/domain/security"
)

type SecurityRepo struct{ c *Client }

func NewSecurityRepo(c *Client) *SecurityRepo { return &SecurityRepo{c: c} }

func (s *SecurityRepo) Ips(ctx context.Context, vdom string) ([]dsec.IpsAnomaly, error) {
	b, err := s.c.do(ctx, "GET", "/api/v2/monitor/ips/anomaly", nil, nil)
	if err != nil {
		return nil, err
	}
	var out []dsec.IpsAnomaly
	if err := decodeResultsAllowEmpty(b, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *SecurityRepo) ListWaf(ctx context.Context, vdom string) ([]dsec.WafClass, error) {
	b, err := s.c.do(ctx, "GET", "/api/v2/cmdb/waf/main-class", vdomQuery(vdom), nil)
	if err != nil {
		return nil, err
	}
	var out []dsec.WafClass
	return out, decodeResults(b, &out)
}

func (s *SecurityRepo) GetWaf(ctx context.Context, id int64, vdom string) (dsec.WafClass, error) {
	b, err := s.c.do(ctx, "GET", "/api/v2/cmdb/waf/main-class/"+strconv.FormatInt(id, 10), vdomQuery(vdom), nil)
	if err != nil {
		return dsec.WafClass{}, err
	}
	var list []dsec.WafClass
	if err := decodeResults(b, &list); err != nil {
		return dsec.WafClass{}, err
	}
	if len(list) == 0 {
		return dsec.WafClass{}, context.DeadlineExceeded
	}
	return list[0], nil
}

func (s *SecurityRepo) ListDlp(ctx context.Context, vdom string) ([]dsec.DlpPattern, error) {
	b, err := s.c.do(ctx, "GET", "/api/v2/cmdb/dlp/filepattern", vdomQuery(vdom), nil)
	if err != nil {
		return nil, err
	}
	var out []dsec.DlpPattern
	return out, decodeResults(b, &out)
}

func (s *SecurityRepo) GetDlp(ctx context.Context, id int64, vdom string) (dsec.DlpPattern, error) {
	b, err := s.c.do(ctx, "GET", "/api/v2/cmdb/dlp/filepattern/"+strconv.FormatInt(id, 10), vdomQuery(vdom), nil)
	if err != nil {
		return dsec.DlpPattern{}, err
	}
	var list []dsec.DlpPattern
	if err := decodeResults(b, &list); err != nil {
		return dsec.DlpPattern{}, err
	}
	if len(list) == 0 {
		return dsec.DlpPattern{}, context.DeadlineExceeded
	}
	return list[0], nil
}

func (s *SecurityRepo) DownloadPac(ctx context.Context, vdom string) ([]byte, error) {
	b, err := s.c.do(ctx, "GET", "/api/v2/monitor/webproxy/pacfile/download", nil, nil)
	if err != nil {
		// HTTP 424: no PAC file configured — empty report, not a failure.
		if strings.Contains(err.Error(), "http 424") {
			return []byte{}, nil
		}
		return nil, err
	}
	return b, nil
}
