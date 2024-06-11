package desired

import (
	"embed"
	"github.com/go-logr/logr"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type StateManager interface {
	Load() error
	Render() error
	Apply() error
}

type AssetsStateManager struct {
	C        client.Client
	s        *runtime.Scheme
	log      logr.Logger
	assets   []*AssetsGroup
	state    []*AssetsGroup
	rootPath string
	fs       embed.FS
	spec     v1.Object
	filter   *LoadFilter
}

func NewAssetsStateManager(
	spec v1.Object,
	c client.Client,
	s *runtime.Scheme,
	log logr.Logger,
	fs embed.FS,
	rootPath string,
	filter *LoadFilter) *AssetsStateManager {

	return &AssetsStateManager{
		C:        c,
		s:        s,
		log:      log,
		rootPath: rootPath,
		fs:       fs,
		spec:     spec,
		filter:   filter,
	}
}

func (m *AssetsStateManager) AddToState(ag *AssetsGroup) {
	if len(ag.Assets) == 0 {
		return
	}
	alreadyInState := false
	for _, s := range m.state {
		if reflect.DeepEqual(s, ag) {
			alreadyInState = true
		}
	}
	if !alreadyInState {
		m.state = append(m.state, ag)
	}
}

func (m *AssetsStateManager) AddToAssets(ag *AssetsGroup) {
	if len(ag.Assets) == 0 {
		return
	}
	alreadyInAssets := false
	for _, a := range m.assets {
		if reflect.DeepEqual(a, ag) {
			alreadyInAssets = true
		}
	}
	if !alreadyInAssets {
		m.assets = append(m.assets, ag)
	}
}

func (m *AssetsStateManager) RootPath() string {
	return m.rootPath
}

func (m *AssetsStateManager) Log() logr.Logger {
	return m.log
}

func (m *AssetsStateManager) Load() error {
	ag := NewAssetsGroup(m.fs, m.rootPath, m.log, m.filter)
	if err := ag.LoadAssets(); err != nil {
		return err
	}
	m.AddToAssets(ag)
	return nil
}

func (m *AssetsStateManager) Render() error {
	for _, ag := range m.assets {
		if err := ag.Render(m.spec); err != nil {
			return err
		}
		m.AddToState(ag)
	}
	return nil
}

func (m *AssetsStateManager) Apply() error {

	if err := m.Load(); err != nil {
		return err
	}

	if err := m.Render(); err != nil {
		return err
	}

	for _, ag := range m.state {
		if err := ag.Apply(m.spec, m.C, m.s, m.log); err != nil {
			return err
		}
	}
	return nil
}
