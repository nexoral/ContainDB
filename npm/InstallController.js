#!/usr/bin/env node

const { spawn } = require('child_process');
const path = require('path');
const os = require('os');

const platform = os.platform();
const arch = os.arch();

const supported = {
  'linux': {
    'x64': 'containdb_linux_amd64'
  },
  'darwin': {
    'x64': 'containdb_darwin_amd64',
    'arm64': 'containdb_darwin_arm64'
  },
  'win32': {
    'x64': 'containdb_windows_amd64.exe'
  }
};

if (!supported[platform] || !supported[platform][arch]) {
  console.error(`Unsupported platform or architecture: ${platform} ${arch}`);
  process.exit(1);
}

const binaryName = supported[platform][arch];
const binaryPath = path.join(__dirname, 'bin', binaryName);

// Copy existing environment and add our custom flag
const env = { ...process.env, CONTAINDB_INSTALL_SOURCE: 'npm' };

const child = spawn(binaryPath, process.argv.slice(2), { stdio: 'inherit', env });

child.on('exit', (code) => {
  process.exit(code);
});
