package strings

import "unsafe"

// GoString copies a NUL-terminated string. p must remain valid until it returns.
func GoString(p *byte) string {
	ptr := unsafe.Pointer(p)
	if ptr == nil {
		return ""
	}
	var length int
	for {
		if *(*byte)(unsafe.Add(ptr, uintptr(length))) == '\x00' {
			break
		}
		length++
	}
	return string(unsafe.Slice((*byte)(ptr), length))
}
