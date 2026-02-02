# WSL Plus (wslp)

A CLI tool that wraps the WSL api, making common and bulk operations easier.

## Documentation

To build and preview the documentation locally:

```bash
cd docs
make run
```

This will:

- Install documentation dependencies
- Auto-generate CLI reference docs from the code
- Start a development server at http://127.0.0.1:8000

The CLI reference documentation in `docs/reference/` is auto-generated from the command definitions and should not be edited manually.

To only generate the CLI reference:

```bash
# /docs/
make cli-ref
```

## Problem

Some WSL tasks are not convenient to perform using the default WSL CLI.

For example:

* Renaming an instance
* Backing up an instance
* Bulk actions on instances (creating, deleting)
* Generating and managing configuration files

## Todo

### Required

- [x] List registered distros
- [x] Show default distro
- [x] Bulk install distros (sequential)
- [ ] Bulk install (concurrent)
- [ ] Rename
- [ ] Register local tarballs
- [ ] Delete distros
- [x] Pretty output

### Stretch

- [ ] Managing configs
- [ ] Managing templates
- [ ] Reporting data
- [ ] Tag management
- [ ] Pro for WSL integration
