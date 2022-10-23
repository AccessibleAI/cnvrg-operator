var _ = require('lodash');
const YAML = require('yaml')


const spec = {
    "clusterDomain": "",
    "clusterInternalDomain": "cluster.local",
    "imageHub": "docker.io/cnvrg",
    "labels": {
        "owner": "cnvrg-control-plane"
    },
    "annotations": null,
    "controlPlane": {
        "image": "core:3.6.99",
        "webapp": {
            "replicas": 1,
            "enabled": false,
            "port": 8080,
            "requests": {
                "cpu": "500m",
                "memory": "4Gi"
            },
            "limits": {
                "cpu": "4",
                "memory": "8Gi"
            },
            "svcName": "app",
            "nodePort": 30080,
            "passengerMaxPoolSize": 50,
            "initialDelaySeconds": 10,
            "readinessPeriodSeconds": 25,
            "readinessTimeoutSeconds": 20,
            "failureThreshold": 5,
            "hpa": {
                "enabled": false,
                "utilization": 85,
                "maxReplicas": 5
            }
        },
        "sidekiq": {
            "enabled": false,
            "split": false,
            "requests": {
                "cpu": "200m",
                "memory": "3750Mi"
            },
            "limits": {
                "cpu": "2",
                "memory": "8Gi"
            },
            "replicas": 2,
            "hpa": {
                "enabled": false,
                "utilization": 85,
                "maxReplicas": 5
            }
        },
        "searchkiq": {
            "enabled": false,
            "requests": {
                "cpu": "200m",
                "memory": "1Gi"
            },
            "limits": {
                "cpu": "2",
                "memory": "8Gi"
            },
            "replicas": 1,
            "hpa": {
                "enabled": false,
                "utilization": 85,
                "maxReplicas": 5
            }
        },
        "systemkiq": {
            "enabled": false,
            "requests": {
                "cpu": "300m",
                "memory": "2Gi"
            },
            "limits": {
                "cpu": "2",
                "memory": "8Gi"
            },
            "replicas": 1,
            "hpa": {
                "enabled": false,
                "utilization": 85,
                "maxReplicas": 5
            }
        },
        "hyper": {
            "enabled": false,
            "image": "hyper-server:latest",
            "port": 5050,
            "replicas": 1,
            "nodePort": 30050,
            "svcName": "hyper",
            "token": "token",
            "requests": {
                "cpu": "100m",
                "memory": "200Mi"
            },
            "limits": {
                "cpu": "2",
                "memory": "4Gi"
            },
            "cpuLimit": "",
            "memoryLimit": "",
            "readinessPeriodSeconds": 100,
            "readinessTimeoutSeconds": 60
        },
        "cnvrgScheduler": {
            "enabled": false,
            "requests": {
                "cpu": "200m",
                "memory": "1000Mi"
            },
            "limits": {
                "cpu": "2",
                "memory": "4Gi"
            },
            "replicas": 1
        },
        "cnvrgClusterProvisionerOperator": {
            "enabled": false,
            "requests": {
                "cpu": "200m",
                "memory": "1Gi"
            },
            "limits": {
                "cpu": "2",
                "memory": "4Gi"
            },
            "image": "cnvrg/ccp-operator:v1",
            "awsCredsRef": ""
        },
        "cnvrgRouter": {
            "enabled": false,
            "image": "nginx:1.21.0",
            "svcName": "cnvrg-router",
            "nodePort": 30081
        },
        "baseConfig": {
            "jobsStorageClass": "",
            "featureFlags": null,
            "sentryUrl": "",
            "runJobsOnSelfCluster": "",
            "agentCustomTag": "latest",
            "intercom": "true",
            "cnvrgJobUid": "0",
            "cnvrgJobRbacStrict": false,
            "cnvrgPrivilegedJob": true,
            "metagpuEnabled": false
        },
        "ldap": {
            "enabled": false,
            "host": "",
            "port": "",
            "account": "userPrincipalName",
            "base": "",
            "adminUser": "",
            "adminPassword": "",
            "ssl": ""
        },
        "smtp": {
            "server": "",
            "port": 587,
            "username": "",
            "password": "",
            "domain": "",
            "opensslVerifyMode": "",
            "sender": "info@cnvrg.io"
        },
        "objectStorage": {
            "type": "minio",
            "bucket": "cnvrg-storage",
            "region": "eastus",
            "accessKey": "",
            "secretKey": "",
            "endpoint": "",
            "azureAccountName": "",
            "azureContainer": "",
            "gcpProject": "",
            "gcpSecretRef": "gcp-storage-secret"
        },
        "mpi": {
            "enabled": false,
            "image": "mpioperator/mpi-operator:v0.2.3",
            "kubectlDeliveryImage": "mpioperator/kubectl-delivery:v0.2.3",
            "extraArgs": null,
            "registry": {
                "name": "mpi-private-registry",
                "url": "docker.io",
                "user": "",
                "password": ""
            },
            "requests": {
                "cpu": "100m",
                "memory": "100Mi"
            },
            "limits": {
                "cpu": "1000m",
                "memory": "1Gi"
            }
        },
        "nomex": {
            "enabled": true,
            "image": "docker.io/cnvrg/nomex:v1.0.0"
        }
    },
    "registry": {
        "name": "cnvrg-app-registry",
        "url": "docker.io",
        "user": "",
        "password": ""
    },
    "dbs": {
        "pg": {
            "enabled": false,
            "serviceAccount": "pg",
            "image": "postgresql-12-centos7:latest",
            "port": 5432,
            "storageSize": "80Gi",
            "svcName": "postgres",
            "storageClass": "",
            "requests": {
                "cpu": "1",
                "memory": "4Gi"
            },
            "limits": {
                "cpu": "12",
                "memory": "32Gi"
            },
            "maxConnections": 500,
            "sharedBuffers": "1024MB",
            "effectiveCacheSize": "2048MB",
            "hugePages": {
                "enabled": false,
                "size": "2Mi",
                "memory": ""
            },
            "nodeSelector": null,
            "credsRef": "pg-creds",
            "pvcName": "pg-storage"
        },
        "redis": {
            "enabled": false,
            "serviceAccount": "redis",
            "image": "cnvrg-redis:v3.0.5.c2",
            "svcName": "redis",
            "port": 6379,
            "storageSize": "10Gi",
            "storageClass": "",
            "requests": {
                "cpu": "100m",
                "memory": "200Mi"
            },
            "limits": {
                "cpu": "1000m",
                "memory": "2Gi"
            },
            "nodeSelector": null,
            "credsRef": "redis-creds",
            "pvcName": "redis-storage"
        },
        "minio": {
            "enabled": false,
            "serviceAccount": "minio",
            "replicas": 1,
            "image": "minio:RELEASE.2021-05-22T02-34-39Z",
            "port": 9000,
            "storageSize": "100Gi",
            "svcName": "minio",
            "nodePort": 30090,
            "storageClass": "",
            "requests": {
                "cpu": "200m",
                "memory": "2Gi"
            },
            "limits": {
                "cpu": "8",
                "memory": "20Gi"
            },
            "nodeSelector": null,
            "pvcName": "minio-storage"
        },
        "es": {
            "enabled": false,
            "serviceAccount": "es",
            "image": "cnvrg-es:v7.8.1.a1-dynamic-indices",
            "port": 9200,
            "storageSize": "80Gi",
            "svcName": "elasticsearch",
            "nodePort": 32200,
            "storageClass": "",
            "requests": {
                "cpu": "500m",
                "memory": "4Gi"
            },
            "limits": {
                "cpu": "4",
                "memory": "8Gi"
            },
            "javaOpts": "",
            "patchEsNodes": false,
            "nodeSelector": null,
            "credsRef": "es-creds",
            "pvcName": "es-storage",
            "cleanupPolicy": {
                "all": "3d",
                "app": "30d",
                "jobs": "14d",
                "endpoints": "1825d"
            },
            "kibana": {
                "enabled": false,
                "serviceAccount": "kibana",
                "svcName": "kibana",
                "port": 8080,
                "image": "kibana-oss:7.8.1",
                "nodePort": 30601,
                "requests": {
                    "cpu": "100m",
                    "memory": "200Mi"
                },
                "limits": {
                    "cpu": "1000m",
                    "memory": "2Gi"
                },
                "credsRef": "kibana-creds"
            },
            "elastalert": {
                "enabled": false,
                "image": "elastalert:3.0.0-beta.1",
                "authProxyImage": "nginx:1.20",
                "credsRef": "elastalert-creds",
                "port": 80,
                "nodePort": 32030,
                "storageSize": "30Gi",
                "svcName": "elastalert",
                "storageClass": "",
                "requests": {
                    "cpu": "100m",
                    "memory": "200Mi"
                },
                "limits": {
                    "cpu": "400m",
                    "memory": "800Mi"
                },
                "nodeSelector": null,
                "pvcName": "elastalert-storage"
            }
        },
        "cvat": {
            "enabled": false,
            "pg": {
                "enabled": false,
                "serviceAccount": "cvat-pg",
                "image": "postgresql-12-centos7:latest",
                "port": 5432,
                "storageSize": "100Gi",
                "svcName": "cvat-postgres",
                "storageClass": "",
                "requests": {
                    "cpu": "1",
                    "memory": "2Gi"
                },
                "limits": {
                    "cpu": "2",
                    "memory": "4Gi"
                },
                "maxConnections": 500,
                "sharedBuffers": "1024MB",
                "effectiveCacheSize": "2048MB",
                "hugePages": {
                    "enabled": false,
                    "size": "2Mi",
                    "memory": ""
                },
                "nodeSelector": null,
                "credsRef": "cvat-pg-secret",
                "pvcName": "cvat-pg-storage"
            },
            "redis": {
                "enabled": false,
                "serviceAccount": "cvat-redis",
                "image": "redis:4.0.5-alpine",
                "svcName": "cvat-redis",
                "port": 6379,
                "storageSize": "10Gi",
                "storageClass": "",
                "requests": {
                    "cpu": "250m",
                    "memory": "500Mi"
                },
                "limits": {
                    "cpu": "1000m",
                    "memory": "2Gi"
                },
                "nodeSelector": null,
                "credsRef": "cvat-redis-secret",
                "pvcName": "cvat-redis-storage"
            }
        },
        "prom": {
            "enabled": false,
            "credsRef": "prom-creds",
            "extraScrapeConfigs": null,
            "image": "prom/prometheus:v2.37.1",
            "grafana": {
                "enabled": false,
                "image": "grafana-oss:9.1.6",
                "svcName": "grafana",
                "port": 8080,
                "nodePort": 30012,
                "credsRef": "grafana-creds"
            }
        }
    },
    "networking": {
        "ingress": {
            "type": "istio",
            "timeout": "18000s",
            "retriesAttempts": 5,
            "perTryTimeout": "3600s",
            "istioGwEnabled": false,
            "istioGwName": ""
        },
        "https": {
            "enabled": false,
            "cert": "",
            "key": "",
            "certSecret": ""
        },
        "proxy": {
            "enabled": false,
            "configRef": "cp-proxy",
            "httpProxy": null,
            "httpsProxy": null,
            "noProxy": null
        }
    },
    "sso": {
        "enabled": false,
        "groups": null,
        "pki": {
            "enabled": false,
            "rootCaSecret": "sso-idp-root-ca",
            "privateKeySecret": "sso-idp-private-key",
            "publicKeySecret": "sso-idp-pki-public-key"
        },
        "jwks": {
            "enabled": false,
            "name": "cnvrg-jwks",
            "image": "cnvrg/jwks:latest",
            "cache": {
                "enabled": true,
                "image": "docker.io/redis"
            }
        },
        "central": {
            "enabled": false,
            "publicUrl": "",
            "cnvrgProxyImage": "docker.io/cnvrg/proxy:v1.0.1",
            "oauthProxyImage": "cnvrg/oauth2-proxy:v7.3.x.ssov3.p2-01",
            "centralUiImage": "docker.io/cnvrg/centralsso:latest",
            "adminUser": "",
            "provider": "",
            "emailDomain": [
                "*"
            ],
            "clientId": "",
            "clientSecret": "",
            "oidcIssuerUrl": "",
            "serviceUrl": "",
            "scope": "openid email profile",
            "insecureOidcAllowUnverifiedEmail": true,
            "whitelistDomain": "",
            "cookieDomain": "",
            "groupsAuth": true
        },
        "authz": {
            "enabled": false,
            "image": "docker.io/cnvrg/proxy:v1.0.0",
            "address": "cnvrg-authz:50052"
        }
    },
    "tenancy": {
        "enabled": false,
        "key": "purpose",
        "value": "cnvrg-control-plane"
    },
    "cnvrgAppPriorityClass": {
        "name": "",
        "value": 0,
        "description": ""
    },
    "cnvrgJobPriorityClass": {
        "name": "",
        "value": 0,
        "description": ""
    }
}

