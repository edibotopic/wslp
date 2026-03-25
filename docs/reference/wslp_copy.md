## wslp copy

Copy a WSL distribution under a new name

### Synopsis

Copy a WSL distribution by exporting it and importing it under a new name.

The new distribution is stored in %USERPROFILE%\WSLCopies\<new-name> by default.
You can override this with the --install-dir flag.

```
wslp copy <source> <new-name> [flags]
```

### Options

```
  -h, --help                 help for copy
  -d, --install-dir string   Directory to store the new distro's virtual disk (overrides default)
```

### SEE ALSO

* [wslp](wslp.md)	 - A tool for managing WSL instances.

