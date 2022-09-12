package v1

var controlPlaneDefault = ControlPlane{
	Image: "core:3.6.99",

	WebApp: WebApp{
		Enabled:  false,
		Replicas: 1,
		Port:     8080,
		Requests: Requests{
			Cpu:    "500m",
			Memory: "4Gi",
		},
		Limits: Limits{
			Cpu:    "4",
			Memory: "8Gi",
		},
		SvcName:                 "app",
		NodePort:                30080,
		PassengerMaxPoolSize:    50,
		InitialDelaySeconds:     10,
		ReadinessPeriodSeconds:  25,
		ReadinessTimeoutSeconds: 20,
		FailureThreshold:        5,
		OauthProxy: OauthProxyServiceConf{
			SkipAuthRegex: []string{
				`\/assets`,
				`\/healthz`,
				`\/public`,
				`\/pack`,
				`\/vscode.tar.gz`,
				`\/jupyter.vsix`,
				`\/gitlens.vsix`,
				`\/ms-python-release.vsix`,
				`\/webhooks`,
				`\/api/v2/metrics`,
				`\/api/v1/events/endpoint_rule_alert`,
			},
			TokenValidationRegex: []string{
				`^\/api`,
			},
		},
		Hpa: Hpa{
			Enabled:     false,
			Utilization: 85,
			MaxReplicas: 5,
		},
	},

	Sidekiq: Sidekiq{
		Enabled: false,
		Split:   false,
		Requests: Requests{
			Cpu:    "200m",
			Memory: "3750Mi",
		},
		Limits: Limits{
			Cpu:    "2",
			Memory: "8Gi",
		},
		Replicas: 2,
		Hpa: Hpa{
			Enabled:     false,
			Utilization: 85,
			MaxReplicas: 5,
		},
	},

	Searchkiq: Searchkiq{
		Enabled: false,
		Requests: Requests{
			Cpu:    "200m",
			Memory: "1Gi",
		},
		Limits: Limits{
			Cpu:    "2",
			Memory: "8Gi",
		},
		Replicas: 1,
		Hpa: Hpa{
			Enabled:     false,
			Utilization: 85,
			MaxReplicas: 5,
		},
	},

	Systemkiq: Systemkiq{
		Enabled: false,
		Requests: Requests{
			Cpu:    "300m",
			Memory: "2Gi",
		},
		Limits: Limits{
			Cpu:    "2",
			Memory: "8Gi",
		},
		Replicas: 1,
		Hpa: Hpa{
			Enabled:     false,
			Utilization: 85,
			MaxReplicas: 5,
		},
	},

	Hyper: Hyper{
		Enabled:  false,
		Image:    "hyper-server:latest",
		Port:     5050,
		Replicas: 1,
		NodePort: 30050,
		SvcName:  "hyper",
		Token:    "token",
		Requests: Requests{
			Cpu:    "100m",
			Memory: "200Mi",
		},
		Limits: Limits{
			Cpu:    "2",
			Memory: "4Gi",
		},
		ReadinessPeriodSeconds:  100,
		ReadinessTimeoutSeconds: 60,
	},

	CnvrgScheduler: CnvrgScheduler{
		Enabled: false,
		Requests: Requests{
			Cpu:    "200m",
			Memory: "1000Mi",
		},
		Limits: Limits{
			Cpu:    "2",
			Memory: "4Gi",
		},
		Replicas: 1,
	},

	CnvrgClusterProvisionerOperator: CnvrgClusterProvisionerOperator{
		Enabled: false,
		Requests: Requests{
			Cpu:    "200m",
			Memory: "1Gi",
		},
		Limits: Limits{
			Cpu:    "2",
			Memory: "4Gi",
		},
		Image:       "cnvrg/ccp-operator:v1",
		AwsCredsRef: "",
	},

	CnvrgRouter: CnvrgRouter{
		Enabled:  false,
		Image:    "nginx:1.21.0",
		SvcName:  "cnvrg-router",
		NodePort: 30081,
	},

	Mpi: Mpi{
		Enabled:              false,
		Image:                "mpioperator/mpi-operator:v0.2.3",
		KubectlDeliveryImage: "mpioperator/kubectl-delivery:v0.2.3",
		ExtraArgs:            nil,
		Requests: Requests{
			Cpu:    "100m",
			Memory: "100Mi",
		},
		Limits: Limits{
			Cpu:    "1000m",
			Memory: "1Gi",
		},
		Registry: Registry{
			Name:     "mpi-private-registry",
			URL:      "docker.io",
			User:     "",
			Password: "",
		},
	},

	BaseConfig: BaseConfig{
		JobsStorageClass:   "",
		FeatureFlags:       nil,
		SentryURL:          "",
		AgentCustomTag:     "latest",
		Intercom:           "true",
		CnvrgJobUID:        "0", // by default cnvrg job is running as root
		CnvrgJobRbacStrict: false,
		CnvrgPrivilegedJob: true,
		MetagpuEnabled:     false,
	},

	ObjectStorage: ObjectStorage{

		Type:             MinioObjectStorageType,
		Bucket:           "cnvrg-storage",
		AccessKey:        "",
		SecretKey:        "",
		Endpoint:         "",
		AzureAccountName: "",
		AzureContainer:   "",
		Region:           "eastus",
		GcpProject:       "",
		GcpSecretRef:     "gcp-storage-secret",
	},

	Ldap: Ldap{
		Enabled:       false,
		Host:          "",
		Port:          "",
		Account:       "userPrincipalName",
		Base:          "", // dc=my-domain,dc=local
		AdminUser:     "",
		AdminPassword: "",
		Ssl:           "", // true/false
	},

	SMTP: SMTP{
		Server:            "",
		Port:              587,
		Username:          "",
		Password:          "",
		Domain:            "",
		OpensslVerifyMode: "",
		Sender:            "info@cnvrg.io",
	},
	Nomex: defaultNomex,
}

