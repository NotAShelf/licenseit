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
3. `{author}` and `{date}` are replaced in the generated config, but they are
   not strictly necessary in the template files.

For example, the MIT license template would look like this.

```markdown
Copyright (c) {date} {author}

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
```

You may find some licenses
[here](https://github.com/licenses/license-templates), should you wish to add
them here.

Feel free to open a pull request to add more licenses.
