package cli

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/apprenda/kismatic/pkg/install"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

type dashboardOpts struct {
	dashboardURLMode   bool
	tokenOnly          bool
	generatedAssetsDir string
	planFilename       string
}

const url = "http://localhost:8001/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy/#!/login"

// NewCmdDashboard opens or displays the dashboard URL
func NewCmdDashboard(in io.Reader, out io.Writer) *cobra.Command {
	opts := &dashboardOpts{}

	cmd := &cobra.Command{
		Use:   "dashboard",
		Short: "Opens/displays the kubernetes dashboard URL of the cluster",
		Long: `
  This command is a convenience command to open the dashbord.
  - Retrieves the token of the secret for ServiceAccount 'kubernetes-dashboard-admin':
  -----------------------------------------------------------------------------------------------------------------------------------------------------
  export SECRET="$(./kubectl get sa kubernetes-dashboard-admin -o 'jsonpath={.secrets[0].name}' -n kube-system --kubeconfig generated/kubeconfig)"
  ./kubectl describe secrets $SECRET -n kube-system --kubeconfig generated/kubeconfig | awk '$1=="token:"{print $2}'
  -----------------------------------------------------------------------------------------------------------------------------------------------------
	
  - Runs 'kubectl proxy':
  -----------------------------------------------------------------------------------------------------------------------------------------------------
  kubectl proxy --kubeconfig generated/kubeconfig
  -----------------------------------------------------------------------------------------------------------------------------------------------------`,
		Example: `  -----------------------------------------------------------------------------------------------------------------------------------------------------
  ./kismatic dashboard
  Opening kubernetes dashboard in default browser...
  Use the kubeconfig in "generated/dashboard-admin-kubeconfig"
  Starting to serve on 127.0.0.1:8001
  -----------------------------------------------------------------------------------------------------------------------------------------------------
  ./kismatic dashboard --url
  http://localhost:8001/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy/#!/login
  -----------------------------------------------------------------------------------------------------------------------------------------------------
  ./kismatic dashboard --token
  Generated kubeconfig in "generated/dashboard-admin-kubeconfig"
  -----------------------------------------------------------------------------------------------------------------------------------------------------`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return fmt.Errorf("Unexpected args: %v", args)
			}
			return doDashboard(in, out, opts)
		},
	}

	cmd.Flags().StringVar(&opts.generatedAssetsDir, "generated-assets-dir", "generated", "path to the directory where assets generated during the installation process will be stored")
	cmd.Flags().BoolVar(&opts.dashboardURLMode, "url", false, "Display the kubernetes dashboard URL instead of opening it in the default browser")
	cmd.Flags().BoolVar(&opts.tokenOnly, "token", false, "Do not open the dashboard, only generate a kubeconfig file with the admin token")
	addPlanFileFlag(cmd.PersistentFlags(), &opts.planFilename)
	return cmd
}

func doDashboard(in io.Reader, out io.Writer, opts *dashboardOpts) error {
	if opts.dashboardURLMode {
		fmt.Fprintln(out, url)
		return nil
	}

	kubeconfig := filepath.Join(opts.generatedAssetsDir, "kubeconfig")
	if stat, err := os.Stat(kubeconfig); os.IsNotExist(err) || stat.IsDir() {
		return fmt.Errorf("Did not find required kubeconfig file %q", kubeconfig)
	}

	var generateErr error
	adminKubeconfig := filepath.Join(opts.generatedAssetsDir, "dashboard-admin-kubeconfig")
	// Generate dashboard admin certificate if it does not exist
	if _, err := os.Stat(adminKubeconfig); os.IsNotExist(err) {
		planner := &install.FilePlanner{File: opts.planFilename}
		plan, err := planner.Read()
		if err != nil {
			return fmt.Errorf("Error reading plan file: %v", err)
		}

		generateErr = generateDashboardAdminKubeconfig(out, opts.generatedAssetsDir, *plan)
	}

	if opts.tokenOnly {
		if generateErr != nil {
			return fmt.Errorf("Error generating a kubeconfig file: %v", generateErr)
		}
		fmt.Fprintf(out, "Generated kubeconfig in %q\n", adminKubeconfig)
		return nil
	}

	if generateErr != nil {
		fmt.Fprintf(out, "Error generating a kubeconfig file, you may still use the dashboard with your own ServiceAccount token\n\n")
	}

	fmt.Fprintf(out, "Opening kubernetes dashboard in default browser...\n")
	if generateErr == nil {
		fmt.Fprintf(out, "Use the kubeconfig in %q\n", adminKubeconfig)
	}
	if err := browser.OpenURL(url); err != nil {
		fmt.Fprintf(out, "Unexpected error opening the kubernetes dashboard: %v. You may access it at %q", err, url)
	}

	cmd := exec.Command("./kubectl", "proxy", "--kubeconfig", kubeconfig)
	cmd.Stdin = in
	cmd.Stdout = out
	cmd.Stderr = out
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Error running kubectl proxy: %v", err)
	}

	return nil
}

func generateDashboardAdminKubeconfig(out io.Writer, generatedAssetsDir string, plan install.Plan) error {
	// All of this is required because cannot set a label on the secret so no selectors
	cmd := exec.Command("./kubectl", "-n", "kube-system", "get", "sa", "kubernetes-dashboard-admin", "-o", "jsonpath={.secrets[0].name}", "--kubeconfig", filepath.Join(generatedAssetsDir, "kubeconfig"))
	sa, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("error getting token secret: %v", err)
	}
	if len(sa) == 0 || !strings.Contains(string(sa), "kubernetes-dashboard-admin-token") {
		return fmt.Errorf("kubernetes-dashboard-admin-token secret not found")
	}

	cmd = exec.Command("./kubectl", "-n", "kube-system", "get", "secrets", string(sa), "-o", "jsonpath={.data.token}", "--kubeconfig", filepath.Join(generatedAssetsDir, "kubeconfig"))
	token, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("error getting the token: %v", err)
	}
	if len(token) == 0 {
		return fmt.Errorf("got an empty token")
	}
	err = install.GenerateDashboardAdminKubeconfig(strings.Trim(string(token), "'"), &plan, generatedAssetsDir)
	if err != nil {
		return fmt.Errorf("error generating dashboard-admin kubeconfig file: %v", err)
	}
	return nil
}
