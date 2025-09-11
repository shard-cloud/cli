#!/usr/bin/env node

import { getExecFile } from "./platform.js";
import { spawnSync } from "node:child_process";
import { existsSync } from "node:fs";
import { installBinaries } from "./install.js";

async function main() {
  const binfile = getExecFile();
  const [, , ...args] = process.argv;
  const isUpdate =
    args[0] === "update" || args[0] === "upgrade" || args[0] === "install";
  if (!existsSync(binfile) || isUpdate) {
    const version = args[0] === "install" ? args[1] : undefined;
    console.log("Installing Shard Cloud CLI, please wait...");
    await installBinaries(version);

    if (isUpdate) {
      return;
    }
  }

  let result = spawnSync(binfile, args, {
    cwd: process.cwd(),
    stdio: "inherit",
  });
  if (result.error) console.log(result.error);

  return result.status;
}

main();