var minioDefaults = Minio{
	Enabled:        false,
	ServiceAccount: "minio",
	Replicas:       1,
	Image:          "minio:RELEASE.2021-05-22T02-34-39Z",
	Port:           9000,
	StorageSize:    "100Gi",
	SvcName:        "minio",
	NodePort:       30090,
	StorageClass:   "",
	Requests: Requests{
		Cpu:    "200m",
		Memory: "2Gi",
	},
	Limits: Limits{
		Cpu:    "8",
		Memory: "20Gi",
	},
	PvcName: "minio-storage",
	SharedStorage: SharedStorage{
		Enabled: false,
		ConsistentHash: ConsistentHash{
			Key:   "httpQueryParameterName",
			Value: "uploadId",
		},
	},
}

var pgDefault = Pg{
	Enabled:        false,
	ServiceAccount: "pg",
	Image:          "postgresql-12-centos7:latest",
	Port:           5432,
	StorageSize:    "80Gi",
	SvcName:        "postgres",
	StorageClass:   "",
	Requests: Requests{
		Cpu:    "1",
		Memory: "4Gi",
	},
	Limits: Limits{
		Cpu:    "12",
		Memory: "32Gi",
	},
	MaxConnections:     500,
	SharedBuffers:      "1024MB", // for the shared_buffers we use 1/4 of given memory
	EffectiveCacheSize: "2048MB", // for the effective_cache_size we set the value to 1/2 of the given memory
	NodeSelector:       nil,
	PvcName:            "pg-storage",
	HugePages: HugePages{
		Enabled: false,
		Size:    "2Mi", // 2Mi/1Gi https://kubernetes.io/docs/tasks/manage-hugepages/scheduling-hugepages/ ,  https://wiki.debian.org/Hugepages#Huge_pages_sizes
		Memory:  "",
	},
	CredsRef: "pg-creds",
}

var redisDefault = Redis{
	Enabled:        false,
	ServiceAccount: "redis",
	Image:          "cnvrg-redis:v3.0.5.c2",
	SvcName:        "redis",
	Port:           6379,
	StorageSize:    "10Gi",
	StorageClass:   "",
	NodeSelector:   nil,
	CredsRef:       "redis-creds",
	PvcName:        "redis-storage",
	Limits: Limits{
		Cpu:    "1000m",
		Memory: "2Gi",
	},
	Requests: Requests{
		Cpu:    "100m",
		Memory: "200Mi",
	},
}

