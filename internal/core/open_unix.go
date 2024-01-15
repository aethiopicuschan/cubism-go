//go:build darwin || freebsd || linux

package core

import "github.com/ebitengine/purego"

func openLibrary(name string) (uintptr, error) {
	return purego.Dlopen(name, purego.RTLD_NOW|purego.RTLD_GLOBAL)
}
