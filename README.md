<!---[![go-build status](https://github.com/rdnt/myst/workflows/go-build/badge.svg)](https://github.com/SHT/myst/actions?query=workflow%3Ago-build)
[![vue-build status](https://github.com/rdnt/myst/workflows/vue-build/badge.svg)](https://github.com/SHT/myst/actions?query=workflow%3Avue-build)-->

# myst 

### Zero-knowledge, end-to-end encrypted password manager

---

#### Development requirements
- Go v1.20
- Node 18

---

#### Setting up

On Windows, you might need to set script execution
policy to unresticted. Open an elevated Powershell terminal and
type the following:

```ps1
Set-ExecutionPolicy Unrestricted
```

Then to build the builder binary, run

```bash
.\script\bootstrap.ps1
```

This will create the `myst.exe` binary
on the root of the project. Use that to build the actual
project:

```bash
.\myst build
```

```bash

npm install -g openapi-typescript-codegen
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
go generate ./...
```
