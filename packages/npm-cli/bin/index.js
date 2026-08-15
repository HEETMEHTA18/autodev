#!/usr/bin/env node

const { spawn, execSync, execFileSync } = require("child_process");
const path = require("path");
const fs = require("fs");
const os = require("os");
const crypto = require("crypto");
const https = require("https");
const http = require("http");

const REPO = "HEETMEHTA18/autodev";
const platformMap = { darwin: "darwin", linux: "linux", win32: "windows" };
const archMap = { x64: "amd64", arm64: "arm64" };
const platform = platformMap[process.platform];
const arch = archMap[process.arch];
const pkgJson = require("../package.json");
const packageVersion = `v${pkgJson.version}`;
const binaryName = platform === "windows" ? "autodev.exe" : "autodev";
const archiveExt = platform === "windows" ? "zip" : "tar.gz";

if (!platform || !arch) {
  console.error(`\nAutoDev cannot run on ${process.platform}/${process.arch}. Supported targets: Linux, macOS and Windows on amd64/arm64.`);
  process.exit(1);
}

const args = process.argv.slice(2);
const offline = args.includes("--offline") || process.env.AUTODEV_OFFLINE === "1";
const noCache = args.includes("--no-cache") || process.env.AUTODEV_NO_CACHE === "1";
const cacheRoot = path.join(
  process.env.AUTODEV_CACHE_DIR || path.join(os.homedir(), ".cache", "autodev"),
  packageVersion,
  `${platform}-${arch}`,
);
const cacheBinary = path.join(cacheRoot, binaryName);
const binDir = __dirname;
const binaryPath = path.join(binDir, binaryName);

function print(message = "") {
  process.stdout.write(`${message}\n`);
}

function fail(message, details = []) {
  console.error(`\n✕ AutoDev bootstrap failed\n  ${message}`);
  for (const detail of details) console.error(`  • ${detail}`);
  process.exit(1);
}

function formatBytes(value) {
  if (!Number.isFinite(value) || value <= 0) return "0 B";
  const units = ["B", "KB", "MB", "GB"];
  let n = value;
  let i = 0;
  while (n >= 1024 && i < units.length - 1) { n /= 1024; i++; }
  return `${n.toFixed(i === 0 ? 0 : 1)} ${units[i]}`;
}

function formatEta(seconds) {
  if (!Number.isFinite(seconds) || seconds <= 0) return "--";
  if (seconds < 60) return `${Math.ceil(seconds)}s`;
  return `${Math.floor(seconds / 60)}m ${Math.ceil(seconds % 60)}s`;
}

function cacheIsValid() {
  try {
    return fs.existsSync(cacheBinary) && fs.statSync(cacheBinary).size > 1000;
  } catch (_) {
    return false;
  }
}

function getLatestReleaseTag() {
  return new Promise((resolve) => {
    const req = https.get(
      {
        hostname: "api.github.com",
        path: `/repos/${REPO}/releases/latest`,
        headers: { "User-Agent": "autodev-npm-cli" },
        timeout: 5000,
      },
      (res) => {
        let body = "";
        res.on("data", (chunk) => (body += chunk));
        res.on("end", () => {
          try {
            const json = JSON.parse(body);
            const valid = /^v?\d+\.\d+\.\d+(-[a-zA-Z0-9.]+)?$/;
            resolve(json.tag_name && valid.test(json.tag_name) ? json.tag_name : null);
          } catch (_) { resolve(null); }
        });
      },
    );
    req.on("error", () => resolve(null));
    req.on("timeout", () => { req.destroy(); resolve(null); });
  });
}

function request(url, redirects = 6) {
  return new Promise((resolve, reject) => {
    if (redirects <= 0) return reject(new Error("Too many redirects"));
    const client = url.startsWith("https:") ? https : http;
    const req = client.get(url, { headers: { "User-Agent": "autodev-npm-cli" } }, (res) => {
      if (res.statusCode >= 300 && res.statusCode < 400 && res.headers.location) {
        res.resume();
        return request(new URL(res.headers.location, url).toString(), redirects - 1).then(resolve, reject);
      }
      resolve(res);
    });
    req.on("error", reject);
    req.setTimeout(15000, () => req.destroy(new Error("request timeout")));
  });
}

