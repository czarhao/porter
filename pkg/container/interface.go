package container

import "porter/pkg/utils"

type Manager interface {
	RunCmd(cmds []string) error
	SetLimit(limit *utils.Limit) error
	WaitProcEnd() error
	DestroyMount() error
	DestroyCgroup() error
}