var esDefault = Es{
	Enabled:        false,
	ServiceAccount: "es",
	Image:          "cnvrg-es:v7.8.1.a1-dynamic-indices",
	Port:           9200,
	StorageSize:    "80Gi",
	SvcName:        "elasticsearch",
	NodePort:       32200,
	StorageClass:   "",
	Requests: Requests{
		Cpu:    "500m",
		Memory: "4Gi",
	},
	Limits: Limits{
		Cpu:    "4",
		Memory: "8Gi",
	},
	JavaOpts:     "",
	PatchEsNodes: false,
	CredsRef:     "es-creds",
	PvcName:      "es-storage",
	CleanupPolicy: CleanupPolicy{
		All:       "3d",
		App:       "30d",
		Jobs:      "14d",
		Endpoints: "1825d",
	},
	Kibana: Kibana{
		Enabled:        false,
		ServiceAccount: "kibana",
		SvcName:        "kibana",
		Port:           8080,
		Image:          "kibana-oss:7.8.1",
		NodePort:       30601,
		Requests: Requests{
			Cpu:    "100m",
			Memory: "200Mi",
		},
		Limits: Limits{
			Cpu:    "1000m",
			Memory: "2Gi",
		},
		OauthProxy: OauthProxyServiceConf{
			SkipAuthRegex:        nil,
			TokenValidationRegex: nil,
		},
		CredsRef: "kibana-creds",
	},
	Elastalert: Elastalert{
		Enabled:        false,
		Image:          "elastalert:3.0.0-beta.1",
		CredsRef:       "elastalert-creds",
		AuthProxyImage: "nginx:1.20",
		Port:           80,
		NodePort:       32030,
		StorageSize:    "30Gi",
		SvcName:        "elastalert",
		StorageClass:   "",
		Requests: Requests{
			Cpu:    "100m",
			Memory: "200Mi",
		},
		Limits: Limits{
			Cpu:    "400m",
			Memory: "800Mi",
		},
		NodeSelector: nil,
		PvcName:      "elastalert-storage",
	},
}

var appDbsDefaults = Dbs{
	Pg:    pgDefault,
	Redis: redisDefault,
	Minio: minioDefaults,
	Es:    esDefault,
	Cvat:  cvatDefault,
	Prom:  promDefaults,
}

var cvatDefault = Cvat{
	Enabled: false,
	Pg:      cvatPgDefault,
	Redis:   cvatRedisDefault,
}

var promDefaults = Prom{
	Enabled:  false,
	CredsRef: "prom-creds",
}

var defaultNomex = Nomex{
	Enabled: true,
	Image:   "docker.io/cnvrg/nomex:v1.0.0",
}

var cvatPgDefault = Pg{
	Enabled:        false,
	ServiceAccount: "cvat-pg",
	Image:          "postgresql-12-centos7:latest",
	Port:           5432,
	StorageSize:    "100Gi",
	SvcName:        "cvat-postgres",
	StorageClass:   "",
	Requests: Requests{
		Cpu:    "1",
		Memory: "2Gi",
	},
	Limits: Limits{
		Cpu:    "2",
		Memory: "4Gi",
	},
	MaxConnections:     500,
	SharedBuffers:      "1024MB", // for the shared_buffers we use 1/4 of given memory
	EffectiveCacheSize: "2048MB", // for the effective_cache_size we set the value to 1/2 of the given memory
	NodeSelector:       nil,
	PvcName:            "cvat-pg-storage",
	HugePages: HugePages{
		Enabled: false,
		Size:    "2Mi", // 2Mi/1Gi https://kubernetes.io/docs/tasks/manage-hugepages/scheduling-hugepages/ ,  https://wiki.debian.org/Hugepages#Huge_pages_sizes
		Memory:  "",
	},
	CredsRef: "cvat-pg-secret",
}

var cvatRedisDefault = Redis{
	Enabled:        false,
	ServiceAccount: "cvat-redis",
	Image:          "redis:4.0.5-alpine",
	SvcName:        "cvat-redis",
	Port:           6379,
	StorageSize:    "10Gi",
	StorageClass:   "",
	NodeSelector:   nil,
	CredsRef:       "cvat-redis-secret",
	PvcName:        "cvat-redis-storage",
	Limits: Limits{
		Cpu:    "1000m",
		Memory: "2Gi",
	},
	Requests: Requests{
		Cpu:    "250m",
		Memory: "500Mi",
	},
}

var networkingDefault = Networking{
	Ingress: Ingress{
		Type:            IstioIngress,
		Timeout:         "18000s",
		RetriesAttempts: 5,
		PerTryTimeout:   "3600s",
		IstioGwEnabled:  false,
		IstioGwName:     "",
	},
	HTTPS: HTTPS{Enabled: false},
	Proxy: Proxy{Enabled: false, ConfigRef: "cp-proxy"},
}

var ssoDefault = SSO{
	Enabled:                          false,
	Image:                            "saas-oauth2-proxy:latest",
	AdminUser:                        "",
	Provider:                         "",
	EmailDomain:                      nil,
	ClientID:                         "",
	ClientSecret:                     "",
	CookieSecret:                     "",
	AzureTenant:                      "", // if IDP is Azure AD
	OidcIssuerURL:                    "", // if IDP oidc
	RealmName:                        "",
	ServiceUrl:                       "",
	InsecureOidcAllowUnverifiedEmail: false,
}
