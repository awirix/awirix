# Welcome to the Vivi's templates contribution guide!

## How to contribute

1. Fork the repository
2. Create a new branch
3. Make your changes
4. Open a pull request
5. You are awesome! :tada:

## What languages are supported?

Any language that compiles to **Lua 5.1** is supported.

## How to add a new language template

1. Create a new folder with the name of your language in the [languages](./languages) folder
2. Create the following files:

### `README.md.tmpl`

This file will be used as the template for the `README.md` file of the new extension.
Try to keep it as simple as possible, and make it look similar to the other README's.

Example:

```markdown
# {{ .Name }} {{ if .NSFW }}🔞{{ end }}

{{ .About }}

> Lua scraper for [vivi](https://github.com/vivi-app/vivi)
```

Available variables for all `*.tmpl` files:

-   `Name`: The name of the extension ~ string
-   `About`: The description of the extension ~ string
-   `NSFW`: Whether the extension is NSFW or not ~ bool
-   `Libs`: Paths to the libraries used by the lua scraper ~ []string

> [Learn golang templates syntax](https://golang.org/pkg/text/template/)

### `scraper.your-language-extension`

This file will be used as the template for the new extension scraper.

It should contain the empty code for the scraper, and it should be a valid file for the language you are using.

Take a look at the [lua](./languages/lua/scraper.lua) and [fennel](./languages/fennel/scraper.fnl) scrapers for example.

### `tester.your-language-extension`

This file will be used as the template for the new extension tester.

It should contain the empty code for the tester, and it should be a valid file for the language you are using.

Take a look at the [lua](./languages/lua/tester.lua) and [fennel](./languages/fennel/tester.fnl) testers for example.

### `scraper.lua` & `tester.lua`

1. If compiler for your language is written in Lua 5.1 without any runtime dependencies, you should call it from these files.
   Take a look at the [fennel](./languages/fennel/scraper.lua) and [teal](./languages/teal/scraper.lua) for examples.
2. Otherwise (non-lua compiler, different version, etc...) you should include an error message in these files
   and point to the compilation instructions in the [README.md.tmpl](#readmemdtmpl) file.
   See the [yue](./languages/yue/scraper.lua) and [typescript](./languages/typescript/scraper.ts) for examples.

### Optional files

It's a good practice to include some formatter configuration files, so the code is consistent across all extensions.
For example, [stylua.toml](./languages/lua/stylua.toml) for Lua
