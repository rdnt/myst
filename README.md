[![go-build status](https://github.com/rdnt/myst/workflows/go-build/badge.svg)](https://github.com/SHT/myst/actions?query=workflow%3Ago-build)
[![vue-build status](https://github.com/rdnt/myst/workflows/vue-build/badge.svg)](https://github.com/SHT/myst/actions?query=workflow%3Avue-build)

Myst: The cloudless password manager.

// keystores
//  - 0000.mst
//  - 0001.mst
// users
//  - user1.json
//  - user2.json

// user1.json example:
// user id:
// user name:
// user argon2id hash encoded
// keystores: [
//  - keystore1: encryption key that is encrypted by : user's master password plus a passphrase
// object with keys: keystore id, value : encrypted encryption key. passphrase ofc not stored
// ]

// if user password is leaked, the passphrase is required to have actual access to the user's keystores
// so 2 credentials
// passphrase is optional