const docs = {
    "clusterDomain": "",
    "clusterInternalDomain": "cluster.local",
    "imageHub": "docker.io/cnvrg",
    "controlPlane": {
        "image": "core:3.6.99",
        "webapp": {
            "replicas": 1,
            "enabled": true,
            "port": 8080,
            "requests": {
                "cpu": "500m",
                "memory": "4Gi"
            },
            "limits": {
                "cpu": "4",
                "memory": "8Gi"
            },
            "svcName": "app",
            "nodePort": 30080,
            "passengerMaxPoolSize": 50,
            "initialDelaySeconds": 10,
            "readinessPeriodSeconds": 25,
            "readinessTimeoutSeconds": 20,
            "failureThreshold": 5,
            "hpa": {
                "enabled": true,
                "utilization": 85,
                "maxReplicas": 5
            }
        },
        "sidekiq": {
            "enabled": true,
            "split": true,
            "requests": {
                "cpu": "200m",
                "memory": "3750Mi"
            },
            "limits": {
                "cpu": "2",
                "memory": "8Gi"
            },
            "replicas": 2,
            "hpa": {
                "enabled": true,
                "utilization": 85,
                "maxReplicas": 5
            }
        },
        "searchkiq": {
            "enabled": true,
            "requests": {
                "cpu": "200m",
                "memory": "1Gi"
            },
            "limits": {
                "cpu": "2",
                "memory": "8Gi"
            },
            "replicas": 1,
            "hpa": {
                "enabled": true,
                "utilization": 85,
                "maxReplicas": 5
            }
        },
        "systemkiq": {
            "enabled": true,
            "requests": {
                "cpu": "300m",
                "memory": "2Gi"
            },
            "limits": {
                "cpu": "2",
                "memory": "8Gi"
            },
            "replicas": 1,
            "hpa": {
                "enabled": true,
                "utilization": 85,
                "maxReplicas": 5
            }
        },
        "hyper": {
            "enabled": true,
            "image": "hyper-server:latest",
            "port": 5050,
            "replicas": 1,
            "nodePort": 30050,
            "svcName": "hyper",
            "token": "token",
            "requests": {
                "cpu": "100m",
                "memory": "200Mi"
            },
            "limits": {
                "cpu": "2",
                "memory": "4Gi"
            },
            "cpuLimit": "",
            "memoryLimit": "",
            "readinessPeriodSeconds": 100,
            "readinessTimeoutSeconds": 60
        },
        "cnvrgScheduler": {
            "enabled": true,
            "requests": {
                "cpu": "200m",
                "memory": "1000Mi"
            },
            "limits": {
                "cpu": "2",
                "memory": "4Gi"
            },
            "replicas": 1
        },
        "cnvrgRouter": {
            "enabled": false,
            "image": "nginx:1.21.0",
            "svcName": "cnvrg-router",
            "nodePort": 30081
        },
        "baseConfig": {
            "jobsStorageClass": "",
            "featureFlags": null,
            "sentryUrl": "",
            "runJobsOnSelfCluster": "",
            "agentCustomTag": "latest",
            "intercom": "true",
            "cnvrgJobUid": "0",
            "cnvrgJobRbacStrict": false,
            "cnvrgPrivilegedJob": true,
            "metagpuEnabled": false
        },
        "ldap": {
            "enabled": false,
            "host": "",
            "port": "",
            "account": "userPrincipalName",
            "base": "",
            "adminUser": "",
            "adminPassword": "",
            "ssl": ""
        },
        "smtp": {
            "server": "",
            "port": 587,
            "username": "",
            "password": "",
            "domain": "",
            "opensslVerifyMode": "",
            "sender": "info@cnvrg.io"
        },
        "objectStorage": {
            "type": "minio",
            "bucket": "cnvrg-storage",
            "region": "eastus",
            "accessKey": "",
            "secretKey": "",
            "endpoint": "",
            "azureAccountName": "",
            "azureContainer": "",
            "gcpProject": "",
            "gcpSecretRef": "gcp-storage-secret"
        },
        "mpi": {
            "enabled": false,
            "image": "mpioperator/mpi-operator:v0.2.3",
            "kubectlDeliveryImage": "mpioperator/kubectl-delivery:v0.2.3",
            "extraArgs": null,
            "registry": {
                "name": "mpi-private-registry",
                "url": "docker.io",
                "user": "",
                "password": ""
            },
            "requests": {
                "cpu": "100m",
                "memory": "100Mi"
            },
            "limits": {
                "cpu": "1000m",
                "memory": "1Gi"
            }
        },
        "nomex": {
            "enabled": true,
            "image": "nomex:v1.0.0"
        }
    },
    "registry": {
        "name": "cnvrg-app-registry",
        "url": "docker.io",
        "user": "",
        "password": ""
    },
    "dbs": {
        "pg": {
            "enabled": true,
            "serviceAccount": "pg",
            "image": "postgresql-12-centos7:latest",
            "port": 5432,
            "storageSize": "80Gi",
            "svcName": "postgres",
            "storageClass": "",
            "requests": {
                "cpu": "1",
                "memory": "4Gi"
            },
            "limits": {
                "cpu": "12",
                "memory": "32Gi"
            },
            "maxConnections": 500,
            "sharedBuffers": "1024MB",
            "effectiveCacheSize": "2048MB",
            "hugePages": {
                "enabled": false,
                "size": "2Mi",
                "memory": ""
            },
            "nodeSelector": null,
            "credsRef": "pg-creds",
            "pvcName": "pg-storage"
        },
        "redis": {
            "enabled": true,
            "serviceAccount": "redis",
            "image": "cnvrg-redis:v3.0.5.c2",
            "svcName": "redis",
            "port": 6379,
            "storageSize": "10Gi",
            "storageClass": "",
            "requests": {
                "cpu": "100m",
                "memory": "200Mi"
            },
            "limits": {
                "cpu": "1000m",
                "memory": "2Gi"
            },
            "nodeSelector": null,
            "credsRef": "redis-creds",
            "pvcName": "redis-storage"
        },
        "minio": {
            "enabled": true,
            "serviceAccount": "minio",
            "replicas": 1,
            "image": "minio:RELEASE.2021-05-22T02-34-39Z",
            "port": 9000,
            "storageSize": "100Gi",
            "svcName": "minio",
            "nodePort": 30090,
            "storageClass": "",
            "requests": {
                "cpu": "200m",
                "memory": "2Gi"
            },
            "limits": {
                "cpu": "8",
                "memory": "20Gi"
            },
            "nodeSelector": null,
            "pvcName": "minio-storage"
        },
        "es": {
            "enabled": true,
            "serviceAccount": "es",
            "image": "cnvrg-es:v7.8.1.a1-dynamic-indices",
            "port": 9200,
            "storageSize": "80Gi",
            "svcName": "elasticsearch",
            "nodePort": 32200,
            "storageClass": "",
            "requests": {
                "cpu": "500m",
                "memory": "4Gi"
            },
            "limits": {
                "cpu": "4",
                "memory": "8Gi"
            },
            "javaOpts": "",
            "patchEsNodes": false,
            "nodeSelector": null,
            "credsRef": "es-creds",
            "pvcName": "es-storage",
            "cleanupPolicy": {
                "all": "3d",
                "app": "30d",
                "jobs": "14d",
                "endpoints": "1825d"
            },
            "kibana": {
                "enabled": true,
                "serviceAccount": "kibana",
                "svcName": "kibana",
                "port": 8080,
                "image": "kibana-oss:7.8.1",
                "nodePort": 30601,
                "requests": {
                    "cpu": "100m",
                    "memory": "200Mi"
                },
                "limits": {
                    "cpu": "1000m",
                    "memory": "2Gi"
                },
                "credsRef": "kibana-creds"
            },
            "elastalert": {
                "enabled": true,
                "image": "elastalert:3.0.0-beta.1",
                "authProxyImage": "nginx:1.20",
                "credsRef": "elastalert-creds",
                "port": 80,
                "nodePort": 32030,
                "storageSize": "30Gi",
                "svcName": "elastalert",
                "storageClass": "",
                "requests": {
                    "cpu": "100m",
                    "memory": "200Mi"
                },
                "limits": {
                    "cpu": "400m",
                    "memory": "800Mi"
                },
                "nodeSelector": null,
                "pvcName": "elastalert-storage"
            }
        },
        "cvat": {
            "enabled": false,
            "pg": {
                "enabled": false,
                "serviceAccount": "cvat-pg",
                "image": "postgresql-12-centos7:latest",
                "port": 5432,
                "storageSize": "100Gi",
                "svcName": "cvat-postgres",
                "storageClass": "",
                "requests": {
                    "cpu": "1",
                    "memory": "2Gi"
                },
                "limits": {
                    "cpu": "2",
                    "memory": "4Gi"
                },
                "maxConnections": 500,
                "sharedBuffers": "1024MB",
                "effectiveCacheSize": "2048MB",
                "hugePages": {
                    "enabled": false,
                    "size": "2Mi",
                    "memory": ""
                },
                "nodeSelector": null,
                "credsRef": "cvat-pg-secret",
                "pvcName": "cvat-pg-storage"
            },
            "redis": {
                "enabled": false,
                "serviceAccount": "cvat-redis",
                "image": "redis:4.0.5-alpine",
                "svcName": "cvat-redis",
                "port": 6379,
                "storageSize": "10Gi",
                "storageClass": "",
                "requests": {
                    "cpu": "250m",
                    "memory": "500Mi"
                },
                "limits": {
                    "cpu": "1000m",
                    "memory": "2Gi"
                },
                "nodeSelector": null,
                "credsRef": "cvat-redis-secret",
                "pvcName": "cvat-redis-storage"
            }
        },
        "prom": {
            "enabled": true,
            "credsRef": "prom-creds",
            "extraScrapeConfigs": null,
            "image": "prometheus:v2.37.1",
            "grafana": {
                "enabled": true,
                "image": "grafana-oss:9.1.6",
                "svcName": "grafana",
                "port": 8080,
                "nodePort": 30012,
                "credsRef": "grafana-creds"
            }
        }
    },
    "networking": {
        "ingress": {
            "type": "istio",
            "timeout": "18000s",
            "retriesAttempts": 5,
            "perTryTimeout": "3600s",
            "istioGwEnabled": false,
            "istioGwName": ""
        },
        "https": {
            "enabled": false,
            "certSecret": ""
        },
        "proxy": {
            "enabled": false,
            "configRef": "cp-proxy",
            "httpProxy": [],
            "httpsProxy": [],
            "noProxy": []
        }
    },
    "sso": {
        "enabled": false,
        "pki": {
            "enabled": false,
            "rootCaSecret": "sso-idp-root-ca",
            "privateKeySecret": "sso-idp-private-key",
            "publicKeySecret": "sso-idp-pki-public-key"
        },
        "jwks": {
            "enabled": false,
            "name": "cnvrg-jwks",
            "image": "jwks:latest",
            "cache": {
                "enabled": true,
                "image": "redis:4.0.5-alpine"
            }
        },
        "central": {
            "enabled": false,
            "publicUrl": "",
            "cnvrgProxyImage": "proxy:v1.0.1",
            "oauthProxyImage": "oauth2-proxy:v7.3.x.ssov3.p2-01",
            "centralUiImage": "centralsso:latest",
            "adminUser": "",
            "provider": "",
            "emailDomain": [
                "*"
            ],
            "clientId": "",
            "clientSecret": "",
            "oidcIssuerUrl": "",
            "serviceUrl": "",
            "scope": "openid email profile",
            "insecureOidcAllowUnverifiedEmail": true,
            "whitelistDomain": [],
            "cookieDomain": [],
            "groupsAuth": false
        },
        "authz": {
            "enabled": false,
            "image": "proxy:v1.0.0",
            "address": "cnvrg-authz:50052"
        }
    },
    "tenancy": {
        "enabled": false,
        "key": "purpose",
        "value": "cnvrg-control-plane"
    },
    "cnvrgAppPriorityClass": {
        "name": "",
        "value": 0,
        "description": ""
    },
    "cnvrgJobPriorityClass": {
        "name": "",
        "value": 0,
        "description": ""
    }
}


const flattenJSON = (obj = {}, res = {}, extraKey = '') => {
    for (key in obj) {
        if (typeof obj[key] !== 'object') {
            res[extraKey + key] = obj[key];
        } else {
            flattenJSON(obj[key], res, `${extraKey}${key}.`);
        }
    }
    return res;
};

// let flat = flattenJSON(spec)
let flat = flattenJSON(docs)
Object.entries(flat).forEach((el, idx) => {
    let val = el[1]
    if (val === "") {
        val = "-"
    }
    console.log("|`" + el[0] + "` | " + val + "|")
    // _.set(spec, el[0], "{{" + ".Values."+ el[0] + "}}")
})
// let yamlSpec = YAML
//     .stringify(spec, null, 2)
//     .replaceAll('"{{', "{{")
//     .replaceAll('}}"', "}}")
// console.log('{{- define "cap_spec" }}')
// console.log(yamlSpec)
// console.log('{{- end }}')

