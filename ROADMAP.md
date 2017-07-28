# Project Status

## Versioning Policy
Kubed __does not follow semver__. Currently Kubed operator implementation is considered alpha. Please report any issues you via Github. Once released, the _major_ version of operator is going to point to the Kubernetes [client-go](https://github.com/kubernetes/client-go#branches-and-tags) version. You can verify this from the `glide.yaml` file. This means there might be breaking changes between point releases of the operator. This generally manifests as changed annotation keys or their meaning. Please always check the release notes for upgrade instructions.

### Release 0.x.0
These are alpha releases done so that users can test Kubed.

### Release 3.0.0
This is going to be the first release of Kubed and uses Kubernetes client-go 3.0.0. We plan to mark the last 0.x.0 release as 3.0.0. This version will support Kubernetes 1.5 & 1.6 .

### Release 4.0.0
This relased will be based on client-go 4.0.0. This is going to include a number of breaking changes (example, turn Kubed into a UAS server) and be supported for Kubernetes 1.7+. Please see the issues in release milestone [here](https://github.com/appscode/kubed/milestone/3).
