package task

import (
	"context"
	"io"
	"log"
	"math"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"github.com/moby/moby/pkg/stdcopy"
)

type Task struct {
	ID            uuid.UUID
	ContainerID   string
	Name          string
	State         State
	Image         string
	Cpu           float64
	Memory        int64
	Disk          int64
	ExposedPorts  nat.PortSet
	PortBindings  map[string]string
	HostPorts     nat.PortMap
	RestartPolicy string
	RestartCount  int
	HealthCheck   string
	StartTime     time.Time
	FinishTime    time.Time
}

type TaskEvent struct {
	ID        uuid.UUID
	State     State
	Timestamp time.Time
	Task      Task
}

type DockerInspectResponse struct {
	Error     error
	Container *struct {
		State struct {
			Status string
		}
		NetworkSettings struct {
			NetworkSettingsBase struct {
				Ports map[string]string
			}
		}
	}
}

// Config struct to hold Docker container config
type Config struct {
	// Name of the task, also used as the container name
	Name string
	// AttachStdin boolean which determines if stdin should be attached
	AttachStdin bool
	// AttachStdout boolean which determines if stdout should be attached
	AttachStdout bool
	// AttachStderr boolean which determines if stderr should be attached
	AttachStderr bool
	// ExposedPorts list of ports exposed
	ExposedPorts nat.PortSet
	// Cmd to be run inside container (optional)
	Cmd []string
	// Image used to run the container
	Image string
	// Cpu
	Cpu float64
	// Memory in MiB
	Memory int64
	// Disk in GiB
	Disk int64
	// Env variables
	Env []string
	// RestartPolicy for the container ["", "always", "unless-stopped", "on-failure"]
	RestartPolicy string
}

func NewConfig(t *Task) *Config {
	return &Config{
		Name:          t.Name,
		ExposedPorts:  t.ExposedPorts,
		Image:         t.Image,
		Cpu:           t.Cpu,
		Memory:        t.Memory,
		Disk:          t.Disk,
		RestartPolicy: t.RestartPolicy,
	}
}

type Docker struct {
	Client *client.Client
	Config Config
}

func NewDocker(c *Config) *Docker {
	dc, _ := client.NewClientWithOpts(
		client.FromEnv,
		client.WithVersion("1.41"),
		client.WithAPIVersionNegotiation(),
	)
	return &Docker{
		Client: dc,
		Config: *c,
	}
}

type DockerResult struct {
	Error       error
	Action      string
	ContainerId string
	Result      string
}

func (d *Docker) Run() DockerResult {
	ctx := context.Background()
	reader, err := d.Client.ImagePull(ctx, d.Config.Image, types.ImagePullOptions{})
	if err != nil {
		log.Printf("Error pulling image %s: %v\n", d.Config.Image, err)
		return DockerResult{Error: err}
	}
	io.Copy(os.Stdout, reader)

	rp := container.RestartPolicy{
		Name: d.Config.RestartPolicy,
	}

	r := container.Resources{
		Memory:   d.Config.Memory,
		NanoCPUs: int64(d.Config.Cpu * math.Pow(10, 9)),
	}

	hc := container.HostConfig{
		RestartPolicy:   rp,
		Resources:       r,
		PublishAllPorts: true,
	}

	resp, err := d.Client.ContainerCreate(ctx, &container.Config{
		Image:        d.Config.Image,
		Tty:          false,
		Env:          d.Config.Env,
		ExposedPorts: d.Config.ExposedPorts,
	}, &hc, nil, nil, d.Config.Name)
	if err != nil {
		log.Printf("Error creating container using image %s: %v\n", d.Config.Image, err)
		return DockerResult{Error: err}
	}

	if err := d.Client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		log.Printf("Error starting container %s: %v\n", resp.ID, err)
		return DockerResult{Error: err}
	}

	out, err := d.Client.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		log.Printf("Error getting logs for container %s: %v\n", resp.ID, err)
		return DockerResult{Error: err}
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	return DockerResult{ContainerId: resp.ID, Action: "start", Result: "success"}
}

func (d *Docker) Stop(id string) DockerResult {
	log.Printf("Attempting to stop container %v", id)
	ctx := context.Background()
	err := d.Client.ContainerStop(ctx, id, container.StopOptions{})
	if err != nil {
		log.Printf("Error stopping container %s: %v\n", id, err)
		return DockerResult{Error: err}
	}

	err = d.Client.ContainerRemove(ctx, id, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   false,
		Force:         false,
	})
	if err != nil {
		log.Printf("Error removing container %s: %v\n", id, err)
		return DockerResult{Error: err}
	}

	return DockerResult{Action: "stop", Result: "success", Error: nil}
}

func (d *Docker) Remove(id string) DockerResult {
	log.Printf("Attempting to remove container %v", id)
	ctx := context.Background()
	err := d.Client.ContainerRemove(ctx, id, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   false,
		Force:         false,
	})
	if err != nil {
		log.Printf("Error removing container %s: %v\n", id, err)
		return DockerResult{Error: err}
	}
	return DockerResult{Action: "remove", Result: "success", Error: nil}
}

func dockerStateToTaskState(dockerState string) State {
	switch dockerState {
	case "created":
		return Pending
	case "running":
		return Running
	case "exited":
		return Completed
	case "dead":
		return Failed
	default:
		return Pending
	}
}

func (d *Docker) Inspect(id string) DockerInspectResponse {
	ctx := context.Background()
	containerInfo, err := d.Client.ContainerInspect(ctx, id)
	if err != nil {
		return DockerInspectResponse{
			Error: err,
		}
	}

	response := DockerInspectResponse{
		Container: &struct {
			State struct {
				Status string
			}
			NetworkSettings struct {
				NetworkSettingsBase struct {
					Ports map[string]string
				}
			}
		}{},
	}

	response.Container.State.Status = containerInfo.State.Status
	response.Container.NetworkSettings.NetworkSettingsBase.Ports = make(map[string]string)

	// Convert port bindings to the format expected by nat.PortMap
	portMap := make(nat.PortMap)
	for port, bindings := range containerInfo.NetworkSettings.Ports {
		if len(bindings) > 0 {
			portMap[nat.Port(port)] = []nat.PortBinding{
				{
					HostIP:   bindings[0].HostIP,
					HostPort: bindings[0].HostPort,
				},
			}
			// Also store in the response's port map for backward compatibility
			response.Container.NetworkSettings.NetworkSettingsBase.Ports[string(port)] = bindings[0].HostPort
		}
	}

	return response
}
