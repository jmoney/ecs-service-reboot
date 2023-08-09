# ECS Service Reboot

## Overview

| Arguemment | Description | Default Value |
| --- | --- | --- |
| cluster | The ECS cluster to look for the service | The environment variable `CLUSTER` |
| region | The region to look for the service in ecs | The environment variable `AWS_REGION` |
| service | The ECS service to reboot | The environment variable `SERVICE` |

## Installation

```bash
brew tap jmoney/aws
brew install reboot-ecs-service
```

## Run Locally

```bash
go run cmd/cli/main.go -cluster dev -region us-east-1 -service my-service
```

Will reboot the service name `my-service` in the `dev cluster` running in the `us-east-1` region.
