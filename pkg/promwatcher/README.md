# Prometheus Watcher

Watch Prometheus TPR instance to deploy proxy for authentication.

## User Guide

To start Prometheus Watcher with `kubed`, use `--enable-prometheus-monitoring` flag.

```
$ kubed --enable-prometheus-monitoring
```

## Architectural Design

If you want to know how this is working, see following workflow

### Workflow
<p align="center">
    <img src="flow.png" width="715">
</p>

* ##### EventType `ADDED`
    1. When a Prometheus object is created, Watcher detects it.
    2. Controller creates deployment for Prometheus Proxy.
    3. Controller creates service for Prometheus Proxy.
    4. Controller adds Ingress rule.
        > Add HTTP rule with path `/prometheus-<prometheus-name>.<namespace>`


* ##### EventType `DELETED`
    1. When a Prometheus object is deleted, Watcher detects it.
    2. Controller deleted deployment of Prometheus Proxy.
    3. Controller deletes service for Prometheus Proxy.
    4. Controller removes Ingress rule.
