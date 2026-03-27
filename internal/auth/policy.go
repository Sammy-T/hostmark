package auth

var policy = map[Role]map[Resource]map[Permission]rule{
	RoleAdmin: {
		ResAcct: {
			PermCreate: allow,
			PermRead:   allow,
			PermUpdate: allow,
			PermDelete: allow,
		},
		ResFile: {
			PermCreate: allow,
			PermRead:   allow,
			PermUpdate: allow,
			PermDelete: allow,
		},
	},
	RoleUser: {
		ResAcct: {
			PermRead:   allowIfOwner,
			PermUpdate: allowIfOwner,
		},
		ResFile: {
			PermCreate: allow,
			PermRead:   allow,
			PermUpdate: allow,
			PermDelete: allow,
		},
	},
}