async function download(url, destination, label) {
  const res = await request(url);
  if (res.statusCode !== 200) {
    res.resume();
    throw new Error(`HTTP ${res.statusCode}`);
  }

  const total = Number(res.headers["content-length"] || 0);
  let received = 0;
  let lastBytes = 0;
  let lastTime = Date.now();
  const started = Date.now();
  const stream = fs.createWriteStream(destination, { flags: "wx" });

  return new Promise((resolve, reject) => {
    const cleanup = (err) => {
      res.destroy();
      stream.destroy();
      try { fs.unlinkSync(destination); } catch (_) {}
      reject(err);
    };

    res.on("data", (chunk) => {
      received += chunk.length;
      const now = Date.now();
      if (now - lastTime >= 250) {
        const speed = (received - lastBytes) / ((now - lastTime) / 1000);
        const elapsed = (now - started) / 1000;
        const eta = total ? (total - received) / Math.max(speed, 1) : 0;
        const percent = total ? Math.min(100, (received / total) * 100) : 0;
        const bar = total ? `${percent.toFixed(0).padStart(3)}%` : formatBytes(received);
        process.stdout.write(`\r  ${label}  ${bar}  ${formatBytes(received)}${total ? ` / ${formatBytes(total)}` : ""}  ${formatBytes(speed)}/s  ETA ${formatEta(eta)}   `);
        lastBytes = received;
        lastTime = now;
      }
    });
    res.on("error", cleanup);
    stream.on("error", cleanup);
    stream.on("finish", () => {
      stream.close(() => {
        process.stdout.write("\r" + " ".repeat(110) + "\r");
        resolve({ bytes: received, elapsed: (Date.now() - started) / 1000 });
      });
    });
    res.pipe(stream);
  });
}

function retryDelay(attempt) {
  return Math.min(8000, 500 * Math.pow(2, attempt - 1)) + Math.floor(Math.random() * 250);
}

async function downloadWithRetry(url, destination, label, attempts = 3) {
  let lastError;
  for (let attempt = 1; attempt <= attempts; attempt++) {
    try {
      return await download(url, destination, `${label} (attempt ${attempt}/${attempts})`);
    } catch (err) {
      lastError = err;
      if (attempt < attempts) {
        const wait = retryDelay(attempt);
        print(`\n  Network error: ${err.message}. Retrying in ${Math.ceil(wait / 1000)}s...`);
        await new Promise((resolve) => setTimeout(resolve, wait));
      }
    }
  }
  throw lastError;
}

async function fetchChecksum(version, archiveFile) {
  const url = `https://github.com/${REPO}/releases/download/${version}/checksums.txt`;
  const res = await request(url);
  if (res.statusCode !== 200) { res.resume(); throw new Error(`checksum file unavailable (HTTP ${res.statusCode})`); }
  let body = "";
  for await (const chunk of res) body += chunk;
  const line = body.split(/\r?\n/).find((entry) => entry.trim().endsWith(`  ${archiveFile}`) || entry.trim().endsWith(` *${archiveFile}`));
  if (!line) throw new Error(`checksum for ${archiveFile} not found`);
  return line.trim().split(/\s+/)[0];
}

function sha256(file) {
  return new Promise((resolve, reject) => {
    const hash = crypto.createHash("sha256");
    const stream = fs.createReadStream(file);
    stream.on("data", (chunk) => hash.update(chunk));
    stream.on("error", reject);
    stream.on("end", () => resolve(hash.digest("hex")));
  });
}

