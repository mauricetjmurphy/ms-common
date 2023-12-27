package mwauthz

import (
	"context"
	"regexp"

	"github.com/NBCUniversal/gvs-ms-common/grpc/auth_context"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/pkg/errors"
)

const tokenSSOKeyMetadata = "sso"

var ssoDigitPattern = regexp.MustCompile(`^[0-9]{9}$`)

type ssoMetadataAuthz struct {
}

func NewSSOMetadataAuthz() Authenticator {
	return &ssoMetadataAuthz{}
}

func (a *ssoMetadataAuthz) HandleAuth(ctx context.Context, _ interface{}) (context.Context, error) {
	ssoId, err := ExtractSSO(ctx)
	if err != nil {
		return ctx, err
	}
	return auth_context.WithSso(ctx, ssoId), nil
}

func ExtractSSO(ctx context.Context) (string, error) {
	val := metautils.ExtractIncoming(ctx).Get(tokenSSOKeyMetadata)
	if val == "" {
		return "", errors.New("authz : unauthorized missing token")
	}

	if !ssoDigitPattern.MatchString(val) {
		return "", errors.New("authz : unauthorized invalid token")
	}

	return val, nil
}
