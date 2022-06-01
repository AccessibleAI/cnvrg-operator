package dumper

import "github.com/AccessibleAI/cnvrg-operator/pkg/desired"

type Dumper interface {
	BuildState() []*desired.State
	Dump() error
	GetCliParams() []*Param
}

type Param struct {
	Name      string
	Shorthand string
	Value     interface{}
	Usage     string
	Required  bool
}
