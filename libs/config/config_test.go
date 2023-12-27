package config_test

import (
	"context"
	"os"
	"testing"

	"github.com/mauricetjmurphy/ms-common/libs/config"
	"github.com/mauricetjmurphy/ms-common/libs/config/testdata"
	"github.com/stretchr/testify/assert"
)

var (
	ctx = context.TODO()
)

func TestLoad_WithYaml(t *testing.T) {
	app := testdata.AppConf{}
	err := config.Load(ctx, &app, config.WithYaml(defYamlConfig))
	assert.Nil(t, err)
	assert.Equal(t, "test", app.DataConf.Data)
	assert.Equal(t, "", app.DataConf.Host)
}

func TestLoad_WithEnv(t *testing.T) {
	//Given
	app := testdata.AppConf{}
	os.Setenv("data_app", "override")
	//When
	err := config.Load(ctx, &app, config.WithEnv(""))
	//Then
	assert.Nil(t, err)
	assert.Equal(t, "override", app.DataConf.Data)
	assert.Equal(t, "localhost", app.DataConf.Host)
	os.Clearenv()
}

func TestLoad_WithCombined(t *testing.T) {
	//Given
	app := testdata.AppConf{}
	os.Setenv("DATA_APP", "override")
	//When
	err := config.Load(ctx, &app,
		config.WithYaml(defYamlConfig),
		config.WithEnv(""),
	)
	assert.Nil(t, err)
	//Then
	assert.Equal(t, "override", app.DataConf.Data)
	os.Clearenv()
}
