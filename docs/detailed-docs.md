# Lupus - detailed docs

This document contains references to detailed documentations/specifications of various Lupus aspects. The references are arranged in kind of coherent story.

Lupus is a platform that let's its users to **design** and **run** any loop workflow. It was designed for Closed Control Loops in telco environment.

Lupus leverages Kubernetes CRD and Operator Pattern. 

You **design** the loop by expressing its workflow in CRD yaml file following the [LupN notation](defs.md#lupn). Then you **run** it, by applying such YAML manifest file and creating [master](defs.md#master) (one loop has one master CR). Master [operator](defs.md#operator) takes care of spawning [lupus-elements](defs.md#lupus-element) (one loop has at least one element CR). The operators developed by Lupus team for these two custom resources take care of running the loop logic.

Loop can contain 2 types of elements: [lupus-elements](defs.md#lupus-element) and [external-elements](defs.md#external-element). [lupus-elements](defs.md#lupus-element) are created via Kubernetes and are responsible for [loop-workflow](defs.md#loop-workflow). 

YAML Manifest file of [lupus-master](defs.md#lupus-master) specifies the [loop-workflow](defs.md#loop-workflow). Detailed documentation of LupN [can be found here](lupn.md).




