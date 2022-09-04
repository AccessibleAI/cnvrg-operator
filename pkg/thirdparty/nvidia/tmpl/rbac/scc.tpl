kind: SecurityContextConstraints
apiVersion: security.openshift.io/v1
metadata:
  annotations:
    mlops.cnvrg.io/default-loader: "false"
    mlops.cnvrg.io/own: "false"
    mlops.cnvrg.io/updatable: "false"
  name: nvidia-device-plugin
allowHostDirVolumePlugin: true
allowHostIPC: true
allowHostNetwork: true
allowHostPID: true
allowHostPorts: true
allowPrivilegeEscalation: true
allowPrivilegedContainer: true
readOnlyRootFilesystem: false
requiredDropCapabilities: null
allowedCapabilities:
- '*'
allowedUnsafeSysctls:
- '*'
fsGroup:
  type: RunAsAny
runAsUser:
  type: RunAsAny
seLinuxContext:
  type: RunAsAny
seccompProfiles:
- '*'
supplementalGroups:
  type: RunAsAny
users:
- system:serviceaccount:{{ .Namespace }}:nvidia
volumes:
- '*'
