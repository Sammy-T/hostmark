package auth

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleEditor Role = "editor"
	RoleUser   Role = "user"
	RoleNone   Role = "none"
)

type Resource string

const (
	ResFile Resource = "file"
	ResNote Resource = "note"
	ResAcct Resource = "acct"
)

type Permission string

const (
	PermCreate Permission = "create"
	PermRead   Permission = "read"
	PermUpdate Permission = "update"
	PermDelete Permission = "delete"
)

// Access processes the corresponding policy rule
// and returns whether access is granted.
func Access(role Role, res Resource, perm Permission, args RuleArgs) bool {
	if role == "" {
		role = RoleNone
	}

	rule, ok := policy[role][res][perm]

	return ok && rule(args)
}
