package auth

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type Resource string

const (
	ResFile Resource = "file"
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
	rule, ok := policy[role][res][perm]

	return ok && rule(args)
}
