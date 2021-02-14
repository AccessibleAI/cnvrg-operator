apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: {{ .Pg.svcName }}
  namespace: "asd"
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: {{ .Pg.storageSize }}
        {% if pg.storageClass != "use-default" %}
        storageClassName: "{{ pg.storageClass }}"
        {% elif storage.ccpStorageClass != "" %}
        storageClassName: "{{ storage.ccpStorageClass }}"
        {% endif %}