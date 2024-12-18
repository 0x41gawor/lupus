# Exemplary use of Lupus for Open5GS

## Live demo
### Prerequisties
- [open5gs-k8s](https://github.com/niloysh/open5gs-k8s) running
- CRD installed
- move to the root dir of this repo

### Steps
#### 1. Main part
Run 4 terminals on MobaXterm and enable split mode:
![](../../upf-net/_img/5.png)

First, at 4 run egress-agent.
```sh
#TODO
```

Then, in 2 run the controller:
```sh
cd lupus
make run
```

In 1, create the Lupus elements
```sh
k apply -f managed-systems/open5gs/sample-loop/master.yaml
```

And finally in 2 run ingress-agent:
```sh
python3 managed-systems/open5gs/sample-loop/ingress-agent.py --interval 30
```