package s3

import (
	"context"
	"io"
	"io/ioutil"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/mauricetjmurphy/ms-common/libs/ptr"
)

//go:generate mockery --output s3mocks --outpkg s3mocks --name Client
type Client interface {
	GetKeys(context.Context, *GetKeysParam) ([]string, error)
	GetObject(context.Context, *GetObjectParam) ([]byte, error)
	GetObjectV2(context.Context, GetObjectRequest) (GetObjectResponse, error)
	GetObjectsAsync(context.Context, *GetObjectsAsyncParam) ([]*Item, error)
	PutObject(context.Context, PutObjectRequest) (PutObjectResponse, error)
	GetPresignURL(context.Context, *GetPresignURLRequest) (GetPresignURLResponse, error)
	GetPresignURLsAsync(context.Context, *GetPresignURLsAsyncRequest) (GetPresignURLsAsyncResponse, error)
}

// clientImpl implements the Client interface.
type clientImpl struct {
	*s3.Client
	*s3.PresignClient
}

// New creates the S3 instance
func New(cfg aws.Config, optFns ...func(*s3.Options)) Client {
	client := s3.NewFromConfig(cfg, optFns...)
	return &clientImpl{
		Client:        client,
		PresignClient: s3.NewPresignClient(client),
	}
}

func (c clientImpl) GetKeys(ctx context.Context, param *GetKeysParam) ([]string, error) {
	resp, err := c.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: ptr.String(param.Bucket),
		Prefix: ptr.String(param.Prefix),
	})
	if err != nil {
		return nil, err
	}

	var keys []string
	if param.Predicate == nil {
		param.Predicate = func(item types.Object) bool {
			return item.Key != nil
		}
	}

	for _, item := range resp.Contents {
		if param.Predicate(item) {
			keys = append(keys, *item.Key)
		}
	}

	return keys, nil
}

func (c clientImpl) GetObject(ctx context.Context, param *GetObjectParam) ([]byte, error) {
	resp, err := c.Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: ptr.String(param.Bucket),
		Key:    ptr.String(param.Key),
	})
	if err != nil {
		return nil, err
	}

	defer func(reader io.ReadCloser) {
		_ = reader.Close()
	}(resp.Body)

	return ioutil.ReadAll(resp.Body)
}

func (c clientImpl) GetObjectsAsync(ctx context.Context, param *GetObjectsAsyncParam) ([]*Item, error) {
	var wg sync.WaitGroup
	s3GetCh := make(chan *Item)
	for _, key := range param.Keys {
		wg.Add(1)
		go func(param *GetObjectParam) {
			defer wg.Done()
			s3GetCh <- c.getItem(ctx, param)
		}(&GetObjectParam{Bucket: param.Bucket, Key: key})
	}

	go func() {
		wg.Wait()
		close(s3GetCh)
	}()

	var resp []*Item
	for elem := range s3GetCh {
		resp = append(resp, elem)
	}
	return resp, nil
}

func (c clientImpl) PutObject(ctx context.Context, param PutObjectRequest) (PutObjectResponse, error) {
	return c.Client.PutObject(ctx, param)
}

func (c clientImpl) GetObjectV2(ctx context.Context, param GetObjectRequest) (GetObjectResponse, error) {
	return c.Client.GetObject(ctx, param)
}

func (c clientImpl) GetPresignURL(ctx context.Context, param *GetPresignURLRequest) (GetPresignURLResponse, error) {
	return c.PresignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: ptr.String(param.Bucket),
		Key:    ptr.String(param.Key),
	}, func(options *s3.PresignOptions) {
		options.Expires = param.LifeTime
	})
}

func (c clientImpl) GetPresignURLsAsync(ctx context.Context,
	param *GetPresignURLsAsyncRequest) (GetPresignURLsAsyncResponse, error) {
	var wg sync.WaitGroup
	s3GetCh := make(chan *PresignItem)
	for _, key := range param.Keys {
		wg.Add(1)
		go func(param *GetPresignURLRequest) {
			defer wg.Done()
			s3GetCh <- c.getPresignURL(ctx, param)
		}(newGetPresignParam(param, key))
	}

	go func() {
		wg.Wait()
		close(s3GetCh)
	}()

	var resp GetPresignURLsAsyncResponse
	for elem := range s3GetCh {
		resp = append(resp, elem)
	}
	return resp, nil
}

func (c clientImpl) getItem(ctx context.Context, param *GetObjectParam) *Item {
	item := &Item{Key: param.Key}
	data, err := c.GetObject(ctx, param)
	if err != nil {
		item.Error = err
		return item
	}
	item.Data = data
	return item
}

func (c clientImpl) getPresignURL(ctx context.Context, param *GetPresignURLRequest) *PresignItem {
	out, err := c.GetPresignURL(ctx, param)
	return &PresignItem{
		PresignedHTTPRequest: out,
		Key:                  param.Key,
		Error:                err,
	}
}

func newGetPresignParam(param *GetPresignURLsAsyncRequest, key string) *GetPresignURLRequest {
	return &GetPresignURLRequest{
		Bucket:   param.Bucket,
		Key:      key,
		LifeTime: param.LifeTime,
	}
}
