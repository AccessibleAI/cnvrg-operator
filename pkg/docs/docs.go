package docs

import (
	"encoding/json"
	"fmt"
	mlopsv1 "github.com/cnvrg-operator/api/v1"
	"github.com/imdario/mergo"
	"github.com/jeremywohl/flatten"
	"os"
	"sort"
)

var fullAppSpec = mlopsv1.CnvrgAppSpec{
	ClusterDomain:    "-",
	NamespaceTenancy: "-",
	ControlPlane: mlopsv1.ControlPlane{
		WebApp: mlopsv1.WebApp{
			Replicas:                0,
			Enabled:                 "-",
			Image:                   "-",
			Port:                    0,
			CPU:                     "-",
			Memory:                  "-",
			SvcName:                 "-",
			NodePort:                0,
			PassengerMaxPoolSize:    0,
			InitialDelaySeconds:     0,
			ReadinessPeriodSeconds:  0,
			ReadinessTimeoutSeconds: 0,
			FailureThreshold:        0,
			OauthProxy:              mlopsv1.OauthProxyServiceConf{},
		},
		Sidekiq: mlopsv1.Sidekiq{
			Enabled:     "-",
			Split:       "-",
			CPU:         "-",
			Memory:      "-",
			Replicas:    0,
			KillTimeout: 0,
		},
		Searchkiq: mlopsv1.Searchkiq{
			Enabled:     "-",
			CPU:         "-",
			Memory:      "-",
			Replicas:    0,
			KillTimeout: 0,
		},
		Systemkiq: mlopsv1.Systemkiq{
			Enabled:     "-",
			CPU:         "-",
			Memory:      "-",
			Replicas:    0,
			KillTimeout: 0,
		},
		Hyper: mlopsv1.Hyper{
			Enabled:                 "-",
			Image:                   "-",
			Port:                    0,
			Replicas:                0,
			NodePort:                0,
			SvcName:                 "-",
			Token:                   "-",
			CPURequest:              "-",
			MemoryRequest:           "-",
			CPULimit:                "-",
			MemoryLimit:             "-",
			EnableReadinessProbe:    "-",
			ReadinessPeriodSeconds:  0,
			ReadinessTimeoutSeconds: 0,
		},
		Seeder: mlopsv1.Seeder{
			Image:           "-",
			SeedCmd:         "-",
			CreateBucketCmd: "-",
		},
		BaseConfig: mlopsv1.BaseConfig{
			JobsStorageClass:     "-",
			FeatureFlags:         nil,
			SentryURL:            "-",
			PassengerAppEnv:      "-",
			RailsEnv:             "-",
			RunJobsOnSelfCluster: "-",
			DefaultComputeConfig: "-",
			DefaultComputeName:   "-",
			UseStdout:            "-",
			ExtractTagsFromCmd:   "-",
			CheckJobExpiration:   "-",
			AgentCustomTag:       "-",
			Intercom:             "-",
			CnvrgJobUID:          "-",
			CcpStorageClass:      "-",
			HostpathNode:         "-",
		},
		Ldap: mlopsv1.Ldap{
			Enabled:       "-",
			Host:          "-",
			Port:          "-",
			Account:       "-",
			Base:          "-",
			AdminUser:     "-",
			AdminPassword: "-",
			Ssl:           "-",
		},
		Rbac: mlopsv1.Rbac{
			Role:               "-",
			ServiceAccountName: "-",
			RoleBindingName:    "-",
		},
		SMTP: mlopsv1.SMTP{
			Server:   "-",
			Port:     "-",
			Username: "-",
			Password: "-",
			Domain:   "-",
		},
		Tenancy: mlopsv1.Tenancy{
			Enabled:        "-",
			DedicatedNodes: "-",
			Key:            "-",
			Value:          "-",
		},
		ObjectStorage: mlopsv1.ObjectStorage{
			CnvrgStorageType:             "-",
			CnvrgStorageBucket:           "-",
			CnvrgStorageAccessKey:        "-",
			CnvrgStorageSecretKey:        "-",
			CnvrgStorageEndpoint:         "-",
			MinioSseMasterKey:            "-",
			CnvrgStorageAzureAccessKey:   "-",
			CnvrgStorageAzureAccountName: "-",
			CnvrgStorageAzureContainer:   "-",
			CnvrgStorageRegion:           "-",
			CnvrgStorageProject:          "-",
			GcpStorageSecret:             "-",
			GcpKeyfileMountPath:          "-",
			GcpKeyfileName:               "-",
			SecretKeyBase:                "-",
			StsIv:                        "-",
			StsKey:                       "-",
		},
		Mpi: mlopsv1.Mpi{
			Enabled:              "-",
			Image:                "-",
			KubectlDeliveryImage: "-",
			ExtraArgs:            nil,
			Registry: mlopsv1.Registry{
				Name:     "-",
				URL:      "-",
				User:     "-",
				Password: "-",
			},
		},
	},
	Registry: mlopsv1.Registry{
		Name:     "-",
		URL:      "-",
		User:     "-",
		Password: "-",
	},
	Dbs: mlopsv1.AppDbs{
		Pg: mlopsv1.Pg{
			Enabled:        "-",
			ServiceAccount: "-",
			SecretName:     "-",
			Image:          "-",
			Port:           0,
			StorageSize:    "-",
			SvcName:        "-",
			Dbname:         "-",
			Pass:           "-",
			User:           "-",
			RunAsUser:      0,
			FsGroup:        0,
			StorageClass:   "-",
			CPURequest:     "-",
			MemoryRequest:  "-",
			MaxConnections: 0,
			SharedBuffers:  "-",
			HugePages: mlopsv1.HugePages{
				Enabled: "-",
				Size:    "-",
				Memory:  "-",
			},
			Fixpg:        "-",
			NodeSelector: nil,
			Tolerations:  nil,
		},
		Redis: mlopsv1.Redis{
			Enabled:        "-",
			ServiceAccount: "-",
			Image:          "-",
			SvcName:        "-",
			Port:           0,
			Appendonly:     "-",
			StorageSize:    "-",
			StorageClass:   "-",
			Limits:         mlopsv1.Limits{},
			Requests:       mlopsv1.Requests{},
			NodeSelector:   nil,
			Tolerations:    nil,
		},
		Minio: mlopsv1.Minio{
			Enabled:        "-",
			ServiceAccount: "-",
			Replicas:       0,
			Image:          "-",
			Port:           0,
			StorageSize:    "-",
			SvcName:        "-",
			NodePort:       0,
			StorageClass:   "-",
			CPURequest:     "-",
			MemoryRequest:  "-",
			SharedStorage:  mlopsv1.SharedStorage{},
			NodeSelector:   nil,
			Tolerations:    nil,
		},
		Es: mlopsv1.Es{
			Enabled:        "-",
			ServiceAccount: "-",
			Image:          "-",
			Port:           0,
			StorageSize:    "-",
			SvcName:        "-",
			RunAsUser:      0,
			FsGroup:        0,
			NodePort:       0,
			StorageClass:   "-",
			CPURequest:     "-",
			MemoryRequest:  "-",
			CPULimit:       "-",
			MemoryLimit:    "-",
			JavaOpts:       "-",
			PatchEsNodes:   "-",
			NodeSelector:   nil,
			Tolerations:    nil,
		},
	},
	Networking: mlopsv1.CnvrgAppNetworking{
		Ingress: mlopsv1.Ingress{
			Enabled:         "-",
			IngressType:     "-",
			Timeout:         "-",
			RetriesAttempts: 0,
			PerTryTimeout:   "-",
			IstioGwName:     "-",
		},
		HTTPS: mlopsv1.HTTPS{
			Enabled:    "-",
			Cert:       "-",
			Key:        "-",
			CertSecret: "-",
		},
	},
	Logging: mlopsv1.CnvrgAppLogging{
		Enabled: "-",
		Elastalert: mlopsv1.Elastalert{
			Enabled:       "-",
			Image:         "-",
			Port:          0,
			NodePort:      0,
			ContainerPort: 0,
			StorageSize:   "-",
			SvcName:       "-",
			StorageClass:  "-",
			CPURequest:    "-",
			MemoryRequest: "-",
			CPULimit:      "-",
			MemoryLimit:   "-",
			RunAsUser:     0,
			FsGroup:       0,
		},
		Kibana: mlopsv1.Kibana{
			Enabled:        "-",
			ServiceAccount: "-",
			SvcName:        "-",
			Port:           0,
			Image:          "-",
			NodePort:       0,
			CPURequest:     "-",
			MemoryRequest:  "-",
			CPULimit:       "-",
			MemoryLimit:    "-",
			OauthProxy:     mlopsv1.OauthProxyServiceConf{},
		},
	},
	Monitoring: mlopsv1.CnvrgAppMonitoring{
		Enabled:            "-",
		UpstreamPrometheus: "-",
		Prometheus: mlopsv1.Prometheus{
			Enabled:       "-",
			Image:         "-",
			CPURequest:    "-",
			MemoryRequest: "-",
			SvcName:       "-",
			Port:          0,
			NodePort:      0,
			StorageSize:   "-",
			StorageClass:  "-",
		},
		Grafana: mlopsv1.Grafana{
			Enabled:    "-",
			Image:      "-",
			SvcName:    "-",
			Port:       0,
			NodePort:   0,
			OauthProxy: mlopsv1.OauthProxyServiceConf{},
		},
	},
	SSO: mlopsv1.SSO{
		Enabled:            "-",
		Image:              "-",
		RedisConnectionUrl: "-",
		AdminUser:          "-",
		Provider:           "-",
		EmailDomain:        "-",
		ClientID:           "-",
		ClientSecret:       "-",
		CookieSecret:       "-",
		AzureTenant:        "-",
		OidcIssuerURL:      "-",
	},
}

