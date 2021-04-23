// Copyright 2021 The Kubernetes Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pkg

import (
	"context"

	radosgwapi "github.com/QuentinPerez/go-radosgw/pkg/api"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog/v2"
	cosi "sigs.k8s.io/container-object-storage-interface-spec"
)

type ProvisionerServer struct {
	provisioner        string
	S3Client           *S3Agent
	radosgwAdminClient *radosgwapi.API
}

// ProvisionerCreateBucket is an idempotent method for creating buckets
// It is expected to create the same bucket given a bucketName and protocol
// If the bucket already exists, then it MUST return codes.AlreadyExists
// Return values
//    nil -                   Bucket successfully created
//    codes.AlreadyExists -   Bucket already exists. No more retries
//    non-nil err -           Internal error                                [requeue'd with exponential backoff]
func (s *ProvisionerServer) ProvisionerCreateBucket(ctx context.Context,
	req *cosi.ProvisionerCreateBucketRequest) (*cosi.ProvisionerCreateBucketResponse, error) {
	klog.Info("Using ceph rgw to create Backend Bucket")
	protocol := req.GetProtocol()
	if protocol == nil {
		klog.ErrorS(errors.New("Invalid Argument"), "Protocol is nil")
		return nil, status.Error(codes.InvalidArgument, "Protocol is nil")
	}
	s3 := protocol.GetS3()
	if s3 == nil {
		klog.ErrorS(errors.New("Invalid Argument"), "S3 protocol is nil")
		return nil, status.Error(codes.InvalidArgument, "S3 Protocol is nil")
	}

	bucketName := s3.BucketName
	klog.V(3).InfoS("Create Bucket", "name", bucketName)

	err := s.S3Client.CreateBucket(s3.BucketName)
	if err != nil {
		// Check to see if the bucket already exists by above api
		klog.ErrorS(err, "Bucket creation failed")
		return nil, status.Error(codes.Internal, "Bucket creation failed")
	}
	klog.Infof("Successfully created Backend Bucket %s", s3.BucketName)

	return &cosi.ProvisionerCreateBucketResponse{
		BucketId: s3.BucketName,
	}, nil
}

func (s *ProvisionerServer) ProvisionerDeleteBucket(ctx context.Context,
	req *cosi.ProvisionerDeleteBucketRequest) (*cosi.ProvisionerDeleteBucketResponse, error) {
	/*	protocol := req.GetProtocol()
		if protocol == nil {
			klog.ErrorS(errors.New("Invalid Argument"), "Protocol is nil")
			return nil, status.Error(codes.InvalidArgument, "Protocol is nil")
		}
		s3 := protocol.GetS3()
		if s3 == nil {
			klog.ErrorS(errors.New("Invalid Argument"), "S3 protocol is nil")
			return nil, status.Error(codes.InvalidArgument, "S3 Protocol is nil")
		}

		if _, err := s.S3Client.DeleteBucket(s3.BucketName); err != nil {
			klog.Info("failed to delete bucket", s3.BucketName)
			return nil, status.Error(codes.Internal, "Bucket deletion failed")
		}
	*/
	return &cosi.ProvisionerDeleteBucketResponse{}, nil
}

func (s *ProvisionerServer) ProvisionerGrantBucketAccess(ctx context.Context,
	req *cosi.ProvisionerGrantBucketAccessRequest) (*cosi.ProvisionerGrantBucketAccessResponse, error) {
	/*	protocol := req.GetProtocol()
		if protocol == nil {
			klog.ErrorS(errors.New("Invalid Argument"), "Protocol is nil")
			return nil, status.Error(codes.InvalidArgument, "Protocol is nil")
		}
		s3 := protocol.GetS3()
		if s3 == nil {
			klog.ErrorS(errors.New("Invalid Argument"), "S3 protocol is nil")
			return nil, status.Error(codes.InvalidArgument, "S3 Protocol is nil")
		}
		user, err := s.radosgwAdminClient.CreateUser(radosgwapi.UserConfig{
			UID:         req.GetAccountName(),
			DisplayName: req.GetAccountName(),
		})
		if err != nil {
			klog.Error("failed to create user", err)
			return nil, status.Error(codes.Internal, "User creation failed")
		}

		policy, err := s.S3Client.GetBucketPolicy(s3.BucketName)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				if aerr.Code() != "NoSuchBucketPolicy" {
					return nil, status.Error(codes.Internal, "fetching policy failed")
				}
			}
		}

		statement := NewPolicyStatement().
			WithSID(req.GetAccountName()).
			ForPrincipals(req.GetAccountName()).
			ForResources(s3.BucketName).
			ForSubResources(s3.BucketName).
			Allows().
			Actions(AllowedActions...)
		if policy == nil {
			policy = NewBucketPolicy(*statement)
		} else {
			policy = policy.ModifyBucketPolicy(*statement)
		}
		out, err := s.S3Client.PutBucketPolicy(s3.BucketName, *policy)
		if err != nil {
			klog.Error("failed to set policy", err)
			return nil, status.Error(codes.Internal, "puting policy failed")
		}
		klog.Infof("failed to set policy %v", out)
		return &cosi.ProvisionerGrantBucketAccessResponse{
			AccountId:               req.GetAccountName(),
			CredentialsFileContents: fmt.Sprintf("[default]\naws_access_key %s\naws_secret_key %s", user.Keys[0].AccessKey, user.Keys[0].SecretKey),
			CredentialsFilePath:     ".aws/credentials",
		}, nil
	*/
	return &cosi.ProvisionerGrantBucketAccessResponse{
		AccountId:               "ceph",
		CredentialsFileContents: "{\"username\":\"ceph\", \"password\": \"ceph123\"}",
	}, nil
}

func (s *ProvisionerServer) ProvisionerRevokeBucketAccess(ctx context.Context,
	req *cosi.ProvisionerRevokeBucketAccessRequest) (*cosi.ProvisionerRevokeBucketAccessResponse, error) {
	// revokes user access to bucket
	/*	if err := s.radosgwAdminClient.RemoveUser(radosgwapi.UserConfig{
			UID:         req.GetAccountName(),
			DisplayName: req.GetAccountName(),
		}); err != nil {
			klog.Error("falied to Revoke Bucket Access")
			return nil, status.Error(codes.Internal, "falied to Revoke Bucket Access")
		}*/
	return &cosi.ProvisionerRevokeBucketAccessResponse{}, nil
}