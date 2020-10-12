package mount

const (
	/* 36 35 98:0 /mnt1 /mnt2 rw,noatime master:1 - ext3 /dev/root rw,errors=continue
	   (1)(2)(3)   (4)   (5)      (6)      (7)   (8) (9)   (10)         (11)

	   (1) mount ID:  unique identifier of the mount (may be reused after umount)
	   (2) parent ID:  ID of parent (or of self for the top of the mount tree)
	   (3) major:minor:  value of st_dev for files on filesystem
	   (4) root:  root of the mount within the filesystem
	   (5) mount point:  mount point relative to the process's root
	   (6) mount options:  per mount options
	   (7) optional fields:  zero or more fields of the form "tag[:value]"
	   (8) separator:  marks the end of the optional fields
	   (9) filesystem type:  name of filesystem of the form "type[.subtype]"
	   (10) mount source:  filesystem specific information or "none"
	   (11) super options:  per super block options*/
	mountinfoFormat = "%d %d %d:%d %s %s %s %s"
)

type Info struct {
	// ID is a unique identifier of the mount (may be reused after umount).
	ID int

	// Parent indicates the ID of the mount parent (or of self for the top of the
	// mount tree).
	Parent int

	// Major indicates one half of the device ID which identifies the device class.
	Major int

	// Minor indicates one half of the device ID which identifies a specific
	// instance of device.
	Minor int

	// Root of the mount within the filesystem.
	Root string

	// Mountpoint indicates the mount point relative to the process's root.
	Mountpoint string

	// Opts represents mount-specific options.
	Opts string

	// Optional represents optional fields.
	Optional string

	// Fstype indicates the type of filesystem, such as EXT3.
	Fstype string

	// Source indicates filesystem specific information or "none".
	Source string

	// VfsOpts represents per super block options.
	VfsOpts string
}

var flags = map[string]struct {
	clear bool
	flag  int
}{
	"defaults":      {false, 0},
	"ro":            {false, RDONLY},
	"rw":            {true, RDONLY},
	"suid":          {true, NOSUID},
	"nosuid":        {false, NOSUID},
	"dev":           {true, NODEV},
	"nodev":         {false, NODEV},
	"exec":          {true, NOEXEC},
	"noexec":        {false, NOEXEC},
	"sync":          {false, SYNCHRONOUS},
	"async":         {true, SYNCHRONOUS},
	"dirsync":       {false, DIRSYNC},
	"remount":       {false, REMOUNT},
	"mand":          {false, MANDLOCK},
	"nomand":        {true, MANDLOCK},
	"atime":         {true, NOATIME},
	"noatime":       {false, NOATIME},
	"diratime":      {true, NODIRATIME},
	"nodiratime":    {false, NODIRATIME},
	"bind":          {false, BIND},
	"rbind":         {false, RBIND},
	"unbindable":    {false, UNBINDABLE},
	"runbindable":   {false, RUNBINDABLE},
	"private":       {false, PRIVATE},
	"rprivate":      {false, RPRIVATE},
	"shared":        {false, SHARED},
	"rshared":       {false, RSHARED},
	"slave":         {false, SLAVE},
	"rslave":        {false, RSLAVE},
	"relatime":      {false, RELATIME},
	"norelatime":    {true, RELATIME},
	"strictatime":   {false, STRICTATIME},
	"nostrictatime": {true, STRICTATIME},
}
