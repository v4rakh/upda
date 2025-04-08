//go:build !prod
// +build !prod

package server

import "embed"

//go:embed web_dev
var embeddedFiles embed.FS
