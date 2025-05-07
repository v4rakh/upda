//go:build prod

package server

import "embed"

//go:embed web/build/*
var webinterfaceFS embed.FS
