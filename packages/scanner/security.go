package scanner

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type AuditPackage struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	Ecosystem string `json:"ecosystem"` // npm, PyPI, Go, etc.
}

type Vulnerability struct {
	ID       string   `json:"id"`
	Summary  string   `json:"summary"`
	Details  string   `json:"details"`
	Severity string   `json:"severity"`
	Aliases  []string `json:"aliases"`
}

type AuditResult struct {
	Package         AuditPackage    `json:"package"`
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
}

// CleanVersion strips common prefix characters from version string
func CleanVersion(v string) string {
	v = strings.TrimSpace(v)
	v = strings.TrimLeft(v, "^~>=<v ")
	// Split by space or comma if there are ranges (e.g., ">=1.0.0 <2.0.0") and take first part
	parts := strings.FieldsFunc(v, func(r rune) bool {
		return r == ' ' || r == ','
	})
	if len(parts) > 0 {
		return parts[0]
	}
	return v
}

func parseGoMod(path string) []AuditPackage {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var pkgs []AuditPackage
	lines := strings.Split(string(data), "\n")

	inRequireBlock := false
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}
		if line == "require (" {
			inRequireBlock = true
			continue
		}
		if line == ")" && inRequireBlock {
			inRequireBlock = false
			continue
		}

		if strings.HasPrefix(line, "require ") && !strings.Contains(line, "(") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				name := parts[1]
				version := CleanVersion(parts[2])
				if version != "" {
					pkgs = append(pkgs, AuditPackage{Name: name, Version: version, Ecosystem: "Go"})
				}
			}
		} else if inRequireBlock {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				name := parts[0]
				version := CleanVersion(parts[1])
				if version != "" {
					pkgs = append(pkgs, AuditPackage{Name: name, Version: version, Ecosystem: "Go"})
				}
			}
		}
	}
	return pkgs
}

func parseRequirements(path string) []AuditPackage {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var pkgs []AuditPackage
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		var name, version string
		if idx := strings.Index(line, "=="); idx != -1 {
			name = strings.TrimSpace(line[:idx])
			version = CleanVersion(line[idx+2:])
		} else if idx := strings.Index(line, ">="); idx != -1 {
			name = strings.TrimSpace(line[:idx])
			version = CleanVersion(line[idx+2:])
		} else if idx := strings.Index(line, "<="); idx != -1 {
			name = strings.TrimSpace(line[:idx])
			version = CleanVersion(line[idx+2:])
		} else {
			continue
		}

		if idx := strings.Index(name, ";"); idx != -1 {
			name = strings.TrimSpace(name[:idx])
		}
		if idx := strings.Index(version, ";"); idx != -1 {
			version = strings.TrimSpace(version[:idx])
		}

		if name != "" && version != "" {
			pkgs = append(pkgs, AuditPackage{Name: name, Version: version, Ecosystem: "PyPI"})
		}
	}
	return pkgs
}

func parsePackageJSON(path string) []AuditPackage {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var pkg map[string]json.RawMessage
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil
	}

	var pkgs []AuditPackage
	extract := func(key string) {
		if raw, ok := pkg[key]; ok {
			var deps map[string]string
			if err := json.Unmarshal(raw, &deps); err == nil {
				for name, version := range deps {
					cleanV := CleanVersion(version)
					if cleanV != "" && cleanV != "latest" && cleanV != "*" {
						pkgs = append(pkgs, AuditPackage{Name: name, Version: cleanV, Ecosystem: "npm"})
					}
				}
			}
		}
	}
	extract("dependencies")
	extract("devDependencies")
	return pkgs
}

type osvQuery struct {
	Package struct {
		Name      string `json:"name"`
		Ecosystem string `json:"ecosystem"`
	} `json:"package"`
	Version string `json:"version"`
}

type osvResponse struct {
	Vulns []struct {
		ID               string   `json:"id"`
		Summary          string   `json:"summary"`
		Details          string   `json:"details"`
		Aliases          []string `json:"aliases"`
		DatabaseSpecific struct {
			Severity string `json:"severity"`
		} `json:"database_specific"`
	} `json:"vulns"`
}

func CheckPackageVulnerabilities(ctx context.Context, client *http.Client, pkg AuditPackage) ([]Vulnerability, error) {
	query := osvQuery{
		Version: pkg.Version,
	}
	query.Package.Name = pkg.Name
	query.Package.Ecosystem = pkg.Ecosystem

	reqBody, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.osv.dev/v1/query", bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
		return nil, fmt.Errorf("OSV API returned status %d: %s", resp.StatusCode, string(body))
	}

	var res osvResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	var vulns []Vulnerability
	for _, v := range res.Vulns {
		sev := v.DatabaseSpecific.Severity
		if sev == "" {
			sev = "MODERATE"
		}
		vulns = append(vulns, Vulnerability{
			ID:       v.ID,
			Summary:  v.Summary,
			Details:  v.Details,
			Severity: sev,
			Aliases:  v.Aliases,
		})
	}
	return vulns, nil
}

func AuditRepository(rootPath string) ([]AuditResult, error) {
	var pkgs []AuditPackage

	packageJsonPath := filepath.Join(rootPath, "package.json")
	if _, err := os.Stat(packageJsonPath); err == nil {
		pkgs = append(pkgs, parsePackageJSON(packageJsonPath)...)
	}

	goModPath := filepath.Join(rootPath, "go.mod")
	if _, err := os.Stat(goModPath); err == nil {
		pkgs = append(pkgs, parseGoMod(goModPath)...)
	}

	requirementsPath := filepath.Join(rootPath, "requirements.txt")
	if _, err := os.Stat(requirementsPath); err == nil {
		pkgs = append(pkgs, parseRequirements(requirementsPath)...)
	}

	if len(pkgs) == 0 {
		return nil, nil
	}

	concurrencyLimit := 10
	semaphore := make(chan struct{}, concurrencyLimit)
	var wg sync.WaitGroup
	var results []AuditResult
	var mu sync.Mutex

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	client := &http.Client{Timeout: 3 * time.Second}

	for _, pkg := range pkgs {
		wg.Add(1)
		go func(p AuditPackage) {
			defer wg.Done()
			select {
			case semaphore <- struct{}{}:
				defer func() { <-semaphore }()
			case <-ctx.Done():
				return
			}

			vulns, err := CheckPackageVulnerabilities(ctx, client, p)
			if err == nil && len(vulns) > 0 {
				mu.Lock()
				results = append(results, AuditResult{
					Package:         p,
					Vulnerabilities: vulns,
				})
				mu.Unlock()
			}
		}(pkg)
	}

	wg.Wait()
	return results, nil
}
