
# Myst 

[![Build status](https://github.com/rdnt/myst/actions/workflows/build.yml/badge.svg)](https://github.com/rdnt/myst/actions/workflows/build.yml)

### Zero-knowledge, end-to-end encrypted password manager.

---

#### Development requirements
- Go v1.20
- Node.JS 18
- golangci-lint v1.53.2 https://golangci-lint.run/usage/install/

---

#### Bootstrap

On Windows, you might need to set script execution
policy to unresticted. Open an elevated Powershell terminal and
type the following:

```powershell
Set-ExecutionPolicy Unrestricted
```

Then to build the builder binary, run

```bash
.\script\bootstrap.ps1
```

This will create the `myst.exe` binary
on the root of the project. You can use that to build the actual
project.

### Building

```bash
.\myst build
```
