# README

Frontend for _upda_.

## Configuration

* During development: `.env.development`
* Production derives the values from their main GoLang application, some environment variables are coupled with their
  backend one, some are exposed via `WEB_*`, but internally mapped to `VITE_` in the `app.go` file

## Development & contribution

Contributions are very welcome!

### Prerequisites

It's probably worth checking out a node environment manager like [nvm manager](https://github.com/nvm-sh/nvm).

Required node and npm versions are outlined in the `package.json`.

### Setup instructions

Run `npm install` which should install all dependencies.

### Start

Use the `npm run start` command to start the development setup. Backend should be running.

### Translation files

Pipeline checks for existing i18n keys, if you need to manually sync any language file, run `npm run i18n-sync`.

### Dependencies

To track unused dependencies

```shell
# install (if not yet present)
npm install -g depcheck typescript

# run
depcheck
```
