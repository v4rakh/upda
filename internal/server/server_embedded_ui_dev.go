//go:build !prod

package server

import "embed"

//go:embed web_dev
var embeddedWebinterfaceFS embed.FS
