package strings

import "unsafe"

func GoString(p uintptr) string {
	ptr := *(*unsafe.Pointer)(unsafe.Pointer(&p))
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
