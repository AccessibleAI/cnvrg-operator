apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: hostpath-provisioner
  namespace: {{ .Namespace }}
  labels:
    k8s-app: hostpath-provisioner
spec:
  selector:
    matchLabels:
      k8s-app: hostpath-provisioner
  template:
    metadata:
      labels:
        k8s-app: hostpath-provisioner
    spec:
      serviceAccountName: hostpath-provisioner-admin
      containers:
        - name: hostpath-provisioner
          image: {{ .Spec.Storage.Hostpath.Image }}
          imagePullPolicy: Always
          env:
            - name: USE_NAMING_PREFIX
              value: "true" # change to true, to have the name of the pvc be part of the directory
            - name: NODE_NAME
              valueFrom:
                fieldRef:
                  fieldPath: spec.nodeName
            - name: PV_DIR
              value: {{ .Spec.Storage.Hostpath.HostPath }}
          volumeMounts:
            - name: pv-volume # root dir where your bind mounts will be on the node
              mountPath: {{ .Spec.Storage.Hostpath.HostPath }}
          resources:
            limits:
              cpu: {{ .Spec.Storage.Hostpath.CPULimit }}
              memory: {{ .Spec.Storage.Hostpath.MemoryLimit }}
            requests:
              cpu: {{ .Spec.Storage.Hostpath.CPURequest }}
              memory: {{ .Spec.Storage.Hostpath.MemoryRequest }}
      volumes:
        - name: pv-volume
          hostPath:
            path: {{ .Spec.Storage.Hostpath.HostPath }}
