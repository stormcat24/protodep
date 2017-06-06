package dependency

import (
	"strings"
)

type ProtoDep struct {
	ProtoOutdir  string               `toml:"proto_outdir"`
	Dependencies []ProtoDepDependency `toml:"dependencies"`
}

type ProtoDepDependency struct {
	Name     string
	Revision string
	Branch   string
}

func (d *ProtoDepDependency) Repository() string {
	tokens := strings.Split(d.Name, "/")
	if len(tokens) > 3 {
		return strings.Join(tokens[0:3], "/")
	} else {
		return d.Name
	}
}

func (d *ProtoDepDependency) Directory() string {
	r := d.Repository()

	if d.Name == r {
		return "."
	} else {
		return "." + strings.Replace(d.Name, r, "", 1)
	}
}
