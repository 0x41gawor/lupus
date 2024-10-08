# MME Network
This document presents the business case for the closed control loop application.
## Real world 

Imagine you are operating a PLMN (Public Land Mobile Network) in Poland with only 4G technology. Your customers are served byMME nodes. We have four **nodes** in four different locations, as shown below:

<img src="img/1.png"/>

Each user generates a **bearer session** (simply called - "session"). We have extremely fast transport network, so from this point of view it is transparent by which node user will be served. We say that at any given moment 'MME servers a number of sessions'. Session generation by users is nearly random, and we cannot control which geographical area will experience the most traffic at any given time. Users are typically served by the nearest MME node. However, some nodes may become significantly more loaded than others. In such cases, given our fast transport network, we can redistribute session loads to other MME nodes.

If session redistribution is possible, how do we determine the split ratio or distribute the traffic? Does each MME hold an equal share of the network? Although we could distribute traffic equally, some nodes may be more efficient than others due to differences in hardware models, available resources, or temporary failures. Thus, the distribution ratio might vary, e.g., 20%-20%-30%-30% instead of 25%-25%-25%-25%

The total number of sessions varies throughout the day, with network usage typically peaking during specific periods. Additionally, as people move around, the number of sessions fluctuates between nodes.

## Mathematically

Bodies:
- **node** - single MME node that can serve a session
- **session** - single unit of demand created by user, need to be served by node

Initially, each node is assigned a certain number of sessions. For example:

```
"Gdansk": 10, "Poznan": 12, "Warsaw": 18,"Krakow": 6
```
> (Note: The specific numbers here are illustrative and could represent counts, megabits, etc.)

During each **round** or **interval**, a random small number of sessions is added or removed from each node. The current session count at a node may be higher, lower, or the same as in the previous round.

Here is an example showing how session counts for each node might change over time:

```sh
# roundNumber. Gdansk-Poznan-Warsaw-Krakow
1. 10-12-18-6
2. 11-12-16-8
3. 10-10-17-9
4. 8-8-20-6
5. 6-7-18-8
```

### The need for our closed control loop

As demonstrated, the overall number of sessions changes over time, and without a control loop, a rational distribution between locations is not maintained. Traffic between nodes and their capacities can vary randomly. What if, at some point, one of the MMEs becomes overloaded, or if one is consistently underutilized over the long term?

For such cause, our closed control loop finds its application. We need to remember that we have the ability to move sessions over nodes to preserve the requested distribution.

This is the use case for reactive closed control loop, it can guarantee that the proper distribution of sessions between nodes will be enforced.

### The need for deliberative closed control loop

Ok, now we understand how our reactive loop will work. It will guarantee that the requested traffic distribution is applied. But what if operator changes its mind and would like to have some other distribution? Or maybe operator will implement some AI algorithms that will analyze and adjust distribution?

This is where the deliberative loop comes into play. It monitors the reactive loop and ensures it adheres to the evolving strategy.



For instance, if a node loses half its resources, we may need to adjust the distribution from `25-25-25-25` to `10-30-30-30`.

> Here, one can remark that for this, only the manual change in configuration file of reactive loop would sufficient. Yes, but:
> - 1st,  this is merely an example to showcase what our platform can do.
> - 2nd, this is simple scenario of more complex abstraction, you don't have to always utilize 100% features of your tool

Or the change in distribution can be result of AI long term analysis.

### Summarizing

We have two loops:
- reactive one enforces the proper distribution on the network
- deliberative comes up with the distribution 