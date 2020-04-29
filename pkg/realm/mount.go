package realm

import (
	"os"
	"syscall"

	log "github.com/sirupsen/logrus"
)

// DefaultMounts will return the defult mounts
func DefaultMounts() *Mounts {
	m := &Mounts{}

	// bin Mount
	bin := Mount{
		CreateMount: false,
		EnableMount: false,
		Name:        "bin",
		Path:        "/bin",
		Mode:        0777,
	}
	m.Mount = append(m.Mount, bin)

	//
	dev := Mount{
		CreateMount: false,
		EnableMount: false,
		Name:        "dev",
		Source:      "devtmpfs",
		Path:        "/dev",
		FSType:      "devtmpfs",
		Flags:       syscall.MS_MGC_VAL,
		Mode:        0777,
	}
	m.Mount = append(m.Mount, dev)

	//
	etc := Mount{
		CreateMount: false,
		EnableMount: false,
		Name:        "etc",
		Path:        "/etc",
		Mode:        0777,
	}
	m.Mount = append(m.Mount, etc)

	//
	home := Mount{
		CreateMount: false,
		EnableMount: false,
		Name:        "home",
		Path:        "/home",
		Mode:        0777,
	}
	m.Mount = append(m.Mount, home)

	//
	mnt := Mount{
		CreateMount: false,
		EnableMount: false,
		Name:        "mnt",
		Path:        "/mnt",
		Mode:        0777,
	}
	m.Mount = append(m.Mount, mnt)

	//
	proc := Mount{
		CreateMount: false,
		EnableMount: false,
		Name:        "proc",
		Source:      "proc",
		Path:        "/proc",
		FSType:      "proc",
		Mode:        0777,
	}
	m.Mount = append(m.Mount, proc)

	//
	sys := Mount{
		CreateMount: false,
		EnableMount: false,
		Name:        "sys",
		Path:        "/sys",
		Mode:        0777,
	}
	m.Mount = append(m.Mount, sys)

	//
	tmp := Mount{
		CreateMount: false,
		EnableMount: false,
		Name:        "tmp",
		Source:      "tmpfs",
		Path:        "/tmp",
		FSType:      "tmpfs",
		Mode:        0777,
	}
	m.Mount = append(m.Mount, tmp)

	//
	usr := Mount{
		CreateMount: false,
		EnableMount: false,
		Name:        "usr",
		Path:        "/usr",
		Mode:        0777,
	}
	m.Mount = append(m.Mount, usr)

	return m
}

// CreateFolder -
func (m *Mounts) CreateFolder() error {

	for x := range m.Mount {
		if m.Mount[x].CreateMount == true {
			err := os.MkdirAll(m.Mount[x].Path, m.Mount[x].Mode)
			if err != nil {
				log.Errorf("Folder[%s] create error [%v]", m.Mount[x].Path, err)
			}
		}
	}
	return nil
}

// CreateMount -
func (m *Mounts) CreateMount() error {
	for x := range m.Mount {
		if m.Mount[x].CreateMount == true {
			err := syscall.Mount(m.Mount[x].Name, m.Mount[x].Path, m.Mount[x].FSType, m.Mount[x].Flags, m.Mount[x].Options)
			if err != nil {
				log.Errorf("Mount [%s] create error [%v]", m.Mount[x].Name, err)
			}
			log.Infof("Mounted [%s] -> [%s]", m.Mount[x].Name, m.Mount[x].Path)
		}
	}
	return nil
}

// CreateNamedMount -
func (m *Mounts) CreateNamedMount(name string, remove bool) error {
	for x := range m.Mount {
		if m.Mount[x].Name == name && m.Mount[x].CreateMount == true {
			err := syscall.Mount(m.Mount[x].Name, m.Mount[x].Path, m.Mount[x].FSType, m.Mount[x].Flags, m.Mount[x].Options)
			if err != nil {
				log.Errorf("Mount [%s] create error [%v]", m.Mount[x].Name, err)
			}
			// Remove this element
			if remove {
				m.Mount = append(m.Mount[:x], m.Mount[x+1:]...)
			}
			log.Infof("Mounted [%s] -> [%s]", m.Mount[x].Name, m.Mount[x].Path)
			return nil
		}
	}
	return nil
}

// GetMount -
func (m *Mounts) GetMount(name string) *Mount {

	for x := range m.Mount {
		if m.Mount[x].Name == name {
			return &m.Mount[x]
		}
	}
	return nil
}
