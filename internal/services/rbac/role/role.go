package role

//go:generate go run ../../../../cmd/injecttrace -file role.go -receiver Role -service Role

import (
	"grubzo/internal/services/rbac/permission"
)

type Role struct {
	Name        string
	Permissions permission.Permissions
}

type Roles map[string]Role

func (r *Role) IsGranted(p permission.Permission) bool {
	return r.Permissions.Contains(p)
}

func (roles Roles) Add(role Role) {
	roles[role.Name] = role
}

func (roles Roles) IsGranted(p permission.Permission) bool {
	for _, v := range roles {
		if v.IsGranted(p) {
			return true
		}
	}
	return false
}

func (roles Roles) HasAndIsGranted(r string, p permission.Permission) bool {
	set, ok := roles[r]
	if !ok {
		return false
	}
	return set.IsGranted(p)
}
