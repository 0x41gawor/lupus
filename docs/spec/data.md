# Data

One of the requirements for Lupus was for it to be [data-driven](../defs.md#data-driven). [Data](../defs.md#data) is the heart and core of fullfillment of this requirement. The runtime of [lupus-element](../defs.md#lupus-element) controller is driven by [data] contents and [actions] chain specified in [LupN] It does not impose any reconciliation logic.

Data is the way in which user can:
- retrive information from [current-state](defs.md#current-state)
- store auxiliary information (as responses from [external-elemetn](../defs.md#external-element)
- store logging/debuggin information
- save information needed to formulate [control-action](defs.md#control-action)
during a single [loop-iteration](../defs.md#loop-iteration).

In each iteration [data](../defs.md#data) resets.

Data is an information carrier. Let's discuss how it stores thie information.



