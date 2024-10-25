//go:build prod
// +build prod

package server

import "embed"

//go:embed web/build/*
var embeddedFiles embed.FS
