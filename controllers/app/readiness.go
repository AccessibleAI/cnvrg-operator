package app

import (
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"k8s.io/apimachinery/pkg/types"
)

func (r *CnvrgAppReconciler) getControlPlaneReadinessStatus(app *mlopsv1.CnvrgApp) (bool, int, map[string]bool, error) {

	readyState := make(map[string]bool)

	// check webapp status
	if app.Spec.ControlPlane.WebApp.Enabled {
		name := types.NamespacedName{Name: app.Spec.ControlPlane.WebApp.SvcName, Namespace: app.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, nil, err
		}
		readyState["webApp"] = ready
	}

	// check sidekiq status
	if app.Spec.ControlPlane.Sidekiq.Enabled {
		name := types.NamespacedName{Name: "sidekiq", Namespace: app.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, nil, err
		}
		readyState["sidekiq"] = ready
	}

	// check searchkiq status
	if app.Spec.ControlPlane.Searchkiq.Enabled {
		name := types.NamespacedName{Name: "searchkiq", Namespace: app.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, nil, err
		}
		readyState["searchkiq"] = ready
	}

	// check systemkiq status
	if app.Spec.ControlPlane.Systemkiq.Enabled {
		name := types.NamespacedName{Name: "systemkiq", Namespace: app.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, nil, err
		}
		readyState["systemkiq"] = ready
	}

	// check postgres status
	if app.Spec.Dbs.Pg.Enabled {
		name := types.NamespacedName{Name: app.Spec.Dbs.Pg.SvcName, Namespace: app.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, nil, err
		}
		readyState["pg"] = ready
	}

	// check minio status
	if app.Spec.Dbs.Minio.Enabled {
		name := types.NamespacedName{Name: app.Spec.Dbs.Minio.SvcName, Namespace: app.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, nil, err
		}
		readyState["minio"] = ready
	}

	// check redis status
	if app.Spec.Dbs.Redis.Enabled {
		name := types.NamespacedName{Name: app.Spec.Dbs.Redis.SvcName, Namespace: app.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, nil, err
		}
		readyState["redis"] = ready
	}

	// check es status
	if app.Spec.Dbs.Es.Enabled {
		name := types.NamespacedName{Name: app.Spec.Dbs.Es.SvcName, Namespace: app.Namespace}
		ready, err := r.CheckStatefulSetReadiness(name)
		if err != nil {
			return false, 0, nil, err
		}
		readyState["es"] = ready
	}

	// check kibana status
	if app.Spec.Dbs.Es.Kibana.Enabled {
		name := types.NamespacedName{Name: app.Spec.Dbs.Es.Kibana.SvcName, Namespace: app.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, nil, err
		}
		readyState["kibana"] = ready
	}

	// check prometheus status
	if app.Spec.Dbs.Prom.Enabled {
		name := types.NamespacedName{Name: "prom", Namespace: app.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, nil, err
		}
		readyState["prometheus"] = ready
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

	return readyCount == len(readyState), percentageReady, readyState, nil
}