function extractArchive(archive, destination) {
  fs.mkdirSync(destination, { recursive: true });
  if (platform === "windows") {
    const escapedArchive = archive.replace(/'/g, "''");
    const escapedDestination = destination.replace(/'/g, "''");
    execSync(`powershell -NoProfile -NonInteractive -Command "Expand-Archive -LiteralPath '${escapedArchive}' -DestinationPath '${escapedDestination}' -Force"`, { stdio: "ignore" });
  } else {
    execFileSync("tar", ["-xzf", archive, "-C", destination], { stdio: "ignore" });
  }
}

function findExtractedBinary(directory) {
  const direct = path.join(directory, binaryName);
  if (fs.existsSync(direct)) return direct;
  const entries = fs.readdirSync(directory, { withFileTypes: true });
  for (const entry of entries) {
    if (entry.isFile() && entry.name === binaryName) return path.join(directory, entry.name);
    if (entry.isDirectory()) {
      const nested = findExtractedBinary(path.join(directory, entry.name));
      if (nested) return nested;
    }
  }
  return null;
}

function atomicInstall(source) {
  fs.mkdirSync(binDir, { recursive: true });
  const staged = path.join(binDir, `.${binaryName}.${process.pid}.tmp`);
  fs.copyFileSync(source, staged);
  if (platform !== "windows") fs.chmodSync(staged, 0o755);
  fs.renameSync(staged, binaryPath);
}

function localDevBinary() {
  const devPaths = [
    path.join(__dirname, "..", "..", "cli", "bin", binaryName),
    path.join(__dirname, "..", "..", "..", "bin", binaryName),
    path.join(__dirname, "..", "..", "..", "packages", "cli", "bin", binaryName),
  ];
  return devPaths.find((candidate) => fs.existsSync(candidate)) || null;
}

async function bootstrap() {
  print("\n  AutoDev bootstrap");
  print("  ─────────────────────────────────────────────────────────");
  print(`  Platform       ${platform}/${arch}`);
  print(`  Package        ${pkgJson.name}@${pkgJson.version}`);
  print(`  Release        ${packageVersion}`);
  print(`  Destination    ${binaryPath}`);
  print(`  Cache          ${noCache ? "disabled" : cacheRoot}`);

  if (offline) {
    if (cacheIsValid()) {
      print("  ✓ Offline mode: using cached binary");
      atomicInstall(cacheBinary);
      return;
    }
    fail("Offline mode was requested, but no valid cached binary is available.", ["Connect to the network and run again, or remove --offline."]);
  }

  if (!noCache && cacheIsValid()) {
    print("  ✓ Valid cached binary found");
    atomicInstall(cacheBinary);
    return;
  }

  const latest = await getLatestReleaseTag();
  if (latest && latest !== packageVersion) {
    print(`  ! Latest GitHub release is ${latest}; using ${packageVersion} to match this npm package.`);
  }

  const archiveName = `autodev_${platform}_${arch}`;
  const archiveFile = `${archiveName}.${archiveExt}`;
  const url = `https://github.com/${REPO}/releases/download/${packageVersion}/${archiveFile}`;
  const root = noCache ? fs.mkdtempSync(path.join(os.tmpdir(), "autodev-")) : cacheRoot;
  fs.mkdirSync(root, { recursive: true });
  const archive = path.join(root, archiveFile);
  const partial = `${archive}.part`;

  try {
    print(`  Source         ${url}`);
    print("\n  Downloading...");
    await downloadWithRetry(url, partial, "  AutoDev");
    fs.renameSync(partial, archive);

    print("  Verifying SHA-256...");
    const expected = await fetchChecksum(packageVersion, archiveFile);
    const actual = await sha256(archive);
    if (expected.toLowerCase() !== actual.toLowerCase()) {
      throw new Error(`checksum mismatch (expected ${expected}, got ${actual})`);
    }
    print("  ✓ Checksum verified");

    const staging = fs.mkdtempSync(path.join(os.tmpdir(), "autodev-extract-"));
    try {
      print("  Extracting safely...");
      extractArchive(archive, staging);
      const extracted = findExtractedBinary(staging);
      if (!extracted) throw new Error(`archive does not contain ${binaryName}`);
      if (fs.statSync(extracted).size < 1000) throw new Error("extracted binary is unexpectedly small");

      if (!noCache) {
        fs.copyFileSync(extracted, cacheBinary);
        if (platform !== "windows") fs.chmodSync(cacheBinary, 0o755);
      }
      atomicInstall(noCache ? extracted : cacheBinary);
    } finally {
      fs.rmSync(staging, { recursive: true, force: true });
    }
    if (noCache) fs.rmSync(root, { recursive: true, force: true });
    else fs.rmSync(archive, { force: true });
    print("\n  ✓ AutoDev is ready");
  } catch (err) {
    try { fs.unlinkSync(partial); } catch (_) {}
    if (noCache) fs.rmSync(root, { recursive: true, force: true });
    fail(err.message, [
      `Requested release: ${packageVersion}`,
      `Target: ${platform}/${arch}`,
      "If you are offline, rerun with --offline after a successful download.",
      "If the release asset is missing, publish the matching CLI release before publishing the npm package.",
    ]);
  }
}

function runBinary(binary, binaryArgs) {
  const child = spawn(binary, binaryArgs, { stdio: "inherit" });
  for (const signal of ["SIGINT", "SIGTERM", "SIGHUP", "SIGQUIT"]) {
    process.on(signal, () => { if (!child.killed) child.kill(signal); });
  }
  child.on("error", (err) => fail(`Unable to start native AutoDev binary: ${err.message}`));
  child.on("close", (code, signal) => {
    if (code !== null) process.exit(code);
    const codes = { SIGINT: 2, SIGTERM: 15, SIGHUP: 1, SIGQUIT: 3 };
    process.exit(128 + (codes[signal] || 0));
  });
}

async function main() {
  const dev = localDevBinary();
  if (dev) return runBinary(dev, args);

  if (args.includes("self") && (args.includes("repair") || args.includes("update"))) {
    if (args.includes("repair")) fs.rmSync(cacheRoot, { recursive: true, force: true });
  }

  if (!fs.existsSync(binaryPath)) await bootstrap();
  if (!fs.existsSync(binaryPath)) fail("Native binary is unavailable after bootstrap.");
  runBinary(binaryPath, args);
}

main().catch((err) => fail(`Unexpected bootstrap error: ${err.message}`));
