<h1 align="center">Lupus</h1>

<p align="center">
  <img src="_img/logo.png" alt="Lupus logo"/>
  <br>
  <i>Lupus - Design and Run Closed Control Loops in Kubernetes</i>
  <br>
</p>

---

Lupus is an open-source platform that enables you to design and run closed control loops within a Kubernetes cluster. The project focuses on managing telco and mobile network systems.

In robotics and automation, a control loop is a non-terminating process that regulates the state of a system. Kubernetes inherently implements the [controller pattern](https://kubernetes.io/docs/concepts/architecture/controller/). We leverage this by combining Kubernetes extensions such as [Custom Resource Definitions (CRDs)](https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/) and the [Operator Pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/) to create a reusable architecture that allows users to design and run any kind of closed control loop for system management.

## ETSI Context

The direct inspiration for this project comes from the work of the [ETSI](https://www.etsi.org) committee, specifically ["ENI - Experiential Networked Intelligence"](https://www.etsi.org/technologies/experiential-networked-intelligence). The document ["Overview of Prominent Control Loop Architectures"](https://www.etsi.org/deliver/etsi_gr/ENI/001_099/017/02.01.01_60/gr_ENI017v020101p.pdf) discusses various control loop architectures. The natural next step is to develop a way to implement and run such loops. Kubernetes was chosen as the runtime due to its widespread adoption within the telco community.

## How to Use It

If you have a system that produces observable data and supports management actions based on that data, Lupus is the perfect platform to begin your automation journey.

<p align="center">
  <img src="_img/readme/1.png" alt="Lupus control loop overview"/>
</p>

To get started:
- Develop Ingress and Egress Agents according to the **Lupin** and **Lupout** interface specifications (detailed documentation coming soon).
- Design your loop workflow using an external tool of your choice (e.g., flowcharts, BPMN diagrams).
- Express the designed logic in Master CR YAML manifest.

Read our [Getting started](getting-started.md) guide for more details.

## Project Status

This project is part of my Master's thesis, supervised by [Dariusz Bursztynowski, Ph.D., Eng.](https://repo.pw.edu.pl/info/author/WEITI-99bdf4cf-dec0-4770-baf2-80874a4d91a0/Profil+osoby+%E2%80%93+Dariusz+Bursztynowski+%E2%80%93+Politechnika+Warszawska), during my second-degree ICT and Cybersecurity course at the [Warsaw University of Technology](https://eng.pw.edu.pl).

The project is under active development, and contributions or feedback are welcome!
