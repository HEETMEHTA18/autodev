#!/usr/bin/env node

const { spawn, execSync, execFileSync } = require("child_process");
const path = require("path");
const fs = require("fs");
const os = require("os");
const https = require("https");
const http = require("http");
const crypto = require("crypto");

const platformMap = { darwin: "darwin", linux: "linux", win32: "windows" };
const archMap = { x64: "amd64", arm64: "arm64" };
const platform = platformMap[process.platform];
const arch = archMap[process.arch];

if (!platform || !arch) {
  console.error(`[autodev] Unsupported platform/architecture: ${process.platform}/${process.arch}`);
  process.exit(1);
}

const ext = platform === "windows" ? "zip" : "tar.gz";
const binaryName = platform === "windows" ? "autodev.exe" : "autodev";
const pkgJson = require("../package.json");
const version = `v${pkgJson.version}`;
const binDir = __dirname;
const binaryPath = path.join(binDir, binaryName);
const bundledBinaryPath = path.join(binDir, "native", `${platform}-${arch}`, binaryName);
const cacheRoot = path.join(
  process.env.XDG_CACHE_HOME || (process.platform === "win32" ? process.env.LOCALAPPDATA || os.homedir() : path.join(os.homedir(), ".cache")),
  "autodev", "binaries", version, `${platform}-${arch}`,
);
const cachedBinaryPath = path.join(cacheRoot, binaryName);

function formatBytes(bytes) {
  if (!Number.isFinite(bytes) || bytes <= 0) return "0 B";
  const units = ["B", "KB", "MB", "GB"];
  const exponent = Math.min(Math.floor(Math.log(bytes) / Math.log(1024)), units.length - 1);
  return `${(bytes / 1024 ** exponent).toFixed(exponent === 0 ? 0 : 1)} ${units[exponent]}`;
}

function isInteractive() {
  return Boolean(process.stderr.isTTY);
}

function progress(downloaded, total, startedAt) {
  if (!isInteractive()) return;
  const elapsed = Math.max((Date.now() - startedAt) / 1000, 0.001);
  const speed = downloaded / elapsed;
  const width = 28;
  const ratio = total > 0 ? Math.min(downloaded / total, 1) : 0;
  const filled = total > 0 ? Math.round(ratio * width) : Math.min(8, width);
  const bar = `${"#".repeat(filled)}${"-".repeat(width - filled)}`;
  const pct = total > 0 ? `${(ratio * 100).toFixed(0)}%` : "...";
  process.stderr.write(`\r[autodev] Downloading [${bar}] ${pct} ${formatBytes(downloaded)}${total > 0 ? ` / ${formatBytes(total)}` : ""} @ ${formatBytes(speed)}/s`);
}

function finishProgress() {
  if (isInteractive()) process.stderr.write("\n");
}

function request(url, redirects = 5) {
  return new Promise((resolve, reject) => {
    if (redirects <= 0) return reject(new Error("Too many redirects"));
    const client = url.startsWith("https") ? https : http;
    const req = client.get(url, { headers: { "User-Agent": "autodev-npm-cli" }, timeout: 15000 }, (res) => {
      if (res.statusCode >= 300 && res.statusCode < 400 && res.headers.location) {
        res.resume();
        resolve(request(res.headers.location, redirects - 1));
        return;
      }
      if (res.statusCode !== 200) {
        res.resume();
        reject(new Error(`HTTP ${res.statusCode} from ${url}`));
        return;
      }
      resolve(res);
    });
    req.on("error", reject);
    req.on("timeout", () => req.destroy(new Error("Request timed out")));
  });
}

async function download(url, destination) {
  const res = await request(url);
  const total = Number(res.headers["content-length"] || 0);
  const startedAt = Date.now();
  let downloaded = 0;
  const file = fs.createWriteStream(destination, { flags: "wx" });

  return new Promise((resolve, reject) => {
    const fail = (error) => {
      file.destroy();
      try { fs.unlinkSync(destination); } catch (_) {}
      reject(error);
    };
    res.on("data", (chunk) => {
      downloaded += chunk.length;
      progress(downloaded, total, startedAt);
    });
    res.on("error", fail);
    file.on("error", fail);
    file.on("finish", () => file.close(() => {
      finishProgress();
      resolve(downloaded);
    }));
    res.pipe(file);
  });
}

async function sha256(filePath) {
  const hash = crypto.createHash("sha256");
  await new Promise((resolve, reject) => {
    const stream = fs.createReadStream(filePath);
    stream.on("data", (chunk) => hash.update(chunk));
    stream.on("error", reject);
    stream.on("end", resolve);
  });
  return hash.digest("hex");
}

