package cli

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/apprenda/kismatic/pkg/install"
)

func TestDashboardCmdMissingPlan(t *testing.T) {
	out := &bytes.Buffer{}
	fp := &fakePlanner{
		exists: false,
	}
	opts := &dashboardOpts{
		planFilename:     "planFile",
		dashboardURLMode: true,
	}
	if err := doDashboard(out, fp, opts); err == nil {
		t.Errorf("dashboard did not return an error when the plan does not exist")
	}
}

func TestDashboardCmdEmptyAddress(t *testing.T) {
	plan := install.Plan{}
	_, err := getDashboardURL(plan)
	if err == nil {
		t.Errorf("dashboard did not return an error when LoadBalancedFQDN is empty")
	}
}

func TestGetDashboardURL(t *testing.T) {
	plan := install.Plan{
		Cluster: install.Cluster{
			AdminPassword: "thePassword",
		},
		Master: install.MasterNodeGroup{
			LoadBalancedFQDN: "cluster.apprenda.local",
		},
	}
	url, err := getDashboardURL(plan)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if strings.Contains(url, plan.Cluster.AdminPassword) {
		t.Errorf("dashboard url contains admin password")
	}
}

func TestGetAuthenticatedDashboardURL(t *testing.T) {
	plan := install.Plan{
		Cluster: install.Cluster{
			AdminPassword: "thePassword",
		},
		Master: install.MasterNodeGroup{
			LoadBalancedFQDN: "cluster.apprenda.local",
		},
	}
	url, err := getAuthenticatedDashboardURL(plan)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !strings.Contains(url, plan.Cluster.AdminPassword) {
		t.Errorf("authenticated dashboard url does not contain admin password")
	}
}

type timeoutHandler struct {
}

func (h timeoutHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	delay, err := strconv.Atoi(strings.Trim(req.URL.Path, "/"))
	if err != nil {
		fmt.Printf("could not parse delay: %s, %v", req.URL.Path, err)
	}
	time.Sleep(time.Duration(delay) * time.Second)
}

func TestVerifyDashboardConnectivity(t *testing.T) {
	server := httptest.NewServer(timeoutHandler{})
	defer server.Close()
	if err := verifyDashboardConnectivity(server.URL + "/1"); err != nil {
		t.Errorf("dashboard returned an error %v", err)
	}
}

func TestVerifyDashboardConnectivityShouldTimeout(t *testing.T) {
	server := httptest.NewServer(timeoutHandler{})
	defer server.Close()
	if err := verifyDashboardConnectivity(server.URL + "/3"); err == nil {
		t.Errorf("ip returned an error %v", err)
	}
}
