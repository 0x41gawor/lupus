# Exemplary use of Lupus for `upf-net`

Let's bring back the general architecture for Lupus application:

![](../../../_img/readme/1.png)

First you need to develop Ingress and Egress Agents:

- [Ingress Agent](ingress-agent.py)
- [Egress Agent](egress-agent.py)

Now design the architecture of Loop:

![](../_img/sample-loop.svg)

Prepare the external elements of Loop:
- [http-bouncer](http-bouncer.go)
- [open-policy-agent](opa.md)

Express the Loop in yaml manifest file of master element.

- [master.yaml](master.yaml)