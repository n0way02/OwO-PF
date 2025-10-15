# Findingway

Lightweight Discord bot that scrapes XIV party finder listings and posts formatted raid listings to configured Discord channels.

This repository contains a fork with a few configuration and build improvements. The original project was created by https://github.com/Veraticus — all credit to the original author; this fork contains minor modifications and added convenience commands.

---

## Overview

Findingway periodically scrapes https://xivpf.com/ (Final Fantasy XIV Party Finder) and posts filtered raid listings to Discord channels. Channels and their configuration are loaded from `config.yaml` and the bot uses a Discord bot token provided via environment variables.

This fork focuses on improving environment handling (no hard-coded tokens), adding a `.env.template`, and making the project easier to build and deploy (local exe, Docker). It does not change the original business logic.

## Features

- Scrapes party-finder listings from xivpf.com
- Cleans and posts formatted listings to one or more Discord channels per configuration
- Command registration for the bot via Discord API
- Configured through `config.yaml` and environment variables

## Short description (for repository listings)

Discord bot to scrape Final Fantasy XIV party finder listings and post them to configured channels. Original project by Veraticus. Modified and maintained by the current repository owner.

## Quickstart (local)

Prerequisites:
- Go (recommended >= 1.20)
- Git

1. Copy the environment template and fill the real values:

```powershell
cp .env.template .env
# Edit .env and set DISCORD_TOKEN and any other values
notepad .env
```

2. Tidy modules (optional but recommended):

```powershell
go mod tidy
```

3. Build (Windows executable):

```powershell
# from project root
go build -o findingway.exe .
```

4. Run:

```powershell
# Ensure DISCORD_TOKEN is set in your environment or loaded from .env by your shell
.
./findingway.exe
```

On Windows you can set the env var inline in PowerShell:

```powershell
$env:DISCORD_TOKEN="your-token-here"; ./findingway.exe
```

## Docker

Build a Docker image locally:

```powershell
docker build -t findingway:local .
```

Run the image (map config or mount .env if needed):

```powershell
docker run --rm -e DISCORD_TOKEN="$env:DISCORD_TOKEN" -v ${PWD}/config.yaml:/findingway/config.yaml findingway:local
```

Notes for hosted builders (Discloud):
- Some hosted builders using older Debian base images (e.g., buster) may fail during apt steps. If you hit errors like "does not have a Release file", update the builder base image used in `Dockerfile` to a newer Debian tag (e.g. `golang:1.22-bullseye`) or ensure `apt-get update` points to valid repositories.
- If the build needs `git` for ldflags (to get commit SHA), ensure the builder image has `git` installed or modify the Dockerfile to `apt-get install -y git` before running `make build`.

## Environment variables

- `DISCORD_TOKEN` (required) — Discord bot token. Do not commit this to GitHub.
- `ONCE` (optional) — set to `true` to run a single iteration and exit.

Use `.env.template` as a starting point. Do NOT commit a populated `.env`.

## Config

- `config.yaml` contains the list of channels and their settings. Edit it to add the channels you want the bot to post to.

## Tests

Run tests:

```powershell
go test ./...
```

## Contributing

If you want to contribute fixes or features, open a pull request. Keep secrets out of commits. If you need CI to use secrets, add them to Actions secrets or the host provider's secret manager.

## Credits

Original project by https://github.com/Veraticus. This repository contains minor modifications (environment handling, .env template, build notes) by the current repository owner.

## License

See the `LICENSE` file in this repository for license details.

---

If you want, I can also:
- add a short repository description to `package.json`/GitHub settings text you can copy, or
- create a small `CONTRIBUTING.md` or `CHANGELOG.md`.
# findingway

