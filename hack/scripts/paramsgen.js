let params = {
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
        "cnvrgRouter": {
            "enabled": false,
            "image": "nginx:1.21.0",
            "svcName": "cnvrg-router",
            "nodePort": 30081
        },
        "baseConfig": {
            "jobsStorageClass": "",
            "featureFlags": {
                "CNVRG_ENABLE_MOUNT_FOLDERS": false,
                "CNVRG_MOUNT_HOST_FOLDERS": false,
                "CNVRG_PROMETHEUS_METRICS": true
            },
            "sentryUrl": "",
            "runJobsOnSelfCluster": "",
            "agentCustomTag": "agnostic-logs",
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
            "image": "cnvrg-redis:v8.0.1",
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
            "image": "minio:RELEASE.2025-04-22T22-12-26Z",
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
            "image": "cnvrg-es:7.17.5",
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
            "patchEsNodes": true,
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
                "authProxyImage": "nginx:1.28.0",
                "credsRef": "elastalert-creds",
                "port": 8080,
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
        "prom": {
            "enabled": true,
            "credsRef": "prom-creds",
            "extraScrapeConfigs": null,
            "image": "prometheus:v2.37.1",
            "grafana": {
                "enabled": true,
                "image": "grafana-oss:9.1.7",
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
            "istioGwEnabled": true,
            "istioGwName": "",
            "istioIngressSelectorKey": "istio",
            "istioIngressSelectorValue": "ingressgateway",
            "ocpSecureRoutes": false
        },
        "https": {
            "enabled": false,
            "certSecret": "",
            "cert": "",
            "key": ""
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
        "version": "v3",
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
            "cacheImage": "redis:v8.0.1"
        },
        "central": {
            "enabled": false,
            "publicUrl": "",
            "oauthProxyImage": "oauth2-proxy:v7.4.ssov3.p6",
            "centralUiImage": "centralsso:0.0.1",
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
            "groupsAuth": false,
            "readiness": true,
            "requests": {
                "cpu": "500m",
                "memory": "1Gi"
            },
            "limits": {
                "cpu": 2,
                "memory": "4Gi"
            }
        },
        "proxy": {
            "enabled": false,
            "image": "cnvrg-proxy:v1.0.15",
            "address": "",
            "readiness": true,
            "requests": {
                "cpu": "500m",
                "memory": "1Gi"
            },
            "limits": {
                "cpu": 2,
                "memory": "4Gi"
            }
        }
    },
    "tenancy": {
        "enabled": false,
        "key": "purpose",
        "value": "cnvrg-control-plane"
    },
    "priorityClass": {
        "appClassRef": "",
        "jobClassRef": ""
    }
};

function toDotNotation(obj, res = {}, current = '') {
    for (const key in obj) {
        let value = obj[key];
        let newKey = (current ? current + "." + key : key);  // joined key with dot
        if (value && typeof value === "object") {
            toDotNotation(value, res, newKey);  // it's a nested object, so do it again
        } else {
            res[newKey] = value;  // it's not an object, so set the property
        }
    }
    return res;
}

let docs = toDotNotation(params)
Object.keys(docs).forEach((key) => {
    if (docs[key] === "") {
        docs[key] = "-"
    }
    console.log(`| \`${key}\` | ${docs[key]} |`)
})
