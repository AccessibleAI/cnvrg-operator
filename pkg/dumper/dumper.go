package dumper

type Dumper interface {
	BuildState() error
	Dump(preserveTmplDirs bool) error
}
