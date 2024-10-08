# Loops description

This document presents more technical (algorithms, data structures) description about loops used in mmet.

Please review [business-case.md](business-case.md) before proceeding.

Status: Doc is still under construction.

## Architecture 
Let's first discuss the components of our system before delving into the algorithms used in the loops.

We can distinguish 4 components:
- `monitored-system` - Represents the MME Network. From the perspective of the closed control loop, it acts as the **managed system**.
- `translation agent` - Serves as a bridge between the loops and the external world; in our case, the external world is `mmet`. It:
  - Retrieves data from `monitored-system`
  - Sends action commands to `monitored-system`
- `reactive loop` - This loop ensures and enforces the demanded traffic distribution.
- `deliberative loop` - Supervising the reactive loop, this component can modify distribution values.

Table below give some insights about technical implementation/placement of these components.

| Component name:   | Implemented as:                                        | Belongs to our platform? |
| ----------------- | ------------------------------------------------------ | ------------------------ |
| monitored-system  | Linux process, likely a Go application                 | No (external)            |
| translation agent | Kubernetes pod, likely a Go application in a container | Yes (entry point)        |
| reactive loop     | Custom Resource API Objects in Kubernetes              | Yes                      |
| deliberative loop | Custom Resource API Objects in Kubernetes              | Yes                      |

## monitored-system

### General descr

This Linux process generates a quadruple each round, where each value represents the number of sessions currently served by the corresponding MME node.
Each round, a random number (which can be negative) is added to each node's count. 

![](img/1.png)

Additionally, the program can receive commands specifying traffic movement from one node to another ("move `x` units of traffic from node `a` to node `b`"), such as:

![](img/2.png)

### Implementation

See [monitored-system/main.go](monitored-system/main.go)

## Translation Agent

### General descr

Translation Agent is a linux process that:

- Operates an HTTP server and listens to requests from `mmet`
- Utilizes an HTTP client to send Move Commands to `mmet`
- Utilizes the kube-api-server client to push data into the `reactive loop`

It is containerized using Docker and deployed as a pod in a Kubernetes cluster (Minikube).

### Implementation

See [translation-agent/main.go](translation-agent/main.go)

## Reactive Loop

It will work as something like this:

![](img/3.png) 
