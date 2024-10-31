## ordo
- ordo is a simplified container orchestrator written in Go. It provides a practical understanding of orchestrator concepts like task management, scheduling, and worker communication.




### Features

#### Task Management

- [x] **Task States**: Implement task states (Pending, Running, Completed, Failed).
- [x] **Docker Integration**: Start/stop Docker containers via Docker API.
- [x] **Task Persistence**: Store tasks in-memory or in BoltDB.

#### Worker Features
- [x] **Task Queue**: FIFO queue for processing tasks.
- [x] **Task Execution**: Run assigned tasks as Docker containers.
- [ ] **Metrics Collection**: Collect CPU, memory, disk usage data. **(Pending)**

#### Manager Features
- [x] **Task Scheduling**: Basic task scheduling (Round-Robin).
- [ ] **Enhanced Scheduler**: Resource-based E-PVM scheduler. **(Pending)**
- [ ] **Health Checks**: Check task health, auto-restart on failure. **(Pending)**

#### General Features
- [x] **Modular Design**: Separate modules for tasks, workers, managers.
- [ ] **Testability**: Set up automated testing. **(Pending)**

#### Limitations & Future Enhancements
- [ ] **Security**: Add security measures. **(Pending)**
- [ ] **Service Discovery**: Enable service discovery for tasks. **(Pending)**
- [ ] **High Availability**: Implement redundancy for manager/workers. **(Pending)**
