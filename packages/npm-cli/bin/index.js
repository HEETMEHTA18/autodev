#!/usr/bin/env node

const { spawn } = require("child_process");
const path = require("path");
const fs = require("fs");

const platformMap = {
  darwin: "darwin",
  linux: "linux",
  win32: "windows",
};

const archMap = {
  x64: "amd64",
  arm64: "arm64",
};

const platform = platformMap[process.platform];
const arch = archMap[process.arch];

if (!platform || !arch) {
  console.error(
    `[autodev] Unsupported platform/architecture: ${process.platform}/${process.arch}`,
  );
  process.exit(1);
}

const binaryName = platform === "windows" ? "autodev.exe" : "autodev";
const pkgJson = require("../package.json");

// Native binaries are bundled into the npm package at publish time.
// This keeps `npx @heetmehta18/autodev` offline-capable after npm install
// and prevents the npm wrapper from silently running a different release.
const binaryPath = path.join(
  __dirname,
  "native",
  `${platform}-${arch}`,
  binaryName,
);

// Development fallback: allow the wrapper to run against a locally built CLI.
const devPaths = [
  path.join(__dirname, "..", "..", "cli", "bin", binaryName),
  path.join(__dirname, "..", "..", "..", "bin", binaryName),
  path.join(__dirname, "..", "..", "..", "packages", "cli", "bin", binaryName),
];

let activeBinaryPath = binaryPath;
if (!fs.existsSync(activeBinaryPath)) {
  for (const devPath of devPaths) {
    if (fs.existsSync(devPath)) {
      activeBinaryPath = devPath;
      break;
    }
  }
}

function failMissingBinary() {
  console.error(
    `[autodev] Native binary missing for ${platform}/${arch} in @heetmehta18/autodev@${pkgJson.version}.`,
  );
  console.error(
    `[autodev] This npm package was published without its platform binary. Please install the latest version or report the release packaging issue.`,
  );
  process.exit(1);
}

function main() {
  if (!fs.existsSync(activeBinaryPath)) {
    failMissingBinary();
  }

  const child = spawn(activeBinaryPath, process.argv.slice(2), {
    stdio: "inherit",
  });

  const signals = ["SIGINT", "SIGTERM", "SIGHUP", "SIGQUIT"];
  signals.forEach((signal) => {
    process.on(signal, () => {
      if (!child.killed) child.kill(signal);
    });
  });

  child.on("close", (code, signal) => {
    if (code !== null) {
      process.exit(code);
      return;
    }

    if (signal) {
      const signalCodes = { SIGINT: 2, SIGTERM: 15, SIGHUP: 1, SIGQUIT: 3 };
      process.exit(128 + (signalCodes[signal] || 0));
      return;
    }

    process.exit(0);
  });

  child.on("error", (err) => {
    console.error(`[autodev] Failed to run native binary: ${err.message}`);
    process.exit(1);
  });
}

main();
