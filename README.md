# Example project in Go using the Docker Engine SDKs and Docker API

This is a test project to play around with the Docker Engine SDKs and Docker API for Go.

## Build and run the project
### Preconditions:
An account for [Docker](https://hub.docker.com/) is recommended to set the parameters in the `example-config.yaml`.
The values to be set are as follows:
```yaml
dockerUsername:      "Your Docker Hub Username"
dockerPassword:      "Your Docker Hub Password or Access Token"
dockerServerAddress: "https://index.docker.io/v1/"
```

### Build the project
Use the following commands:
```sh
# Build the project
go build
```

### Run the project
```
# Run the main.go
go run main.go

# Or, run the built service
./crispy-garbanzo
```
