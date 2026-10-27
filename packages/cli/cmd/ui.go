package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/autodev-sh/autodev/core/osinfo"
	"github.com/spf13/cobra"
)

func newUICmd() *cobra.Command {
	var port int

	cmd := &cobra.Command{
		Use:   "ui",
		Short: "Start the local AutoDev interactive web dashboard",
		Long:  `Launches a web-based dashboard on your local machine to check environment health, system specs, and manage installation states in a premium user interface.`,
		Example: `  autodev ui
  autodev ui --port 8080`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUI(port)
		},
	}

	cmd.Flags().IntVarP(&port, "port", "p", 8080, "port to run the web server on")
	return cmd
}

type uiToolCheck struct {
	Name    string `json:"name"`
	Cmd     string `json:"cmd"`
	Status  string `json:"status"` // OK | MISSING
	Version string `json:"version"`
	Hint    string `json:"hint"`
}

type uiStatusResponse struct {
	OS             string        `json:"os"`
	Arch           string        `json:"arch"`
	CPUCores       int           `json:"cpu_cores"`
	RAM            string        `json:"ram"`
	PackageManager string        `json:"package_manager"`
	Checks         []uiToolCheck `json:"checks"`
}

func runUI(port int) error {
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		// Port taken, fallback to random available port
		listener, err = net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			return fmt.Errorf("failed to start web server: %w", err)
		}
		addr = listener.Addr().String()
	}

	url := fmt.Sprintf("http://%s", addr)
	fmt.Printf("\n  AutoDev UI dashboard starting at %s ...\n", url)

	// Register Handlers
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/api/status", handleAPIStatus)

	// Start browser in a goroutine
	go func() {
		time.Sleep(500 * time.Millisecond)
		openBrowser(url)
	}()

	server := &http.Server{}
	return server.Serve(listener)
}

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	}
	_ = err
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, err := template.New("dashboard").Parse(dashboardHTML)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = tmpl.Execute(w, nil)
}

func handleAPIStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Detect system info
	info, _ := osinfo.Detect()

	// Tool checks run in parallel (800ms select limit to keep UI responsive)
	var wg sync.WaitGroup
	var mu sync.Mutex

	type rawCheck struct {
		name string
		cmd  string
		args []string
		hint string
	}

	rawChecks := []rawCheck{
		{name: "Git", cmd: "git", args: []string{"--version"}, hint: "https://git-scm.com"},
		{name: "Node.js", cmd: "node", args: []string{"--version"}, hint: "autodev install nodejs"},
		{name: "npm", cmd: "npm", args: []string{"--version"}, hint: "Comes with Node.js"},
		{name: "pnpm", cmd: "pnpm", args: []string{"--version"}, hint: "npm install -g pnpm"},
		{name: "yarn", cmd: "yarn", args: []string{"--version"}, hint: "npm install -g yarn"},
		{name: "Bun", cmd: "bun", args: []string{"--version"}, hint: "https://bun.sh"},
		{name: "Go", cmd: "go", args: []string{"version"}, hint: "autodev install go"},
		{name: "Python 3", cmd: "python3", args: []string{"--version"}, hint: "autodev install python"},
		{name: "pip", cmd: "pip3", args: []string{"--version"}, hint: "Comes with Python 3"},
		{name: "Rust", cmd: "rustc", args: []string{"--version"}, hint: "autodev install rust"},
		{name: "Docker", cmd: "docker", args: []string{"--version"}, hint: "autodev install docker"},
		{name: "docker compose", cmd: "docker", args: []string{"compose", "version"}, hint: "Upgrade Docker"},
		{name: "kubectl", cmd: "kubectl", args: []string{"version", "--client", "--short"}, hint: "autodev install kubectl"},
		{name: "Terraform", cmd: "terraform", args: []string{"version"}, hint: "autodev install terraform"},
		{name: "Flutter", cmd: "flutter", args: []string{"--version"}, hint: "autodev install flutter"},
		{name: "Java", cmd: "java", args: []string{"-version"}, hint: "autodev install java"},
		{name: "PHP", cmd: "php", args: []string{"--version"}, hint: "autodev install php"},
		{name: "Ruby", cmd: "ruby", args: []string{"--version"}, hint: "autodev install ruby"},
	}

	results := make(map[int]uiToolCheck)
	wg.Add(len(rawChecks))

	for i, c := range rawChecks {
		go func(idx int, rc rawCheck) {
			defer wg.Done()

			_, err := exec.LookPath(rc.cmd)
			if err != nil {
				mu.Lock()
				results[idx] = uiToolCheck{
					Name:    rc.name,
					Cmd:     rc.cmd,
					Status:  "MISSING",
					Version: "",
					Hint:    rc.hint,
				}
				mu.Unlock()
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), 700*time.Millisecond)
			defer cancel()

			cmd := exec.CommandContext(ctx, rc.cmd, rc.args...)
			out, err := cmd.CombinedOutput()

			versionStr := ""
			statusStr := "MISSING"
			if err == nil {
				versionStr = strings.TrimSpace(strings.Split(string(out), "\n")[0])
				if len(versionStr) > 40 {
					versionStr = versionStr[:40]
				}
				statusStr = "OK"
			}

			mu.Lock()
			results[idx] = uiToolCheck{
				Name:    rc.name,
				Cmd:     rc.cmd,
				Status:  statusStr,
				Version: versionStr,
				Hint:    rc.hint,
			}
			mu.Unlock()
		}(i, c)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(800 * time.Millisecond):
	}

	// Assemble response
	resChecks := make([]uiToolCheck, len(rawChecks))
	for i := range rawChecks {
		mu.Lock()
		tc, exists := results[i]
		mu.Unlock()
		if !exists {
			tc = uiToolCheck{
				Name:    rawChecks[i].name,
				Cmd:     rawChecks[i].cmd,
				Status:  "MISSING",
				Version: "",
				Hint:    rawChecks[i].hint,
			}
		}
		resChecks[i] = tc
	}

	resp := uiStatusResponse{
		OS:             info.Version,
		Arch:           info.Arch,
		CPUCores:       info.CPUCores,
		RAM:            osinfo.FormatRAM(info.RAMBytes),
		PackageManager: info.PackageManager,
		Checks:         resChecks,
	}

	_ = json.NewEncoder(w).Encode(resp)
}

const dashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>AutoDev — Interactive Dashboard</title>
  <link href="https://fonts.googleapis.com/css2?family=Space+Grotesk:wght@400;700;900&family=JetBrains+Mono:wght@400;700&display=swap" rel="stylesheet">
  <style>
    :root {
      --bg: #000000;
      --card-bg: #0a0a0a;
      --border: #2a2a2a;
      --text: #ffffff;
      --text-dim: #888888;
      --yellow: #0070f3;
      --green: #00c853;
      --red: #ee0000;
    }

    * {
      box-sizing: border-box;
      margin: 0;
      padding: 0;
    }

    body {
      background-color: var(--bg);
      color: var(--text);
      font-family: 'Space Grotesk', sans-serif;
      padding: 2rem;
      min-height: 100vh;
    }

    .container {
      max-width: 1200px;
      margin: 0 auto;
    }

    header {
      border: 4px solid var(--text);
      background-color: var(--yellow);
      color: var(--bg);
      padding: 2rem;
      margin-bottom: 2rem;
      display: flex;
      justify-content: space-between;
      align-items: center;
      box-shadow: 8px 8px 0px var(--text);
    }

    h1 {
      font-size: 3rem;
      font-weight: 900;
      letter-spacing: -2px;
      line-height: 0.9;
    }

    .subtitle {
      font-weight: 700;
      text-transform: uppercase;
      font-size: 0.9rem;
      letter-spacing: 2px;
      margin-top: 5px;
    }

    /* Grid layout */
    .grid {
      display: grid;
      grid-template-columns: 1fr;
      gap: 2rem;
    }

    @media (min-width: 768px) {
      .grid {
        grid-template-columns: 1fr 3fr;
      }
    }

    .panel {
      border: 4px solid var(--border);
      background-color: var(--card-bg);
      padding: 1.5rem;
      box-shadow: 6px 6px 0px var(--border);
    }

    .panel-title {
      font-size: 1.25rem;
      font-weight: 900;
      text-transform: uppercase;
      border-bottom: 2px solid var(--border);
      padding-bottom: 0.75rem;
      margin-bottom: 1rem;
      color: var(--yellow);
    }

    /* System info list */
    .spec-item {
      margin-bottom: 1rem;
    }

    .spec-label {
      font-size: 0.8rem;
      text-transform: uppercase;
      color: var(--text-dim);
      font-weight: 700;
    }

    .spec-val {
      font-size: 1.15rem;
      font-weight: 700;
      margin-top: 2px;
    }

    /* Tool Grid */
    .tools-grid {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
      gap: 1rem;
    }

    .tool-card {
      border: 2px solid var(--border);
      background-color: #080808;
      padding: 1.25rem;
      transition: all 0.2s ease;
      position: relative;
    }

    .tool-card:hover {
      border-color: var(--text);
      transform: translate(-2px, -2px);
      box-shadow: 4px 4px 0px var(--border);
    }

    .tool-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 0.5rem;
    }

    .tool-name {
      font-size: 1.1rem;
      font-weight: 700;
    }

    .badge {
      font-size: 0.75rem;
      font-weight: 900;
      padding: 2px 8px;
      text-transform: uppercase;
      border: 2px solid;
    }

    .badge-ok {
      background-color: rgba(0, 255, 135, 0.1);
      color: var(--green);
      border-color: var(--green);
    }

    .badge-missing {
      background-color: rgba(255, 95, 86, 0.1);
      color: var(--red);
      border-color: var(--red);
    }

    .tool-ver {
      font-family: 'JetBrains Mono', monospace;
      font-size: 0.85rem;
      color: var(--text-dim);
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
      margin-bottom: 0.75rem;
    }

    .tool-action {
      font-size: 0.8rem;
      color: var(--yellow);
      text-decoration: none;
      font-weight: 700;
      display: inline-block;
    }

    .tool-action:hover {
      text-decoration: underline;
    }

    .loading {
      text-align: center;
      padding: 4rem;
      font-size: 1.5rem;
      font-weight: 700;
      color: var(--yellow);
    }
  </style>
