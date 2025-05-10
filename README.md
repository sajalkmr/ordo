# Ordo

Ordo is a container orchestrator written in Go that provides a practical implementation of container orchestration concepts. It allows you to manage and schedule Docker containers across multiple nodes.

## Features

- **Task Management**: Create, stop, and monitor container tasks
- **Node Management**: Add and manage worker nodes
- **Task Scheduling**: Basic task scheduling across available nodes
- **Status Monitoring**: View task and node status
- **Docker Integration**: Direct integration with Docker for container management

## Requirements

- Go 1.16 or later
- Docker
- BoltDB (v1.3.1)
- chi (v5.0.3)

## Project Structure

- `cmd/`: Command-line interface implementation
- `manager/`: Manager node implementation
- `worker/`: Worker node implementation
- `task/`: Task management and scheduling
- `store/`: Data persistence layer
- `utils/`: Utility functions
- `node/`: Node management
- `scheduler/`: Task scheduling logic
- `stats/`: System statistics collection

## Getting Started

1. Clone the repository
2. Install dependencies: `go mod download`
3. Build the project: `go build`
4. Run the manager: `./ordo manager`
5. Add worker nodes: `./ordo node add <node-address>`
6. Create tasks: `./ordo run <task-config>`

## Example Task Configuration

```json
{
    "name": "example-task",
    "image": "nginx:latest",
    "ports": ["80:80"],
    "env": ["ENV=production"]
}
```

## License

MIT License
