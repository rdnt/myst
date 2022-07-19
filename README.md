<!---[![go-build status](https://github.com/rdnt/myst/workflows/go-build/badge.svg)](https://github.com/SHT/myst/actions?query=workflow%3Ago-build)
[![vue-build status](https://github.com/rdnt/myst/workflows/vue-build/badge.svg)](https://github.com/SHT/myst/actions?query=workflow%3Avue-build)-->

# myst 

#### Zero-knowledge, end-to-end encrypted password manager


### Development requirements
- Go v1.18
- Node 16+

# Bootstrap

```bash
npm install -g openapi-typescript-codegen
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
go generate ./...
```
