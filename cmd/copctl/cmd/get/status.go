package get

import (
	"context"
	"encoding/json"
	"fmt"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/cmd/copctl/utils"
	"github.com/AccessibleAI/cnvrg-shim/apis/metacloud/v1alpha1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"istio.io/istio/pkg/config/protocol"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"os/signal"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"syscall"
	"time"
)

type Status struct {
	NamespacedName      types.NamespacedName
	StatusConfigmapName string
	Interval            time.Duration
}

func init() {
	statusCmd.PersistentFlags().StringP("namespace", "", "", "current namespace")
	statusCmd.PersistentFlags().StringP("name", "", "", "name of the CnvrgApp CR")
	statusCmd.PersistentFlags().StringP("status-configmap", "", "service-instance-status", "the status cm name")
	statusCmd.PersistentFlags().IntP("interval", "i", 1, "status generation interval")

	viper.BindPFlag("namespace", statusCmd.PersistentFlags().Lookup("namespace"))
	viper.BindPFlag("name", statusCmd.PersistentFlags().Lookup("name"))
	viper.BindPFlag("status-configmap", statusCmd.PersistentFlags().Lookup("status-configmap"))
	viper.BindPFlag("interval", statusCmd.PersistentFlags().Lookup("interval"))

	Cmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:     "status",
	Aliases: []string{"s"},
	Short:   "Get cnvrg app status",
	Run: func(cmd *cobra.Command, args []string) {

		NewStatus(
			viper.GetString("namespace"),
			viper.GetString("name"),
			viper.GetString("status-configmap"),
			viper.GetInt("interval"),
		).run()

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		for {
			select {
			case s := <-sigCh:
				zap.S().Infof("signal: %s, shutting down", s)
				zap.S().Info("bye bye 👋")
				os.Exit(0)
			}
		}

	},
}

func NewStatus(ns, name, statusConfigmap string, interval int) *Status {
	if name == "" || ns == "" {
		zap.S().Fatal("name and namespace must be set")
	}
	return &Status{
		NamespacedName:      types.NamespacedName{Namespace: ns, Name: name},
		StatusConfigmapName: statusConfigmap,
		Interval:            time.Second * time.Duration(interval),
	}
}

func (s *Status) loadCnvrgApp(cap *mlopsv1.CnvrgApp) error {
	return utils.
		Kubecrudclient().
		Get(context.Background(),
			s.NamespacedName,
			cap,
			[]client.GetOption{}...,
		)

}

func (s *Status) run() {
	for {
		cnvrgAppInstance := &mlopsv1.CnvrgApp{}
		// fetch cnvrgapp instance
		if err := s.loadCnvrgApp(cnvrgAppInstance); err != nil {
			zap.S().Errorf("failed to fetch cnvrgapp instance, err: %s", err.Error())
			continue
		}

		siStatus, err := json.Marshal(s.generateStatus(cnvrgAppInstance))
		if err != nil {
			zap.S().Errorf("failed to marshal service instance, err: %s", err.Error())
			continue
		}

		if err := s.writeStatus(siStatus); err != nil {
			zap.S().Errorf("failed to create status configmap, err: %s", err.Error())
			continue
		}
		zap.S().Info("service instance configmap updated")
		time.Sleep(s.Interval)
	}
}

func (s *Status) writeStatus(payload []byte) error {

	// get status cm
	cm, err := utils.
		Clientset().
		CoreV1().
		ConfigMaps(s.NamespacedName.Namespace).
		Get(context.Background(), s.StatusConfigmapName, v1.GetOptions{})

	// if not found create it
	if errors.IsNotFound(err) {
		// construct configmap
		statusCm := &corev1.ConfigMap{
			ObjectMeta: v1.ObjectMeta{
				Namespace: s.NamespacedName.Namespace,
				Name:      s.StatusConfigmapName,
			},
			Data: map[string]string{"serviceInstanceStatus": string(payload)},
		}
		//create configmap
		_, err = utils.
			Clientset().
			CoreV1().
			ConfigMaps(s.NamespacedName.Namespace).
			Create(context.Background(), statusCm, v1.CreateOptions{})

		return err

	} else if err != nil {
		// error fetching configmap
		return err
	}

	// status configmap exists, update it
	cm.Data = map[string]string{"serviceInstanceStatus": string(payload)}
	_, err = utils.
		Clientset().
		CoreV1().
		ConfigMaps(s.NamespacedName.Namespace).
		Update(context.Background(), cm, v1.UpdateOptions{})

	return err

}

func (s *Status) generateStatus(cnvrgApp *mlopsv1.CnvrgApp) *v1alpha1.ServiceInstanceStatus {

	port := 80
	proto := protocol.HTTP
	status := v1alpha1.StatusHealthy

	if cnvrgApp.Spec.Networking.HTTPS.Enabled {
		proto = protocol.HTTPS
		port = 443
	}

	if cnvrgApp.Status.Status != "READY" {
		status = v1alpha1.StatusReconciling
	}

	return &v1alpha1.ServiceInstanceStatus{
		Status: status,
		Sins: []v1alpha1.Sins{
			{
				Name: cnvrgApp.Name,
				IngressEndpoints: []v1alpha1.IngressEndpoint{
					{
						Protocol: proto,
						Address:  []string{fmt.Sprintf("app.%s", cnvrgApp.Spec.ClusterDomain)},
						Port:     uint32(port),
					},
				},
			},
		},
	}

}
