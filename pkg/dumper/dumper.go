package dumper

type Dumper interface {
	BuildState() error
	Dump() error
}
