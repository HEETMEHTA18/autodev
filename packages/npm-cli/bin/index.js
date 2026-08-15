#!/usr/bin/env node

const { spawn, execSync, execFileSync } = require("child_process");
const path = require("path");
const fs = require("fs");
const os = require("os");
const https = require("https");
const http = require("http");

const platformMap = { darwin: "darwin", linux: "linux", win32: "windows" };
const archMap = { x64: "amd64", arm64: "arm64" };
const platform = platformMap[process.platform];
const arch = archMap[process.arch];

if (!platform || !arch) {
  console.error(
    `[autodev] Unsupported platform/architecture: ${process.platform}/${process.arch}`,
  );
  process.exit(1);
}

const ext = platform === "windows" ? "zip" : "tar.gz";
const binaryName = platform === "windows" ? "autodev.exe" : "autodev";
const pkgJson = require("../package.json");
const version = `v${pkgJson.version}`;

const binDir = __dirname;
const binaryPath = path.join(binDir, binaryName);

// Development fallback paths. These are intentionally checked before downloading.
const devPaths = [
  path.join(__dirname, "..", "..", "cli", "bin", binaryName),
  path.join(__dirname, "..", "..", "..", "bin", binaryName),
  path.join(__dirname, "..", "..", "..", "packages", "cli", "bin", binaryName),
];

let activeBinaryPath = binaryPath;
for (const devPath of devPaths) {
  if (fs.existsSync(devPath)) {
    activeBinaryPath = devPath;
    break;
  }
}

function download(url, destPath, maxRedirects = 5) {
  return new Promise((resolve, reject) => {
    if (maxRedirects <= 0) return reject(new Error("Too many redirects"));

    const client = url.startsWith("https") ? https : http;
    client
      .get(url, { headers: { "User-Agent": "autodev-npm-cli" } }, (res) => {
        if (
          res.statusCode >= 300 &&
          res.statusCode < 400 &&
          res.headers.location
        ) {
          return download(res.headers.location, destPath, maxRedirects - 1)
            .then(resolve)
            .catch(reject);
        }

        if (res.statusCode !== 200) {
          res.resume();
          return reject(new Error(`HTTP ${res.statusCode} from ${url}`));
        }

        const fileStream = fs.createWriteStream(destPath);
        res.pipe(fileStream);
        fileStream.on("finish", () => {
          fileStream.close();
          resolve();
        });
        fileStream.on("error", reject);
        res.on("error", reject);
      })
      .on("error", reject);
  });
}

async function downloadBinary() {
  const archiveName = `autodev_${platform}_${arch}`;
  const archiveFile = `${archiveName}.${ext}`;
  const url = `https://github.com/HEETMEHTA18/autodev/releases/download/${version}/${archiveFile}`;
  const tempFile = path.join(
    os.tmpdir(),
    `autodev_download_${Date.now()}.${ext}`,
  );

  console.log(
    `[autodev] Native binary not bundled for ${platform}/${arch}. Downloading AutoDev ${version}...`,
  );
  console.log(`[autodev] Downloading from: ${url}`);

  try {
    await download(url, tempFile);

    if (!fs.existsSync(tempFile)) {
      throw new Error("Download completed but archive was not found on disk.");
    }

    const stat = fs.statSync(tempFile);
    if (stat.size < 1000) {
      throw new Error(
        `Downloaded archive is too small (${stat.size} bytes), likely an error page.`,
      );
    }

    console.log(`[autodev] Downloaded ${(stat.size / 1024 / 1024).toFixed(1)} MB`);
    console.log(`[autodev] Extracting binary...`);

    if (ext === "zip") {
      if (process.platform === "win32") {
        const escapedTempFile = tempFile.replace(/'/g, "''");
        const escapedBinDir = binDir.replace(/'/g, "''");
        execSync(
          `powershell -NoProfile -Command "Expand-Archive -LiteralPath '${escapedTempFile}' -DestinationPath '${escapedBinDir}' -Force"`,
          { stdio: "inherit" },
        );
      } else {
        execFileSync("unzip", ["-o", tempFile, "-d", binDir], {
          stdio: "inherit",
        });
      }
    } else {
      execFileSync("tar", ["-xzf", tempFile, "-C", binDir], {
        stdio: "inherit",
      });
    }

    if (!fs.existsSync(binaryPath)) {
      throw new Error(`Archive extracted successfully but ${binaryName} was not found.`);
    }

    if (process.platform !== "win32") {
      fs.chmodSync(binaryPath, 0o755);
    }

    console.log(`[autodev] Installation successful.`);
  } catch (err) {
    console.error(`\n[autodev] Error installing native binary: ${err.message}`);
    console.error(`[autodev] Release asset: ${url}`);
    process.exit(1);
  } finally {
    if (fs.existsSync(tempFile)) fs.unlinkSync(tempFile);
  }
}

async function main() {
  if (activeBinaryPath === binaryPath && !fs.existsSync(binaryPath)) {
    await downloadBinary();
  }

  if (!fs.existsSync(activeBinaryPath)) {
    console.error(`[autodev] Native binary is unavailable: ${activeBinaryPath}`);
    process.exit(1);
  }

  const args = process.argv.slice(2);
  const child = spawn(activeBinaryPath, args, { stdio: "inherit" });

  const signals = ["SIGINT", "SIGTERM", "SIGHUP", "SIGQUIT"];
  signals.forEach((signal) => {
    process.on(signal, () => {
      if (!child.killed) child.kill(signal);
    });
  });

  child.on("close", (code, signal) => {
    if (code !== null) {
      process.exit(code);
    } else if (signal) {
      const signalCodes = { SIGINT: 2, SIGTERM: 15, SIGHUP: 1, SIGQUIT: 3 };
      process.exit(128 + (signalCodes[signal] || 0));
    } else {
      process.exit(0);
    }
  });

  child.on("error", (err) => {
    console.error(`[autodev] Failed to run binary: ${err.message}`);
    process.exit(1);
  });
}

main().catch((err) => {
  console.error(`[autodev] Fatal error: ${err.message}`);
  process.exit(1);
});
