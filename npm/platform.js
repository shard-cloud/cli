import { join, resolve } from "node:path";
import url from "node:url";

const __dirname = url.fileURLToPath(new URL(".", import.meta.url));

export const ARCH_MAPPING = {
  ia32: "386",
  x64: "amd64",
  arm: "arm",
};

export const PLATFORM_MAPPING = {
  darwin: "darwin",
  linux: "linux",
  win32: "windows",
};

export function getPlatform() {
  const platform = PLATFORM_MAPPING[process.platform];
  if (!platform) {
    throw new Error("Unsupported Platform: " + process.platform);
  }
  return platform;
}

export function getFileExtension(platform) {
  return platform === "windows" ? "zip" : "tar.gz";
}

export function getDownloadURL(version) {
  const platform = getPlatform();
  const arch = getArch();
  return `https://github.com/shard-cloud/cli/releases/download/${version}/shardcloud_${platform}_${arch}.${getFileExtension(
    platform
  )}`;
}

export function getArch() {
  const arch = ARCH_MAPPING[process.arch];
  if (!arch) {
    throw new Error("Unsupported Arch: " + process.arch);
  }
  return arch;
}

export const BIN_NAME = "shardcloud";

export const BIN_DIR = join(__dirname, "bin");

export function getExecFile() {
  const extension = process.platform === "win32" ? ".exe" : "";
  const execfile = resolve(BIN_DIR, `shardcloud${extension}`);

  return execfile;
}
