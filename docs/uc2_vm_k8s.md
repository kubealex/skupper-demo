## Use case 2: VM service to k8s/OCP Cluster communication

In this use case, we are going to expose a simple golang application that is showing the hostname of the machine where it is running.
At the end of the demo, we will see that the application is accessed using an OCP Route, but it is **actually served by the VM**.

For this section, we are going to use **cluster-A** and we will use the _skupper gateway_ configuration.

### Configure Skupper on the VM

```shell
skupper init --ingress none --platform podman
skupper gateway init --type podman --platform kubernetes --context cluster-A
```

This will configure our VM as a gateway peer, where traffic will be routed.

---

### Configure Skupper on OCP cluster

We will then create our OCP service:

```shell
skupper service create hello-skupper 8080 --platform kubernetes --context cluster-A
```

And configure the gateway so that our port 8080 on localhost will be an endpoint for the created service and expose the route to reach it using our cluster:

```shell
skupper gateway bind hello-skupper localhost 8080 --platform kubernetes --context cluster-A
oc expose svc hello-skupper --context cluster-A
```

---

### Workload creation

The last step is to initialize the mini-application, a simple REST endpoint that listens on port 8080 of our host and just prints a welcome message from the host's FQDN.

The binary is in the _hello-skupper_ folder and you can launch it like this:

```shell
./hello-skupper/hello-skupper
```

You can then browse or cURL to _[http://hello-skupper-skupper-demo.apps.OCP_DOMAIN/](http://hello-skupper-skupper-demo.apps.OCP_DOMAIN/)_ and if everything is correctly configured, you'll see the following, confirming we are accessing our local application from OCP:

![](../_assets/hello-skupper.png)

### Cleanup

To cleanup, run the following commands:

```shell
oc delete route hello-skupper --context cluster-A
skupper service delete hello-skupper --platform kubernetes --context cluster-A
skupper gateway delete --platform kubernetes --context cluster-A
skupper delete --platform podman
```

---
