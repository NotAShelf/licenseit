# Licenseit

Small, simple and streamlined utility for generating license files based on
customizable templates, for easily and quickly creating license files.

## Why?

I used to have aliases that `curl` project sites, or my own repository for
storing license files. Licenseit embeds any templates that use so that I can
call it anywhere to quickly generate a license.

## Usage

You must pass the name of the license template you wish to use, and the name of
the author that will be put in the license file.

```bash
licenseit gpl3 -author "John Doe"
```

Alternatively, if you don't want to pass the name of the author each time, you
can have licenseit read the author from a configuration file.

```json
{
  "author": "Your Name"
}
```

## Adding New Templates

License templates are embedded into the program using the
[embed package](https://pkg.go.dev/embed). You may add new licenses to the
[templates directory](./templates) using a very simple format.

1. The license must have `.txt` or `.md` as its extension. Licenseit will look
   for `name.txt`, `name.md` or `name`.
2. License name must be unique. There can't be duplicates of the same license
   with different extensions.
3. `{author}` and `{year}` are replaced in the generated config, but they are
   not strictly necessary in the template files.

Feel free to open a pull request to add more licenses.
