package cgetopt

// #include <unistd.h>
// #include <stdlib.h>
//
// void pushToArray(char **arr, size_t ind, char *val) {
//	arr[ind] = val;
// }
import "C"
import "unsafe"

var Optarg string
var (
	Optind,
	Opterr,
	Optopt int
)

func Cgetopt(args []string, optstring string) rune {
	ln := len(args)

	var cargs **C.char
	cargs = (**C.char)(C.malloc((C.ulong)(uintptr(ln) * unsafe.Sizeof(cargs))))
	defer C.free(unsafe.Pointer(cargs))

	for i, j := range args {
		cj := C.CString(j)
		defer C.free(unsafe.Pointer(cj))

		C.pushToArray(cargs, C.ulong(i), cj)
	}

	coptstring := C.CString(optstring)
	defer C.free(unsafe.Pointer(coptstring))

	ret := rune(C.getopt(C.int(ln), cargs, coptstring))
	if ret == -1 {
		return ret
	}

	if C.optarg != nil {
		Optarg = C.GoString(C.optarg)
	}
	Opterr = int(C.opterr)
	Optind = int(C.optind)
	Optopt = int(C.optopt)

	return ret
}
