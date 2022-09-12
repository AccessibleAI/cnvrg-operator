package thirdparty

import (
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func applyCtpDefaults(ctp *mlopsv1.CnvrgThirdParty, desiredCtp *mlopsv1.CnvrgThirdPartySpec, c client.Client) error {
	if ctp.Spec.Metagpu.Enabled && len(ctp.Spec.Metagpu.NodeSelector) == 0 {
		desiredCtp.Metagpu.NodeSelector = map[string]string{"accelerator": "nvidia"}
	}
	return nil
}
