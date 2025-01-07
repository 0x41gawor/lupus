# Getting started

## Target audience

Lupus is a platform targeted for management of [systems](defs.md#managed-system) in telco/mobile networks industry. Lupus proposes such management in an automated way by [Closed Control Loops](defs.md#closed-control-loops). Actually, Lupus does not managed itself. It rather provides a framework that let's its user to design and run loop workflow. The actual runtime components of Lupus loop are external to it (and called [external blocks](defs.md#external-block)). Lupus is implemented on top of Kubernetes and leverages its [controller pattern](defs.md#controller-pattern). Thus, due to the speed of Kubernetes itself, it is destined for non-realtime parts of telco systems.

Here are some requirements and assumptions for Lupus platform:
- Targeted for telco industry (especially mobile networks).
- Will enable to design and run any closed control loop architecture proposed in [Overview of Prominent Control Loop Architectures](https://www.etsi.org/deliver/etsi_gr/ENI/001_099/017/02.01.01_60/gr_ENI017v020101p.pdf).
- Processes managed by Lupus will be non-realtime.
- Implemented on top of Kubernetes, leveraging its controller pattern.
- Actual runtime components of Lupus loop are external to it (e.g. These are some HTTP Servers, especially [Open Policy Agents](https://www.openpolicyagent.org))

## What Lupus can do

Here we will discuss what is in the scope of Lupus responsibilities and what is not. The best idea to explain such topic is to provide some simple example. The example will be simplified as much as it can be and not applicable in real world.

Let's say you are monitoring a temperature in 4 rooms named "A", "B", "C", "D". Each room has a knob on its radiator that can be controlled remotely. You have a database with records of desired temperature for each room. Building administrator from time to time changes this values based on the current season and weather. In your headquarters you have a HTTP server that exposes an endpoint in which you can present room name and it current temperature and endpoint will answer with instruction how to move the radiator knob to get closer to desired temperature. As for now, each room has a microcontroller connected the thermometer and radiator.

Your mission is to keep the desired temperature in each room.

![](../_img/39.png)

As for now a microcontroller in each room periodically quries HTTP server in HQ. In this simple case this is OK, but [reconcillation](defs.md#reconcilliation) logic is performed by controller, locally. It can cause several problems:
- What if logic changes? you will have to reinstall program in each controller
- What if logic is more complex? And either it is hard to code it in controller's native language or controller has no sufficient computing power to perform it

The next problem is that. Imagine company where such "temperature management" is only one aspect for management. One among of thousands more. It would be nice to have a single "place" of management in company.

It would be better that microcontrolle's program servers only as interface between its devices and external world, and that the loop workflow is derived somewhere else. 

**Here is where the Lupus comes into play!**

First let's deprive microcontrollers of the loop workflow. Reprogram them so they only can be queried for current temperature and instructed how to move radiator knob. 

Let's say that you have some datacenter building with Kubernetes cluster. Here we can install Lupus components that will perform loop workflow. 

New architecture would look like this:

![](../_img/40.png)

The [Ingress-Agent](defs.md#ingress-agent) gathers temperature information from all rooms, send it to Lupus. Lupus performs loop logic (query HQ server in this case) and sends its result to [Egress-Agent](defs.md#egress-agent) which then sends appropriate actions to microcontrollers in rooms.

It is worth to note that loop workflow is designed in Kubernetes YAML definition files and Kubernetes serves as loop runtime. 

Such approach gives several advantages:
- Single place (Kubernetes cluster) for all management loops in company
- If the loop workdflow will change, the change has to be done only here
- The loop workflow is coded in [LupN](defs.md#lupn) a notation dedicated to loops

### What Lupus cannot do?

As you can note in the example above we still need the HQ server that performs the loop logic. It is not the Lupus mission to perform actual calculations. Lupus controls only the loop workflow. Here the example is very easy, but imagine the case, where loop has to query multiple servers and make decision based on the responses along the way.

> Of course, sometimes it is not it is overkill to deploy dedicated server with some little loop logic on it. In such case Lupus offers [go code snippets](defs.md#user-functions) possible to be inserted in loop workflow.

### Lupus mission

Lupus mission is more to design loop workflow and introcude single point of management of such loops in an organisation, rather than stand as loop logic runtime (this is done in [external-blocks](defs.md#external-block)).

## Step by step guide

Ok, let's use the same example to show step by step process with detailed description how to use Lupus for closed control loop.

The most general architecture of Closed Control Loop looks like this:

![](../managed-systems/_img/1.svg)

The four rooms state here as [managed-system](defs.md#managed-system) and the [management-problem](defs.md#management-problem) is to keep their temperatures at desired level.

![](../_img/41.png)

### 1. Develop Ingress Agent

General architecture of Lupus as controller looks like this:

![](../_img/42.png)

Ingress and Egress Agents are input and output points to the Lupus system. They are external to Kubernetes and have to be developed by [designer](defs.md#designer). The only requirement for them is to be compatible (stick along) with [Lupin and Lupout interfaces specificatin](lupin-lupout-if.md). 

Our loop will have only one Lupus Element, thus it will be both the [Ingress Element](defs.md#ingress-element) and [Egress Element](defs.md#egress-element).

Our Ingress Agent on one side has to communicate with managed-system, and on the other has to implement [Lupin interface](lupin-lupout-if.md).

![](../_img/43.png)

On the left side it will act as a MQTT Broker and will gather temperature measurement from microcontrollers. Periodically it will change the Status of Element CR* with such json input:

```json
{
    "A": {
        "temp": 12
    },
    "B": {
        "temp": 10
    },
    "C": {
        "temp": 8
    },
    "D": {
        "temp": 14
    },
}
```

> *Python offers great library for interworking with kube-api-server

### 2. Design the Lupus Elements

