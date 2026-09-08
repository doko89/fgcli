package user

import (
	"errors"
	"strings"
)

// UserGroup is the cmdb/user/group subset fgcli manages: the group name
// plus its member login names (wire shape: member:[{name:...}]).
type UserGroup struct {
	Name    string        `json:"name"`
	Members []GroupMember `json:"member,omitempty"`
}

type GroupMember struct {
	Name string `json:"name"`
}

func (g UserGroup) Validate() error {
	if strings.TrimSpace(g.Name) == "" {
		return errors.New("user group: name is required")
	}
	return nil
}

// MemberNames returns the member login names in order.
func (g UserGroup) MemberNames() []string {
	out := make([]string, 0, len(g.Members))
	for _, m := range g.Members {
		out = append(out, m.Name)
	}
	return out
}

// AddUnique appends name when absent (case-insensitive); added=false when
// the member is already present.
func AddUnique(ms []GroupMember, name string) ([]GroupMember, bool) {
	for _, m := range ms {
		if strings.EqualFold(m.Name, name) {
			return ms, false
		}
	}
	return append(ms, GroupMember{Name: name}), true
}

// RemoveName drops name (case-insensitive); removed=false when absent.
func RemoveName(ms []GroupMember, name string) ([]GroupMember, bool) {
	kept := ms[:0]
	removed := false
	for _, m := range ms {
		if strings.EqualFold(m.Name, name) {
			removed = true
			continue
		}
		kept = append(kept, m)
	}
	if !removed {
		return ms, false
	}
	return kept, true
}

// LocalUser is the cmdb/user/local subset fgcli manages.
type LocalUser struct {
	Name   string `json:"name"`
	Passwd string `json:"passwd,omitempty"`
	Status string `json:"status,omitempty"`
	Type   string `json:"type,omitempty"`
}

func (u LocalUser) Validate() error {
	if strings.TrimSpace(u.Name) == "" {
		return errors.New("user local: name is required")
	}
	return nil
}
