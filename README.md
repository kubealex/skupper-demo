# Skupper demo 101

This repo contains the basic setup to run a quick demo of how Red Hat Service Interconnect (based on project [skupper.io](skupper.io)) works.

The repo has a dedicated readme for each use case, instructions expects you cloned this repo and you are working in the root folder of the project.

- [Podman to k8s/OCP Cluster communication](./docs/uc1_podman_k8s.md)
- [VM service to k8s/OCP Cluster communication](./docs/uc2_vm_k8s.md)
- [OCP Cluster to OCP Cluster communication](./docs/uc3_k8s_k8s.md)

#### What you need to run the demo

- An OCP Cluster (two, if you want to test cluster to cluster communication), we will call them **cluster-A** and **cluster-B**
- A VM running podman
- 10 minutes of your time
- Skupper CLI, of course.

You can follow the [skupper.io](https://skupper.io) instructions to install the CLI.

TL:DR

```shell
curl https://skupper.io/install.sh | sh
```

## Preliminary steps

In our OCP Clusters, we will use the Red Hat Service Interconnect operator to handle our resources and the console.

For k8s clusters you can use the skupper CLI to create a site using the **skupper init** directly.

Let's set up contexts to avoid confusion and allow using a single terminal session to work with both:

_*on cluster-A*_:

```shell
    oc login -u OCP_USERNAME -p OCP_PASSWORD https://OCP_API_URL_CLUSTER_A:6443
    oc config set-cluster cluster-A --server=https://OCP_API_URL_CLUSTER_A:6443 --insecure-skip-tls-verify=true
    oc config set-credentials admin-cluster-A --token=$(oc whoami -t)
    oc config set-context cluster-A --cluster=cluster-A --namespace=skupper-demo --user=admin-cluster-A
```

_*on cluster-B*_:

```shell
    oc login -u OCP_USERNAME -p OCP_PASSWORD https://OCP_API_URL_CLUSTER_B:6443
    oc config set-cluster cluster-B --server=https://OCP_API_URL_CLUSTER_B:6443 --insecure-skip-tls-verify=true
    oc config set-credentials admin-cluster-B --token=$(oc whoami -t)
    oc config set-context cluster-B --cluster=cluster-B --namespace=skupper-demo --user=admin-cluster-B
```

Install the Skupper operator on (each) OCP Cluster by running:

```shell
    oc apply -f resources/cluster-A/skupper-project.yml -f resources/cluster-A/skupper-operator-subscription.yml --context cluster-A
    oc apply -f resources/cluster-B/skupper-project.yml -f resources/cluster-B/skupper-operator-subscription.yml --context cluster-B
```

---
