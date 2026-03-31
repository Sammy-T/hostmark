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
		ResNote: {
			PermCreate: allowIfOwner,
			PermRead:   allowWithVisibility,
			PermUpdate: allowIfOwner,
			PermDelete: allowIfOwner,
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
		ResNote: {
			PermCreate: allowIfOwner,
			PermRead:   allowWithVisibility,
			PermUpdate: allowIfOwner,
			PermDelete: allowIfOwner,
		},
	},
}
