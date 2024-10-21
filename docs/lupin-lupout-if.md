# Spec of Lupin and Lupout interfaces

## Intro

![](../_img/readme/1.png)

It is impossible to develop Lupus in way that it can be plugged into any managed system as it is (en soi) and ready to go without any integration effort.

To adjust any of plethora of managed-systems the concept of Ingress and Egress Agent was born. Both of them work as translation agents. They translate communication from one system (Lupus) to another (managed-system) and vice versa. 


![](../managed-systems/_img/1.svg)

According to the image above:
- Lupin interface is the interface of Feedback, any signal at this interface will trigger an iteration of closed control loop
- Lupout interafce is the interface of Action, Lupus transmits signals at this interface to trigger changes in managed-system

Any piece of software that implements Lupin interface can be named Ingress Agent.

Any piecie of software that implements Lupout interface can be named Egress Agent.

## Lupin interface

The entry point of Lupus is Observe CR. 

If Ingress Agent wants to indicate that a new state of managed-system can be observed, it has to modify the Observe CR Status `input` field. The value placed in this field will represent the new observed state.

The `Status.Input` field of Observe CR is of type [RawExtension](https://github.com/kubernetes/apimachinery/blob/829ed199f4e0454344a5bc5ef7859a01ef9b8e22/pkg/runtime/types.go#L94) which can receive any json object.

The root fields of json sent here, will be the subject of tagging system in Lupus.

Ingress Agent implements the Lupin interface if at some point of its code it sends a HTTP request to the kube-api-server that updates the status of Observe CR. More precisely the `input` field. The value has to be json object that represent current observed state of a managed-system.

## Lupout interface

The exit point from Lupus (where Lupus gives out the control (pl. "oddaje sterowanie")) is Execute CR controller.

At this point, Lupus, based on its input (feedback) on Lupin interface, has decided what needs to be done. 

Lupus has prepared some set of commands that have to be executed on managed-system.

Egress Agent implements the Lupout interface if it exposes a HTTP Server listening for these commands. Commands will arrive in body in json format with the root object named `commands`. The mission of Egress Agent is to translate this json into set of actions that can be performed with usage of managed-system exposed API.