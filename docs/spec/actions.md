# Actions
Detailed specification of actions (especially their notation in LupN) is present in [lupn specifcation](lupn.md). This document presents some examples of how actions work.

Actions were developed as the most atomic operations that when combined in the appropriate manner will deliver a tool to fully operate on [data](../defs.md#data).

Sometimes some opeartion that seems to be atomic at the first glance will require usage of two actions combined. On the other hand, sometimes an operation that seems to be atomic at the first glance, appears to be a specific case of more general operation. A good example here is an extinct action of concat. This action was designed to merge two fields into one, but it turned out that it is a specific case of nest, where the `InputKey` list has only 2 elements. 

## Examples

### Send


