package assets

import (
	"github.com/gobuffalo/packr/v2"

	"qmetry_uploader/modules/osx"
)

// Load ...
func Load() error {
	box := packr.New("assets", "./assets")
	return osx.LoadAssets(box)
}
