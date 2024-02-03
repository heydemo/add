package addmain

import (
	"embed"
)

//go:embed templates/*.yml
var Templates embed.FS
