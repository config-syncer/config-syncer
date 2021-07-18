---
title: Contributing | Kubed
description: Contributing
menu:
  product_kubed_{{ .version }}:
    identifier: contributing-kubed
    name: Contributing
    parent: welcome
    weight: 15
product_name: kubed
menu_name: product_kubed_{{ .version }}
section_menu_id: welcome
url: /products/kubed/{{ .version }}/welcome/contributing/
aliases:
  - /products/kubed/{{ .version }}/CONTRIBUTING/
---

# Contribution Guidelines
Want to hack on Kubed?

AppsCode projects are [Apache 2.0 licensed](https://github.com/kubeops/kubed/blob/master/LICENSE) and accept contributions via
GitHub pull requests.  This document outlines some of the conventions on
development workflow, commit message formatting, contact points and other
resources to make it easier to get your contribution accepted.

## Certificate of Origin

By contributing to this project you agree to the Developer Certificate of
Origin (DCO). This document was created by the Linux Kernel community and is a
simple statement that you, as a contributor, have the legal right to make the
contribution. See the [DCO](https://github.com/kubeops/kubed/blob/master/DCO) file for details.

## Developer Guide

We have a [Developer Guide](/docs/setup/developer-guide/overview.md) that outlines everything you need to know from setting up your
dev environment to how to build and test Kubed. If you find something undocumented or incorrect along the way,
please feel free to send a Pull Request.

## Getting Help

If you have a question about Kubed or having problem using it, you can contact us on the [AppsCode Slack team](https://appscode.slack.com/messages/C6HSHCKBL/details/) channel `#kubed`. Follow [this link](https://slack.appscode.com) to get invitation to our Slack channel.

## Bugs/Feature request

If you have found a bug with Kubed or want to request for new features, please [file an issue](https://github.com/kubeops/kubed/issues/new).

## Submit PR

If you fix a bug or developed a new feature, feel free to submit a PR. In either case, please file a [Github issue](https://github.com/kubeops/kubed/issues/new) first, so that we can have a discussion on it. This is a rough outline of what a contributor's workflow looks like:

- Create a topic branch from where you want to base your work (usually master).
- Make commits of logical units.
- Push your changes to a topic branch in your fork of the repository.
- Make sure the tests pass, and add any new tests as appropriate.
- Submit a pull request to the original repository.

Thanks for your contributions!

## Spread the word

If you have written blog post or tutorial on Kubed, please share it with us on [Twitter](https://twitter.com/AppsCodeHQ) or [Slack](https://slack.appscode.com).
