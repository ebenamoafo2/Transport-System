//go:build integration_test

package test_containers

import (
	"context"
	"fmt"
	"sync"

	"github.com/moby/moby/api/types/container"
	mobycontainer "github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	"github.com/rs/zerolog/log"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	mysqlOnce sync.Once
	mysqlHost string
	mysqlPort string
)

type testContainerRunner struct {
	servicePort       int
	name              string
	image             string
	exposedPorts      []string
	env               map[string]string
	hostConfigModifier func(hostConfig *mobycontainer.HostConfig)
}

func (r testContainerRunner) Run(ctx context.Context) (testcontainers.Container, error) {
	containerReq := testcontainers.ContainerRequest{
		Name:         r.name,
		Image:        r.image,
		ExposedPorts: r.exposedPorts,
		Env:          r.env,
		WaitingFor:   wait.ForListeningPort(fmt.Sprintf("%d/tcp", r.servicePort)),
		HostConfigModifier: func(hostConfig *mobycontainer.HostConfig) {
			if r.hostConfigModifier != nil {
				r.hostConfigModifier(hostConfig)
			}
		},
	}

	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: containerReq,
		Started:          true,
	})
}

func GetMySqlContainer(ctx context.Context, db, user, pass string, port *int) (string, string) {
	mysqlOnce.Do(func() {
		c, err := testContainerRunner{
			servicePort:  3306,
			name:         "mysql",
			image:        "mysql:8.0.44",
			exposedPorts: []string{"3306/tcp"},
			env: map[string]string{
				"MYSQL_ROOT_PASSWORD": pass,
				"MYSQL_DATABASE":      db,
				"MYSQL_USER":          user,
				"MYSQL_PASSWORD":      pass,
			},
			hostConfigModifier: mysqlHostConfigModifier(port),
		}.Run(ctx)

		if err != nil {
			log.Fatal().Err(err).Msg("failed to start mysql container")
		}

		mp, err := c.MappedPort(ctx, "3306/tcp")
		if err != nil {
			log.Fatal().Err(err).Msg("failed to get mapped port for mysql container")
		}

		h, err := c.Host(ctx)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to get host for mysql container")
		}

		mysqlHost = h
		mysqlPort = mp.Port()
	})
	return mysqlHost, mysqlPort
}

func mysqlHostConfigModifier(port *int) func(hostConfig *container.HostConfig) {
	return func(hostConfig *container.HostConfig) {
		hostConfig.AutoRemove = true
		if port != nil {
			hostConfig.PortBindings = network.PortMap{
				"3306/tcp
				": []network.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: fmt.Sprintf("%d", *port),
					},
				},
			}
		}
	}
}