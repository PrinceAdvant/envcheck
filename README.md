# envcheck

A small Go library and CLI that compares a `.env` file with its template, usually `.env.example`, and reports keys that are missing, empty, or undocumented. Run it before `docker compose up` or as a CI step so a forgotten variable fails the deploy instead of crashing the app at runtime.

## Install

```
go install github.com/PrinceAdvant/envcheck/cmd/envcheck@latest
```

## Usage

```
envcheck -template .env.example -env .env
```

Output:

```
missing: MAIL_FROM
empty:   REDIS_URL
extra:   DEBUG
```

Exit code 1 when any key is missing or empty. Add `-strict` to also fail on keys that are not in the template, which keeps the template honest as documentation.

The parser understands comments, blank lines, an optional `export ` prefix and quoted values. It does not expand `${VAR}` references; that is the job of Docker Compose or your shell.

## In a deploy script

```
envcheck -template .env.example -env .env && docker compose up -d
```

## Why this matters with Compose

Docker Compose substitutes an unset variable with an empty string and only prints a warning, which is easy to miss in CI output. Values can also come from several places (shell, `.env`, `env_file`, `environment:`) with a precedence order that surprises people. This guide to [environment variables in Docker Compose](https://peon.sh/blogs/environment-variables-docker-compose) explains which value wins.

Maintained by people who work on [Peon](https://peon.sh/), an open-source platform for deploying Docker apps to your own servers. MIT licensed.
