# managed-system
Any system that can be managed by Lupus in a [Closed Control Loop](#closed-control-loop) manner. Typically we aim into systems from telco and mobile networks industry.

# control-loop
Is a fundametal building block of [control systems (controller)](#controller) (which manages, commands, direct or regulates the bahavior or other devies or systems) in industrial control systems.

We divide control-loops in two groups wheter they incorporate feedback mechanism or not:
- open-control-loops: the control action (input to the managed-system) is indepentent of the managed-system output.
- closed-control-loops: in this the output of managed-system is "fed back" to the [controller](#controller) and influences the control-action

# current-state
The observed, present at the moment and real state of a [managed-system](#managed-system). In contrary to an abstract term of [desired-state](#desired-state).

# desired-state
The state of [managed-system](#managed-system) that we would like to achieve from the management point of view. Typically it is derived from the [management-problem](#management-problem). 

# closed-control-loop
Here, we adopt the defintion of 'control-loop' from Kubernetes. https://kubernetes.io/docs/concepts/architecture/controller/.  Kubernetes concept of sending [current-state](#current-state) as input to the [controller](#current-state) is de facto a feedback mechanism, hence every control-loop in Kubernetes is closed-control-loop.

Defition: non-terminating loop that regulates the state of a system, by bringing the [current-state](#current-state) of [managed-system](#managed-system) closer to the [desired-state](#desired-state) in each iteration.

![](../_img/46.png)

# controller
A brain of closed-control-loop. It reads the [control-feedback](#control-feedback) and comes up with [control-action](#control-action). 

Lupus aspires to be the one. 

The ultimate goal of controller is to solve or keep in harmony the [management-problem](#management-problem).

# control-feedback
A representation of [current-state](#current-state) sent by the [managed-system](#managed-system) to the [controller](#controller).

# control-action
An action send to or performed on [managed-system](#managed-system) by the [controller](#controller) that has to bring it closer to [desired-state](#desired-state).


# reconciliation
An act of bring the [current-state](#current-state) of a [managed-system](#managed-system) closer to the [desired-state](#desired-state). 

# ingress-agent
When [user](#user) wants to integrate his [managed-system](#managed-system) with Lupus - as in any other case of combining two softwares together -  he needs some API Integration. Sometimes [user](#user) cannot modify the [managed-system](#managed-system). For such more general case and more secure option we came up with ingress-agent. Ingress-Agent links Lupus and [managed-system] together. Typical operations for Ingress-Agent is to:
- receive or gather itself iformation from [managed-system](#managed-system) (sometimes from multiple sources)
- watch for changes in [managed-system]
- translate [managed-system](#managed-system) nomenclature and the Lupus one
- etc..

Communication of Ingress-Agent and [managed-system](#managed-system) is case-specific, out of scope of this repo and to be done by the [user](#user). But communication of Ingres-Agent and Lupus is standardized as [lupin-interface](#lupin-interface). 

Ingress-Agent communicates with [managed-system](#managed-system) and provides Lupus with [control-feedback](#control-feedback) via [lupin-interface](#lupin-interface).

# egress-agent
Same as [Ingress-Agent](#ingress-agent) but at the output of Lupus. 

On one side Egress-Agent receives the [control-action](#control-action) from Lupus. This [control-action] is represented as [resulting-data](#resulting-data) json. The mission of Egress-Agent is to translate this into real actions performed on [managed-system] and send them on the other side.

Interface on which Lupus send its [control-action](#control-action) is named [lupout-interface](#lupout-interface). 

# reconcilliation-logic
It comprises of operations that are mandatory for [reconciliation](#reconciliation) goal. 

# loop logic 
Same as [reconciliation-logic](#reconcilliation-logic), but reconciliation logic describes more what has to be done in a single iteration and loop-logic is more of a general loop attribute.

In context of this repo, we can highlight a computing-part of the loop-logic. Computing part is responsible for the actual calculations, while other parts of loop-logic contain sequence of these calculations or data fed from/to them.

# loop-workflow
It is a [workflow](#workflow) present in one iteration of loop. Loop workflow should deliver the proper [loop-logic](#loop-logic). It is build from [loop-elements](#loop-element). In Lupus, loop-workflow is expressed in [LupN](#lupn).

# loop-element
Building entity of Loop Workflow. Can be [lupus-element] or [external-element].

# lupus-element
A [loop-element] that runs in Kubernetes Cluster. Its mission is to express [loop workflow], but delegate loop-logic to external elements. Lupus-element is implemented as [custrom-resource](#custom-resources) and thus its [k8s-controller](#k8s-controller) can't have any loop-logic implemented. Loop-logic has to be derived somewhere else.

In Kubernetes it is a [custom-resource](#custom-resources) named `elements.lupus.gawor.io`. 

# external-element
A [loop-element] that runs outside of Kubernetes Cluster. Its mission is to deliver computation part for [loop-logic](#loop-logic). 

As for now 3 types of external-elements are supported:
- HTTP servers
- Open Policy Agent (subtype of HTTP Server)
- [user-functions](#user-functions)

# user-functions
It was assumed that [k8s-controller](#k8s-controller) of [lupus-element](#lupus-element) should not contain any computation part of [loop-logic](#loop-logic). But sometimes it is non-sense to deploy HTTP Server solely for the reason of it serving as [external-element](#external-element) for Lupus (e.g. in case if [computation-part](#computation-part) of [loop-logic](#loop-logic contains few lines of code)). We allow to deploy such logic in [k8s-controller](#k8s-controller) of [lupus-element](#lupus-element) by introduction of "user-functions". These functions act logically as [external-elemenet](#external-element) with the difference that physically they are deployed in Kubernetes cluster, not outside of it. 

//TODO link for docs

# lupin-interface
Interface implemented by the [Ingress-Agent](#ingress-agent) on which it sends [control-feedback](#control-feedback) to the Lupus.

# lupout-interface
Interface implemented by the [Egress-Agent](#egress-agent) on which it receiver [control-action](#control-action) from the Lupus.

# management-problem
A problem that occurs in given [managed-system](#managed-system) and has to be solved by [controller](#controller). 

> E.g. keeping the constant and desired temperature in a room.

# user

A person or organisation willing to use Lupus project for automation of its [managed-systems](#managed-system). 

Adjusting a [managed-system](#managed-system) with [Ingress](#ingress-agent) or [Egress Agents](#egress-agent) will always require software developers and such people are included in "user" definition, but we distinguish the role of a [designer](#designer). 

# designer

Our assumption was that design of [loop-workflow](#loop-workflow) will require the minimal set of technical/engineering skills. Thus the person who solely operates on [LupN](#lupn) and delegates work related to [Ingress](#ingress-agent) and [Egress Agents] or [external-elements](#external-element) to Software Engineers is denoted by the name "loop-designer". Its mission is to design a loop.

# data

Data is an information carrier implemented in JSON format.

Data is the way in which user can:
- retrive information from [current-state](#current-state)
- store auxiliary information
- save information need to formulate [control-action](#control-action)
during a [loop-iteration](#loop-iteration).

The intial form of [Data](#data) in each [loop-iteration](#loop-iteration) is given by the [current-state](#current-state) input on [lupin-interface](#lupin-interface). Then [lupus-element](#lupus-element) performs modifications on it by [actions](#action). The resulting form of Data is called [final-data](#final-data) and in this form is sent on [lupout-interface](#lupout-interface). 

# ingress-element
A [loop-element](#loop-element) which terminates the [lupin-interface](#lupin-interface). The first element of the [loop-workflow](#loop-workflow).

# egress-element

A [loop-element](#loop-element) that utilizes the [lupout-interface](#lupout-interface). 

# master
Master is a type of [custom-resource] named `masters.lupus.gawor.io`. It mission is to spawn [lupus-elements](#lupus-element). Its YAML file includes the LupN notation.

# lupn
A special notation in YAML, developed in purpose of expressing a [loop-workflow]. A LupN notation is included in [master](#master) YAML specification and later interpretted by its controller to create and run defined loop. 

LupN supports the expression of sequential loop workflow, flow-control or immediate exits. Also, by the mechanism of [actions](#action) it allows to reach out and modify information stored in [data](#data).

LupN introduces few abstract terms/mechanisms as actions, data, destination, next for the cause of workflow expression. (Same as programming languages come with their classes, objects etc..). 

[Full documentation is available here!](lupn.md)

# destination
Destination represents adress infromation of [external-element](#external-element).

//TODO link for its docs in LupN

# action

Do not confuse with [control-action](#control-action).

It is a LupN term.

Actions performs modifications od Data. This is the way in which we can interfact with Data, to:
- retrive information from [current-state](#current-state)
- store auxiliary information
- save information need to formulate [control-action](#control-action)
during a [loop-iteration](#loop-iteration).

// TODO link for docs

# final-data
In each iteration when the Data will be modified with all the elements's actions we name such Data a "resulting". It is simply the Data at the output of element.

# workflow
Workflow is a chain of actions connected in a single direction. Sometimes the next action of chain results from condition expression.

# data-driven
Lupus primary goal was to be data-driven, do not comprise any [reconciliation logic](#reconcilliation-logic) on its own. Lupus internal elements had only to express the loop workflow, do not perform any computation needed in [reconciliation-logic](#reconcilliation-logic) (as calculating the diffrences between [current](#current-state) and [desired](#desired-state)) or reaching out for some information to external databases. Only in this way it is possible to propose a framework that can run **ANY** loop. The [k8s-controller](#k8s-controller) of loop element in Kubernetes can't have any of the [loop-logic](#loop-logic), the logic has to be delegated somehwere else, but still belong to the loop. Thus we divide loop-elements in two groups:
- [lupus-element](#lupus-element) - These run in Kubernetes cluster, they serve to express loop workflow and delegate actual computing to [external-elements]
- [external-element](#external-element) - These run outside of Kubernetes cluster, typically as HTTP servers (especially [Open Policy Agent](#open-policy-agent))

# custom-resources
It is a Kubernetes term: https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/

This is a mechanism to extend Kubernetes. Sometimes built-in (default) resources like Pods, Deployments, Services are not sufficient for our goals. We can define a new resource type and register it in a Cluster. Typically these are created for managemenet of complex application configurations or automated deployment of stateful applications. To register a new resource kind in cluster you need to define it with [Custom Resource Definition file](#crd) and apply it.

# crd
It is a Kubernetes term: https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/

Custom Resource Defintions (CRD) is a file that specifies custom resource. You can apply such file to register new resource type.

# k8s-controller
It is a Kubernetes term: https://kubernetes.io/docs/concepts/architecture/controller/

Every resource has its k8s-controller. K8s-controller runs control loop for them. If anything changes in an API Object like Pod, Deployment, Service etc. Controller gets notified and compares resource [current-state](#current-state) (`status` stored in [etcd](#etcd)) with [desired-state](#desired-state) (`spec` specified e.g. in YAML). Controller for built-in types (kinds) are developed by the Kubernetes team.

# operator-pattern
This is a Kubernetes term: https://kubernetes.io/docs/concepts/extend-kubernetes/operator/

This is a mechanism of extending Kubernetes. When we create [Custom Resources](#custom-resources), we can also implement [controllers](#controller) for them to apply a [control-loop](#closed-control-loop) for them. Such behavior is named "operator-pattern". Sometimes we simply call such controllers "operators". The name "operator" has its genesis in the meaning that such controller typically replaces real human operator of some application (which deployment required custom resource to be defined).

# operator
Colloqually a [k8s-controller](#k8s-controller) for [custom-resources](#custom-resources) is called an "operator".

# etcd
This is a Kubernetes term: https://kubernetes.io/docs/tasks/administer-cluster/configure-upgrade-etcd/, https://etcd.io

This is a database that stores the current state of a cluster.
