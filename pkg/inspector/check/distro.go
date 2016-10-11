package check

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
)

const (
	Ubuntu      Distro = "ubuntu"
	RHEL        Distro = "rhel"
	CentOS      Distro = "centos"
	Darwin      Distro = "darwin"
	Unsupported Distro = ""
)

// Distro is a Linux distribution that the inspector supports
type Distro string

// DetectDistro uses the /etc/os-release file to get distro information.
func DetectDistro() (Distro, error) {
	if runtime.GOOS == "darwin" {
		return Darwin, nil
	}
	f, err := os.Open("/etc/os-release")
	if err != nil {
		return Unsupported, fmt.Errorf("error reading /etc/os-release file: %v", err)
	}
	defer f.Close()
	return detectDistroFromOSRelease(f)
}

func detectDistroFromOSRelease(r io.Reader) (Distro, error) {
	idLine := ""
	s := bufio.NewScanner(r)
	for s.Scan() {
		l := s.Text()
		if strings.HasPrefix(l, "ID=") {
			idLine = l
			break
		}
	}
	if idLine == "" {
		return Unsupported, errors.New("/etc/os-release file does not contain ID= field")
	}
	fields := strings.Split(idLine, "=")
	if len(fields) != 2 {
		return Unsupported, fmt.Errorf("Unknown format of /etc/os-release file. ID line was: %s", idLine)
	}

	// Remove double-quotes from field value
	switch strings.Replace(fields[1], "\"", "", -1) {
	case "centos":
		return CentOS, nil
	case "rhel":
		return RHEL, nil
	case "ubuntu":
		return Ubuntu, nil
	default:
		return Unsupported, fmt.Errorf("Unsupported distribution detected: %s", fields[1])
	}
}
