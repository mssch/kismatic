package inspector

import (
	"strings"
	"testing"
)

func TestDetectDistroFromOSRelease(t *testing.T) {
	tests := []struct {
		osReleaseFile  string
		expectedDistro Distro
		expectErr      bool
	}{
		{
			osReleaseFile:  centos7ReleaseFile,
			expectedDistro: CentOS,
			expectErr:      false,
		},
		{
			osReleaseFile:  rhel7ReleaseFile,
			expectedDistro: RHEL,
			expectErr:      false,
		},
		{
			osReleaseFile:  ubuntu1604ReleaseFile,
			expectedDistro: Ubuntu,
			expectErr:      false,
		},
		{
			osReleaseFile:  "",
			expectedDistro: Unsupported,
			expectErr:      true,
		},
		{
			osReleaseFile:  missingIDFieldOSReleaseFile,
			expectedDistro: Unsupported,
			expectErr:      true,
		},
	}

	for _, test := range tests {
		d, err := detectDistroFromOSRelease(strings.NewReader(test.osReleaseFile))
		if test.expectErr && err == nil {
			t.Error("expected an error, but didn't get one")
		}

		if !test.expectErr && err != nil {
			t.Errorf("unexpected error occurred when running test: %v", err)
		}

		if d != test.expectedDistro {
			t.Errorf("failed to detect distro. expected %s, found %s", test.expectedDistro, d)
		}
	}
}

var centos7ReleaseFile = `NAME="CentOS Linux"
VERSION="7 (Core)"
ID="centos"
ID_LIKE="rhel fedora"
VERSION_ID="7"
PRETTY_NAME="CentOS Linux 7 (Core)"
ANSI_COLOR="0;31"
CPE_NAME="cpe:/o:centos:centos:7"
HOME_URL="https://www.centos.org/"
BUG_REPORT_URL="https://bugs.centos.org/"

CENTOS_MANTISBT_PROJECT="CentOS-7"
CENTOS_MANTISBT_PROJECT_VERSION="7"
REDHAT_SUPPORT_PRODUCT="centos"
REDHAT_SUPPORT_PRODUCT_VERSION="7"`

var rhel7ReleaseFile = `NAME="Red Hat Enterprise Linux Server"
VERSION="7.2 (Maipo)"
ID="rhel"
ID_LIKE="fedora"
VERSION_ID="7.2"
PRETTY_NAME="Red Hat Enterprise Linux Server 7.2 (Maipo)"
ANSI_COLOR="0;31"
CPE_NAME="cpe:/o:redhat:enterprise_linux:7.2:GA:server"
HOME_URL="https://www.redhat.com/"
BUG_REPORT_URL="https://bugzilla.redhat.com/"

REDHAT_BUGZILLA_PRODUCT="Red Hat Enterprise Linux 7"
REDHAT_BUGZILLA_PRODUCT_VERSION=7.2
REDHAT_SUPPORT_PRODUCT="Red Hat Enterprise Linux"
REDHAT_SUPPORT_PRODUCT_VERSION="7.2"`

var ubuntu1604ReleaseFile = `NAME="Ubuntu"
VERSION="16.04.1 LTS (Xenial Xerus)"
ID=ubuntu
ID_LIKE=debian
PRETTY_NAME="Ubuntu 16.04.1 LTS"
VERSION_ID="16.04"
HOME_URL="http://www.ubuntu.com/"
SUPPORT_URL="http://help.ubuntu.com/"
BUG_REPORT_URL="http://bugs.launchpad.net/ubuntu/"
UBUNTU_CODENAME=xenial`

var missingIDFieldOSReleaseFile = `NAME="Ubuntu"
VERSION="16.04.1 LTS (Xenial Xerus)"
ID_LIKE=debian
PRETTY_NAME="Ubuntu 16.04.1 LTS"
VERSION_ID="16.04"
HOME_URL="http://www.ubuntu.com/"
SUPPORT_URL="http://help.ubuntu.com/"
BUG_REPORT_URL="http://bugs.launchpad.net/ubuntu/"
UBUNTU_CODENAME=xenial`
