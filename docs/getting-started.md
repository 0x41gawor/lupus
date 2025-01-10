# Getting started

## Target audience

Lupus is a platform targeted for autoamted management of [systems](defs.md#managed-system) in telco/mobile networks industry. Lupus proposes such management by the means of [Closed Control Loops](defs.md#closed-control-loops). Actually, Lupus does not manage itself. It rather provides a framework that let's its users to design and run loop workflow. The actual processing components of Lupus loop are external to it (and called [external blocks](defs.md#external-block)). Lupus is implemented on top of Kubernetes and leverages its [controller pattern](defs.md#controller-pattern). Thus, due to the speed of Kubernetes itself, it is destined for non-realtime parts of telco systems.

Here are some requirements and assumptions for Lupus platform:
- Targeted for telco industry (especially mobile networks).
- Will enable to design and run any closed control loop architecture, especially the ones proposed in [Overview of Prominent Control Loop Architectures](https://www.etsi.org/deliver/etsi_gr/ENI/001_099/017/02.01.01_60/gr_ENI017v020101p.pdf).
- Processes managed by Lupus will be non-realtime.
- Implemented on top of Kubernetes, leveraging its controller pattern.
- Actual processing components of Lupus loop are external to it (e.g. These are some HTTP Servers, especially [Open Policy Agents](https://www.openpolicyagent.org))

## What Lupus can do

Here we will discuss what is in the scope of Lupus responsibilities and what is not. The best idea to explain such topic is to provide some simple example. The example will be simplified as much as it can be and is not applicable in real world.

Let's say you are monitoring a temperature in 4 rooms named "A", "B", "C", "D". Each room has a knob on its radiator that can be controlled remotely. You have a database with records of desired temperature for each room. Building administrator from time to time changes this values based on the current season and weather. In your headquarters you have a HTTP server that exposes an endpoint for which you can present room name and it current temperature and endpoint will answer with instruction how to move the radiator knob to get closer to the desired temperature. As for now, each room has a microcontroller connected to the thermometer and radiator.

Your mission is to keep the desired temperature in each room.

![](../_img/39.png)

As for now, microcontroller in each room periodically quries HTTP server in HQ. In this simple case this is OK, but [reconcillation](defs.md#reconcilliation) logic is performed by controller, locally. It can cause several problems:
- What if logic changes? You will have to reinstall program in each controller
- What if logic is more complex? Either it is hard to code it in controller's native language or controller has no sufficient computing power to perform it

Imagine company where such "temperature management" is only one aspect for management. One among of thousands more. It would be nice to have a single "place" of automated management in company.

**It would be better that microcontroller's program serves only as interface between its devices and external world, and that the loop workflow is expressed somewhere else.** 

**Here is where the Lupus comes into play!**

First let's deprive microcontrollers of the loop workflow. Reprogram them so they only can be queried for current temperature and instructed how to move radiator knob. 

Next, let's say that you have some datacenter building with Kubernetes cluster. Here we can install Lupus components that will perform loop workflow. 

New architecture would look like this:

![](../_img/40.png)

The [Ingress-Agent](defs.md#ingress-agent) gathers temperature information from all rooms, then sends it to Lupus. Lupus performs loop logic (query HQ server in this case) for each room and sends results to [Egress-Agent](defs.md#egress-agent). Egress Agent the sends appropriate set of actions to each microcontroller.

It is worth to note that loop workflow is designed in Kubernetes YAML definition files and Kubernetes serves as loop runtime. 

Such approach gives several advantages:
- Single place (Kubernetes cluster) for all management loops in company
- If the loop workdflow will change, the change has to be done only here
- The loop workflow is coded in [LupN](defs.md#lupn) a notation dedicated to loops

### What Lupus cannot do?

As you can note in the example above we still need the HQ server that performs the reconcilliation logic. It is not the Lupus mission to perform actual calculations. Lupus controls only the loop workflow. Here the example is very easy, but imagine the case, where loop has to query multiple servers and make decision based on the responses along the way. This is what Lupus takes care of.

> Of course, sometimes it is not it is overkill to deploy dedicated server with some little reconcilliation logic on it. In such case Lupus offers [go code snippets](defs.md#user-functions) possible to be inserted in loop workflow.

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

Ingress and Egress Agents are input and output points to the Lupus system. They are external to Kubernetes cluster and have to be developed by [designer](defs.md#designer). The only requirement for them is to be compatible with [Lupin and Lupout interfaces specification](lupin-lupout-if.md). 

Our loop will have only one [Lupus Element](defs.md#lupus-element), thus it will be both the [Ingress](defs.md#ingress-element) and [Egress](defs.md#egress-element) Element.

[Ingress Agent](defs.md#ingress-agent) on one side has to communicate with [managed-system](defs.md#managed-system), and on the other it has to implement [Lupin interface](lupin-lupout-if.md).

![](../_img/43.png)

On the left side it will act as a MQTT Broker and will gather temperature measurement from microcontrollers. On the right side (the one implemeting Lupin interface) it will change the Status of [Ingress Element](defs.md#ingress-element) CR* (i na a periodical manner) with such json input:

```json
{
    "A": {
        "room": "A"
        "temp": 12
    },
    "B": {
        "room": "B"
        "temp": 10
    },
    "C": {
        "room": "C"
        "temp": 8
    },
    "D": {
        "room": "D"
        "temp": 14
    },
}
```

> *Python offers great library for interworking with kube-api-server

### 2. Design the Lupus Elements

Ok, we've received [managed-system-state](defs.md#managed-system-state), now we have to come up with Loop workflow. What has to be done in each loop iteration?

Looking back at how it was done previously, we need to query HTTP server in HQ for each room.

How to express such logic in Lupus?

We need to prepare YAML config file for Kubernetes Object [`masters.lupus.gawor.io`](defs.md#masterslupusgaworio). The spec of this object expresses the loop workflow.

```yaml
apiVersion: lupus.gawor.io/v1
kind: Master
metadata:
  labels:
    app.kubernetes.io/name: lupus
    app.kubernetes.io/managed-by: kustomize
  name: Temp
spec:
  name: Temp
  elements:
    - name: "Main"
      descr: "Queries HQ Server for each room"
      actions:
        - name: "room A"
          type: send
          inputKey: "A"
          outputKey: "A"
          destination:
            type: http
            http: 
              path: "http://hq.server.corpnet/administration/rooms/temp"
              method: GET
          next: "room B"
        - name: "room B"
          type: send
          inputKey: "B"
          outputKey: "B"
          destination:
            type: http
            http: 
              path: "http://hq.server.corpnet/administration/rooms/temp"
              method: GET
          next: "room C"
        - name: "room C"
          type: send
          inputKey: "C"
          outputKey: "C"
          destination:
            type: http
            http: 
              path: "http://hq.server.corpnet/administration/rooms/temp"
              method: GET
          next: "room D"
        - name: "room D"
          type: send
          inputKey: "D"
          outputKey: "D"
          destination:
            type: http
            http: 
              path: "http://hq.server.corpnet/administration/rooms/temp"
              method: GET
          next: final
      next: 
        - type: destination
          destination:
            type: http
            http:
              path: "http://hq.server.corpnet/lupus/temp/egress-agent"  
              method: POST
            keys: ["*"]
```

The general structure of YAML config file for Kubernetes resources is assumed to be known by reader.

One [Master](defs.md#master) API Object corresponds to one management loop. Let's focus on its `spec`.

`spec` of Master is expressed in [LupN](defs.md#lupn). LupN notation will be interpreted by Master [operator](defs.md#operator). Loop workflow is based on [Data-Concept](defs.md#data). This is the information carrier for single loop-iteration. Iteration at any time can retrieve or save operational information in [Data](defs.md#data). Data is simply an information coded in Json format. In each iteration its initial form is exact the same json that [Ingress Element](defs.md#ingress-element) has got from [Ingress-Agent](defs.md#ingress-agent).

Notation above compiles to the following loop:

1. Retrieve "A" field from Data.
2. Send its content to HQ Server.
3. Save response in "A" field in Data.
4. Retrieve "B" field from Data.
5. nd its content to HQ Server.
6. Save response in "B" field in Data.
7. Retrieve "C" field from Data.
8. Send its content to HQ Server.
9. Save response in "C" field in Data.
10. Retrieve "D" field from Data.
11. Send its content to HQ Server.
12. Save response in "D" field in Data.
13. Send resulted Data to the Egress-Agent.

Which converts to such Diagram:
![](../_img/44.png)

> Solid line means that one entity sends something to another, while stripped line shows only the sequence of execution.

Let's examine the [Master's](defs.md#master) spec. Full documentation of LupN is [here](lupn.md). Master has 2 fields in its spec:
- `name` - used to uniquely identify the loop and its elements (elements are prefixed with loop name). In our example the name is "Temp", which will tell us that from among of plethora of Closed Control Loops, this one is responsible for automated management of Temperature.
- `elements` - list of loops elements. Element is the first level of logical division of loop. In case of [OODA loop](https://www.etsi.org/deliver/etsi_gr/ENI/001_099/017/02.01.01_60/gr_ENI017v020101p.pdf) elements could be: {"Observe", "Orient", "Decide", "Act"}. In our simple example we have only one element - "Main".

Element spec contains:
- `name` that uniquely identifies the element in scope of its loop. In our example we have only one element hence the name "Main".
- `descr` - description of element e.g. its goal or summary of its logic. It is for faciliation of [desginer's](defs.md#designer) work. In our example it simply describes its operation as this may be what desinger assumed as convention.
- `actions` - [Action](defs.md#action) is the second level of logical division of loop. Actually, actions are performed on the [Data](defs.md#data) object. Action can be of several types. Each actions takes one field from Data (identified by `inputKey`), performs action dependable on its type, and, if applicable, saves the result in field indicated by `outputKey`. In our example we can see only actions of type `send`. This action takes some Data field, sends it as input to indicated [destination](defs.md#destination) (e.g. HTTP Endpoint) and saves received json response in Data field indicated by `outputKey`. Full list of available actions types and their documentation is [here](actions.md). <br>Each action defines the next action in the workflow* as the `next` field. Special keyword `final` exists to terminate the actions chain.
- `next` - next specifies the next [element](defs.md#loop-element) in loop workflow and resulting Data fields that has to be sent there. It can be any [destination](defs.md#destination) or a [lupus-element](defs.md#lupus-element). 

> *Control Flow with conditions (if expression) or sudden exits are possible.

At the end of the loop iteration we have such [Resulting data](defs.md#resulting-data):
```json
{
    "A": {
        "move_gauge": 1
    },
    "B": {
        "move_gauge": 0
    },
    "C": {
        "move_gauge": 2
    },
    "D": {
        "move_gauge": -2
    }
}
```
Note that these are the response from HQ Server for each room. Instead of Data representing the [managed-system-state](defs.md#managed-system-state) we have now Data that represents the [management-action](defs.md#management-action).

We send such Data to Egress-Agent. It will be his mission to translate it to [Management-Action] and execute such action.

### 3. Design the Egress-Agent

[Egress-Agent](defs.md#egress-agent) has to implement [Lupout interface](defs.md#lupout-interface) on its left side. It can be HTTP server (as it is in our case). Then, on the right side it has to perform [Management-Action](defs.md#management-action), which is out of scope of Lupus specification and specific to [managed-system](defs.md#managed-system). 

![](../_img/45.png)

In our case Egress-Agent exposes endpoint "http://hq.server.corpnet/lupus/temp/egress-agent, which accepts json in the same format as resulting data above. Then it translates this json to 4 separate management-actions. Each involves sending "move gauge" instruction to respective microcontroller. Exact way of communication here is not relevant and out of scope of this tutorial.

### 4. Closing word

Now the loop is ready. Its workflow runs in Kubernetes as Master CR and Element CR, and its logic runs in HQ Server.

## Summary

This was just quick walktrhrough. To go further explore:
- detailed documentation of Lupus -> [detailed-doc.md](detailed-doc.md), or
- examples directory -> [examples](/examples/)