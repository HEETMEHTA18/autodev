// Package catalog loads the YAML-driven package registry and resolves
// dependency graphs. Adding a new tool means adding a YAML entry only.
package catalog

import (
	"context"
	_ "embed"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

//go:embed catalog.yaml
var catalogData []byte

// InstallSteps holds platform-specific installation instructions.
type InstallSteps struct {
	Method   string   `yaml:"method"`   // apt | brew | winget | script | npm | pip | cargo
	Packages []string `yaml:"packages"` // for apt/brew/winget/npm/pip
	Script   []string `yaml:"script"`   // for script method
}

// Package represents one installable entry in the catalog.
type Package struct {
	ID          string   `yaml:"id"`
	Name        string   `yaml:"name"`
	Category    string   `yaml:"category"`
	Description string   `yaml:"description"`
	Icon        string   `yaml:"icon"`
	Deps        []string `yaml:"deps"`
	PostInstall []string `yaml:"post_install"`
	Verify      string   `yaml:"verify"`
	Install     struct {
		Linux   InstallSteps `yaml:"linux"`
		Darwin  InstallSteps `yaml:"darwin"`
		Windows InstallSteps `yaml:"windows"`
	} `yaml:"install"`
}

// Profile is a named set of packages for a developer role.
type Profile struct {
	ID          string   `yaml:"id"`
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Icon        string   `yaml:"icon"`
	Packages    []string `yaml:"packages"`
}

// Catalog is the parsed registry.
type Catalog struct {
	Packages []Package `yaml:"packages"`
	Profiles []Profile `yaml:"profiles"`
	byID     map[string]*Package
}

// Load parses the embedded catalog.yaml.
func Load() (*Catalog, error) {
	var c Catalog
	if err := yaml.Unmarshal(catalogData, &c); err != nil {
		return nil, fmt.Errorf("failed to parse catalog: %w", err)
	}
	c.byID = make(map[string]*Package, len(c.Packages))
	for i := range c.Packages {
		c.byID[c.Packages[i].ID] = &c.Packages[i]
	}
	return &c, nil
}

// GetPackage returns a package by ID.
func (c *Catalog) GetPackage(id string) (*Package, bool) {
	p, ok := c.byID[id]
	return p, ok
}

// GetProfile returns a profile by ID.
func (c *Catalog) GetProfile(id string) (*Profile, bool) {
	for i := range c.Profiles {
		if c.Profiles[i].ID == id {
			return &c.Profiles[i], true
		}
	}
	return nil, false
}

// ByCategory returns packages grouped by category in order.
func (c *Catalog) ByCategory() map[string][]*Package {
	result := make(map[string][]*Package)
	for i := range c.Packages {
		p := &c.Packages[i]
		result[p.Category] = append(result[p.Category], p)
	}
	return result
}

// CategoryOrder is the preferred display order of categories.
var CategoryOrder = []string{
	"Languages",
	"Frameworks",
	"Databases",
	"DevOps",
	"Mobile",
	"AI/ML",
	"Tools",
}

// Resolve returns the full ordered install list for a set of package IDs,
// including all transitive dependencies (topological sort).
func (c *Catalog) Resolve(ids []string) ([]*Package, error) {
	visited := map[string]bool{}
	var order []*Package

	var visit func(id string) error
	visit = func(id string) error {
		if visited[id] {
			return nil
		}
		visited[id] = true
		pkg, ok := c.byID[id]
		if !ok {
			return fmt.Errorf("unknown package: %q — check catalog.yaml", id)
		}
		for _, dep := range pkg.Deps {
			if err := visit(dep); err != nil {
				return err
			}
		}
		order = append(order, pkg)
		return nil
	}

	for _, id := range ids {
		if err := visit(id); err != nil {
			return nil, err
		}
	}
	return order, nil
}

// ResolveProfile returns the resolved install plan for a profile.
func (c *Catalog) ResolveProfile(profileID string) ([]*Package, error) {
	prof, ok := c.GetProfile(profileID)
	if !ok {
		return nil, fmt.Errorf("unknown profile: %q", profileID)
	}
	return c.Resolve(prof.Packages)
}

// IsInstalled checks if the package is already installed on the system using its Verify command.
func (p *Package) IsInstalled() bool {
	if p.Verify == "" {
		return false
	}
	parts := strings.Fields(p.Verify)
	if len(parts) == 0 {
		return false
	}

	// Check if the binary executable exists in PATH
	_, err := exec.LookPath(parts[0])
	if err != nil {
		return false
	}

	// Reject complex shell expressions in the verify command. If the verify
	// string contains pipes, redirection, command substitution or boolean
	// operators, we avoid running it automatically to reduce command-injection
	// risk and return false so the caller can perform manual verification.
	unsafeTokens := []string{"|", "$(", "&&", "||", ";", "`", ">", "<"}
	for _, t := range unsafeTokens {
		if strings.Contains(p.Verify, t) {
			return false
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
