# Deployment Workflow

This repository ships with a one-command deploy script: [`deploy`](/Users/herock/Workspace/Dev/pureapi/new-api/deploy).

## Goals

- Refuse to deploy if the git worktree is dirty.
- Show the latest commit message and wait for confirmation before publishing.
- Push the current branch to `origin`.
- Build the frontend with `bun run build`.
- Build a Linux binary into `build/<version>/pureapi-server`.
- Upload the release to the server.
- Switch the live binary to the new release.
- Restart `pureapi` and run a health check.
- Roll back automatically if the health check fails.

## Version Format

Releases use `YYYYMMDDNN`, for example:

- `2026041101`
- `2026041102`

The script calculates the next sequence by checking both:

- local `build/`
- remote `releases/`

This avoids collisions when you deploy from different machines or after cleaning local artifacts.

## First-Time Setup

### 1. Local config

Copy the example config:

```bash
cd /Users/herock/Workspace/Dev/pureapi/new-api
cp deploy.env.example .deploy.env
```

Edit `.deploy.env` for your server.

The deploy workflow requires `rsync` on both the local machine and the remote server. Missing `rsync` is treated as a hard failure.

### 2. Remote directory layout

The script expects or creates these paths:

```text
/home/pureapi/releases/
/home/pureapi/current
/home/pureapi/shared/
```

Each release lives in its own directory under `/home/pureapi/releases/<version>/`, and `/home/pureapi/current` is a symlink to the active release.

### 3. Non-interactive sudo

The deploy user must be able to restart the service without typing a password, for example:

```text
pureapi ALL=(root) NOPASSWD:/usr/bin/systemctl restart pureapi
```

Without this, one-command deploy is not possible.

## Recommended systemd Unit

```ini
[Unit]
Description=PureAPI New-API Server
After=network.target postgresql.service

[Service]
Type=exec
User=pureapi
WorkingDirectory=/home/pureapi
EnvironmentFile=/home/pureapi/shared/.env
ExecStart=/home/pureapi/current/pureapi-server --log-dir /home/pureapi/logs
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

`Type=exec` is the recommended mode for long-running services because startup failures such as a missing executable are surfaced more reliably in systemd's own lifecycle handling. Source: [systemd.service](https://www.freedesktop.org/software/systemd/man/255/systemd.service.html)

`WorkingDirectory=/home/pureapi` is intentional. This project still has relative-path defaults such as `one-api.db`, so using the stable base directory is safer than pointing the working directory at the current release.

`--log-dir /home/pureapi/logs` is included on purpose so logs stay in one stable shared directory instead of being written inside the currently active release directory.

## Usage

Interactive deploy:

```bash
cd /Users/herock/Workspace/Dev/pureapi/new-api
./deploy
```

Non-interactive deploy:

```bash
cd /Users/herock/Workspace/Dev/pureapi/new-api
./deploy --yes
```

## Failure Behavior

If deployment fails, the script prints the step that failed. Common examples:

- uncommitted local changes
- SSH connection failure
- `sudo -n` not allowed on the server
- frontend or Go build failure
- upload failure
- remote health check failure

If the health check fails after restart, the script automatically switches the live binary back to the previous target and restarts the service.
