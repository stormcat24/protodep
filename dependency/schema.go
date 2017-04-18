package dependency

type ProtoDep struct {
	Dependencies []ProtoDepDependency
}

type ProtoDepDependency struct {
	Name     string
	Revision string
	Branch   string
}
