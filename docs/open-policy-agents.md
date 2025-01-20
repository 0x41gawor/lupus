# Open Policy Agents

Again we are starting from [data-driven](defs.md#data-driven) design requirement. Since operators could have no computing-part of the loop-logic. We needed some other place to express it.
Anyway it is a good practive so have loop-workflow defined somewhere and actual computing-part (or policies) somewhere else.

Open Policy Agent is a great place to store polcicies. Homepage: https://www.openpolicyagent.org

We recommend to store computing parts (policies) in OPA since it is convenient. 

In Lupus, OPA is supported as external systems one of Destinations. 

Actually at first OPA was the sole external-system under consideration. It came with time that we could broaden the scope with any HTTP server.

