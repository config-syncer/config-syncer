---
title: Changelog | Kubed
description: Changelog
menu:
  product_kubed_{{ .version }}:
    identifier: changelog-kubed
    name: Changelog
    parent: welcome
    weight: 10
product_name: kubed
menu_name: product_kubed_{{ .version }}
section_menu_id: welcome
url: /products/kubed/{{ .version }}/welcome/changelog/
aliases:
  - /products/kubed/{{ .version }}/CHANGELOG/
---

# Change Log

## [Unreleased](https://github.com/kubeops/kubed/tree/HEAD)

[Full Changelog](https://github.com/kubeops/kubed/compare/v0.12.0...HEAD)

**Merged pull requests:**

- Update to Kubernetes v1.18.3 [\#398](https://github.com/kubeops/kubed/pull/398) ([tamalsaha](https://github.com/tamalsaha))

## [v0.12.0](https://github.com/kubeops/kubed/tree/v0.12.0) (2020-05-18)
[Full Changelog](https://github.com/kubeops/kubed/compare/v0.12.0-rc.3...v0.12.0)

**Closed issues:**

- unknown rbac rules [\#384](https://github.com/kubeops/kubed/issues/384)

**Merged pull requests:**

- Auto generate chart readme file [\#397](https://github.com/kubeops/kubed/pull/397) ([tamalsaha](https://github.com/tamalsaha))
- Clean up CI pipeline [\#396](https://github.com/kubeops/kubed/pull/396) ([tamalsaha](https://github.com/tamalsaha))

## [v0.12.0-rc.3](https://github.com/kubeops/kubed/tree/v0.12.0-rc.3) (2020-04-25)
[Full Changelog](https://github.com/kubeops/kubed/compare/v0.12.0-rc.2...v0.12.0-rc.3)

**Closed issues:**

- Context not found in kubeconfig file [\#389](https://github.com/kubeops/kubed/issues/389)
- Crash when the namespace is terminating [\#380](https://github.com/kubeops/kubed/issues/380)
- Don't regenerate a certificate authority and a certificate everytime you deploy the helm chart [\#371](https://github.com/kubeops/kubed/issues/371)
- Kubed shows abnormal high I/O and memory usage [\#357](https://github.com/kubeops/kubed/issues/357)

**Merged pull requests:**

- Prepare v0.12.0-rc.3 release [\#395](https://github.com/kubeops/kubed/pull/395) ([tamalsaha](https://github.com/tamalsaha))
- Publish Helm chart from release workflow [\#394](https://github.com/kubeops/kubed/pull/394) ([tamalsaha](https://github.com/tamalsaha))
- Apply various fixes to chart [\#393](https://github.com/kubeops/kubed/pull/393) ([tamalsaha](https://github.com/tamalsaha))
- Custom securityContext in template [\#392](https://github.com/kubeops/kubed/pull/392) ([jsrolon](https://github.com/jsrolon))
- Parameterizes run command and secure port in helm [\#390](https://github.com/kubeops/kubed/pull/390) ([masstamike](https://github.com/masstamike))
- Clean up Helm chart's README with removed attributes [\#387](https://github.com/kubeops/kubed/pull/387) ([olivierlemasle](https://github.com/olivierlemasle))
- Allow specifying rather than generating certs [\#385](https://github.com/kubeops/kubed/pull/385) ([tamalsaha](https://github.com/tamalsaha))

## [v0.12.0-rc.2](https://github.com/kubeops/kubed/tree/v0.12.0-rc.2) (2020-01-10)
[Full Changelog](https://github.com/kubeops/kubed/compare/v0.12.0-rc.1...v0.12.0-rc.2)

**Merged pull requests:**

- Permit kubed to write events [\#379](https://github.com/kubeops/kubed/pull/379) ([tamalsaha](https://github.com/tamalsaha))

## [v0.12.0-rc.1](https://github.com/kubeops/kubed/tree/v0.12.0-rc.1) (2020-01-10)
[Full Changelog](https://github.com/kubeops/kubed/compare/v0.12.0-rc.0...v0.12.0-rc.1)

**Closed issues:**

- Kubernetes 1.16 and "extensions/v1beta1" [\#369](https://github.com/kubeops/kubed/issues/369)
- Project roadmap [\#363](https://github.com/kubeops/kubed/issues/363)
- kubed is messing up api-resources [\#351](https://github.com/kubeops/kubed/issues/351)
- Don't regenerate a certificate authority and a certificate everytime you deploy the helm chart [\#346](https://github.com/kubeops/kubed/issues/346)
- liveness check does not detect when kubed is not responding to API requests anymore [\#343](https://github.com/kubeops/kubed/issues/343)
- kube-controller-manager errors [\#340](https://github.com/kubeops/kubed/issues/340)
- Performance optimizing syncer [\#335](https://github.com/kubeops/kubed/issues/335)
- How to properly restore snapshot [\#303](https://github.com/kubeops/kubed/issues/303)
- Use audit policy server to forward events [\#194](https://github.com/kubeops/kubed/issues/194)
- Auto mount image pull secret for docker registry [\#191](https://github.com/kubeops/kubed/issues/191)
- Need automatic clean up of backed up yamls. [\#169](https://github.com/kubeops/kubed/issues/169)
- Use new Events api in 1.9 [\#148](https://github.com/kubeops/kubed/issues/148)
- watches for changes to ConfigMap objects and performs rolling upgrades on their associated deployments [\#145](https://github.com/kubeops/kubed/issues/145)
- Watch cloud provider specific forwarder. [\#108](https://github.com/kubeops/kubed/issues/108)
- Include processes running on host when OOM is reported. [\#107](https://github.com/kubeops/kubed/issues/107)
- Notifier routing [\#101](https://github.com/kubeops/kubed/issues/101)
- Show event UID in SMS/Chat version [\#93](https://github.com/kubeops/kubed/issues/93)
- Setup retention policy for snapshot operation [\#75](https://github.com/kubeops/kubed/issues/75)
- Add Dry Run option for janitors [\#60](https://github.com/kubeops/kubed/issues/60)
- Explore generic DELETE watcher [\#41](https://github.com/kubeops/kubed/issues/41)

**Merged pull requests:**

- Prepare v0.12.0-rc.1 release [\#378](https://github.com/kubeops/kubed/pull/378) ([tamalsaha](https://github.com/tamalsaha))
- Remove Kubernetes dependency in chart [\#377](https://github.com/kubeops/kubed/pull/377) ([tamalsaha](https://github.com/tamalsaha))

## [v0.12.0-rc.0](https://github.com/kubeops/kubed/tree/v0.12.0-rc.0) (2020-01-10)
[Full Changelog](https://github.com/kubeops/kubed/compare/v0.11.0...v0.12.0-rc.0)

**Merged pull requests:**

- Update syncer docs [\#376](https://github.com/kubeops/kubed/pull/376) ([tamalsaha](https://github.com/tamalsaha))
- Fix intra cluster sync docs [\#375](https://github.com/kubeops/kubed/pull/375) ([tamalsaha](https://github.com/tamalsaha))
- Update docs for v0.12.0-rc.0 [\#374](https://github.com/kubeops/kubed/pull/374) ([tamalsaha](https://github.com/tamalsaha))
- Reboot kubed project [\#373](https://github.com/kubeops/kubed/pull/373) ([tamalsaha](https://github.com/tamalsaha))
- Delete script based installer [\#372](https://github.com/kubeops/kubed/pull/372) ([tamalsaha](https://github.com/tamalsaha))
- Change old "extensions/v1beta1" to new "apps/v1" [\#370](https://github.com/kubeops/kubed/pull/370) ([ruzickap](https://github.com/ruzickap))
- helm chart: Add kube-config into secret if needed [\#368](https://github.com/kubeops/kubed/pull/368) ([lonwern](https://github.com/lonwern))
- Download onessl version v0.13.1 for Kubernetes 1.16 fix [\#367](https://github.com/kubeops/kubed/pull/367) ([tamalsaha](https://github.com/tamalsaha))
- Templatize front matter [\#366](https://github.com/kubeops/kubed/pull/366) ([tamalsaha](https://github.com/tamalsaha))

## [v0.11.0](https://github.com/kubeops/kubed/tree/v0.11.0) (2019-09-10)
[Full Changelog](https://github.com/kubeops/kubed/compare/0.10.0...v0.11.0)

**Closed issues:**

- Feature request:  specify source namespaces via configuration \(one or all\) [\#355](https://github.com/kubeops/kubed/issues/355)

**Merged pull requests:**

- Use v0.11.0 instead of 0.11.0 [\#365](https://github.com/kubeops/kubed/pull/365) ([tamalsaha](https://github.com/tamalsaha))
- Prepare docs for 0.11.0 release  [\#364](https://github.com/kubeops/kubed/pull/364) ([tamalsaha](https://github.com/tamalsaha))
- Use osm package from kmodules.xyz/objectstore-api [\#362](https://github.com/kubeops/kubed/pull/362) ([tamalsaha](https://github.com/tamalsaha))
- Fix nil pointer exception [\#361](https://github.com/kubeops/kubed/pull/361) ([tamalsaha](https://github.com/tamalsaha))
- Update dependencies [\#360](https://github.com/kubeops/kubed/pull/360) ([tamalsaha](https://github.com/tamalsaha))
- Add Makefile [\#359](https://github.com/kubeops/kubed/pull/359) ([tamalsaha](https://github.com/tamalsaha))
- Implementation of  \#355.  [\#356](https://github.com/kubeops/kubed/pull/356) ([gralfca](https://github.com/gralfca))
- Use absolute path as aliases for reference docs [\#353](https://github.com/kubeops/kubed/pull/353) ([tamalsaha](https://github.com/tamalsaha))
- Update to k8s 1.14.0 client libraries using go.mod [\#352](https://github.com/kubeops/kubed/pull/352) ([tamalsaha](https://github.com/tamalsaha))
- Update go-notify & envconfig packages [\#350](https://github.com/kubeops/kubed/pull/350) ([tamalsaha](https://github.com/tamalsaha))
- Add notifier secret to chart [\#349](https://github.com/kubeops/kubed/pull/349) ([tamalsaha](https://github.com/tamalsaha))
- Remove notifier instructions for Hipchat and Stride [\#348](https://github.com/kubeops/kubed/pull/348) ([tamalsaha](https://github.com/tamalsaha))
- Use Backend api objects from kmodules/objectstore-api. [\#347](https://github.com/kubeops/kubed/pull/347) ([tamalsaha](https://github.com/tamalsaha))

## [0.10.0](https://github.com/kubeops/kubed/tree/0.10.0) (2019-04-23)
[Full Changelog](https://github.com/kubeops/kubed/compare/0.9.0...0.10.0)

**Closed issues:**

- helm deployment doesn't use pre-existing kubed secret [\#312](https://github.com/kubeops/kubed/issues/312)
- kubed is using too much memory [\#304](https://github.com/kubeops/kubed/issues/304)

**Merged pull requests:**

- Prepare docs for 0.10.0 release [\#345](https://github.com/kubeops/kubed/pull/345) ([tamalsaha](https://github.com/tamalsaha))
- Update Kubernetes client libraries to 1.13.5 [\#344](https://github.com/kubeops/kubed/pull/344) ([tamalsaha](https://github.com/tamalsaha))
- Improved code style and saying hello to good practices! [\#342](https://github.com/kubeops/kubed/pull/342) ([AnikHasibul](https://github.com/AnikHasibul))
- Update Kubernetes client libraries to 1.13.0 [\#341](https://github.com/kubeops/kubed/pull/341) ([tamalsaha](https://github.com/tamalsaha))
- Pass pod annotation to deployment [\#339](https://github.com/kubeops/kubed/pull/339) ([tamalsaha](https://github.com/tamalsaha))
- Don't use priority class when operator namespace is not kube-system [\#338](https://github.com/kubeops/kubed/pull/338) ([tamalsaha](https://github.com/tamalsaha))
- Use onessl 0.10.0 [\#337](https://github.com/kubeops/kubed/pull/337) ([tamalsaha](https://github.com/tamalsaha))
- Fix the case for deploying using MINGW64 for windows [\#336](https://github.com/kubeops/kubed/pull/336) ([tamalsaha](https://github.com/tamalsaha))
- Adds option to allocate/use pvc with helm installation [\#334](https://github.com/kubeops/kubed/pull/334) ([DerekHeldtWerle](https://github.com/DerekHeldtWerle))
- Add the ability to insert config sections directly [\#333](https://github.com/kubeops/kubed/pull/333) ([pirtoo](https://github.com/pirtoo))
- Remove apiserver.ca from chart and update onessl [\#332](https://github.com/kubeops/kubed/pull/332) ([tamalsaha](https://github.com/tamalsaha))
- Add certificate health checker [\#331](https://github.com/kubeops/kubed/pull/331) ([tamalsaha](https://github.com/tamalsaha))

## [0.9.0](https://github.com/kubeops/kubed/tree/0.9.0) (2018-12-17)
[Full Changelog](https://github.com/kubeops/kubed/compare/0.8.0...0.9.0)

**Fixed bugs:**

- Fix analytics flag [\#305](https://github.com/kubeops/kubed/pull/305) ([tamalsaha](https://github.com/tamalsaha))

**Closed issues:**

- Dependabot couldn't find a Gopkg.toml for this project [\#320](https://github.com/kubeops/kubed/issues/320)
- New version release [\#318](https://github.com/kubeops/kubed/issues/318)
- Kubed fails to start on initial start [\#316](https://github.com/kubeops/kubed/issues/316)
- Resource Requests and Limits in helm chart [\#315](https://github.com/kubeops/kubed/issues/315)

**Merged pull requests:**

- Update osm version to 0.9.1 [\#329](https://github.com/kubeops/kubed/pull/329) ([tamalsaha](https://github.com/tamalsaha))
- Update dependencies [\#328](https://github.com/kubeops/kubed/pull/328) ([tamalsaha](https://github.com/tamalsaha))
- Permit specifying compute resources for the kubed container. [\#327](https://github.com/kubeops/kubed/pull/327) ([niclic](https://github.com/niclic))
- Use rbac/v1 api [\#325](https://github.com/kubeops/kubed/pull/325) ([tamalsaha](https://github.com/tamalsaha))
- Prepare docs for 0.9.0 release [\#324](https://github.com/kubeops/kubed/pull/324) ([tamalsaha](https://github.com/tamalsaha))
- Update osm version to 0.9.0 [\#323](https://github.com/kubeops/kubed/pull/323) ([tamalsaha](https://github.com/tamalsaha))
- Use flags.DumpAll to dump flags [\#322](https://github.com/kubeops/kubed/pull/322) ([tamalsaha](https://github.com/tamalsaha))
- Set periodic analytics [\#321](https://github.com/kubeops/kubed/pull/321) ([tamalsaha](https://github.com/tamalsaha))
- Update Kubernetes client libraries to 1.12.0 [\#319](https://github.com/kubeops/kubed/pull/319) ([tamalsaha](https://github.com/tamalsaha))
- Update kubernetes client libraries to 1.12.0 [\#314](https://github.com/kubeops/kubed/pull/314) ([tamalsaha](https://github.com/tamalsaha))
- Check if Kubernetes version is supported before running operator [\#313](https://github.com/kubeops/kubed/pull/313) ([tamalsaha](https://github.com/tamalsaha))
- Use kubernetes-1.11.3 [\#311](https://github.com/kubeops/kubed/pull/311) ([tamalsaha](https://github.com/tamalsaha))
-  Update pipeline [\#310](https://github.com/kubeops/kubed/pull/310) ([tahsinrahman](https://github.com/tahsinrahman))
- fix uninstall for concourse [\#309](https://github.com/kubeops/kubed/pull/309) ([tahsinrahman](https://github.com/tahsinrahman))
- Improve Helm chart options [\#308](https://github.com/kubeops/kubed/pull/308) ([tamalsaha](https://github.com/tamalsaha))
- Revendor apis [\#307](https://github.com/kubeops/kubed/pull/307) ([tamalsaha](https://github.com/tamalsaha))
- Use concourse scripts from libbuild [\#306](https://github.com/kubeops/kubed/pull/306) ([tahsinrahman](https://github.com/tahsinrahman))
- Fix extended apiserver issues with Kubernetes 1.11 [\#302](https://github.com/kubeops/kubed/pull/302) ([tamalsaha](https://github.com/tamalsaha))

## [0.8.0](https://github.com/kubeops/kubed/tree/0.8.0) (2018-07-10)
[Full Changelog](https://github.com/kubeops/kubed/compare/0.7.0...0.8.0)

**Fixed bugs:**

- Handle syncing for updated namespaces [\#299](https://github.com/kubeops/kubed/pull/299) ([tamalsaha](https://github.com/tamalsaha))
- Remove infinite spin loop from operator [\#294](https://github.com/kubeops/kubed/pull/294) ([tamalsaha](https://github.com/tamalsaha))

**Merged pull requests:**

- Prepare docs for 0.8.0 [\#301](https://github.com/kubeops/kubed/pull/301) ([tamalsaha](https://github.com/tamalsaha))
- Add chart config for event forwarder and recycle bin [\#300](https://github.com/kubeops/kubed/pull/300) ([tamalsaha](https://github.com/tamalsaha))
- Improve logging for syncer [\#298](https://github.com/kubeops/kubed/pull/298) ([tamalsaha](https://github.com/tamalsaha))
- Expose webhook server to expose operator metrics [\#297](https://github.com/kubeops/kubed/pull/297) ([tamalsaha](https://github.com/tamalsaha))
- Remove outdated installer links [\#296](https://github.com/kubeops/kubed/pull/296) ([tamalsaha](https://github.com/tamalsaha))
- Use yaml file to create service account in installer script [\#295](https://github.com/kubeops/kubed/pull/295) ([tamalsaha](https://github.com/tamalsaha))
- Use osm 0.7.1 [\#293](https://github.com/kubeops/kubed/pull/293) ([tamalsaha](https://github.com/tamalsaha))
- Deploy in kube-system namespace using Helm [\#292](https://github.com/kubeops/kubed/pull/292) ([tamalsaha](https://github.com/tamalsaha))
- Update client-go to v8.0.0 [\#291](https://github.com/kubeops/kubed/pull/291) ([tamalsaha](https://github.com/tamalsaha))
- Format shell script [\#290](https://github.com/kubeops/kubed/pull/290) ([tamalsaha](https://github.com/tamalsaha))
- Fix openapi schema for metav1.Duration [\#289](https://github.com/kubeops/kubed/pull/289) ([tamalsaha](https://github.com/tamalsaha))
- Move openapi-spec to api folder [\#288](https://github.com/kubeops/kubed/pull/288) ([tamalsaha](https://github.com/tamalsaha))
- Add togglable tabs for Installation: Script & Helm [\#287](https://github.com/kubeops/kubed/pull/287) ([sajibcse68](https://github.com/sajibcse68))

## [0.7.0](https://github.com/kubeops/kubed/tree/0.7.0) (2018-06-01)
[Full Changelog](https://github.com/kubeops/kubed/compare/0.7.0-rc.2...0.7.0)

**Fixed bugs:**

- secrets syncing is not proper [\#233](https://github.com/kubeops/kubed/issues/233)
- Ensure bad backups are not used to overwrite last good backup [\#176](https://github.com/kubeops/kubed/issues/176)

**Closed issues:**

- kubectl returns results super slow after installing kubed [\#279](https://github.com/kubeops/kubed/issues/279)
- Event Forwarder Hipchat notifier sends messages not notifications [\#260](https://github.com/kubeops/kubed/issues/260)
- Fix backup manage RBAC issue [\#256](https://github.com/kubeops/kubed/issues/256)
- Fix tests [\#240](https://github.com/kubeops/kubed/issues/240)
- invalid header field value error when setting up with S3. [\#161](https://github.com/kubeops/kubed/issues/161)
- Restart kubed in e2e tests when config.yaml changes [\#158](https://github.com/kubeops/kubed/issues/158)

**Merged pull requests:**

- Prepare 0.7.0 release [\#286](https://github.com/kubeops/kubed/pull/286) ([tamalsaha](https://github.com/tamalsaha))

## [0.7.0-rc.2](https://github.com/kubeops/kubed/tree/0.7.0-rc.2) (2018-05-31)
[Full Changelog](https://github.com/kubeops/kubed/compare/0.7.0-rc.1...0.7.0-rc.2)

**Merged pull requests:**

- Update changelog [\#285](https://github.com/kubeops/kubed/pull/285) ([tamalsaha](https://github.com/tamalsaha))
- Add document for Stride [\#284](https://github.com/kubeops/kubed/pull/284) ([tamalsaha](https://github.com/tamalsaha))
- Prepare docs for 0.7.0-rc.2 release [\#283](https://github.com/kubeops/kubed/pull/283) ([tamalsaha](https://github.com/tamalsaha))
- Disable api server by default in 1.8 cluster. [\#282](https://github.com/kubeops/kubed/pull/282) ([tamalsaha](https://github.com/tamalsaha))
- Fix grammar [\#281](https://github.com/kubeops/kubed/pull/281) ([tamalsaha](https://github.com/tamalsaha))
- Allow setting cluster-name during installation [\#280](https://github.com/kubeops/kubed/pull/280) ([tamalsaha](https://github.com/tamalsaha))

## [0.7.0-rc.1](https://github.com/kubeops/kubed/tree/0.7.0-rc.1) (2018-05-30)
[Full Changelog](https://github.com/kubeops/kubed/compare/0.7.0-rc.0...0.7.0-rc.1)

**Fixed bugs:**

- ConfigSyncer does not sync configmap/secret when new namespace is created [\#266](https://github.com/kubeops/kubed/issues/266)

**Merged pull requests:**

- Prepare docs for 0.7.0-rc.1 [\#278](https://github.com/kubeops/kubed/pull/278) ([tamalsaha](https://github.com/tamalsaha))
- Fixed secret type of synced secret [\#277](https://github.com/kubeops/kubed/pull/277) ([hossainemruz](https://github.com/hossainemruz))
- concourse - delete cluster on exit [\#275](https://github.com/kubeops/kubed/pull/275) ([tahsinrahman](https://github.com/tahsinrahman))

## [0.7.0-rc.0](https://github.com/kubeops/kubed/tree/0.7.0-rc.0) (2018-05-28)
[Full Changelog](https://github.com/kubeops/kubed/compare/0.6.0-rc.0...0.7.0-rc.0)

**Fixed bugs:**

- ClusterRole kubed-operator needs 'patch' and 'delete' permissions for configmaps/secrets resources [\#267](https://github.com/kubeops/kubed/issues/267)
- Cron lib keeps running every 1h [\#83](https://github.com/kubeops/kubed/issues/83)
- Fix backup command [\#254](https://github.com/kubeops/kubed/pull/254) ([tamalsaha](https://github.com/tamalsaha))

**Closed issues:**

- Check that client-ca and requestheader-ca are not same [\#242](https://github.com/kubeops/kubed/issues/242)
- Support self-signed CA for Minio [\#241](https://github.com/kubeops/kubed/issues/241)
- List and delete all old indices matching prefix [\#177](https://github.com/kubeops/kubed/issues/177)
- Certificate signer [\#147](https://github.com/kubeops/kubed/issues/147)
- Extract docker LABELS [\#139](https://github.com/kubeops/kubed/issues/139)
- Enforce Pod policy via admission webhook [\#118](https://github.com/kubeops/kubed/issues/118)
- Rethink copying config/secret to kube-public namespace [\#113](https://github.com/kubeops/kubed/issues/113)
- Kubed api features [\#86](https://github.com/kubeops/kubed/issues/86)
- Log warnings against kubed-config [\#81](https://github.com/kubeops/kubed/issues/81)
- Perform as a Image review process [\#72](https://github.com/kubeops/kubed/issues/72)
- k8sguard [\#22](https://github.com/kubeops/kubed/issues/22)

**Merged pull requests:**

- Update changelog [\#276](https://github.com/kubeops/kubed/pull/276) ([tamalsaha](https://github.com/tamalsaha))
- Use same config for chart and script installer [\#274](https://github.com/kubeops/kubed/pull/274) ([tamalsaha](https://github.com/tamalsaha))
- Prepare docs for 7.0.0-rc.0 [\#273](https://github.com/kubeops/kubed/pull/273) ([tamalsaha](https://github.com/tamalsaha))
- Add concourse test [\#272](https://github.com/kubeops/kubed/pull/272) ([tahsinrahman](https://github.com/tahsinrahman))
- Improve installer [\#271](https://github.com/kubeops/kubed/pull/271) ([tamalsaha](https://github.com/tamalsaha))
- Improve e2e test [\#270](https://github.com/kubeops/kubed/pull/270) ([hossainemruz](https://github.com/hossainemruz))
- Revendor dependencies [\#269](https://github.com/kubeops/kubed/pull/269) ([tamalsaha](https://github.com/tamalsaha))
- Add missing RBAC rules [\#268](https://github.com/kubeops/kubed/pull/268) ([hossainemruz](https://github.com/hossainemruz))
- Don't panic if admission options is nil [\#264](https://github.com/kubeops/kubed/pull/264) ([tamalsaha](https://github.com/tamalsaha))
- Disable admission controllers for webhook server [\#263](https://github.com/kubeops/kubed/pull/263) ([tamalsaha](https://github.com/tamalsaha))
- Sync secret Kind [\#262](https://github.com/kubeops/kubed/pull/262) ([farcaller](https://github.com/farcaller))
- Update client-go to 7.0.0 [\#261](https://github.com/kubeops/kubed/pull/261) ([tamalsaha](https://github.com/tamalsaha))
- Support private registry for chart [\#259](https://github.com/kubeops/kubed/pull/259) ([tamalsaha](https://github.com/tamalsaha))
- Improve installer [\#258](https://github.com/kubeops/kubed/pull/258) ([tamalsaha](https://github.com/tamalsaha))
- Add support for SSL certificate for S3 compatible custom server \(i.e. Minio\) [\#257](https://github.com/kubeops/kubed/pull/257) ([hossainemruz](https://github.com/hossainemruz))
- Rename snapshot command to backup [\#255](https://github.com/kubeops/kubed/pull/255) ([tamalsaha](https://github.com/tamalsaha))
- Correctly load default config [\#253](https://github.com/kubeops/kubed/pull/253) ([tamalsaha](https://github.com/tamalsaha))
- Add RBAC instructions for GKE cluster [\#252](https://github.com/kubeops/kubed/pull/252) ([tamalsaha](https://github.com/tamalsaha))
- Update chart repository location [\#251](https://github.com/kubeops/kubed/pull/251) ([tamalsaha](https://github.com/tamalsaha))
- Support installing from local installer scripts [\#250](https://github.com/kubeops/kubed/pull/250) ([tamalsaha](https://github.com/tamalsaha))
- Move swagger.json to openapi-spec/v2 [\#249](https://github.com/kubeops/kubed/pull/249) ([tamalsaha](https://github.com/tamalsaha))
- Generate swagger.json [\#248](https://github.com/kubeops/kubed/pull/248) ([tamalsaha](https://github.com/tamalsaha))
- Generate openapi spec [\#247](https://github.com/kubeops/kubed/pull/247) ([tamalsaha](https://github.com/tamalsaha))
- Delete internal clientset [\#246](https://github.com/kubeops/kubed/pull/246) ([tamalsaha](https://github.com/tamalsaha))
- Revendor dependencies [\#245](https://github.com/kubeops/kubed/pull/245) ([tamalsaha](https://github.com/tamalsaha))
- Skip downloading onessl if already exists [\#244](https://github.com/kubeops/kubed/pull/244) ([tamalsaha](https://github.com/tamalsaha))
- Rename --analytics to --enable-analytics [\#243](https://github.com/kubeops/kubed/pull/243) ([tamalsaha](https://github.com/tamalsaha))
- Add travis yaml [\#239](https://github.com/kubeops/kubed/pull/239) ([tahsinrahman](https://github.com/tahsinrahman))
- Update chart to match new config format [\#238](https://github.com/kubeops/kubed/pull/238) ([tamalsaha](https://github.com/tamalsaha))
- Remove reference to Voyager [\#237](https://github.com/kubeops/kubed/pull/237) ([tamalsaha](https://github.com/tamalsaha))
- Make it clear that installer is a single command [\#236](https://github.com/kubeops/kubed/pull/236) ([tamalsaha](https://github.com/tamalsaha))
- Fix installer [\#235](https://github.com/kubeops/kubed/pull/235) ([tamalsaha](https://github.com/tamalsaha))
- Update chart to match RBAC best practices for charts [\#234](https://github.com/kubeops/kubed/pull/234) ([tamalsaha](https://github.com/tamalsaha))
- Add checks to installer script [\#232](https://github.com/kubeops/kubed/pull/232) ([tamalsaha](https://github.com/tamalsaha))

## [0.6.0-rc.0](https://github.com/kubeops/kubed/tree/0.6.0-rc.0) (2018-03-03)
[Full Changelog](https://github.com/kubeops/kubed/compare/0.5.0...0.6.0-rc.0)

**Closed issues:**

- Cleanup indexed data [\#212](https://github.com/kubeops/kubed/issues/212)
- Delete search index when namespace is deleted [\#210](https://github.com/kubeops/kubed/issues/210)
- Allow configuring resource types for Add/Update notification [\#192](https://github.com/kubeops/kubed/issues/192)
- Used shared informer and queue [\#152](https://github.com/kubeops/kubed/issues/152)
- Use separate bleve index for Events and other kube api objects [\#106](https://github.com/kubeops/kubed/issues/106)
- Support soft loading of ClusterConfig [\#51](https://github.com/kubeops/kubed/issues/51)
- Expose Kubed api server as a UAS [\#19](https://github.com/kubeops/kubed/issues/19)

**Merged pull requests:**

- Fix docs [\#231](https://github.com/kubeops/kubed/pull/231) ([tamalsaha](https://github.com/tamalsaha))
- Update docs for 0.6.0-rc.0 release [\#230](https://github.com/kubeops/kubed/pull/230) ([tamalsaha](https://github.com/tamalsaha))
- Upgrade github.com/blevesearch/bleve to 0.7.0 [\#229](https://github.com/kubeops/kubed/pull/229) ([tamalsaha](https://github.com/tamalsaha))
- Use github.com/json-iterator/go [\#228](https://github.com/kubeops/kubed/pull/228) ([tamalsaha](https://github.com/tamalsaha))
- Remove unused options field [\#227](https://github.com/kubeops/kubed/pull/227) ([tamalsaha](https://github.com/tamalsaha))
- Sync chart to stable charts repo [\#226](https://github.com/kubeops/kubed/pull/226) ([tamalsaha](https://github.com/tamalsaha))
- Generate internal types [\#225](https://github.com/kubeops/kubed/pull/225) ([tamalsaha](https://github.com/tamalsaha))
- Use rbac/v1 apis [\#224](https://github.com/kubeops/kubed/pull/224) ([tamalsaha](https://github.com/tamalsaha))
- Create user facing aggregate roles [\#223](https://github.com/kubeops/kubed/pull/223) ([tamalsaha](https://github.com/tamalsaha))
- Use official code generator scripts [\#222](https://github.com/kubeops/kubed/pull/222) ([tamalsaha](https://github.com/tamalsaha))
- Update charts to support api registration [\#221](https://github.com/kubeops/kubed/pull/221) ([tamalsaha](https://github.com/tamalsaha))
- Use ${} form for onessl envsubst [\#220](https://github.com/kubeops/kubed/pull/220) ([tamalsaha](https://github.com/tamalsaha))
- Update .gitignore file [\#219](https://github.com/kubeops/kubed/pull/219) ([tamalsaha](https://github.com/tamalsaha))
- Rename Stuff back to SearchResult [\#218](https://github.com/kubeops/kubed/pull/218) ([tamalsaha](https://github.com/tamalsaha))
- Fix locks in resource indexer [\#217](https://github.com/kubeops/kubed/pull/217) ([tamalsaha](https://github.com/tamalsaha))
- Move apis out of pkg package [\#216](https://github.com/kubeops/kubed/pull/216) ([tamalsaha](https://github.com/tamalsaha))
- Document recent changes [\#215](https://github.com/kubeops/kubed/pull/215) ([tamalsaha](https://github.com/tamalsaha))
- Rename searchresult to stuff [\#214](https://github.com/kubeops/kubed/pull/214) ([tamalsaha](https://github.com/tamalsaha))
- Add installer script [\#211](https://github.com/kubeops/kubed/pull/211) ([tamalsaha](https://github.com/tamalsaha))
- Add tests for RestMapper [\#209](https://github.com/kubeops/kubed/pull/209) ([tamalsaha](https://github.com/tamalsaha))
- Set GroupVersionKind for event handlers [\#208](https://github.com/kubeops/kubed/pull/208) ([tamalsaha](https://github.com/tamalsaha))
- Rename api package by version [\#207](https://github.com/kubeops/kubed/pull/207) ([tamalsaha](https://github.com/tamalsaha))
- Properly handle update events for trashcan [\#206](https://github.com/kubeops/kubed/pull/206) ([tamalsaha](https://github.com/tamalsaha))
- Use fsnotify from kutil [\#205](https://github.com/kubeops/kubed/pull/205) ([tamalsaha](https://github.com/tamalsaha))
- Fix NPE [\#204](https://github.com/kubeops/kubed/pull/204) ([tamalsaha](https://github.com/tamalsaha))
- Generate DeepCopy methods for ClusterConfig [\#203](https://github.com/kubeops/kubed/pull/203) ([tamalsaha](https://github.com/tamalsaha))
- Fix config validator for event forwarder [\#202](https://github.com/kubeops/kubed/pull/202) ([tamalsaha](https://github.com/tamalsaha))
- Transform event forwarder rules to rules format [\#201](https://github.com/kubeops/kubed/pull/201) ([tamalsaha](https://github.com/tamalsaha))
- Split Setup\(\) into New\(\) and Configure\(\) [\#199](https://github.com/kubeops/kubed/pull/199) ([tamalsaha](https://github.com/tamalsaha))
- Remove reverse index [\#198](https://github.com/kubeops/kubed/pull/198) ([tamalsaha](https://github.com/tamalsaha))
- Update bleve to v0.6.0 [\#197](https://github.com/kubeops/kubed/pull/197) ([tamalsaha](https://github.com/tamalsaha))
- Turn kubed api server into an EAS [\#196](https://github.com/kubeops/kubed/pull/196) ([tamalsaha](https://github.com/tamalsaha))
- Allow configuring resource types for Add/Update notification [\#195](https://github.com/kubeops/kubed/pull/195) ([tamalsaha](https://github.com/tamalsaha))
- Use SharedInformerFactory [\#193](https://github.com/kubeops/kubed/pull/193) ([tamalsaha](https://github.com/tamalsaha))
- Support soft loading of ClusterConfig [\#125](https://github.com/kubeops/kubed/pull/125) ([tamalsaha](https://github.com/tamalsaha))

## [0.5.0](https://github.com/kubeops/kubed/tree/0.5.0) (2018-01-17)
[Full Changelog](https://github.com/kubeops/kubed/compare/0.4.0...0.5.0)

**Closed issues:**

- Support syncing config across clusters [\#144](https://github.com/kubeops/kubed/issues/144)

**Merged pull requests:**

- Prepare docs for 0.5.0 [\#190](https://github.com/kubeops/kubed/pull/190) ([tamalsaha](https://github.com/tamalsaha))
- Update changelog for 0.5.0 [\#189](https://github.com/kubeops/kubed/pull/189) ([tamalsaha](https://github.com/tamalsaha))
- Document valid time units for janitor TTL [\#188](https://github.com/kubeops/kubed/pull/188) ([tamalsaha](https://github.com/tamalsaha))
- Reset shard duration for influx janitor [\#187](https://github.com/kubeops/kubed/pull/187) ([tamalsaha](https://github.com/tamalsaha))
- Set min retention policy for kubed influx janitor [\#186](https://github.com/kubeops/kubed/pull/186) ([tamalsaha](https://github.com/tamalsaha))
- Log influx janitor result [\#185](https://github.com/kubeops/kubed/pull/185) ([tamalsaha](https://github.com/tamalsaha))
- Update github.com/influxdata/influxdb to v1.3.3 [\#184](https://github.com/kubeops/kubed/pull/184) ([tamalsaha](https://github.com/tamalsaha))
- Increase burst and qps for kube client [\#183](https://github.com/kubeops/kubed/pull/183) ([tamalsaha](https://github.com/tamalsaha))
- Update github.com/influxdata/influxdb to v1.1.1 [\#182](https://github.com/kubeops/kubed/pull/182) ([tamalsaha](https://github.com/tamalsaha))
- Update Elasticsearch client to olivere/elastic.v5 [\#181](https://github.com/kubeops/kubed/pull/181) ([tamalsaha](https://github.com/tamalsaha))
- Add Telegram as notifier [\#180](https://github.com/kubeops/kubed/pull/180) ([tamalsaha](https://github.com/tamalsaha))
- Delete all older indices prior to a date [\#179](https://github.com/kubeops/kubed/pull/179) ([aerokite](https://github.com/aerokite))
- Ensure bad backups are not used to overwrite last good backup [\#178](https://github.com/kubeops/kubed/pull/178) ([tamalsaha](https://github.com/tamalsaha))

## [0.4.0](https://github.com/kubeops/kubed/tree/0.4.0) (2018-01-08)
[Full Changelog](https://github.com/kubeops/kubed/compare/0.3.1...0.4.0)

**Closed issues:**

- Config/Secret Target selected namespaces via Annotation [\#150](https://github.com/kubeops/kubed/issues/150)

**Merged pull requests:**

- Fixed docs for syncer [\#175](https://github.com/kubeops/kubed/pull/175) ([diptadas](https://github.com/diptadas))
- Update changelog [\#174](https://github.com/kubeops/kubed/pull/174) ([tamalsaha](https://github.com/tamalsaha))
- Reorganize docs for hosting on product site [\#173](https://github.com/kubeops/kubed/pull/173) ([tamalsaha](https://github.com/tamalsaha))
- Add support for new DB types [\#172](https://github.com/kubeops/kubed/pull/172) ([tamalsaha](https://github.com/tamalsaha))
- Rename `kubeConfig` -\> `kubeConfigFile` [\#171](https://github.com/kubeops/kubed/pull/171) ([tamalsaha](https://github.com/tamalsaha))
- Update docs for syncer [\#170](https://github.com/kubeops/kubed/pull/170) ([diptadas](https://github.com/diptadas))
- Fix analytics client-id detection [\#168](https://github.com/kubeops/kubed/pull/168) ([tamalsaha](https://github.com/tamalsaha))
- Auto detect AWS bucket region [\#166](https://github.com/kubeops/kubed/pull/166) ([tamalsaha](https://github.com/tamalsaha))
- Support hipchat server [\#165](https://github.com/kubeops/kubed/pull/165) ([tamalsaha](https://github.com/tamalsaha))
- Write event for syncer origin conflict [\#164](https://github.com/kubeops/kubed/pull/164) ([diptadas](https://github.com/diptadas))
- Fix Syncer [\#163](https://github.com/kubeops/kubed/pull/163) ([diptadas](https://github.com/diptadas))
- Remove unnecessary IsPreferredAPIResource api calls [\#162](https://github.com/kubeops/kubed/pull/162) ([tamalsaha](https://github.com/tamalsaha))
- Sync configmap/secret to selected namespaces/contexts [\#154](https://github.com/kubeops/kubed/pull/154) ([diptadas](https://github.com/diptadas))

## [0.3.1](https://github.com/kubeops/kubed/tree/0.3.1) (2017-12-21)
[Full Changelog](https://github.com/kubeops/kubed/compare/0.3.0...0.3.1)

**Fixed bugs:**

- Support region for s3 backend [\#159](https://github.com/kubeops/kubed/pull/159) ([tamalsaha](https://github.com/tamalsaha))

**Closed issues:**

- Audit report [\#136](https://github.com/kubeops/kubed/issues/136)
- s3 snapshotter try to list all buckets [\#133](https://github.com/kubeops/kubed/issues/133)

**Merged pull requests:**

- Prepare docs for 0.3.1 [\#160](https://github.com/kubeops/kubed/pull/160) ([tamalsaha](https://github.com/tamalsaha))
- Fixed e2e tests [\#157](https://github.com/kubeops/kubed/pull/157) ([diptadas](https://github.com/diptadas))
- Set ClientID for analytics [\#156](https://github.com/kubeops/kubed/pull/156) ([tamalsaha](https://github.com/tamalsaha))
- notifier doc fixes [\#155](https://github.com/kubeops/kubed/pull/155) ([kargakis](https://github.com/kargakis))
- Cleanup object versions [\#153](https://github.com/kubeops/kubed/pull/153) ([tamalsaha](https://github.com/tamalsaha))
- Add front matter for docs 0.3.0 [\#151](https://github.com/kubeops/kubed/pull/151) ([sajibcse68](https://github.com/sajibcse68))
- Add front matter for kubed cli [\#149](https://github.com/kubeops/kubed/pull/149) ([tamalsaha](https://github.com/tamalsaha))
- Revendor dependemcies [\#146](https://github.com/kubeops/kubed/pull/146) ([tamalsaha](https://github.com/tamalsaha))
- Add config file in chart [\#143](https://github.com/kubeops/kubed/pull/143) ([tamalsaha](https://github.com/tamalsaha))
- Use BackupManager from kutil [\#142](https://github.com/kubeops/kubed/pull/142) ([tamalsaha](https://github.com/tamalsaha))
- Avoid listing buckets [\#141](https://github.com/kubeops/kubed/pull/141) ([tamalsaha](https://github.com/tamalsaha))
- Make chart namespaced [\#140](https://github.com/kubeops/kubed/pull/140) ([tamalsaha](https://github.com/tamalsaha))
- Add test event forward [\#138](https://github.com/kubeops/kubed/pull/138) ([a8uhnf](https://github.com/a8uhnf))
- Use client-go 5.x [\#137](https://github.com/kubeops/kubed/pull/137) ([tamalsaha](https://github.com/tamalsaha))
- Add test for Kubed [\#135](https://github.com/kubeops/kubed/pull/135) ([a8uhnf](https://github.com/a8uhnf))
- This should be enableSearchIndex [\#134](https://github.com/kubeops/kubed/pull/134) ([a8uhnf](https://github.com/a8uhnf))

## [0.3.0](https://github.com/kubeops/kubed/tree/0.3.0) (2017-09-26)
[Full Changelog](https://github.com/kubeops/kubed/compare/0.2.0...0.3.0)

**Implemented enhancements:**

- Annotate replicated objects indicating they are a replica and the source [\#112](https://github.com/kubeops/kubed/issues/112)
- Support TLS for elasticsearch connection [\#126](https://github.com/kubeops/kubed/pull/126) ([aerokite](https://github.com/aerokite))

**Fixed bugs:**

- Installing kubed fails due to missing service account [\#121](https://github.com/kubeops/kubed/issues/121)
- Cleanup search index when a namespace is deleted. [\#109](https://github.com/kubeops/kubed/issues/109)

**Closed issues:**

- Vault Integration [\#119](https://github.com/kubeops/kubed/issues/119)
- Notify about new CSR requests [\#73](https://github.com/kubeops/kubed/issues/73)
- Support auth for Elasticsearch janitor [\#64](https://github.com/kubeops/kubed/issues/64)
- Support CRD [\#53](https://github.com/kubeops/kubed/issues/53)

**Merged pull requests:**

- Update docs for 0.3.0 [\#132](https://github.com/kubeops/kubed/pull/132) ([tamalsaha](https://github.com/tamalsaha))
- Prepare docs for 0.3.0 release [\#131](https://github.com/kubeops/kubed/pull/131) ([tamalsaha](https://github.com/tamalsaha))
- Revendor dependencies. [\#130](https://github.com/kubeops/kubed/pull/130) ([tamalsaha](https://github.com/tamalsaha))
- Install kubed as a critical addon [\#129](https://github.com/kubeops/kubed/pull/129) ([tamalsaha](https://github.com/tamalsaha))
- Add changelog [\#128](https://github.com/kubeops/kubed/pull/128) ([tamalsaha](https://github.com/tamalsaha))
- Revendor kutil [\#127](https://github.com/kubeops/kubed/pull/127) ([tamalsaha](https://github.com/tamalsaha))
- Revendor for generator clients. [\#124](https://github.com/kubeops/kubed/pull/124) ([tamalsaha](https://github.com/tamalsaha))
- Update chart to match recent convention [\#123](https://github.com/kubeops/kubed/pull/123) ([tamalsaha](https://github.com/tamalsaha))
- Use correct service account for RBAC installer [\#122](https://github.com/kubeops/kubed/pull/122) ([tamalsaha](https://github.com/tamalsaha))
- Fix command in Developer-guide doc [\#120](https://github.com/kubeops/kubed/pull/120) ([the-redback](https://github.com/the-redback))
- Forward CSR approved/denied events [\#117](https://github.com/kubeops/kubed/pull/117) ([tamalsaha](https://github.com/tamalsaha))
- Use kutil package for utils [\#116](https://github.com/kubeops/kubed/pull/116) ([tamalsaha](https://github.com/tamalsaha))
- Annotate copied configmaps & secrets with kubed.appscode.com/origin [\#115](https://github.com/kubeops/kubed/pull/115) ([tamalsaha](https://github.com/tamalsaha))
- Use client-go 4.0.0 [\#114](https://github.com/kubeops/kubed/pull/114) ([tamalsaha](https://github.com/tamalsaha))
- Fix config object. [\#105](https://github.com/kubeops/kubed/pull/105) ([tamalsaha](https://github.com/tamalsaha))

## [0.2.0](https://github.com/kubeops/kubed/tree/0.2.0) (2017-08-04)
[Full Changelog](https://github.com/kubeops/kubed/compare/0.1.0...0.2.0)

**Closed issues:**

- Use Title case for notifier names [\#99](https://github.com/kubeops/kubed/issues/99)
- Add pushover [\#98](https://github.com/kubeops/kubed/issues/98)
- Indicate cluster name in the CHAT/SMS version [\#94](https://github.com/kubeops/kubed/issues/94)

**Merged pull requests:**

- Prepare docs for 0.2.0 release. [\#104](https://github.com/kubeops/kubed/pull/104) ([tamalsaha](https://github.com/tamalsaha))
- Add support for new ByPush interface. [\#103](https://github.com/kubeops/kubed/pull/103) ([tamalsaha](https://github.com/tamalsaha))
- Add cluster name [\#102](https://github.com/kubeops/kubed/pull/102) ([tamalsaha](https://github.com/tamalsaha))
- Support pushover.net [\#100](https://github.com/kubeops/kubed/pull/100) ([tamalsaha](https://github.com/tamalsaha))
- Add DCO [\#97](https://github.com/kubeops/kubed/pull/97) ([tamalsaha](https://github.com/tamalsaha))
- Use robfig/cron master since we don't need delete feature. [\#96](https://github.com/kubeops/kubed/pull/96) ([tamalsaha](https://github.com/tamalsaha))
- Fix domains to namespaces. [\#95](https://github.com/kubeops/kubed/pull/95) ([tamalsaha](https://github.com/tamalsaha))

## [0.1.0](https://github.com/kubeops/kubed/tree/0.1.0) (2017-08-01)
**Implemented enhancements:**

- Fix API Response Types [\#77](https://github.com/kubeops/kubed/pull/77) ([sadlil](https://github.com/sadlil))
- WIP: Enable Kubed Health API [\#66](https://github.com/kubeops/kubed/pull/66) ([sadlil](https://github.com/sadlil))
- Add Service to ServiceMonitor Reverse index [\#46](https://github.com/kubeops/kubed/pull/46) ([sadlil](https://github.com/sadlil))
- Remove voyager and searchlight controller from Kubed [\#10](https://github.com/kubeops/kubed/pull/10) ([sadlil](https://github.com/sadlil))
- Reverse index for ServiceMonitor to Prometheus [\#47](https://github.com/kubeops/kubed/pull/47) ([sadlil](https://github.com/sadlil))
- Pod to service Reverse index and Full Text search [\#21](https://github.com/kubeops/kubed/pull/21) ([sadlil](https://github.com/sadlil))

**Fixed bugs:**

- Restarting kubed did not fix existing namespaces [\#13](https://github.com/kubeops/kubed/issues/13)
- invalid memory address or nil pointer dereference [\#59](https://github.com/kubeops/kubed/issues/59)
- Assign TypeMeta [\#40](https://github.com/kubeops/kubed/issues/40)
- Make snapshotter storage inline properly [\#84](https://github.com/kubeops/kubed/pull/84) ([tamalsaha](https://github.com/tamalsaha))
- Forwarding events only if recently added [\#67](https://github.com/kubeops/kubed/pull/67) ([tamalsaha](https://github.com/tamalsaha))

**Closed issues:**

- Remove Voyager & Searchlight from Kubed [\#7](https://github.com/kubeops/kubed/issues/7)
- Move prometheus YAML here? [\#5](https://github.com/kubeops/kubed/issues/5)
- Add README.md for promwatcher [\#2](https://github.com/kubeops/kubed/issues/2)
- Local volumes does not work for cluster snapshot [\#52](https://github.com/kubeops/kubed/issues/52)
- Tutorial.md -\> 404 [\#49](https://github.com/kubeops/kubed/issues/49)
- Install as critical addon [\#36](https://github.com/kubeops/kubed/issues/36)
- Sync configmap/secret based on label [\#27](https://github.com/kubeops/kubed/issues/27)
- Use Kubernetes style response objects [\#26](https://github.com/kubeops/kubed/issues/26)
- Support RBAC [\#25](https://github.com/kubeops/kubed/issues/25)
- Send email for Warning events [\#24](https://github.com/kubeops/kubed/issues/24)
- Create full-text search index for Pharm [\#17](https://github.com/kubeops/kubed/issues/17)
- Keep backup of deleted or updated objects [\#16](https://github.com/kubeops/kubed/issues/16)
- Notify cluster admin about soon to be expired certs. [\#15](https://github.com/kubeops/kubed/issues/15)
- Notify cluster admin when some resource is deleted [\#11](https://github.com/kubeops/kubed/issues/11)
- Backup etcd [\#9](https://github.com/kubeops/kubed/issues/9)
- Turn kubed into a reverse index [\#8](https://github.com/kubeops/kubed/issues/8)
- Pass configurations in a secret [\#6](https://github.com/kubeops/kubed/issues/6)

**Merged pull requests:**

- Upload snapshot file in .tar.gz form [\#92](https://github.com/kubeops/kubed/pull/92) ([tamalsaha](https://github.com/tamalsaha))
- Take the first backup using go routine. [\#91](https://github.com/kubeops/kubed/pull/91) ([tamalsaha](https://github.com/tamalsaha))
- Fix Hipchat notifications [\#90](https://github.com/kubeops/kubed/pull/90) ([tamalsaha](https://github.com/tamalsaha))
- Only watch for warning events [\#89](https://github.com/kubeops/kubed/pull/89) ([tamalsaha](https://github.com/tamalsaha))
- Support overwriting old snapshot files. [\#88](https://github.com/kubeops/kubed/pull/88) ([tamalsaha](https://github.com/tamalsaha))
- Test osm credential using `osm lc` [\#87](https://github.com/kubeops/kubed/pull/87) ([tamalsaha](https://github.com/tamalsaha))
- Support multiple receivers for each notification [\#85](https://github.com/kubeops/kubed/pull/85) ([tamalsaha](https://github.com/tamalsaha))
- Fix panic: check reverse index enabled. [\#80](https://github.com/kubeops/kubed/pull/80) ([tamalsaha](https://github.com/tamalsaha))
- Rename kubed-notifier to notifier-info [\#79](https://github.com/kubeops/kubed/pull/79) ([tamalsaha](https://github.com/tamalsaha))
- Update local snapshotter installer scripts. [\#78](https://github.com/kubeops/kubed/pull/78) ([tamalsaha](https://github.com/tamalsaha))
- Show how to use multiple notifiers [\#76](https://github.com/kubeops/kubed/pull/76) ([tamalsaha](https://github.com/tamalsaha))
- Document config [\#74](https://github.com/kubeops/kubed/pull/74) ([tamalsaha](https://github.com/tamalsaha))
- User docs - part 15 [\#71](https://github.com/kubeops/kubed/pull/71) ([tamalsaha](https://github.com/tamalsaha))
- Obfuscate secrets in index and recycle bin [\#69](https://github.com/kubeops/kubed/pull/69) ([tamalsaha](https://github.com/tamalsaha))
- Update apiServer config [\#65](https://github.com/kubeops/kubed/pull/65) ([tamalsaha](https://github.com/tamalsaha))
- Document janitors [\#62](https://github.com/kubeops/kubed/pull/62) ([tamalsaha](https://github.com/tamalsaha))
- Update event forwarder docs [\#61](https://github.com/kubeops/kubed/pull/61) ([tamalsaha](https://github.com/tamalsaha))
- Document event forwarder [\#58](https://github.com/kubeops/kubed/pull/58) ([tamalsaha](https://github.com/tamalsaha))
- User docs - recycle bin [\#57](https://github.com/kubeops/kubed/pull/57) ([tamalsaha](https://github.com/tamalsaha))
- Use docs - part 2 [\#56](https://github.com/kubeops/kubed/pull/56) ([tamalsaha](https://github.com/tamalsaha))
- User Docs - part 1 [\#50](https://github.com/kubeops/kubed/pull/50) ([tamalsaha](https://github.com/tamalsaha))
- Require config to pass notification receiver address [\#48](https://github.com/kubeops/kubed/pull/48) ([tamalsaha](https://github.com/tamalsaha))
- Cleanup Reverse Index [\#44](https://github.com/kubeops/kubed/pull/44) ([tamalsaha](https://github.com/tamalsaha))
- Assign TypeKind [\#43](https://github.com/kubeops/kubed/pull/43) ([tamalsaha](https://github.com/tamalsaha))
- Add event forwarder. [\#38](https://github.com/kubeops/kubed/pull/38) ([tamalsaha](https://github.com/tamalsaha))
- Index resources for searching [\#37](https://github.com/kubeops/kubed/pull/37) ([tamalsaha](https://github.com/tamalsaha))
- Refine cluster config [\#34](https://github.com/kubeops/kubed/pull/34) ([tamalsaha](https://github.com/tamalsaha))
- Watch everything [\#31](https://github.com/kubeops/kubed/pull/31) ([tamalsaha](https://github.com/tamalsaha))
- Add docs from stash [\#29](https://github.com/kubeops/kubed/pull/29) ([tamalsaha](https://github.com/tamalsaha))
- Generate reference docs [\#28](https://github.com/kubeops/kubed/pull/28) ([tamalsaha](https://github.com/tamalsaha))
- Notify admin exp certs [\#23](https://github.com/kubeops/kubed/pull/23) ([ashiquzzaman33](https://github.com/ashiquzzaman33))
- Pass configurations in a secret [\#20](https://github.com/kubeops/kubed/pull/20) ([ashiquzzaman33](https://github.com/ashiquzzaman33))
- Remove provider name flag [\#12](https://github.com/kubeops/kubed/pull/12) ([tamalsaha](https://github.com/tamalsaha))
- Add documentation for Prometheus Watcher [\#4](https://github.com/kubeops/kubed/pull/4) ([aerokite](https://github.com/aerokite))
- Change package to kubeops.dev/kubed [\#3](https://github.com/kubeops/kubed/pull/3) ([tamalsaha](https://github.com/tamalsaha))
- Add Prometheus TPR watcher [\#1](https://github.com/kubeops/kubed/pull/1) ([aerokite](https://github.com/aerokite))
- Use docs - part 13 [\#68](https://github.com/kubeops/kubed/pull/68) ([tamalsaha](https://github.com/tamalsaha))
- Turn janitors into an array [\#55](https://github.com/kubeops/kubed/pull/55) ([tamalsaha](https://github.com/tamalsaha))
- Add kubed check command to verify cluster config [\#54](https://github.com/kubeops/kubed/pull/54) ([tamalsaha](https://github.com/tamalsaha))
- Update Elastic to Elasticsearch [\#45](https://github.com/kubeops/kubed/pull/45) ([tamalsaha](https://github.com/tamalsaha))
- Sync configmap & secret with annotation kubernetes.appscode.com/sync [\#42](https://github.com/kubeops/kubed/pull/42) ([tamalsaha](https://github.com/tamalsaha))
- Various bug fixes [\#39](https://github.com/kubeops/kubed/pull/39) ([tamalsaha](https://github.com/tamalsaha))
- Organize operator [\#35](https://github.com/kubeops/kubed/pull/35) ([tamalsaha](https://github.com/tamalsaha))
- Cleanup config format [\#33](https://github.com/kubeops/kubed/pull/33) ([tamalsaha](https://github.com/tamalsaha))
- Allow recovering deleted Kube objects [\#32](https://github.com/kubeops/kubed/pull/32) ([tamalsaha](https://github.com/tamalsaha))
- Add cluster backup command from appctl [\#30](https://github.com/kubeops/kubed/pull/30) ([tamalsaha](https://github.com/tamalsaha))
- Use client-go [\#18](https://github.com/kubeops/kubed/pull/18) ([tamalsaha](https://github.com/tamalsaha))



\* *This Change Log was automatically generated by [github_changelog_generator](https://github.com/skywinder/Github-Changelog-Generator)*