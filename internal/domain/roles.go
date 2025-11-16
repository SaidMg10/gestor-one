package domain

// Role represents a role in the system.
const (
	RoleSuperAdmin = "superadmin"
	RoleAdmin      = "admin"
	RoleEmployee   = "employee"
	RoleAccountant = "accountant"
)

var AllowedRoles = map[string]bool{
	RoleSuperAdmin: true,
	RoleAdmin:      true,
	RoleEmployee:   true,
	RoleAccountant: true,
}

func IsValidRole(role string) bool {
	_, ok := AllowedRoles[role]
	return ok
}
