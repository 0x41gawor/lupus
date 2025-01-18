# Getting started
## Target audience
Lupus is a platform targeted for automated management of [systems](defs.md#managed-system) in mobile networks industry, but it can be used with any system that exposes API for observation of its current state and for management or control actions.

## Initial development requirements and assumptions
- Targeted for telco industry, especially Mobile Networks
- Enables to design and run any closed control loop architecture (workflow), especially the ones proposed in [Overview of Prominent Control Loop Architectures](https://www.etsi.org/deliver/etsi_gr/ENI/001_099/017/02.01.01_60/gr_ENI017v020101p.pdf).
- Processes managed by Lupus will be non-realtime.
- Implemented on top of Kubernetes, leveraging its controller pattern.
- Lupus is data-driven, which means that it does not impose and does not have any loop logic sewn in by default
- Actual processing components of Lupus loop are external to it (e.g. These are some HTTP Servers, especially [Open Policy Agents](https://www.openpolicyagent.org))

## What problems Lupus can solve

They example provided below will not apply to telco/mobile systems. Instead it will be a simple management-problem that every engineer should understand. For more advanced and close to the real world examples explore the [examples](../examples/) directory.

Imagine you are monitoring the temperature in four office rooms named "A," "B," "C," and "D." Each room has a radiator with a knob that can be controlled remotely. 
Additionaly, each room is equipped with a microcontroller connected to a thermometer and the radiator.

![](../_img/47.png)

A database stores the desired temperature for each room, and the building administrator periodically updates these values based on the current season, weather conditions and time of day.
In your headquarters, there is an HTTP server with an endpoint. You can send the room name and its current temperature to this endpoint, and it will respond with instructions on how to adjust the radiator knob to bring the temperature closer to the desired level.

<img src="../_img/48.png" style="zoom:50%">


Your goal is to maintain the desired temperature in each room.


Currently, the microcontroller in each room periodically queries the HTTP server in the headquarters and adjusts temperature according to received instruction. While this simple setup works for now, the reconciliation logic is handled locally by the microcontroller. This approach can lead to several issues:
- Updating reconciliation logic: If the logic changes, you'll need to reinstall the program on every microcontroller.
- Complex logic: For more complex logic, it may be difficult to implement in the microcontroller's native language, or the microcontroller might lack the necessary computing power to handle it.
- Scalability: In a company where "temperature management" is just one of thousands of management tasks, having a single centralized system for automated management would be far more efficient.

It would be better that microcontroller's program serves only as interface between its devices and external world, and that the loop workflow is expressed somewhere else.

**Here is where the Lupus comes into play!**

First, let's deprive microcontrollers of the reconciliaton logic. We will reprogram them so they only:
- can be queried for the current temperature
- can be instructed how to move radiator's knob

Next, let's assume that you have some server with available computing power and a running Kubernetes cluster. We can use this to deploy a Closed Control Loop for rooms temperature management.

New architecture would look like this:

![](../_img/49.png)


In this setup, [Ingress-Agent](defs.md#ingress-agent) periodically gathers temperature from all romms, then sends it to Lupus. Lupus performs the [reconciliation-logic](defs.md#reconciliation-logic) for each  (which in this case is to query HQ Server) and sends results to [Egress-Agent](defs.md#egress-agent). Eggress Agent then translates it into appropriate set of actions for each microcontroller.

Such approach gives several advantages:
- Kubernetes cluster acts as a single point of management for all closed control loops in company
- The [loop-workflow](defs.md#loop-workflow) is coded in [LupN](defs.md#lupn) a notation dedicated for it
- If the [loop-workflow](defs.md#loop-workflow) will change, changes has to be made only in Lupus

## Lupus constraints
As you can note in the example above, we still need the HQ Server to be running. It is not the Lupus mission to perform [computing part of loop-logic](defs.md#computing-part). Lupus controls only the [loop-workflow](defs.md#loop-workflow) and delegates such workload to [external-elements](defs.md#external-element). Here, the example is very easy. But imagine the case, where loop has to query multiple servers and control its flow based on the responses along the way. This is what Lupus can take care of.

## Lupus mission
Lupus mission is to express and run a [loop-workflow](defs.md#loop-workflow) and introduce a single point of management of such loops in an organisation.

## Closing words

This was just a quick overview. To go further explore:
- detailed documentation of Lupus -> [deatailed-docs](detailed-docs.md)
- examples directory -> [examples](../examples/) to see Lupus application in various examples