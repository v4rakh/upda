//go:build !prod

package server

import "embed"

// fake location for dev build
//
//go:embed web/index.html
var webinterfaceFS embed.FS