Inspired and indebted to the [Rust project of the same name](https://github.com/epitaque/findingway/) by [epitaque](https://github.com/epitaque). You can see it in action in [Aether PUG DSR](https://discord.gg/aetherpugdsr) in the #pf-checks channel.

findingway scrapes https://xivpf.com/listings every 3 minutes, collects the resulting listings, and posts them onto a Discord channel of your choice. Note that xivpf.com is not particularly accurate and includes private listings; there does not seem to be a way to segment them out at the current time.

## Running

findingway ingests its configuration file at `./config.yaml` to determine what to parse.

findingway requires one environment variable to start:

* **DISCORD_TOKEN**: You have to create a [Discord bot for findingway](https://discord.com/developers/applications). Once you've done so, you can add the bot token here.

findingway also accepts one optional environment variable:

* **ONCE**: If present, findingway will run only once and then exit successfully. Otherwise it will run perpetually and update the target channel every three minutes.

I'm not totally sure if findingway can "just run" in other Discords, even if added. The emojis it uses are present only in APD, and bots can't always use emojis across Discords. If it can't be run in other Discords, I can create a configuration file for mapping roles and jobs to emojis -- someone just open an issue and let me know.

## Deployment

The repository automatically builds Docker images; you can access them if you want to run findingway in your own Discord.

I run this in Fargate for Aether PUG DSR. Here's a task definition you might find useful:

```
{
  "ipcMode": null,
  "executionRoleArn": "arn:aws:iam::AWS_ACCOUNT_ID:role/ecsTaskExecutionRole",
  "containerDefinitions": [
    {
      "dnsSearchDomains": null,
      "environmentFiles": null,
      "logConfiguration": {
        "logDriver": "awslogs",
        "secretOptions": null,
        "options": {
          "awslogs-group": "/ecs/findingway",
          "awslogs-region": "us-east-1",
          "awslogs-stream-prefix": "ecs"
        }
      },
      "entryPoint": null,
      "portMappings": [],
      "command": null,
      "linuxParameters": null,
      "cpu": 0,
      "environment": [
        {
          "name": "DATA_CENTRE",
          "value": "Aether"
        },
        {
          "name": "DISCORD_CHANNEL_ID",
          "value": "your discord channel ID"
        },
        {
          "name": "DISCORD_TOKEN",
          "value": "your discord token"
        },
        {
          "name": "DUTY",
          "value": "Dragonsong's Reprise (Ultimate)"
        }
      ],
      "resourceRequirements": null,
      "ulimits": null,
      "dnsServers": null,
      "mountPoints": [],
      "workingDirectory": null,
      "secrets": null,
      "dockerSecurityOptions": null,
      "memory": null,
      "memoryReservation": null,
      "volumesFrom": [],
      "stopTimeout": null,
      "image": "ghcr.io/veraticus/findingway:main",
      "startTimeout": null,
      "firelensConfiguration": null,
      "dependsOn": null,
      "disableNetworking": null,
      "interactive": null,
      "healthCheck": null,
      "essential": true,
      "links": null,
      "hostname": null,
      "extraHosts": null,
      "pseudoTerminal": null,
      "user": null,
      "readonlyRootFilesystem": null,
      "dockerLabels": null,
      "systemControls": null,
      "privileged": null,
      "name": "findingway"
    }
  ],
  "placementConstraints": [],
  "memory": "512",
  "taskRoleArn": "arn:aws:iam::AWS_ACCOUNT_ID:role/ecsTaskExecutionRole",
  "compatibilities": [
    "EC2",
    "FARGATE"
  ],
  "taskDefinitionArn": "arn:aws:ecs:us-east-1:AWS_ACCOUNT_ID:task-definition/findingway:1",
  "family": "findingway",
  "requiresAttributes": [
    {
      "targetId": null,
      "targetType": null,
      "value": null,
      "name": "com.amazonaws.ecs.capability.logging-driver.awslogs"
    },
    {
      "targetId": null,
      "targetType": null,
      "value": null,
      "name": "ecs.capability.execution-role-awslogs"
    },
    {
      "targetId": null,
      "targetType": null,
      "value": null,
      "name": "com.amazonaws.ecs.capability.docker-remote-api.1.19"
    },
    {
      "targetId": null,
      "targetType": null,
      "value": null,
      "name": "com.amazonaws.ecs.capability.task-iam-role"
    },
    {
      "targetId": null,
      "targetType": null,
      "value": null,
      "name": "com.amazonaws.ecs.capability.docker-remote-api.1.18"
    },
    {
      "targetId": null,
      "targetType": null,
      "value": null,
      "name": "ecs.capability.task-eni"
    }
  ],
  "pidMode": null,
  "requiresCompatibilities": [
    "FARGATE"
  ],
  "networkMode": "awsvpc",
  "runtimePlatform": {
    "operatingSystemFamily": "LINUX",
    "cpuArchitecture": null
  },
  "cpu": "256",
  "revision": 1,
  "status": "ACTIVE",
  "inferenceAccelerators": null,
  "proxyConfiguration": null,
  "volumes": []
}
```
