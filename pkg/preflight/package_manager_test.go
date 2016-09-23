package preflight

import (
	"fmt"
	"testing"
)

func TestYumPackageManager(t *testing.T) {
	m := yumManager{
		run: func(string, ...string) ([]byte, error) {
			out := `Installed Packages
potentiallySomeGarbageData
NetworkManager.x86_64           1:1.0.6-30.el7_2               @koji-override-1`
			return []byte(out), nil
		},
	}
	p := pkg{"NetworkManager.x86_64", "1:1.0.6-30.el7_2"}
	ok, _ := m.isInstalled(p)
	if !ok {
		t.Error("expected true, but got false")
	}
}

func TestYumPackageManagerPackageNotFound(t *testing.T) {
	m := yumManager{
		run: func(string, ...string) ([]byte, error) {
			out := `Error: No matching Packages to list`
			return []byte(out), nil
		},
	}
	p := pkg{"NonExistent", "1.0"}
	ok, err := m.isInstalled(p)
	if ok {
		t.Error("expected false, but got true")
	}
	if err != nil {
		t.Errorf("got an unexpected error: %v", err)
	}
}

func TestYumPackageManagerPackageNotFound2(t *testing.T) {
	m := yumManager{
		run: func(string, ...string) ([]byte, error) {
			out := `Installed Packages
NetworkManager.x86_64           1:1.0.6-30.el7_2               @koji-override-1`
			return []byte(out), nil
		},
	}
	p := pkg{"NonExistent", "1.0"}
	ok, err := m.isInstalled(p)
	if ok {
		t.Error("expected false, but got true")
	}
	if err != nil {
		t.Errorf("got an unexpected error: %v", err)
	}
}

func TestYumPackageManagerExecError(t *testing.T) {
	m := yumManager{
		run: func(string, ...string) ([]byte, error) {
			return nil, fmt.Errorf("some error")
		},
	}
	p := pkg{"SomePkg", "1.0"}
	ok, err := m.isInstalled(p)
	if ok {
		t.Error("expected false, but got true")
	}
	if err == nil {
		t.Error("expected an error, but didn't get one")
	}
}

func TestAptPackageManagerExecError(t *testing.T) {
	m := aptManager{
		run: func(string, ...string) ([]byte, error) {
			out := `Desired=Unknown/Install/Remove/Purge/Hold
| Status=Not/Inst/Conf-files/Unpacked/halF-conf/Half-inst/trig-aWait/Trig-pend
|/ Err?=(none)/Reinst-required (Status,Err: uppercase=bad)
||/ Name                                                  Version                         Architecture                    Description
+++-=====================================================-===============================-===============================-================================================================================================================
ii  libc6:amd64                                           2.23-0ubuntu3                   amd64                           GNU C Library: Shared libraries`
			return []byte(out), nil
		},
	}
	p := pkg{"libc6:amd64", "2.23-0ubuntu3"}
	ok, _ := m.isInstalled(p)
	if !ok {
		t.Errorf("expected true, but got false")
	}
}
