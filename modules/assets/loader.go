package assets

import (
	"github.com/gobuffalo/packr/v2"
)

// Load ...
func Load() *packr.Box {
	return packr.New("assets", "./assets")
}
