import axios from "axios";
import {
  getDownloadURL,
  BIN_DIR,
  getPlatform,
  getFileExtension,
} from "./platform.js";
import { join } from "node:path";
import {
  mkdirSync,
  createWriteStream,
  createReadStream,
  unlinkSync,
} from "node:fs";
import { promisify } from "node:util";
import { pipeline } from "node:stream";
import { extract as extractTar } from "tar-fs";
import { createGunzip } from "node:zlib";

import AdmZip from "adm-zip";

const pipelineAsync = promisify(pipeline);

export async function installBinaries(version) {
  const latestVersion = version || (await getLatestVersion());
  const url = getDownloadURL(latestVersion);
  const platform = getPlatform();
  const extension = getFileExtension(platform);
  mkdirSync(BIN_DIR, { recursive: true });
  const installPath = join(BIN_DIR, `shardcloud.${extension}`);
  const response = await axios.get(url, {
    timeout: 30000,
    responseType: "stream",
  });

  const writeStream = createWriteStream(installPath);

  await pipelineAsync(response.data, writeStream);
  await extract(installPath, BIN_DIR);
  unlinkSync(installPath);
}

function extract(path, dest) {
  return new Promise((resolve, reject) => {
    if (path.endsWith(".zip")) {
      const zip = new AdmZip(path);
      zip.extractAllToAsync(dest, true, false, (err) =>
        err ? reject(err) : resolve()
      );
      return;
    }

    const extract = extractTar(dest);
    const stream = createReadStream(path).pipe(createGunzip()).pipe(extract);

    stream.on("finish", resolve);
    stream.on("error", reject);
  });
}

async function getLatestVersion() {
  try {
    const url = "https://api.github.com/repos/shard-cloud/cli/releases/latest";

    const response = await axios.get(url, {
      timeout: 30000,
      headers: {
        Accept: "application/vnd.github.v3+json",
      },
    });

    const { tag_name } = response.data;

    if (!tag_name) {
      throw new Error("No tag_name found in release");
    }

    return tag_name;
  } catch (error) {
    throw new Error(`Failed to get latest version: ${error.message}`);
  }
}
