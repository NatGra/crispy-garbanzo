package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	yaml "gopkg.in/yaml.v2"
)

type config struct {
	DockerUser          string `yaml:"dockerUsername"`
	DockerPwd           string `yaml:"dockerPassword"`
	DockerServerAddress string `yaml:"dockerServerAddress"`
}

func main() {
	c, err := readConfig("example-config.yaml")
	if err != nil {
		panic(err)
	}
	print("c.DockerUser: ", c.DockerUser, "\n")

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	// Listing images
	listImages(ctx, cli)

	// Listing containers
	listContainers(ctx, cli)

	// Retagging and pushing an image to an image repository
	retagAndPushImage(ctx, cli, c, "docker.io/library/alpine", "natgra/sample-image-repo:latest")
}

func readConfig(filename string) (*config, error) {
	yamlConfig, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("error while reading config: %v ", err)
	}

	c := &config{}
	err = yaml.Unmarshal(yamlConfig, c)
	if err != nil {
		log.Fatalf("error while unmarshalling: %v", err)
	}

	return c, nil
}

func listImages(ctx context.Context, cli *client.Client) {
	print("Listing available images:\n")

	images, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	for _, image := range images {
		fmt.Println(image.ID)
	}
}

func listContainers(ctx context.Context, cli *client.Client) {
	print("Listing available containers:\n")

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Println(container.ID)
	}
}

func retagAndPushImage(ctx context.Context, cli *client.Client, config *config, source string, target string) {
	print("Pulling image from original container registry.\n")
	out, err := cli.ImagePull(ctx, source, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	defer out.Close()

	print("Retagging image.\n")
	if err := cli.ImageTag(ctx, source, target); err != nil {
		panic(err)
	}

	print("Pushing image to another container registry.\n")
	authConfigBytes, _ := json.Marshal(
		types.AuthConfig{
			Username:      config.DockerUser,
			Password:      config.DockerPwd,
			ServerAddress: config.DockerServerAddress,
		})
	authConfigEncoded := base64.URLEncoding.EncodeToString(authConfigBytes)
	out, err = cli.ImagePush(ctx, target, types.ImagePushOptions{RegistryAuth: authConfigEncoded})
	if err != nil {
		panic(err)
	}

	io.Copy(os.Stdout, out)
}
