# Connecting AutoDevs × DevMentor

This guide explains how to connect the **AutoDev CLI** with your **DevMentor Prompt Intelligence** system so that your prompts are captured, refined, scored, and synchronized in one go.

---

## 🔗 How the Connection Works

AutoDev communicates directly with the DevMentor API hosted at `https://devmentor-jmjh.onrender.com/` using a secure telemetry tunnel. The integration is governed by the core prompt-capture driver in [capture.go](file:///media/heet18/Futuristic1/Heet/Github/Autodev/packages/core/promptcapture/capture.go).

### 1. Authentication Token
The CLI checks for the presence of the `DEVMENTOR_TOKEN` environment variable in your terminal environment.

```bash
export DEVMENTOR_TOKEN="your_devmentor_api_token_here"
```

### 2. Endpoint Payload
When active, AutoDev automatically intercepts user prompts and sends them to the DevMentor prompt intelligence API:
* **URL**: `POST https://devmentor-jmjh.onrender.com/api/v1/prompts/event`
* **Content-Type**: `application/json`
* **Authorization**: `Bearer <DEVMENTOR_TOKEN>`
* **Payload Structure**:
  ```json
  {
    "original_prompt": "Install a React TypeScript app with Tailwind",
    "project_name": "my-app",
    "file_context": "AutoDev Prompt Tracking System active"
  }
  ```

### 3. Refinement Response
The DevMentor API returns a JSON payload containing details for local recording and AI rule updates:
```json
{
  "original_prompt": "Install a React TypeScript app with Tailwind",
  "refined_prompt": "Scaffold a new React project using TypeScript and configure Tailwind CSS automatically.",
  "score": 85,
  "workflow": "scaffolding",
  "technologies": ["react", "typescript", "tailwind"]
}
```

### 4. Local Workspace Update
AutoDev automatically writes these results back to your local environment under the `.autodevs` rules folder:
* **Refined Prompts**: Appends prompt details and score evaluation to `.autodevs/refined-prompts.md`.
* **Workflows**: Logs workflow categories and original queries to `.autodevs/workflows.md`.
* **Project Metadata**: Merges newly detected technologies into `.autodevs/metadata.json`.

---

## ⚡ Synchronizing in One Go

To activate the real-time capture and sync process, choose one of the three workflows:

### Workflow A: The Background Daemon (Recommended)
Runs a background loop to automatically extract prompts from active developer sessions:
```bash
autodev daemon
```
* **Mechanism**: The daemon monitors your running AI terminal processes. It dynamically locates and imports prompt history lines directly from your active Antigravity conversations, cleanses them, sends them to the DevMentor API, and refines them into your workspace files.

### Workflow B: Terminal Session Interceptor
If you want to wrap a specific terminal AI assistant session and sync every prompt immediately:
```bash
autodev capture gemini
# or
autodev capture claude
```
* **Mechanism**: Intercepts your input lines in real-time, forwarding them to the underlying AI CLI while automatically logging and syncing clean captures to DevMentor.

### Workflow C: Manual Batch Sync
If you have offline prompts or queued events stored under `.autodevs/analytics/queue.json`, push them all to the DevMentor API at once:
```bash
autodev sync
```
* **Mechanism**: Scans the queue file, uploads all unsynced events in a single batch, and clears the local queue file upon successful response.
