package auth

type RuleArgs struct {
	User       string
	Owner      string
	Visibility string
}

type rule func(args RuleArgs) bool

var allow rule = func(args RuleArgs) bool { return true }

var allowIfMember rule = func(args RuleArgs) bool { return args.User != "" }

var allowIfOwner rule = func(args RuleArgs) bool { return allowIfMember(args) && args.Owner != "" && args.User == args.Owner }

var allowWithVisibility rule = func(args RuleArgs) bool {
	switch args.Visibility {
	case "public":
		return allow(args)
	case "protected":
		return allowIfMember(args)
	case "private":
		return allowIfOwner(args)
	default:
		return false
	}
}
