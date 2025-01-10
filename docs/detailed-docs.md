# Lupus - detailed docs

This document contains references to detailed documentations/specifications of various Lupus aspects. The references are arranged in kind of coherent story.

Lupus is a platform that let's its users to **design** and **run** any loop workflow. It was designed for Closed Control Loops in telco environment.

Lupus leverages Kubernetes CRD and Operator Pattern. 

You **design** the loop by expressing its workflow in CRD yaml file following the [LupN notation](defs.md#lupn). Then you **run** it, by applying such YAML file and creating resource instances of type `masters.lupus.gawor.io` (one loop has one master CR) and `elements.lupus.gawor.io` (one loop has at least one element CR). The operators developed by Lupus team for these two custom resources take care of running the loop logic.


The most important attribute to design a loop is its workflow. Loop workflow can contain 2 types of elements:
- [lupus-element](defs.md#lupus-element) - This element states as [building-entity](defs.md#building-entity), it runs in a Kubernetes Cluster as `elements.lupus.gawor.io` CR instance. 
- [external-element](defs.md#external-element) - This kind of element is external to Kubernetes Cluster. It serves as [computing-entity](defs.md#computing-entity).