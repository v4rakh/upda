# README

## Configuration

* During development: `.env.development`
* Production derives the values from their main GoLang application, some environment variables are coupled with their
  backend one, some are exposed via `WEB_INTERFACE_*`, but internally mapped to `VITE_` in the `server.go` file

## Development & contribution

Contributions are very welcome!

### Prerequisites

It's probably worth checking out a node environment manager like [nvm manager](https://github.com/nvm-sh/nvm).

Required node and npm versions are outlined in the `package.json` in the `"engines"` section.

The application uses `pnpm`. Make sure to install it. For this you can use your global npm installation and invoke
`npm install --global pnpm` or have it as system package.

### Setup instructions

Run `pnpm install` which should install all dependencies.

### Start

Use the `pnpm start` command to start the development setup. Backend should be running.

### Translation files

Pipeline checks for existing i18n keys, if you need to manually sync any language file, run `pnpm i18n-sync`.

#### New translation languages

Adapt `package.json` and add the respective key to `i18n-sync` and `checkstyle:i18n` for the `--languages` argument,
e.g. `--languages en,new`.

Run `pnpm i18n-sync` to create and sync new language

Adapt `languages.ts`

```typescript
enum Languages {
    // ...
    NEW = 'new'
}
```

Adapt `locales.ts`

```typescript
enum Locales {
    // ...
    NEW = 'new-NEW'
}
```

Adapt `LocaleProvider.tsx` to detect the new language, assign proper _momentjs_ locale.

```typescript
// add import for Ant Design
import new_NEW from 'antd/lib/locale/new_NEW';

// add import for moment
// @ts-ignore
import new_localization from 'moment/dist/locale/new.js';

// getLocaleFromLanguage
if (LANGUAGES.NEW === language) {
    // ...
}

// getLanguageFromLocale
if (LANGUAGES.NEW === language) {
    // ...
}

// setLocales (uses import of moment)
// ...
else if (Languages.NEW === language) {
    // use import for Ant Design
    setAntLocale(new_NEW);
    // use import for moment
    moment.updateLocale(language, de_localization);
}
```

Adapt `i18n/index.ts` and import new language and make i18n detect it (locales are derived from it)

```typescript jsx
import newTranslation from './translations/new.json';

const resources = {
    //...
    new: {
        ...newTranslation
    }
};
```

### Dependencies

To track unused dependencies

```shell
# install (if not yet present)
npm install --global depcheck typescript

# run
depcheck
```

To manage (outdated) dependencies, `pnpm` helps:

```shell
# shows outdated dependencies
pnpm outdated

# upgrades to latest in range definition of package.json
pnpm up

# upgrades to latest available, this includes major upgrades
pnpm up --latest
```
