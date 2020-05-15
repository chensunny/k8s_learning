package errors

import (
	"fmt"
	"runtime"
	"strings"
)

func SmallerStacktrace(skip, depth int) string {
	callers := make([]uintptr, 128)
	//skip SmallerStacktrace、runtime.Callers，so skip+2
	n := runtime.Callers(skip+2, callers)
	frames := runtime.CallersFrames(callers[:n-1])
	var callerStrings []string
	for {
		frame, ok := frames.Next()
		if !ok {
			break
		}
		fullFunctionName := frame.Function
		callerStr := fmt.Sprintf("%s %s %d", frame.File, fullFunctionName[strings.LastIndex(fullFunctionName, "/")+1:], frame.Line)
		callerStrings = append(callerStrings, callerStr)
	}
	return strings.Join(callerStrings[:depth], "\n")
}
