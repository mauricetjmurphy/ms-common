package banner_test

import (
	_ "embed"
	"testing"

	"github.com/mauricetjmurphy/ms-common/libs/banner"
)

var (
	//go:embed "banner.txt"
	bannerData []byte
)

func TestLoadBanner(t *testing.T) {
	banner.Load(bannerData)
}
