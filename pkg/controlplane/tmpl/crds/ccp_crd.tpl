
---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.4.1
  creationTimestamp: null
  name: cnvrgclusterprovisioners.mlops.cnvrg.io
spec:
  group: mlops.cnvrg.io
  names:
    kind: CnvrgClusterProvisioner
    listKind: CnvrgClusterProvisionerList
    plural: cnvrgclusterprovisioners
    shortNames:
    - ccp
    singular: cnvrgclusterprovisioner
  scope: Namespaced
  versions:
  - name: v1
    schema:
      openAPIV3Schema:
        description: CnvrgClusterProvisioner is the Schema for the cnvrgclusterprovisioners
          API
        properties:
          apiVersion:
            description: 'APIVersion defines the versioned schema of this representation
              of an object. Servers should convert recognized schemas to the latest
              internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
            type: string
          kind:
            description: 'Kind is a string value representing the REST resource this
              object represents. Servers may infer this from the endpoint the client
              submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
            type: string
          metadata:
            type: object
          spec:
            description: CnvrgClusterProvisionerSpec defines the desired state of
              CnvrgClusterProvisioner
            properties:
              aws:
                properties:
                  availabilityZones:
                    items:
                      type: string
                    type: array
                  metadata:
                    properties:
                      region:
                        type: string
                      version:
                        type: string
                    type: object
                  nodeGroups:
                    items:
                      properties:
                        availabilityZones:
                          items:
                            type: string
                          type: array
                        iamAddonPolicies:
                          properties:
                            AWSLoadBalancerController:
                              type: boolean
                            XRay:
                              type: boolean
                            appMesh:
                              type: boolean
                            appMeshPreview:
                              type: boolean
                            autoScaler:
                              type: boolean
                            certManager:
                              type: boolean
                            cloudWatch:
                              type: boolean
                            ebs:
                              type: boolean
                            efs:
                              type: boolean
                            externalDNS:
                              type: boolean
                            fsx:
                              type: boolean
                            imageBuilder:
                              type: boolean
                          type: object
                        iamAttachPolicyARNs:
                          items:
                            type: string
                          type: array
                        metadata:
                          properties:
                            autoScaling:
                              type: boolean
                            desiredCapacity:
                              type: integer
                            instanceType:
                              type: string
                            isGpu:
                              type: boolean
                            isHpu:
                              type: boolean
                            labels:
                              additionalProperties:
                                type: string
                              type: object
                            maxSize:
                              type: integer
                            minSize:
                              type: integer
                            name:
                              type: string
                            privateNetworking:
                              type: boolean
                            tags:
                              additionalProperties:
                                type: string
                              type: object
                            taints:
                              items:
                                properties:
                                  effect:
                                    type: string
                                  key:
                                    type: string
                                  value:
                                    type: string
                                required:
                                - effect
                                - key
                                - value
                                type: object
                              type: array
                            volumeSize:
                              type: integer
                          type: object
                        securityGroups:
                          items:
                            type: string
                          type: array
                        spotInstances:
                          type: boolean
                      type: object
                    type: array
                  vpc:
                    properties:
                      clusterEndpoints:
                        properties:
                          privateAccess:
                            type: boolean
                          publicAccess:
                            type: boolean
                        type: object
                      id:
                        type: string
                      privateSubnets:
                        items:
                          type: string
                        type: array
                      publicAccessCIDRs:
                        items:
                          type: string
                        type: array
                      publicSubnets:
                        items:
                          type: string
                        type: array
                      securityGroup:
                        type: string
                    type: object
                type: object
              clusterSlug:
                type: string
              cnvrgDedicatedNodeGroup:
                properties:
                  spotInstances:
                    type: boolean
                type: object
              metadata:
                additionalProperties:
                  type: string
                type: object
              name:
                description: Foo is an example field of CnvrgClusterProvisioner. Edit
                  cnvrgclusterprovisioner_types.go to remove/update Foo string `json:"foo,omitempty"`
                type: string
              role:
                additionalProperties:
                  type: string
                type: object
              type:
                type: string
            required:
            - clusterSlug
            - role
            type: object
          status:
            description: CnvrgClusterProvisionerStatus defines the observed state
              of CnvrgClusterProvisioner
            properties:
              cluster:
                additionalProperties:
                  properties:
                    message:
                      type: string
                    status:
                      type: string
                  type: object
                type: object
              general:
                description: 'INSERT ADDITIONAL STATUS FIELD - define observed state
                  of cluster Important: Run "make" to regenerate code after modifying
                  this file'
                properties:
                  message:
                    type: string
                  status:
                    type: string
                type: object
              postActions:
                properties:
                  message:
                    type: string
                  status:
                    type: string
                type: object
            type: object
        type: object
    served: true
    storage: true
    subresources:
      status: {}
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
