# REAPER

[English](README.md) | 简体中文

REpository ArchivER（REAPER）是一个用于从任何Git服务器归档 Git 仓库的工具。

- [功能](#功能)
- [安装](#安装)
- [使用方法](#使用方法)
  - [rip](#rip)
  - [run](#run)
  - [daemon](#daemon)
  - [bury](#bury)
- [配置](#配置)
- [存储](#存储)
- [使用 Docker 运行](#使用-docker-运行)
  - [Docker CLI](#docker-cli)
  - [Docker Compose](#docker-compose)

## 功能

- 从任何Git服务器归档 Git 仓库
- 归档用户/组织的仓库（见 [配置](https://github.com/LeslieLeung/reaper/wiki/Configuration#repository))
- 定时任务
- 多种存储类型（见 [存储](#存储)）
- Docker 支持（见 [使用 Docker 运行](#使用-docker-运行)）

## 安装

```bash
curl -sSfL https://raw.githubusercontent.com/LeslieLeung/reaper/main/install.sh | sh -s -- -b /usr/local/bin
```

或从 [Release](https://github.com/LeslieLeung/reaper/releases) 获取。

## 使用方法

你需要创建一个配置文件来使用REAPER。

```yaml
repository:
  - name: reaper
    url: github.com/leslieleung/reaper
    cron: "0 * * * *"
    storage:
      - localFile
      - backblaze
    useCache: True
    allBranches: True
    depth: 0
    downloadReleases: True
    downloadIssues: True
    downloadWiki: True
    downloadDiscussion: True

storage:
  - name: localFile
    type: file
    path: ./repo
  - name: backblaze
    type: s3
    endpoint: s3.us-west-000.backblazeb2.com
    region: us-west-000
    bucket: your-bucket-name
    accessKeyID: your-access-key-id
    secretAccessKey: your-secret-access-key
```

然后，你可以使用配置文件运行REAPER。

```bash
reaper -c config.yaml
# 或者如果你的配置文件名为config.yaml，只需调用reaper
reaper
```

### rip

`rip`命令会归档在配置中定义的单个 Git 仓库。

```bash
reaper rip reaper
```

### run

`run`命令会归档在配置中定义的所有 Git 仓库。

```bash
reaper run
```

结合cron，你可以定期归档 Git 仓库。

### bury

`bury`命令会归档指定 Git 仓库的所有发布产物。

```bash
reaper bury reaper
```

### daemon

`daemon`命令会启动一个守护进程，它会在后台运行，归档在配置中定义的所有 Git 仓库。

```bash
reaper daemon
# 使用 nohup 后台运行
nohup reaper daemon &
```

## 配置

有关配置，你可以查看此[示例](config/example.config.yaml)。

更多细节，可查看[配置文档](https://github.com/LeslieLeung/reaper/wiki/Configuration)。

## 存储

REAPER支持多种存储类型。

- [x] 文件
- [x] AWS S3

## 使用 Docker 运行

### Docker CLI

一次性运行。 
- 修改 `${pwd}/config/example.config.yaml` 为你的配置文件本地路径。
- 自定义 `${pwd}/repo:/repo` 为你需要的存储路径。容器内路径需要与配置文件中的路径一致。

```bash
docker run --rm \
    -v ${pwd}/config/example.config.yaml:/config.yaml \
    -v ${pwd}/repo:/repo \
    leslieleung/reaper:latest \
    run
```

### Docker Compose

示例Compose配置，见 [docker-compose.yml](docker-compose.yml)。

```bash
git clone https://github.com/leslieleung/reaper.git
docker compose up -d
```

## 常见问题

见 [FAQ](https://github.com/LeslieLeung/reaper/wiki/FAQ)。

## Stargazers over time

[![Stargazers over time](https://starchart.cc/LeslieLeung/reaper.svg)](https://starchart.cc/LeslieLeung/reaper)
