package uuid_test

import (
	"testing"

	"github.com/NBCUniversal/gvs-ms-common/libs/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUuid(t *testing.T) {
	assert.NotNil(t, uuid.New())
	assert.NotEqualf(t, uuid.New(), uuid.New(), "unique on each generation")
}

func TestNewString(t *testing.T) {
	assert.NotNil(t, uuid.NewString())
	assert.NotEqualf(t, uuid.NewString(), uuid.NewString(), "unique on each generation")
}
