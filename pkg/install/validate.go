package install

import (
	"errors"
	"fmt"
	"net"
	"os"
	"time"
)

// ValidatePlan runs validation against the installation plan to ensure
// that the plan contains valid user input. Returns true, nil if the validation
// is successful. Otherwise, returns false and a collection of validation errors.
func ValidatePlan(p *Plan) (bool, []error) {
	v := newValidator()
	v.validate(p)
	return v.valid()
}

type validatable interface {
	validate() (bool, []error)
}

type validator struct {
	errs []error
}

func newValidator() *validator {
	return &validator{
		errs: []error{},
	}
}

func (v *validator) addError(err ...error) {
	v.errs = append(v.errs, err...)
}

func (v *validator) validate(obj validatable) {
	if ok, err := obj.validate(); !ok {
		v.addError(err...)
	}
}

func (v *validator) validateWithErrPrefix(prefix string, obj validatable) {
	if ok, err := obj.validate(); !ok {
		newErrs := make([]error, len(err), len(err))
		for i, err := range err {
			newErrs[i] = fmt.Errorf("%s: %v", prefix, err)
		}
		v.addError(newErrs...)
	}
}

func (v *validator) valid() (bool, []error) {
	if len(v.errs) > 0 {
		return false, v.errs
	}
	return true, nil
}

func (p *Plan) validate() (bool, []error) {
	v := newValidator()

	v.validate(&p.Cluster)
	v.validateWithErrPrefix("Etcd nodes", &p.Etcd)
	v.validateWithErrPrefix("Master nodes", &p.Master)
	v.validateWithErrPrefix("Worker nodes", &p.Worker)

	return v.valid()
}

func (c *Cluster) validate() (bool, []error) {
	v := newValidator()
	if c.Name == "" {
		v.addError(errors.New("Cluster name cannot be empty"))
	}
	if c.AdminPassword == "" {
		v.addError(errors.New("Admin password cannot be empty"))
	}
	v.validate(&c.Networking)
	v.validate(&c.Certificates)
	v.validate(&c.SSH)

	return v.valid()
}

func (n *NetworkConfig) validate() (bool, []error) {
	v := newValidator()
	if n.Type == "" {
		v.addError(errors.New("Networking type cannot be empty"))
	}
	if n.Type != "bridged" && n.Type != "overlay" {
		v.addError(fmt.Errorf("Invalid networking type %q was provided", n.Type))
	}
	if n.PodCIDRBlock == "" {
		v.addError(errors.New("Pod CIDR block cannot be empty"))
	}
	if _, _, err := net.ParseCIDR(n.PodCIDRBlock); n.PodCIDRBlock != "" && err != nil {
		v.addError(fmt.Errorf("Invalid Pod CIDR block provided: %v", err))
	}

	if n.ServiceCIDRBlock == "" {
		v.addError(errors.New("Service CIDR block cannot be empty"))
	}
	if _, _, err := net.ParseCIDR(n.ServiceCIDRBlock); n.ServiceCIDRBlock != "" && err != nil {
		v.addError(fmt.Errorf("Invalid Service CIDR block provided: %v", err))
	}
	return v.valid()
}

func (c *CertsConfig) validate() (bool, []error) {
	v := newValidator()
	if _, err := time.ParseDuration(c.Expiry); err != nil {
		v.addError(fmt.Errorf("Invalid certificate expiry %q provided: %v", c.Expiry, err))
	}

	if c.LocationCity == "" {
		v.addError(errors.New("Certificate location_city field is required"))
	}
	if c.LocationState == "" {
		v.addError(errors.New("Certificate location_state field is required"))
	}
	if c.LocationCountry == "" {
		v.addError(errors.New("Certificate location_country field is required"))
	}
	return v.valid()
}

func (s *SSHConfig) validate() (bool, []error) {
	v := newValidator()
	if s.User == "" {
		v.addError(errors.New("SSH user field is required"))
	}
	if s.Key == "" {
		v.addError(errors.New("SSH key field is required"))
	}
	if _, err := os.Stat(s.Key); os.IsNotExist(err) {
		v.addError(fmt.Errorf("SSH Key file was not found at %q", s.Key))
	}
	if s.Port < 1 || s.Port > 65535 {
		v.addError(fmt.Errorf("SSH port %d is invalid. Port must be in the range 1-65535", s.Port))
	}
	return v.valid()
}

func (ng *NodeGroup) validate() (bool, []error) {
	v := newValidator()
	if ng == nil || len(ng.Nodes) <= 0 {
		v.addError(fmt.Errorf("Atleast one node is required"))
	}
	if ng.ExpectedCount <= 0 {
		v.addError(fmt.Errorf("Node count must be greater than 0"))
	}
	if len(ng.Nodes) != ng.ExpectedCount && (len(ng.Nodes) > 0 && ng.ExpectedCount > 0) {
		v.addError(fmt.Errorf("Expected node count (%d) does not match the number of nodes provided (%d)", ng.ExpectedCount, len(ng.Nodes)))
	}
	for i, n := range ng.Nodes {
		v.validateWithErrPrefix(fmt.Sprintf("Node #%d", i+1), &n)
	}
	return v.valid()
}

func (mng *MasterNodeGroup) validate() (bool, []error) {
	v := newValidator()
	v.validate(&mng.NodeGroup)

	if mng.LoadBalancedFQDN == "" {
		v.addError(fmt.Errorf("Load balanced FQDN is required"))
	}

	if mng.LoadBalancedShortName == "" {
		v.addError(fmt.Errorf("Load balanced shortname is required"))
	}

	return v.valid()
}

func (n *Node) validate() (bool, []error) {
	v := newValidator()
	if n.Host == "" {
		v.addError(fmt.Errorf("Node host field is required"))
	}
	if n.IP == "" {
		v.addError(fmt.Errorf("Node IP field is required"))
	}
	if ip := net.ParseIP(n.IP); ip == nil {
		v.addError(fmt.Errorf("Invalid IP provided"))
	}
	return v.valid()
}