function copyAtomic(source, destination) {
  fs.mkdirSync(path.dirname(destination), { recursive: true });
  const temp = `${destination}.tmp-${process.pid}-${Date.now()}`;
  fs.copyFileSync(source, temp);
  if (process.platform !== "win32") fs.chmodSync(temp, 0o755);
  fs.renameSync(temp, destination);
}

function extract(tempArchive) {
  if (ext === "zip") {
    if (process.platform === "win32") {
      const escapedTemp = tempArchive.replace(/'/g, "''");
      const escapedBin = binDir.replace(/'/g, "''");
      execSync(`powershell -NoProfile -Command "Expand-Archive -LiteralPath '${escapedTemp}' -DestinationPath '${escapedBin}' -Force"`, { stdio: "inherit" });
    } else {
      execFileSync("unzip", ["-o", tempArchive, "-d", binDir], { stdio: "inherit" });
    }
  } else {
    execFileSync("tar", ["-xzf", tempArchive, "-C", binDir], { stdio: "inherit" });
  }
}

async function downloadBinary() {
  const asset = `autodev_${platform}_${arch}.${ext}`;
  const url = `https://github.com/HEETMEHTA18/autodev/releases/download/${version}/${asset}`;
  const tempRoot = fs.mkdtempSync(path.join(os.tmpdir(), "autodev-"));
  const tempArchive = path.join(tempRoot, asset);

  console.error("");
  console.error("[autodev] AutoDev binary bootstrap");
  console.error(`[autodev] Platform : ${platform}/${arch}`);
  console.error(`[autodev] Version  : ${version}`);
  console.error(`[autodev] Source   : GitHub Releases`);
  console.error(`[autodev] Cache    : ${cacheRoot}`);

  try {
    fs.mkdirSync(cacheRoot, { recursive: true });
    console.error(`[autodev] Downloading ${asset}...`);
    const bytes = await download(url, tempArchive);
    if (bytes < 1000) throw new Error(`Downloaded archive is unexpectedly small (${bytes} bytes)`);
    console.error(`[autodev] Download complete (${formatBytes(bytes)})`);

    console.error("[autodev] Extracting binary...");
    extract(tempArchive);
    if (!fs.existsSync(binaryPath)) throw new Error(`Archive did not contain expected ${binaryName}`);

    const digest = await sha256(binaryPath);
    console.error(`[autodev] Binary SHA-256: ${digest}`);
    copyAtomic(binaryPath, cachedBinaryPath);
    console.error("[autodev] Binary ready.\n");
  } catch (error) {
    finishProgress();
    console.error(`\n[autodev] Bootstrap failed: ${error.message}`);
    console.error(`[autodev] Release: ${version}`);
    console.error(`[autodev] Asset  : ${asset}`);
    console.error("[autodev] Try again, check your network connection, or use the standalone release binary.");
    process.exit(1);
  } finally {
    fs.rmSync(tempRoot, { recursive: true, force: true });
  }
}

function resolveBinary() {
  const candidates = [
    binaryPath,
    bundledBinaryPath,
    cachedBinaryPath,
    path.join(__dirname, "..", "..", "cli", "bin", binaryName),
    path.join(__dirname, "..", "..", "..", "bin", binaryName),
    path.join(__dirname, "..", "..", "..", "packages", "cli", "bin", binaryName),
  ];
  return candidates.find((candidate) => fs.existsSync(candidate)) || null;
}

async function main() {
  let activeBinaryPath = resolveBinary();
  if (!activeBinaryPath) {
    await downloadBinary();
    activeBinaryPath = resolveBinary();
  }
  if (!activeBinaryPath) {
    console.error(`[autodev] Binary resolution failed for ${platform}/${arch}.`);
    process.exit(1);
  }

  if (process.platform !== "win32") {
    try { fs.chmodSync(activeBinaryPath, 0o755); } catch (_) {}
  }

  const child = spawn(activeBinaryPath, process.argv.slice(2), { stdio: "inherit" });
  ["SIGINT", "SIGTERM", "SIGHUP", "SIGQUIT"].forEach((signal) => {
    process.on(signal, () => { if (!child.killed) child.kill(signal); });
  });
  child.on("close", (code, signal) => {
    if (code !== null) return process.exit(code);
    if (signal) {
      const codes = { SIGINT: 2, SIGTERM: 15, SIGHUP: 1, SIGQUIT: 3 };
      return process.exit(128 + (codes[signal] || 0));
    }
    return process.exit(0);
  });
  child.on("error", (error) => {
    console.error(`[autodev] Failed to run native binary: ${error.message}`);
    process.exit(1);
  });
}

main().catch((error) => {
  console.error(`[autodev] Fatal error: ${error.message}`);
  process.exit(1);
});
