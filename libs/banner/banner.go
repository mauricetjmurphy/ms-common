package banner

import (
	"bytes"
	"os"

	"github.com/dimiro1/banner"
)

// Load provides a banner from array of bytes.
func Load(data []byte) {
	reader := bytes.NewReader(data)
	banner.Init(os.Stdout, true, true, reader)
}
