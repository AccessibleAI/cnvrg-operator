kind: SecurityContextConstraints
apiVersion: security.openshift.io/v1
metadata:
  name: "cnvrg-scc-control-plane"
allowHostDirVolumePlugin: false
allowHostIPC: false
allowHostNetwork: false
allowHostPID: false
allowHostPorts: false
allowPrivilegeEscalation: true
allowPrivilegedContainer: false
allowedCapabilities: null
readOnlyRootFilesystem: false
runAsUser:
  type: MustRunAsRange
  uidRangeMin: 26
  uidRangeMax: 1000
seLinuxContext:
  type: RunAsAny
fsGroup:
  type: MustRunAs
  ranges:
  - min: 26
    max: 1000
supplementalGroups:
  type: MustRunAs
  ranges:
  - min: 26
    max: 1000
groups:
- "system:serviceaccounts:cnvrg"
