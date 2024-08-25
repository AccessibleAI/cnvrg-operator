package app

import (
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"k8s.io/apimachinery/pkg/types"
)

const (
	webApp     = "webapp"
	sidekiq    = "sidekiq"
	searchkiq  = "searchkiq"
	systemkiq  = "systemkiq"
	pg         = "pg"
	minio      = "minio"
	redis      = "redis"
	es         = "es"
	kibana     = "kibana"
	prometheus = "prometheus"
	scheduler  = "scheduler"
)

type StackReadiness struct {
	readyState      map[string]bool
	isReady         bool
	percentageReady int
}

func (r *CnvrgAppReconciler) getControlPlaneReadinessStatus(app *mlopsv1.CnvrgApp) (*StackReadiness, error) {
	readyState := make(map[string]bool)

	// check webapp status
	if app.Spec.ControlPlane.WebApp.Enabled {
		name := types.NamespacedName{Name: app.Spec.ControlPlane.WebApp.SvcName, Namespace: app.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return nil, err
		}
		readyState[webApp] = ready
	}

	// check sidekiq status
	if app.Spec.ControlPlane.Sidekiq.Enabled {
		name := types.NamespacedName{Name: sidekiq, Namespace: app.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return nil, err
		}
		readyState[sidekiq] = ready
	}

	// check searchkiq status
	if app.Spec.ControlPlane.Searchkiq.Enabled {
		name := types.NamespacedName{Name: searchkiq, Namespace: app.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return nil, err
		}
		readyState[searchkiq] = ready
	}

	// check systemkiq status
	if app.Spec.ControlPlane.Systemkiq.Enabled {
		name := types.NamespacedName{Name: systemkiq, Namespace: app.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return nil, err
		}
		readyState[systemkiq] = ready
	}

	// check postgres status
	if app.Spec.Dbs.Pg.Enabled {
		name := types.NamespacedName{Name: app.Spec.Dbs.Pg.SvcName, Namespace: app.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return nil, err
		}
		readyState[pg] = ready
	}

	// check minio status
	if app.Spec.Dbs.Minio.Enabled {
		name := types.NamespacedName{Name: app.Spec.Dbs.Minio.SvcName, Namespace: app.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return nil, err
		}
		readyState[minio] = ready
	}

	// check redis status
	if app.Spec.Dbs.Redis.Enabled {
		name := types.NamespacedName{Name: app.Spec.Dbs.Redis.SvcName, Namespace: app.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return nil, err
		}
		readyState[redis] = ready
	}

	// check es status
	if app.Spec.Dbs.Es.Enabled {
		name := types.NamespacedName{Name: app.Spec.Dbs.Es.SvcName, Namespace: app.Namespace}
		ready, err := r.CheckStatefulSetReadiness(name)
		if err != nil {
			return nil, err
		}
		readyState[es] = ready
	}

	// check kibana status
	if app.Spec.Dbs.Es.Kibana.Enabled {
		name := types.NamespacedName{Name: app.Spec.Dbs.Es.Kibana.SvcName, Namespace: app.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return nil, err
		}
		readyState[kibana] = ready
	}

	// check prometheus status
	if app.Spec.Dbs.Prom.Enabled {
		name := types.NamespacedName{Name: app.Spec.Dbs.Prom.SvcName, Namespace: app.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return nil, err
		}
		readyState[prometheus] = ready
	}

	// check scheduler status
	if app.Spec.ControlPlane.CnvrgScheduler.Enabled {
		name := types.NamespacedName{Name: scheduler, Namespace: app.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return nil, err
		}
		readyState[scheduler] = ready
	}

	percentageReady := 0

	readyCount := 0

	for _, ready := range readyState {
		if ready {
			readyCount++
		}
	}

	if len(readyState) > 0 {
		percentageReady = readyCount * 100 / len(readyState)
	}
	if len(readyState) == 0 {
		percentageReady = 100
	}

	return &StackReadiness{
		isReady:         readyCount == len(readyState),
		percentageReady: percentageReady,
		readyState:      readyState,
	}, nil
}

func (s *StackReadiness) webAppReady() bool {
	if r, ok := s.readyState[webApp]; ok {
		return r
	}
	return false
}

func (s *StackReadiness) sidekiqReady() bool {
	if r, ok := s.readyState[sidekiq]; ok {
		return r
	}
	return false
}

func (s *StackReadiness) searchkiqReady() bool {
	if r, ok := s.readyState[searchkiq]; ok {
		return r
	}
	return false
}

func (s *StackReadiness) systemkiqReady() bool {
	if r, ok := s.readyState[systemkiq]; ok {
		return r
	}
	return false
}

func (s *StackReadiness) pgReady() bool {
	if r, ok := s.readyState[pg]; ok {
		return r
	}
	return false
}

func (s *StackReadiness) minioReady() bool {
	if r, ok := s.readyState[minio]; ok {
		return r
	}
	return false
}

func (s *StackReadiness) redisReady() bool {
	if r, ok := s.readyState[redis]; ok {
		return r
	}
	return false
}

func (s *StackReadiness) esReady() bool {
	if r, ok := s.readyState[es]; ok {
		return r
	}
	return false
}

func (s *StackReadiness) kibanaReady() bool {
	if r, ok := s.readyState[kibana]; ok {
		return r
	}
	return false
}

func (s *StackReadiness) prometheusReady() bool {
	if r, ok := s.readyState[prometheus]; ok {
		return r
	}
	return false
}

func (s *StackReadiness) schedulerReady() bool {
	if r, ok := s.readyState[scheduler]; ok {
		return r
	}
	return false
}
