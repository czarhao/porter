package actions

import (
	"fmt"
	"runtime"
)

const version = "0.1"

func Version() string {
	v := fmt.Sprintf("Porter version: %v  &&  ", version)
	v += fmt.Sprintf("GO version: %v", runtime.Version())
	return v
}
