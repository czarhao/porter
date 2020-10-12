package container

import "syscall"

const (
	CLONE_FLAGS = syscall.CLONE_NEWUTS |
				  syscall.CLONE_NEWPID |
				  syscall.CLONE_NEWNS |
				  syscall.CLONE_NEWNET |
				  syscall.CLONE_NEWIPC
)
