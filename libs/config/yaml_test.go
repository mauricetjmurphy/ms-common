package config_test

import (
	"context"
	"testing"

	"github.com/mauricetjmurphy/ms-common/libs/config"
	"github.com/mauricetjmurphy/ms-common/libs/config/testdata"
	"github.com/stretchr/testify/assert"
)

var (
	defYamlConfig = "testdata/app-conf.yaml"
)

func TestLoadYamlConfig(t *testing.T) {
	app := testdata.AppConf{}
	if err := config.NewYaml(defYamlConfig).Load(context.TODO(), &app); err != nil {
		assert.NotNil(t, err)
	}
	assert.Equal(t, "test", app.DataConf.Data)
}
