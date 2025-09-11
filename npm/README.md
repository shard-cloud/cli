Shard Cloud CLI.
This package provides a direct way to interact with the official ShardCloud API.

### Installation

To install the CLI, just run the following command in your terminal:

macOS, Linux, and WSL:

```bash
curl -fsSL https://cli.shardcloud.app | bash
```

Windows | need [npm](https://www.npmjs.com/) installed:

```bash
npm install -g shard-cloud-cli
```

### Commands

List of all commands available:

| Command | Description                                             |
| ------- | ------------------------------------------------------- |
| me      | Print the current logged-in user                        |
| backup  | Manage app backups                                      |
| commit  | Create a commit in the current application              |
| create  | Create an application with the current .shardcloud file |
| login   | Log in to Shard Cloud                                   |
| logout  | Log out                                                 |
| logs    | View the logs of your applications                      |
| status  | View the status of your applications                    |
| delete  | Delete your applications                                |
| start   | Start a stopped app                                     |
| stop    | Stop a running app                                      |
| restart | Restart a running app                                   |
| restore | Restore a backup                                        |

### Update

To update the CLI, just run the following command in your terminal:

macOS, Linux, and WSL:

```bash
curl -fsSL https://cli.shardcloud.app | bash
```

Windows | need [npm](https://www.npmjs.com/) installed:

```bash
shardcloud update
```

### Install a specific version (NPM only)

```bash
shardcloud install 1.0.0
```
