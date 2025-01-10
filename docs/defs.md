# managed-system

# current-state

# desired-state

# closed-control-loop

# external-element

# controller-pattern

# reconciliation
An act of bring the [current-state](#current-state) of a [managed-system](#managed-system) closer to the [desired-state](#desired-state). 

# ingress-agent

# egress-agent

# reconcilliation-logic
It comprises of operations that are mandatory to [reconciliation](#reconciliation) goal. 

# loop logic 
Same as [reconciliation-logic](#reconcilliation-logic), while reconciliation logic describes more what has to be done in single iteration and loop-logic is more of a general loop attribute.

# loop-workflow
It is a [workflow](#workflow) present in one iteration of loop. Loop workflow should deliver the proper [loop-lofic](#loop-logic).

# loop-element
Building entity of Loop Workflow. Can be Lupus Element or External Block.

# user-functions

# lupin-interface

# lupout-interface

# management-problem

# lupus-element

# designer

# data

# ingress-element

# egress-element

# managed-system-state

# elements.lupus.gawor.io

# masters.lupus.gawor.io

# master

# lupn

# operator

# loop-iteration

# destination

# action

# resulting-data
In each iteration when the Data will be modified with all the elements's actions we name such Data a "resulting". It is simply the Data at the output of element.

# management-action

# workflow
Workflow is a chain of actions connected in a single direction. Sometimes the next action of chain results from condition expression.


# data-driven
Lupus primary goal was to be data-driven, do not comprise any reconciliation logic on its own. Lupus internal elements had only to express the loop workflow, do not perform any computation needed in [reconciliation-logic](#reconcilliation-logic) (as calculating the diffrences between [current](#current-state) and [desired](#desired-state)) or reach out for some information needed by [reconciliation-logic](#reconcilliation-logic). Only in this way it is possible to propose a framework that can run **ANY** loop. The [controller](#controller) of loop element can't have any of the [loop-logic](#loop-logic), the logic has to be delegated somehwere else, but still belong to the loop. Thus we divide loop-elements in two groups:
- [lupus-element](#lupus-element) - These run in Kubernetes cluster, they serve to express loop workflow and delegate actual computing to [external-elements]
- [external-element](#external-element) - These run outside of Kubernetes cluster, typically as HTTP servers (especially [Open Policy Agent](#open-policy-agent))


# controller

