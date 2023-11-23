package shellcode

import (
	"syscall"
	"unsafe"
)

func runShellcode(shellcode []byte, bg bool) error {
	sc_addr := uintptr(unsafe.Pointer(&shellcode[0]))
	page := (*(*[0xFFFFFF]byte)(unsafe.Pointer(sc_addr & ^uintptr(syscall.Getpagesize()-1))))[:syscall.Getpagesize()]
	syscall.Mprotect(page, syscall.PROT_READ|syscall.PROT_EXEC)
	spointer := unsafe.Pointer(&shellcode)
	sc_ptr := *(*func())(unsafe.Pointer(&spointer))
	if bg {
		go sc_ptr()
	} else {
		sc_ptr()
	}

	return nil
}
