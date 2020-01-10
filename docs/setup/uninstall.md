---
title: Kubed Uninstall
description: Kubed Uninstall
menu:
  product_kubed_{{ .version }}:
    identifier: kubed-uninstall
    name: Uninstall
    parent: setup
    weight: 20
product_name: kubed
menu_name: product_kubed_{{ .version }}
section_menu_id: setup
---

> New to Kubed? Please start [here](/docs/concepts/README.md).

# Uninstall Kubed

Please follow the steps below to uninstall Kubed:

<ul class="nav nav-tabs" id="installerTab" role="tablist">
  <li class="nav-item">
    <a class="nav-link active" id="helm3-tab" data-toggle="tab" href="#helm3" role="tab" aria-controls="helm3" aria-selected="true">Helm 3</a>
  </li>
  <li class="nav-item">
    <a class="nav-link" id="helm2-tab" data-toggle="tab" href="#helm2" role="tab" aria-controls="helm2" aria-selected="false">Helm 2</a>
  </li>
  <li class="nav-item">
    <a class="nav-link" id="script-tab" data-toggle="tab" href="#script" role="tab" aria-controls="script" aria-selected="false">YAML</a>
  </li>
</ul>
<div class="tab-content" id="installerTabContent">
  <div class="tab-pane fade show active" id="helm3" role="tabpanel" aria-labelledby="helm3-tab">

## Using Helm 3

In Helm 3, release names are [scoped to a namespace](https://v3.helm.sh/docs/faq/#release-names-are-now-scoped-to-the-namespace). So, provide the namespace you used to install the operator when installing.

```console
$ helm uninstall kubed --namespace kube-system
```

</div>
<div class="tab-pane fade" id="helm2" role="tabpanel" aria-labelledby="helm2-tab">

## Using Helm 2

```console
$ helm delete kubed
```

</div>
<div class="tab-pane fade" id="script" role="tabpanel" aria-labelledby="script-tab">

## Using YAML (with helm 3)

If you prefer to not use Helm, you can generate YAMLs from Kubed chart and uninstall using `kubectl`.

```console
$ helm template kubed appscode/kubed --namespace kube-system | kubectl delete -f -
```

</div>
</div>
