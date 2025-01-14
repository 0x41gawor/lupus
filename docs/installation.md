# Installation

Installation of Lupus requires technical skills, some knowledge of programming, Kubernetes, Linux system etc..

Lupus is implemented as a [kubebuilder](https://book.kubebuilder.io) project. The recommended way of using Lupus is to clone this repo and embody the role as if you were the developer of this kubebuilder project. 

There is no thing as installation of Lupus (e.g. in your Operating System or so). You can install [Custom Resource Definitions](defs.md#crd) of [Lupus elements](defs.md#lupus-element) in your cluster and run [controllers](defs.md#controller) for them. And this is what this document is about.

## Steps

### 1. Prerequisities

Of course you need to have a running Kubernetes cluster (can be [Minikube](https://minikube.sigs.k8s.io/docs/)), some container engine (like [Docker](https://docs.docker.com)) and [Go language](https://go.dev) installed. 

#### 1.1 Install Kubebuilder

Follow [this tutorial](https://book.kubebuilder.io/quick-start).

### 2. Clone this repo

```sh
git clone  https://github.com/0x41gawor/lupus
cd lupus
```

### 3. Install CRD into the cluster

This will apply [crd](defs.md#crd) of [master](defs.md#master) and [element](defs.md#lupus-element) to the cluster enabling you to use them.

```sh
make install
```

### 4. Run controllers for master and element custom resources

This will run [controller](defs.md#controller) for [master](defs.md#master) and [element](defs.md#lupus-element) custom resources. 

```sh
make run
```

If you want to run these controllers as pods in you Kubernetes cluster get familiar with [kubebuilder](https://book.kubebuilder.io) platform.

But I do not recommend such approach unless you know you won't modify your [lupus-controllers](defs.md#lupus-controllers). 

## Next steps

Get familiar with [getting-started](getting-started.md) guide.
