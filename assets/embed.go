// Package assets embeds the production's graphics, bitmap fonts and music.
package assets

import "embed"

// Files contains application-owned media; reusable effects live in DCK.
//
//go:embed */*
var Files embed.FS
