# Lupus - detailed docs (The Lupus book)

This document describes Lupus in a story-like manner. It explains Lupus concepts (sometimes containing links to a discussion about taken decisions or change history), contains links to specifications but at the same time follows an example, so the explained&specified concepts can be related to an example.


## 1. Management problem

The first thing that have to be described is [management-problem](defs.md#management-problem) and [managed-system](defs.md#managed-system). 

In a real word, we continuosly encounter a situations when it would be nice if the work of some system could be contantly regulated. For example:
- we would like a feature in cars that would regulate the work of engine in order to keep stedy velocity
- it would be nice if the refiregerator could keep a cool but not below-zero temperature regardless of how often the door is opened or what the outside temperature is on a given day
- it would be nice if cloud server could quarantee that an application with sufficient resources to cover users needs will be up and running

The problems listed above can be regarded as [management-problem](defs.md#management-problem). There are no technical limitations to achieve these goals, all of the systems above have the proper equipment e.g. we can add more gas to the engine or deliver more power to the compressor in refrigerator. The issue lies in actual doing it in a appropriate moments. E.g. adding more gas, when car slows down or delivering more power when the temperature in refrigerator goes up. Thus, this is the problem of pure management. 

And the system that we aim to manage is [managed-system](defs.md#managed-system). 

## 2. Control System

A [control-system](defs.md#control-system) is a system that regulates the work of [managed-system](defs.md#managed-system).
For example:
- a Cruising Control System regulates the work of engine in order to keep steady pace,
- refiregerator who controls the work of compressor to maintain a cool temperature,
- Kubernetes that keeps an eye if the desired number of Pods is running

In other words, it solves the [management-problem](defs.md#management-problem). 

But how? What is the general architecture to approach each such problem?

A reply to this question first arosed in Robotics and Automation for [Industrial Process Control](https://en.wikipedia.org/wiki/Industrial_process_control). The answer is to use [control loops](defs.md#control-loop).

## 3. Control loop

It consists of the:
- process sensor, 
- the controller function, 
- and final control element (FCE) which controls the process necessary to automaically adjust the value of a measured process variable (PV) to equal the value of a deisred set-point. 

Control loops are categorized based on whether they incorporate feedback mechanisms:
- **Open Control Loops**: The control action (input to the managed system) is independent of the managed system's output.
- **Closed Control Loops**: The output of the managed system is "fed back" to the [control-system](defs.md#control-system) and influences the control action.

In Lupus we focus only on the [closed-control-loops](defs.md#closed-control-loop). 

## 4. Closed control loop

This is the starting point for our reference architecture:

![](../_img/46.png)

and definitions as such:
- [control-system](defs.md#control-system) - A system that solves the [management-problem](defs.md#management-problem) of a [managed-system](defs.md#managed-system) by the means of [closed-control-loop](defs.md#closed-control-loop).
- [closed-control-loop](defs.md#closed-control-loop) - A non-terminating loop that regulates the state of a [managed-system](defs.md#managed-system) by iteratively bringing the [current-state](defs.md#current-state) closer to the [desired-state](defs.md#desired-state).
- [control-action](defs.md#control-action) - An action sent to or performed on a [managed-system](defs.md#managed-system) that brings it closer to the [desired-state](defs.md#desired-state).
- [control-feedback](defs.md#control-feedback) - A representation of the [current-state](defs.md#current-state) sent from [managed-system](defs.md#managed-system). 

//TODO a good example will explain it really good

