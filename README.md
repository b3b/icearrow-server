# IceArrow server

This repository represents the backend component of the [IceArrow](https://github.com/kkomelin/icearrow) app.

The IceArrow backend is built on top of:

- [Yopass](https://github.com/jhaals/yopass/tree/7e50bef6aacc5b401149914fd3472404f1b65e5c): A platform for the secure sharing of secrets, passwords, and files.
- [Walrus](https://docs.walrus.site/): A decentralized storage and data availability protocol optimized for large binary files.

## Getting Started

To run the project locally using Docker:

```bash
docker-compose up --build
```

### Accessing the application


- Storage API: Available at http://127.0.0.1/secret
- Walrus HTTP API: Available at http://127.0.0.1:31415

### API Usage Example

For examples of using the service and interacting with the Walrus API, refer to [walrus.ipynb](./notebooks/walrus.ipynb).

## Go Packages

- `walrus/walrus.go`: Implements `WalrusClient`, a client to interact with the Walrus HTTP API.
- `hippo/hippo.go`: Defines the `Hippo` struct, connecting with the Walrus storage via `WalrusClient`.

## Commands

- `icearrow-server`: Command to run the server.
- `walrus-client`: Wrapper for the Walrus package.
