package dockersetupcommand

import (
	"fmt"
	"io"
	"log/slog"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

const (
	DEFAULT_NET_INTERFACE = "eth0"
	CLEANUP_CONTAINER_CMD = "tc qdisc del dev %s root"

	IMAGE_NAME          = "alpine:3.24.1"
	START_CONTAINER_CMD = "apk add --no-cache iproute2 && tc qdisc add dev %s root %s && sleep infinity"
)

var _ ContainerThrottler = (*DockerContainerThrottler)(nil)

type ContainerThrottler interface {
	Throttle(cmd string, netInterface string, targetContainerId string) error
	Unthrottle(netInterface string, targetContainerId string) error
}

type DockerContainerThrottler struct {
	helperContainerId string // name or id of the helper container, it will be set after the Throttle method (so the Build method) is called
	builder           HelperContainerBuilder
	starter           ContainerStarter
	stopper           ContainerStopper
	executor          ContainerCommandExecutor
}

func NewDockerContainerThrottler(builder HelperContainerBuilder, starter ContainerStarter, stopper ContainerStopper, executor ContainerCommandExecutor) *DockerContainerThrottler {
	return &DockerContainerThrottler{
		builder:  builder,
		starter:  starter,
		stopper:  stopper,
		executor: executor,
	}
}

func (t *DockerContainerThrottler) Throttle(cmd string, netInterface string, targetContainerId string) error {
	helperContainerId, err := t.builder.Build(targetContainerId, netInterface, cmd)
	if err != nil {
		return fmt.Errorf("failed to build helper container for target container '%s': %w", targetContainerId, err)
	}

	t.helperContainerId = helperContainerId
	err = t.starter.Start(t.helperContainerId)
	if err != nil {
		return fmt.Errorf("failed to start helper container '%s': %w", t.helperContainerId, err)
	}

	return nil
}

func (t *DockerContainerThrottler) Unthrottle(netInterface string, targetContainerId string) error {
	if t.helperContainerId == "" {
		return fmt.Errorf("helper container id is empty, cannot unthrottle")
	}

	cleanupCmd := fmt.Sprintf(CLEANUP_CONTAINER_CMD, netInterface)
	err := t.executor.Execute(t.helperContainerId, cleanupCmd)
	if err != nil {
		return fmt.Errorf("failed to execute cleanup command in helper container '%s': %w", t.helperContainerId, err)
	}

	err = t.stopper.Stop(t.helperContainerId)
	if err != nil {
		return fmt.Errorf("failed to stop helper container for target container '%s' during cleanup: %w", t.helperContainerId, err)
	}

	return nil
}

// Highly optimizable implementation
func (e *DockerHandler) Build(targetContainerId string, netInterface string, cmd string) (string, error) {
	// Pull the alpine image
	resp, err := e.client.ImagePull(e.ctx, IMAGE_NAME, client.ImagePullOptions{}) // check if the image is already present locally, if not pull it. Also it always checks for updates to the image
	if err != nil {
		return "", fmt.Errorf("failed to pull image '%s': %w", IMAGE_NAME, err)
	}
	defer resp.Close() // nolint: errcheck

	// Wait for the image pull to complete and by reading and discarding the response body
	_, err = io.Copy(io.Discard, resp)
	if err != nil {
		return "", fmt.Errorf("failed to read image pull response: %w", err)
	}

	// Configure the container with the specified image, with iproute2 installed and running indefinitely (sleep infinity)
	netIssueCmd := fmt.Sprintf(START_CONTAINER_CMD, netInterface, cmd)

	config := &container.Config{
		Image: IMAGE_NAME,
		Cmd:   []string{"sh", "-c", netIssueCmd},
	}

	// Net admin permission, to manipulate net interface
	hostConfig := &container.HostConfig{
		CapAdd:      []string{"NET_ADMIN"},
		AutoRemove:  true, // Automatically remove the container when it exits
		NetworkMode: container.NetworkMode("container:" + targetContainerId),
	}

	opts := client.ContainerCreateOptions{
		Config:     config,
		HostConfig: hostConfig,
	}

	containerResp, err := e.client.ContainerCreate(
		e.ctx,
		opts,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}

	slog.Info("Helper container created successfully", "containerId", containerResp.ID)

	return containerResp.ID, nil
}
