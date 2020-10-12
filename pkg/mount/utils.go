package mount

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"
)

// cat uid_map, if show:
//          0          0 4294967295
// proc run in userNS
func isRunningInUserNS() bool {
	file, err := os.Open("/proc/self/uid_map")
	if err != nil {
		return false
	}
	defer file.Close()

	buf := bufio.NewReader(file)
	l, _, err := buf.ReadLine()
	if err != nil {
		return false
	}

	var a, b, c int64
	_, _ = fmt.Sscanf(string(l), "%d %d %d", &a, &b, &c)
	if a == 0 && b == 0 && c == 4294967295 {
		return false
	}
	return true
}

func makeRSlave(mountPoint string) error {
	mounted, err := isMounted(mountPoint)
	if err != nil {
		return err
	}
	if !mounted {
		if err := mount(mountPoint, mountPoint, "none", "bind,rw"); err != nil {
			return err
		}
	}
	if _, err = isMounted(mountPoint); err != nil {
		return err
	}
	return forceMount("", mountPoint, "none", "rslave")
}

func isMounted(mountPoint string) (bool, error) {
	entries, err := parseMountTable()
	if err != nil {
		return false, err
	}

	// Search the table for the mountpoint
	for _, e := range entries {
		if e.Mountpoint == mountPoint {
			return true, nil
		}
	}
	return false, nil
}

func mount(device, target, mType, options string) error {
	flag, _ := parseOptions(options)
	if flag&REMOUNT != REMOUNT {
		if mounted, err := isMounted(target); err != nil || mounted {
			return err
		}
	}
	return forceMount(device, target, mType, options)
}

func parseMountTable() ([]*Info, error) {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var (
		s   = bufio.NewScanner(f)
		out []*Info
	)

	for s.Scan() {
		if err := s.Err(); err != nil {
			return nil, err
		}

		var (
			p              = &Info{}
			text           = s.Text()
			optionalFields string
		)

		if _, err := fmt.Sscanf(text, mountinfoFormat,
			&p.ID, &p.Parent, &p.Major, &p.Minor,
			&p.Root, &p.Mountpoint, &p.Opts, &optionalFields); err != nil {
			return nil, fmt.Errorf("scanning '%s' failed: %s", text, err)
		}
		// Safe as mountinfo encodes mountpoints with spaces as \040.
		index := strings.Index(text, " - ")
		postSeparatorFields := strings.Fields(text[index+3:])
		if len(postSeparatorFields) < 3 {
			return nil, fmt.Errorf("Error found less than 3 fields post '-' in %q", text)
		}

		if optionalFields != "-" {
			p.Optional = optionalFields
		}

		p.Fstype = postSeparatorFields[0]
		p.Source = postSeparatorFields[1]
		p.VfsOpts = strings.Join(postSeparatorFields[2:], " ")
		out = append(out, p)
	}
	return out, nil
}

func parseOptions(options string) (uintptr, string) {
	var (
		flag int
		data []string
	)
	for _, o := range strings.Split(options, ",") {
		if f, exists := flags[o]; exists && f.flag != 0 {
			if f.clear {
				flag &= ^f.flag
			} else {
				flag |= f.flag
			}
		} else {
			data = append(data, o)
		}
	}
	return uintptr(flag), strings.Join(data, ",")
}

func forceMount(device, target, mType, options string) error {
	flag, data := parseOptions(options)
	if err := syscall.Mount(device, target, mType, flag, data); err != nil {
		return err
	}
	if flag&syscall.MS_BIND == syscall.MS_BIND && flag&syscall.MS_RDONLY == syscall.MS_RDONLY {
		return syscall.Mount(device, target, mType, flag|syscall.MS_REMOUNT, data)
	}
	return nil
}
