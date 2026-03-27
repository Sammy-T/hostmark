package auth

type RuleArgs struct {
	User  string
	Owner string
}

type rule func(args RuleArgs) bool

var allow rule = func(args RuleArgs) bool { return true }

var allowIfOwner rule = func(args RuleArgs) bool { return args.User != "" && args.Owner != "" && args.User == args.Owner }
