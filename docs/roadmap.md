---
title: Roadmap | Kubed
description: Roadmap of kubed
menu:
  product_kubed_0.9.0:
    identifier: roadmap-kubed
    name: Roadmap
    parent: welcome
    weight: 20
product_name: kubed
menu_name: product_kubed_0.9.0
section_menu_id: welcome
url: /products/kubed/0.9.0/welcome/roadmap/
aliases:
  - /products/kubed/0.9.0/roadmap/
---

# Project Status

## Versioning Policy
Kubed __does not follow semver__. Currently Kubed operator implementation is considered alpha. Please report any issues you via Github. Once released, the _major_ version of operator is going to point to the Kubernetes [client-go](https://github.com/kubernetes/client-go#branches-and-tags) version. You can verify this from the `glide.yaml` file. This means there might be breaking changes between point releases of the operator. This generally manifests as changed annotation keys or their meaning. Please always check the release notes for upgrade instructions.
