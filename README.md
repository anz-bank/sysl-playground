# sysl-playground

Sysl Playground is the live playground for [Sysl](http://sysl.io/) which is running in https://sysl-playground.herokuapp.com/.

## Features

- Examples for most features of sysl, include diagram generation, import, export and code generation etc.


## Run locally
### Requirements
- The backend service use  [Sysl](http://sysl.io/) internally to compile the input. Need to install `Sysl` to run locally, the installation guide can be found in
[Sysl docs](https://sysl.io/docs/installation).

### Run
```bash
go run main.go
```
Open http://localhost:3030/ to see the playground website.

## Deployment
- The `Dockerfile` build the application including the dependencies.
- The current pipeline is use [Heroku](https://www.heroku.com/) container registry for the deployment, refer to `.github/workflows/workflow.yml` for the workflow. Just need to configure the `HEROKU_API_KEY` and `HEROKU_APP_NAME` to the secret.
