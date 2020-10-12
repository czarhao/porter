package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func PathExists(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

func Path2name(path string, subSuffix bool) string {
	s, e := 0, 0
	for i, v := range path {
		switch v {
		case '.':
			e = i
		case '/':
			s = i
		}
	}
	if s < e && subSuffix {
		return path[s+1 : e]
	}
	return path[s:]
}

func ReadUserCommand() ([]string, error) {
	pipe := os.NewFile(uintptr(3), "pipe")
	defer pipe.Close()
	msg, err := ioutil.ReadAll(pipe)
	if err != nil {
		return nil, fmt.Errorf("init proc read pipe error %v", err)
	}
	msgStr := string(msg)
	return strings.Split(msgStr, " "), nil
}