package docs

import (
	mlopsv1 "github.com/cnvrg-operator/api/v1"
)

var defaultEnabled = false

var controlPlane = mlopsv1.ControlPlane{
	WebApp: mlopsv1.WebApp{
		Enabled:  &defaultEnabled,
		Replicas: 0,
		Image:    "",
	},
	Sidekiq: mlopsv1.Sidekiq{
		Enabled: &defaultEnabled,
	},
	Searchkiq: mlopsv1.Searchkiq{
		Enabled: &defaultEnabled,
	},
	Systemkiq: mlopsv1.Systemkiq{
		Enabled: &defaultEnabled,
	},
	Hyper: mlopsv1.Hyper{
		Enabled: &defaultEnabled,
		Image:   "",
	},
	BaseConfig: mlopsv1.BaseConfig{
		AgentCustomTag:  "",
		Intercom:        "",
		CcpStorageClass: "",
	},
	Ldap: mlopsv1.Ldap{
		Enabled:       &defaultEnabled,
		Host:          "",
		Port:          "",
		Account:       "",
		Base:          "",
		AdminUser:     "",
		AdminPassword: "",
		Ssl:           "",
	},
	SMTP: mlopsv1.SMTP{
		Server:   "",
		Port:     0,
		Username: "",
		Password: "",
		Domain:   "",
	},
	ObjectStorage: mlopsv1.ObjectStorage{
		CnvrgStorageType:             "",
		CnvrgStorageBucket:           "",
		CnvrgStorageAccessKey:        "",
		CnvrgStorageSecretKey:        "",
		CnvrgStorageEndpoint:         "",
		CnvrgStorageAzureAccessKey:   "",
		CnvrgStorageAzureAccountName: "",
		CnvrgStorageAzureContainer:   "",
		CnvrgStorageRegion:           "",
		CnvrgStorageProject:          "",
		GcpStorageSecret:             "",
		GcpKeyfileMountPath:          "",
		GcpKeyfileName:               "",
	},
	Mpi: mlopsv1.Mpi{
		Enabled:              &defaultEnabled,
		Image:                "",
		KubectlDeliveryImage: "",
		ExtraArgs:            nil,
		Registry:             mlopsv1.Registry{},
	},
}

var registry = mlopsv1.Registry{
	URL:      "",
	User:     "",
	Password: "",
}

var dbs = mlopsv1.AppDbs{
	Pg: mlopsv1.Pg{
		Enabled:     nil,
		Image:       "",
		StorageSize: "",
		HugePages: mlopsv1.HugePages{
			Enabled: nil,
			Size:    "",
			Memory:  "",
		},
		NodeSelector: nil,
		CredsRef:     "",
	},
	Redis: mlopsv1.Redis{
		Enabled:      nil,
		Image:        "",
		StorageSize:  "",
		NodeSelector: nil,
		CredsRef:     "",
	},
	Minio: mlopsv1.Minio{
		Enabled:      nil,
		Image:        "",
		StorageSize:  "",
		NodeSelector: nil,
	},
	Es: mlopsv1.Es{
		Enabled:      nil,
		Image:        "",
		StorageSize:  "",
		NodeSelector: nil,
		CredsRef:     "",
	},
}

var networking = mlopsv1.CnvrgInfraNetworking{
	Ingress: mlopsv1.Ingress{
		Type:            "",
		Timeout:         "",
		RetriesAttempts: 0,
		PerTryTimeout:   "",
		IstioGwName:     "",
	},
	HTTPS: mlopsv1.HTTPS{
		Enabled:    nil,
		Cert:       "",
		Key:        "",
		CertSecret: "",
	},
	Istio: mlopsv1.Istio{
		Enabled:               &defaultEnabled,
		ExternalIP:            nil,
		IngressSvcAnnotations: nil,
		LBSourceRanges:        nil,
	},
}
