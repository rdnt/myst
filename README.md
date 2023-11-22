
# Myst 

[![Build status](https://github.com/rdnt/myst/actions/workflows/build.yml/badge.svg)](https://github.com/rdnt/myst/actions/workflows/build.yml)

![image](https://github.com/rdnt/myst/assets/17600197/d4cabfca-ba75-4426-92fe-704c9ab4ea36)


### Zero-knowledge, end-to-end encrypted password manager.

# /!\ Notice

Myst is a proof of concept that was built as part of my thesis, and although it could be a usable app, it was built in isolation
and without much user feedback.  
The security and encryption practices could have edge-cases or be completely broken, since they have not been validated by
external contributors and the test coverage is not as high as I would want it to be.

I have not yet decided what the future of the project will be, although I would hope it could be useful to some people.  
A useful at least to me path is for it to become a CLI password manager with password sharing within the local network.

If you have any ideas/suggestions or are interested in collaborating, please either reach out to me or fork away, whatever you prefer.

---

#### Development requirements
- Go v1.21
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
