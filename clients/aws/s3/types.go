package s3

import (
	"time"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type Predicable interface {
	Predicate(item types.Object) bool
}

type PredicateFunc func(item types.Object) bool

func (fn PredicateFunc) Predicate(item types.Object) bool {
	return fn(item)
}

type GetKeysParam struct {
	Bucket    string
	Prefix    string
	Predicate PredicateFunc
}

type Item struct {
	Key   string
	Data  []byte
	Error error
}

type GetObjectParam struct {
	Bucket string
	Key    string
}

type GetObjectsAsyncParam struct {
	Bucket string
	Keys   []string
}

type PutObjectRequest *s3.PutObjectInput

type PutObjectResponse *s3.PutObjectOutput

type GetObjectRequest *s3.GetObjectInput

type GetObjectResponse *s3.GetObjectOutput

type GetPresignURLRequest struct {
	Bucket   string
	Key      string
	LifeTime time.Duration
}

type GetPresignURLResponse *v4.PresignedHTTPRequest

type GetPresignURLsAsyncRequest struct {
	GetObjectsAsyncParam
	LifeTime time.Duration
}

type PresignItem struct {
	*v4.PresignedHTTPRequest
	Key   string
	Error error
}

type GetPresignURLsAsyncResponse []*PresignItem
