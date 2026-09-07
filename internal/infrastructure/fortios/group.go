package fortios

import (
	"context"
	"net/url"

	dgrp "github.com/local/fgcli/internal/domain/group"
)

// GroupRepo is path-parameterized: vipgrp and addrgrp share the shape.
type GroupRepo struct {
	c    *Client
	path string
}

func NewVipGrpRepo(c *Client) *GroupRepo {
	return &GroupRepo{c: c, path: "/api/v2/cmdb/firewall/vipgrp"}
}
func NewAddrGrpRepo(c *Client) *GroupRepo {
	return &GroupRepo{c: c, path: "/api/v2/cmdb/firewall/addrgrp"}
}

func (s *GroupRepo) List(ctx context.Context, vdom string) ([]dgrp.Group, error) {
	b, err := s.c.do(ctx, "GET", s.path, vdomQuery(vdom), nil)
	if err != nil {
		return nil, err
	}
	var out []dgrp.Group
	return out, decodeResults(b, &out)
}

func (s *GroupRepo) Get(ctx context.Context, name, vdom string) (dgrp.Group, error) {
	b, err := s.c.do(ctx, "GET", s.path+"/"+url.PathEscape(name), vdomQuery(vdom), nil)
	if err != nil {
		return dgrp.Group{}, err
	}
	var list []dgrp.Group
	if err := decodeResults(b, &list); err != nil {
		return dgrp.Group{}, err
	}
	if len(list) == 0 {
		return dgrp.Group{}, context.DeadlineExceeded
	}
	return list[0], nil
}

func (s *GroupRepo) Create(ctx context.Context, g dgrp.Group, vdom string) (dgrp.Group, error) {
	if _, err := s.c.do(ctx, "POST", s.path, vdomQuery(vdom), g); err != nil {
		return dgrp.Group{}, err
	}
	return g, nil
}

func (s *GroupRepo) Update(ctx context.Context, g dgrp.Group, vdom string) (dgrp.Group, error) {
	if _, err := s.c.do(ctx, "PUT", s.path+"/"+url.PathEscape(g.Name), vdomQuery(vdom), g); err != nil {
		return dgrp.Group{}, err
	}
	return g, nil
}

func (s *GroupRepo) Delete(ctx context.Context, name, vdom string) error {
	_, err := s.c.do(ctx, "DELETE", s.path+"/"+url.PathEscape(name), vdomQuery(vdom), nil)
	return err
}
