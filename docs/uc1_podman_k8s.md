## Use case 1: Podman to k8s/OCP Cluster communication

In this use case, we are going to expose an nginx server running in a container on the VM in OCP Cluster.
At the end of the demo, we will see that the webserver is accessed using an OCP Route, but it is actually served by the **container running in our VM**.

For this section, we are going to use **cluster-A**.

### Skupper initialization on OCP cluster-A

The needed steps to start operating Skupper on the cluster are:

- Initializing the site
- Generate the secret with link configuration
- Retrieve the secret's content to proceed on Podman site

#### Site initialization

We start bycreating the site for cluster-A after jumping in the _skupper-demo_ project:

```shell
    oc apply -f resources/cluster-A/skupper-site-config.yml --context cluster-A
```

The site is configured to also expose the Red Hat Service Interconnect console, with credentials **admin/redhat** that you can use to visualize the connections.

#### Generate the secret with link configuration

We then generate the required secret for site-to-site link configuration:

```shell
    oc apply -f resources/cluster-A/skupper-site-link-config.yml --context cluster-A
```

#### Retrieve the secret's content to proceed on Podman site

Proceed retrievieng the secret and we will use it to create the link with the Podman site we created before:

```shell
    skupper token create token-cluster-A.yml --context cluster-A
```

---

### Skupper initialization on the VM

On the VM, we will proceed by:

- Initializing the Podman site
- Establish the link with the OCP site

#### Initializing the Podman site

Let's start initializing Skupper on the VM, but first we need to enable Podman socket for the user to allow API communication:

```shell
    systemctl --user start podman.socket
    skupper init --platform podman --site-name podman-site --ingress none
```

#### Establish the link with the OCP site

And finally we establish the link between the two sites:

```shell
    skupper link create --platform podman --name ocp-link token-cluster-A.yml
```

We will get a confirmation that the link is up and running:

```shell
    [sysadmin@rhel9-vm ~]$ skupper link status --platform podman
    Links created from this site:
            Link ocp-link is connected
```

On OCP, if you go in the Skupper console (https://skupper-skupper-demo.apps.CLUSTER*DOMAIN/) and log in with \_admin/redhat\* credentials, you can confirm the link is up:

![](../_assets/podman-link.png)

---

### Workload creation and configuration

Now we can start our container. We will use a standard _nginx_ that creates a running instance on port 80 of the container.

```shell
podman run --name nginx -d --network skupper docker.io/nginx
```

Let's expose the container _nginx_ it on the skupper network with the address _podman-nginx_ listening on port _8080_ linked to container port _80_:

```shell
skupper expose host nginx --platform podman --address podman-nginx --port 8080 --target-port 80
```

The following command will create an OCP service (note the **--platform kubernetes** to switch from podman context and use the connection to the cluster) and we will expose the service, creating a Route on OCP:

```shell
skupper service create  podman-nginx 8080 --platform kubernetes --context cluster-A
oc expose svc/podman-nginx --context cluster-A
```

If you now navigate to [http://hello-skupper-skupper-demo.apps.OCP_DOMAIN/](http://hello-skupper-skupper-demo.apps.OCP_DOMAIN/) you should be able to access the nginx instance created on the VM with Podman:

![](../_assets/hello-nginx.png)

---

### Cleanup

To cleanup, run the following commands:

```shell
podman rm -f nginx
oc delete route podman-nginx --context cluster-A
skupper service delete podman-nginx --platform kubernetes --context cluster-A
skupper service delete podman-nginx --platform podman
skupper link delete ocp-link --platform podman
skupper delete --platform podman
skupper delete --platform kubernetes --context cluster-A
```