</head>
<body>
  <div class="container">
    <header>
      <div>
        <h1>AUTODEV</h1>
        <div class="subtitle">Clone. Scan. Install. Build.</div>
      </div>
      <div>
        <span class="badge badge-ok" style="background: black; color: #0070f3; font-size: 0.9rem; padding: 6px 12px;">LOCAL SERVER RUNNING</span>
      </div>
    </header>

    <div id="loader" class="loading">DIAGNOSING SYSTEM ENVIRONMENT...</div>

    <div id="dashboard-content" class="grid" style="display: none;">
      <!-- Left Panel: Specs -->
      <div class="panel">
        <div class="panel-title">System Specs</div>
        <div class="spec-item">
          <div class="spec-label">Operating System</div>
          <div id="spec-os" class="spec-val">-</div>
        </div>
        <div class="spec-item">
          <div class="spec-label">Architecture</div>
          <div id="spec-arch" class="spec-val">-</div>
        </div>
        <div class="spec-item">
          <div class="spec-label">CPU Cores</div>
          <div id="spec-cpu" class="spec-val">-</div>
        </div>
        <div class="spec-item">
          <div class="spec-label">System RAM</div>
          <div id="spec-ram" class="spec-val">-</div>
        </div>
        <div class="spec-item">
          <div class="spec-label">Package Manager</div>
          <div id="spec-pm" class="spec-val">-</div>
        </div>
      </div>

      <!-- Right Panel: Managed Tools -->
      <div class="panel">
        <div class="panel-title" style="display: flex; justify-content: space-between;">
          <span>Managed Dev Tools</span>
          <span id="install-summary" style="font-size: 0.9rem; color: var(--text-dim); text-transform: none;">-</span>
        </div>
        <div id="tools-container" class="tools-grid"></div>
      </div>
    </div>
  </div>

  <script>
    async function loadStatus() {
      try {
        const response = await fetch('/api/status');
        const data = await response.json();

        // System specs
        document.getElementById('spec-os').innerText = data.os;
        document.getElementById('spec-arch').innerText = data.arch;
        document.getElementById('spec-cpu').innerText = data.cpu_cores + ' cores';
        document.getElementById('spec-ram').innerText = data.ram;
        document.getElementById('spec-pm').innerText = data.package_manager;

        // Tools
        const container = document.getElementById('tools-container');
        container.innerHTML = '';

        let installedCount = 0;
        data.checks.forEach(tool => {
          const isOK = tool.status === 'OK';
          if (isOK) installedCount++;

          const card = document.createElement('div');
          card.className = 'tool-card';
          
          card.innerHTML = 
            '<div class="tool-header">' +
            '  <span class="tool-name">' + tool.name + '</span>' +
            '  <span class="badge ' + (isOK ? 'badge-ok' : 'badge-missing') + '">' + tool.status + '</span>' +
            '</div>' +
            '<div class="tool-ver">' + (isOK ? tool.version : 'Not found') + '</div>' +
            '<a href="' + (tool.hint.startsWith('http') ? tool.hint : '#') + '" target="' + (tool.hint.startsWith('http') ? '_blank' : '_self') + '" class="tool-action">' +
            '  ' + (isOK ? 'Documentation' : 'Install Guide') + '' +
            '</a>';
          container.appendChild(card);
        });

        document.getElementById('install-summary').innerText = '(' + installedCount + ' / ' + data.checks.length + ' installed)';

        // Show content
        document.getElementById('loader').style.display = 'none';
        document.getElementById('dashboard-content').style.display = 'grid';

      } catch (err) {
        document.getElementById('loader').innerText = 'ERROR CONNECTING TO LOCAL SERVER';
        console.error(err);
      }
    }

    loadStatus();
  </script>
</body>
</html>
`
