# Brigade GitLab Gateway

Send [GitLab events](https://gitlab.com/help/user/project/integrations/webhooks) into a [Brigade](https://github.com/Azure/brigade) pipeline. 

This is a Brigade gateway that listens to GitLab webhooks event stream and triggers events inside of Brigade.

## Prerequisites

1. Have a running [Kubernetes](https://kubernetes.io/docs/setup/) environment
2. Setup [Helm](https://github.com/kubernetes/helm)
3. Setup [Brigade](https://github.com/Azure/brigade) core

## Install

### From File
Clone Brigade GitLab Gateway and change directory
```bash
$ git clone https://github.com/lukepatrick/brigade-gitlab-gateway
$ cd brigade-gitlab-gateway
```
Helm install brigade-gitlab-gateway
> note name and namespace (something important about brigade core)
```bash
$ helm install --name brigade-gl ./charts/brigade-gitlab-gateway
```

### From Repo
Add this project as a helm repo

```bash
$ helm repo add glgw https://lukepatrick.github.io/brigade-gitlab-gateway
$ helm install -n brig-gl glgw/brigade-gitlab-gateway
```

## Building from Source
You must have the Go toolchain, make, and dep installed. For Docker support, you will need to have Docker installed as well. 
See more at [Brigade Developers Guide](https://github.com/Azure/brigade/blob/master/docs/topics/developers.md) 
From there:

```bash
$ make build
```
To build a Docker image
```bash
$ make docker-build
```

## GitLab Integration
The Default URL for the GitLab Gateway is at `:7446/events/gitlab/`. In your GitLab project, go to Settings -> Integrations. Depending on how you set up 
your Kubernetes and the GitLab Gateway will determine your externally accessable host/IP/Port. Out of the box the gateway sets up as LoadBalancer; use the host/Cluster IP and check the GitLab Gateway Kubernetes Service for the external port (something like 30001).

Enter that IP/Port and URL at the Webhook Integration URL. The Secret Token will be the same string used in the Brigade Project *values.yaml* `sharedSecret` property.

Check the boxes for the Trigger events to publish from the GitLab instance. SSL is optional.

## [Scripting Guide](docs/scripting.md)
tl;dr: GitLab Gateway produces 8 events: `push`, `tag`, `issue`, `comment`, `mergerequest`, `wikipage`, `pipeline`, `build`.


# Contributing

This project welcomes contributions and suggestions.

# License

MIT