package duties

import "fmt"

var LogInfoFunc func(msg string)
var LogErrorFunc func(msg string, err error)

var loggingDisabled = false

func logInfo(format string, data ...interface{}) {
	if LogInfoFunc != nil && !loggingDisabled {
		LogInfoFunc(fmt.Sprintf(format, data...))
	}
}

func logError(format string, err error, data ...interface{}) {
	if LogErrorFunc != nil && !loggingDisabled {
		LogErrorFunc(fmt.Sprintf(format, data...), err)
	}
}
