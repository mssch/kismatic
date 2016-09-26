package inspector

import (
	"errors"
	"fmt"
	"testing"
)

func runMock(out string, err error) func(string, ...string) ([]byte, error) {
	return func(string, ...string) ([]byte, error) {
		return []byte(out), err
	}
}

func TestRPMPackageManager(t *testing.T) {
	out := `Installed Packages
potentiallySomeGarbageData
NetworkManager.x86_64           1:1.0.6-30.el7_2               @koji-override-1`
	m := rpmManager{
		run: runMock(out, nil),
	}
	p := packageQuery{"NetworkManager.x86_64", "1:1.0.6-30.el7_2"}
	ok, _ := m.isInstalled(p)
	if !ok {
		t.Error("expected true, but got false")
	}
}

func TestRPMPackageManagerPackageNotFound(t *testing.T) {
	out := `Error: No matching Packages to list`
	m := rpmManager{
		run: runMock(out, nil),
	}
	p := packageQuery{"NonExistent", "1.0"}
	ok, err := m.isInstalled(p)
	if ok {
		t.Error("expected false, but got true")
	}
	if err != nil {
		t.Errorf("got an unexpected error: %v", err)
	}
}

func TestRPMPackageManagerPackageNotFound2(t *testing.T) {
	out := `Installed Packages
NetworkManager.x86_64           1:1.0.6-30.el7_2               @koji-override-1`
	m := rpmManager{
		run: runMock(out, nil),
	}
	p := packageQuery{"NonExistent", "1.0"}
	ok, err := m.isInstalled(p)
	if ok {
		t.Error("expected false, but got true")
	}
	if err != nil {
		t.Errorf("got an unexpected error: %v", err)
	}
}

func TestRPMPackageManagerExecError(t *testing.T) {
	m := rpmManager{
		run: runMock("", fmt.Errorf("some error")),
	}
	p := packageQuery{"SomePkg", "1.0"}
	ok, err := m.isInstalled(p)
	if ok {
		t.Error("expected false, but got true")
	}
	if err == nil {
		t.Error("expected an error, but didn't get one")
	}

	ok, err = m.isAvailable(p)
	if ok {
		t.Error("expected false, but got true")
	}
	if err == nil {
		t.Error("expected an error, but didn't get one")
	}
}

func TestRPMPackageManagerIsAvailable(t *testing.T) {
	out := `Available Packages
389-ds-base.x86_64                        1.3.4.0-33.el7_2               updates
389-ds-base-devel.x86_64                  1.3.4.0-33.el7_2               updates
389-ds-base-libs.x86_64                   1.3.4.0-33.el7_2               updates
Cython.x86_64                             0.19-3.el7                     base
ElectricFence.i686                        2.2.2-39.el7                   base`
	m := rpmManager{
		run: runMock(out, nil),
	}
	p := packageQuery{"Cython.x86_64", "0.19-3.el7"}
	ok, _ := m.isAvailable(p)
	if !ok {
		t.Error("expected true, but got false")
	}
}

func TestRPMPackageManagerIsNotAvailable(t *testing.T) {
	out := `Available Packages
389-ds-base.x86_64                        1.3.4.0-33.el7_2               updates
389-ds-base-devel.x86_64                  1.3.4.0-33.el7_2               updates
389-ds-base-libs.x86_64                   1.3.4.0-33.el7_2               updates
Cython.x86_64                             0.19-3.el7                     base
ElectricFence.i686                        2.2.2-39.el7                   base`
	m := rpmManager{
		run: runMock(out, nil),
	}
	p := packageQuery{"NonExistent", "1.0.0"}
	ok, _ := m.isAvailable(p)
	if ok {
		t.Error("expected false, but got true")
	}
}

func TestDebPackageManager(t *testing.T) {
	out := `Desired=Unknown/Install/Remove/Purge/Hold
| Status=Not/Inst/Conf-files/Unpacked/halF-conf/Half-inst/trig-aWait/Trig-pend
|/ Err?=(none)/Reinst-required (Status,Err: uppercase=bad)
||/ Name                                                  Version                         Architecture                    Description
+++-=====================================================-===============================-===============================-================================================================================================================
ii  libc6:amd64                                           2.23-0ubuntu3                   amd64                           GNU C Library: Shared libraries`
	m := debManager{
		run: runMock(out, nil),
	}
	p := packageQuery{"libc6:amd64", "2.23-0ubuntu3"}
	ok, _ := m.isInstalled(p)
	if !ok {
		t.Errorf("expected true, but got false")
	}
}

func TestDebPackageManagerPackageNotFound(t *testing.T) {
	m := debManager{
		run: runMock("dpkg-query: no packages found matching libc6a", nil),
	}
	p := packageQuery{"", ""}
	ok, err := m.isInstalled(p)
	if ok {
		t.Errorf("expected false, but got true")
	}
	if err != nil {
		t.Errorf("unexpected error returned: %v", err)
	}
}

func TestDebPackageManagerExecError(t *testing.T) {
	m := debManager{
		run: runMock("", errors.New("Some error happened...")),
	}
	p := packageQuery{"", ""}
	ok, err := m.isInstalled(p)
	if ok {
		t.Error("expected false, but got true")
	}
	if err == nil {
		t.Error("expected an error, but didn't get one")
	}

	ok, err = m.isAvailable(p)
	if ok {
		t.Error("expected false, but got true")
	}
	if err == nil {
		t.Error("expected an error, but didn't get one")
	}
}

func TestDebPackageManagerIsAvailable(t *testing.T) {
	m := debManager{
		run: runMock("we don't really care about the output, just the exit status", nil),
	}
	p := packageQuery{"", ""}
	ok, _ := m.isAvailable(p)
	if !ok {
		t.Error("expected true, but got false")
	}
}

func TestDebPackageManagerIsNotAvailable(t *testing.T) {
	m := debManager{
		run: runMock("", errors.New("package not found")),
	}
	p := packageQuery{"", ""}
	ok, err := m.isAvailable(p)
	if ok {
		t.Error("expected false, but got true")
	}
	if err == nil {
		t.Error("expected an error, but didn't get one")
	}
}
