kind: SecurityContextConstraints
apiVersion: security.openshift.io/v1
metadata:
  name: cnvrg-scc-control-plane
  annotations:
    mlops.cnvrg.io/default-loader: "false"
    mlops.cnvrg.io/own: "false"
    mlops.cnvrg.io/updatable: "false"
allowHostDirVolumePlugin: false
allowHostIPC: false
allowHostNetwork: false
allowHostPID: false
allowHostPorts: false
allowPrivilegeEscalation: true
allowPrivilegedContainer: true
allowedCapabilities: null
readOnlyRootFilesystem: false
runAsUser:
  type: RunAsAny
seLinuxContext:
  type: RunAsAny
fsGroup:
  type: RunAsAny
supplementalGroups:
  type: RunAsAny
groups:
- "system:serviceaccounts:{{ .Namespace }}"
