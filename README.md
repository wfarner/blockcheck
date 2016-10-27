# blockcheck
A tool to help keep Markdown code blocks in sync with a file contents

## Overview

Keeping documentation in sync with code and resource files is hard.  It's even harder since most common Markdown
processors (like GitHub's) doesn't allow you to include other files.

Blockcheck works around this by making it easy to check that a code block matches another file.  You'll still have the
file contents repeated, but if you include blockcheck in CI you can prevent mismatches.  With this in place you can
confidently do things like have inline source files without worrying about maintaining both; or write unit tests against
code blocks in your docs!

## Usage

Blockcheck works by scanning Markdown files for sequences like:

<pre>
&lt;!-- blockcheck file.txt --&gt;
```
A code block, which
must match file.txt
```
</pre>

Try it!  The block below is kept in sync with `file.txt` in this repository.

```shell
$ blockcheck README.md
```

<!-- blockcheck file.txt -->
```
A code block, which
must match file.txt
```

You can also pipe file names for convenient chaining with `find`:
```shell
$ find . -name '*.md' | blockcheck
```

That's it!

## Installing

```shell
$ go get -u github.com/wfarner/blockcheck
```

## Contributing

Issues and pull requests welcome!  Before submitting patches, please verify that `make` passes.