var fullInfraSpec = mlopsv1.CnvrgInfraSpec{
	ClusterDomain:     "-",
	InfraNamespace:    "-",
	InfraReconcilerCm: "-",
	Monitoring: mlopsv1.CnvrgInfraMonitoring{
		Enabled: "-",
		PrometheusOperator: mlopsv1.PrometheusOperator{
			Enabled: "-",
			Images: mlopsv1.Images{
				OperatorImage:                 "-",
				ConfigReloaderImage:           "-",
				PrometheusConfigReloaderImage: "-",
				KubeRbacProxyImage:            "-",
			},
		},
		Prometheus: mlopsv1.Prometheus{
			Enabled:       "-",
			Image:         "-",
			CPURequest:    "-",
			MemoryRequest: "-",
			SvcName:       "-",
			Port:          0,
			NodePort:      0,
			StorageSize:   "-",
			StorageClass:  "-",
		},
		KubeletServiceMonitor: "-",
		NodeExporter: mlopsv1.NodeExporter{
			Enabled: "-",
			Image:   "-",
		},
		KubeStateMetrics: mlopsv1.KubeStateMetrics{
			Enabled: "-",
			Image:   "-",
		},
		Grafana: mlopsv1.Grafana{
			Enabled:    "-",
			Image:      "-",
			SvcName:    "-",
			Port:       0,
			NodePort:   0,
			OauthProxy: mlopsv1.OauthProxyServiceConf{},
		},
		DcgmExporter: mlopsv1.DcgmExporter{
			Enabled: "-",
			Image:   "-",
		},
	},
	Networking: mlopsv1.CnvrgInfraNetworking{
		Ingress: mlopsv1.Ingress{
			Enabled:         "-",
			IngressType:     "-",
			Timeout:         "-",
			RetriesAttempts: 0,
			PerTryTimeout:   "-",
			IstioGwName:     "-",
		},
		HTTPS: mlopsv1.HTTPS{
			Enabled:    "-",
			Cert:       "-",
			Key:        "-",
			CertSecret: "-",
		},
		Istio: mlopsv1.Istio{
			Enabled:                  "-",
			OperatorImage:            "-",
			Hub:                      "-",
			Tag:                      "-",
			ProxyImage:               "-",
			MixerImage:               "-",
			PilotImage:               "-",
			ExternalIP:               "-",
			IngressSvcAnnotations:    "-",
			IngressSvcExtraPorts:     "-",
			LoadBalancerSourceRanges: "-",
		},
	},
	Logging: mlopsv1.CnvrgInfraLogging{
		Enabled:   "-",
		Fluentbit: mlopsv1.Fluentbit{Image: ""},
	},
	Registry: mlopsv1.Registry{
		Name:     "-",
		URL:      "-",
		User:     "-",
		Password: "-",
	},
	Storage: mlopsv1.Storage{
		Enabled: "-",
		Hostpath: mlopsv1.Hostpath{
			Enabled:          "-",
			Image:            "-",
			HostPath:         "-",
			StorageClassName: "-",
			NodeName:         "-",
			CPURequest:       "-",
			MemoryRequest:    "-",
			CPULimit:         "-",
			MemoryLimit:      "-",
			ReclaimPolicy:    "-",
			DefaultSc:        "-",
		},
		Nfs: mlopsv1.Nfs{
			Enabled:          "-",
			Image:            "-",
			Provisioner:      "-",
			StorageClassName: "-",
			Server:           "-",
			Path:             "-",
			CPURequest:       "-",
			MemoryRequest:    "-",
			CPULimit:         "-",
			MemoryLimit:      "-",
			ReclaimPolicy:    "-",
			DefaultSc:        "-",
		},
	},
	Dbs: mlopsv1.InfraDbs{Redis: mlopsv1.Redis{
		Enabled:        "-",
		ServiceAccount: "-",
		Image:          "-",
		SvcName:        "-",
		Port:           0,
		Appendonly:     "-",
		StorageSize:    "-",
		StorageClass:   "-",
		Limits:         mlopsv1.Limits{},
		Requests:       mlopsv1.Requests{},
		NodeSelector:   nil,
		Tolerations:    nil,
	}},
	SSO: mlopsv1.SSO{
		Enabled:            "-",
		Image:              "-",
		RedisConnectionUrl: "-",
		AdminUser:          "-",
		Provider:           "-",
		EmailDomain:        "-",
		ClientID:           "-",
		ClientSecret:       "-",
		CookieSecret:       "-",
		AzureTenant:        "-",
		OidcIssuerURL:      "-",
	},
	Gpu: mlopsv1.Gpu{NvidiaDp: mlopsv1.NvidiaDp{
		Enabled: "-",
		Image:   "-",
	}},
}

