package ptr_test

import (
	"testing"
	"time"

	"github.com/NBCUniversal/gvs-ms-common/libs/ptr"
	"github.com/stretchr/testify/assert"
)

func TestBool(t *testing.T) {
	ptrVal := ptr.Bool(false)
	assert.Equal(t, false, *ptrVal)
}

func TestInt(t *testing.T) {
	ptrVal := ptr.Int(123)
	assert.Equal(t, 123, *ptrVal)
}

func TestInt32(t *testing.T) {
	ptrVal := ptr.Int32(123)
	assert.Equal(t, int32(123), *ptrVal)
}

func TestInt64(t *testing.T) {
	ptrVal := ptr.Int64(123)
	assert.Equal(t, int64(123), *ptrVal)
}

func TestString(t *testing.T) {
	ptrVal := ptr.String("foobar")
	assert.Equal(t, "foobar", *ptrVal)
}

func TestTime(t *testing.T) {
	ptrVal := ptr.Time(time.Unix(123456, 0))
	assert.Equal(t, time.Unix(123456, 0), *ptrVal)
}

func TestDuration(t *testing.T) {
	ptrVal := ptr.Duration(4 * time.Minute)
	assert.Equal(t, 4*time.Minute, *ptrVal)
}
