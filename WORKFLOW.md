_below is more like open discussion than a design doc_
### Introduction

The Ceph COSI driver is similar to Ceph CSI, it implements interface between COSI enabled contianerised environment and Ceph RGW. It dynamically provision RGW buckets and access for that bucket, can attach them to workloads.

The following CRDs are defined for managing the lifecycle of Buckets:

 - BucketRequest(BR) - Represents a request to provision a Bucket
 - BucketClass(BC) - Represents a class of Buckets with similar characteristics
 - Bucket(B) - Represents a Bucket or its equivalent in the storage backend

 The following CRDs are defined for managing the lifecycle of workloads accessing the Bucket:

 - BucketAccessRequest(BAR) - Represents a request to access a Bucket
 - BucketAccessClass(BAC) - Represents a class of accessors with similar access requirements
 - BucketAccess(BA) - Represents a access token or service account in the storage backend


Similar to CSI, it has three running components :
- Controller - manages life cycle of COSI objects, watches for the requests for above mentioned objects
- Node adapter - bootstraps bucket access request to applications for their workloads
- Sidecar/Driver - receives the request and calls corresponding API in the driver

More details and definitions can be found in the [COSI Spec](https://github.com/kubernetes-sigs/container-object-storage-interface-spec/blob/master/spec.md). The original COSI [KEP](https://github.com/kubernetes/enhancements/tree/master/keps/sig-storage/1979-object-storage-support) may also be a useful reference.

### Workflow
1. User makes a BucketRequest (BR) (namespace-scoped) with reference to a BucketClass (BC)
1. The COSI controller makes a request to Ceph COSI driver (via gRPC) to create a bucket.
1. Ceph COSI driver will create the bucket and return the result
1. The COSI controller creates a Bucket (cluster-scoped) for the BucketRequest (BR)
1. User makes a BucketAccessRequest (BAR) (namespace-scoped) with reference to a BucketAccessClass (BAC) and BucketRequest (BR)
1. The COSI controller makes a request to Ceph COSI driver (via gRPC) to grant access to the bucket
1. Ceph COSI driver returns S3 bucket access credentials to COSI (via gRPC)
1. The COSI controller creates a BucketAccess (BA) (cluster-scoped) for the BucketAccessRequest (BAR)

The BC and BAC are expected to be created by the admin, and users are expected to create BRs and BARs. When an application Pod is given access to a BAR, both need to be running in the same namespace. Even though a BR is namespace-scoped, access to it is limited to a selection of allowed namespaces by the COSI controller.

### COSI Driver
The user credentials and endpoint need to be passed as a Kubernetes Secret to the driver. Using these details, the driver will perform bucket operations via S3 APIs and user operations via RADOS Gateway administrator operations (RGW admin ops) APIs. The following are the API's defined by the driver:

- ProvisionerCreateBucket() - called upon create BucketRequest, the driver will create user with help of RGW admin ops API and then create a bucket using S3 API using that user.
- ProvisionerDeleteBucket() - called upon delete BucketRequest, the driver will fetch BucketInfo using RGW admin ops API, then fetch the owner info, delete bucket via S3 API using the owner.
- ProvisionerGrantBucketAccess() - called upon create BucketAccessRequest, create user with RGW admin ops API and attach the user to required bucket with help of bucket policy.
- ProvisionerRevokeBucketAccess() - called upon delete BucketAccessRequest, revoke the permission on with user and then delete the user.

```
P.S Since there can be multiple users(BARs) attaching to the same bucket, it is not tested how much RGW scale with COSI.
```

### Consuming at application
Users can attach buckets to application Pods by referencing a BucketAccessRequest and the COSI node adapter details as volume mount as follows:

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
The COSI adapter sidecar injects the Ceph driver's credentials into the application pod into the mount path given (/data/cosi in the above example).

### Integrating with Rook
The controller, node adapter, and sidecar/driver need to be deployed. This design proposes that Rook deploy these components for a Rook-Ceph cluster. This design also proposes that Rook create a `CephObjectStoreUser` and refer to the credentials secret created by that controller.
