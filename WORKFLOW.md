_below is more like open disscussion than a design doc_
### Introduction

The following CRDs are defined for managing the lifecycle of Buckets:

 - BucketRequest - Represents a request to provision a Bucket
 - BucketClass - Represents a class of Buckets with similar characteristics
 - Bucket - Represents a Bucket or its equivalent in the storage backend

 The following CRDs are defined for managing the lifecycle of workloads accessing the Bucket:

 - BucketAccessRequest - Represents a request to access a Bucket
 - BucketAccessClass - Represents a class of accessors with similar access requirements
 - BucketAccess - Represents a access token or service account in the storage backend

Simialr to CSI, it has three running components :
- Controller - mananges life cycle of COSI objects, watches for the requests for above mentioned objects
- Node adapter - bootstraps bucket access request to applications for their worloads
- Sidecar/Driver - receives the request and calls corresponding api in the driver

More details and defintions can be found the [COSI Spec](https://github.com/kubernetes-sigs/container-object-storage-interface-spec/blob/master/spec.md)

### Workflow

- BucketRequest(namespace scoped) can be created with reference to a BucketClass, so user will get a Bucket(clusterscoped). The COSI driver will create a `backend bucket` in RGW for a Bucket Request
- BucketAccessRequest(namespace scoped) can be create with reference to a BucketAccessClass and BucketRequest, so user will get a BucketAccess(clusterscoped).The COSI driver will create a `s3 user` with proper access
- 
### COSI Driver
The `user credentials` and `endpoint` need to passed as `Kubernetes Secret` to the Driver. Using that details, the driver will perform bucket operations via S3 apis and user operations via rados-admin Ops apis. 

### Consuming at application
The app pod can refer BucketAccessRequest and Node adapter details as volume mount as follows:
```yaml
spec:
  containers:
      volumeMounts:
        - name: cosi-secrets
          mountPath: /data/cosi
  volumes:
  - name: cosi-secrets
    csi:
      driver: objectstorage.k8s.io
      volumeAttributes:
        bar-name: sample-bar
        bar-namespace: default
```
The adapter bootstrap the credentials into application pod, for Ceph COSI Driver credentials passed as `.aws/credentials` file into the pod.

### Integrating with Rook
The controller, node adapter and sidecar/driver need to deployed via Rook. Rook also need to create a `cephobjectstoreuser` and refer the secret in the deployment.