func GenerateDocs() {

	// app params
	app := mlopsv1.DefaultCnvrgAppSpec()
	if err := mergo.Merge(&app, &fullAppSpec); err != nil {
		fmt.Fprintf(os.Stderr, "failed generate docs: %v", err)
		os.Exit(1)
	}
	b, _ := json.Marshal(app)
	flatAppParams, _ := flatten.FlattenString(string(b), "", flatten.DotStyle)
	appParams := make(map[string]interface{})
	_ = json.Unmarshal([]byte(flatAppParams), &appParams)

	// infra params
	infra := mlopsv1.DefaultCnvrgInfraSpec()
	if err := mergo.Merge(&infra, &fullInfraSpec); err != nil {
		fmt.Fprintf(os.Stderr, "failed generate docs: %v", err)
		os.Exit(1)
	}
	b, _ = json.Marshal(infra)
	flatInfraParams, _ := flatten.FlattenString(string(b), "", flatten.DotStyle)
	infraParams := make(map[string]interface{})
	_ = json.Unmarshal([]byte(flatInfraParams), &infraParams)

	finalParams := make(map[string]interface{})
	skipKeys := []string{
		"controlPlane.baseConfig.sentryUrl",
		"controlPlane.objectStorage.stsIv",
		"controlPlane.objectStorage.stsKey",
		"controlPlane.objectStorage.secretKeyBase",
		"controlPlane.objectStorage.minioSseMasterKey",
	}
	for key, value := range appParams {
		skipKey := false
		for _, item := range skipKeys {
			if item == key {
				skipKey = true
			}
		}
		if !skipKey {
			finalParams[key] = value
		}
	}

	for key, value := range infraParams {
		finalParams[key] = value
	}
	keys := make([]string, 0, len(finalParams))
	for k := range finalParams {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Println(fmt.Sprintf("|`%v`|%v", k, finalParams[k]))
	}
}
