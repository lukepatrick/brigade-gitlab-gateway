# Brigade GitLab Gateway

Send [GitLab events](https://gitlab.com/help/user/project/integrations/webhooks) into a [Brigade](https://github.com/Azure/brigade) pipeline. 

This is a Brigade gateway that listens to GitLab webhooks event stream and triggers events inside of Brigade.

## Prerequisites

1. Have a running [Kubernetes](https://kubernetes.io/docs/setup/) environment
2. Setup [Helm](https://github.com/kubernetes/helm)
3. Setup [Brigade](https://github.com/Azure/brigade) core

## Install

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

# Contributing

This project welcomes contributions and suggestions.

# License

MIT