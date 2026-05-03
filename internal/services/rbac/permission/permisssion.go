package permission

//go:generate go run ../../../../cmd/injecttrace -file permisssion.go -receiver Permissions -service Permissions

type Permission string

type Permissions map[Permission]struct{}

func (set Permissions) Add(p Permission) {
	set[p] = struct{}{}
}

func (set Permissions) Remove(p Permission) {
	delete(set, p)
}

func (set Permissions) Contains(p Permission) bool {
	_, ok := set[p]
	return ok
}

func (set Permissions) Array() []Permission {
	result := make([]Permission, 0, len(set))
	for k := range set {
		result = append(result, k)
	}
	return result
}

func PermissionsFromArray(perms []Permission) Permissions {
	res := Permissions{}
	for _, perm := range perms {
		res.Add(perm)
	}
	return res
}

const (
	Dashboard Permission = "dashboard"
	Orders    Permission = "orders"
	Items     Permission = "items"
	Employee  Permission = "employee"
	Location  Permission = "location"
	RBAC      Permission = "rbac"
)

var (
	DashboardPermissions = []Permission{Dashboard}
	OrderPermissions     = []Permission{Orders}
	MenuPermissions      = []Permission{Items}
	EmployeePermissions  = []Permission{Employee}
	LocationPermissions  = []Permission{Location}
	RBACPermissions      = []Permission{RBAC}

	List = append(
		append(
			append(
				append(
					append(DashboardPermissions, OrderPermissions...),
					MenuPermissions...),
				EmployeePermissions...),
			LocationPermissions...),
		RBACPermissions...,
	)
)
