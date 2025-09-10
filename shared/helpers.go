package shared

import (
	"runtime"
)

func GetFunctionName(depthList ...int) string {
	var depth int
	if depthList == nil {
		depth = 1
	} else {
		depth = depthList[0]
	}
	function, _, _, _ := runtime.Caller(depth)

	return runtime.FuncForPC(function).Name()
}

// func shortFuncName(full string) string {
// 	// remove path up to last slash
// 	if idx := strings.LastIndex(full, "/"); idx != -1 {
// 		full = full[idx+1:]
// 	}
// 	// Optionally remove package prefix before dot, keeping Type.Method or func
// 	if i := strings.Index(full, "."); i != -1 {
// 		return full[i+1:]
// 	}
// 	return full
// }
