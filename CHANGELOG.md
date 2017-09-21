# Change Log

## [Unreleased](https://github.com/appscode/kubed/tree/HEAD)

[Full Changelog](https://github.com/appscode/kubed/compare/4.0.0-alpha.0...HEAD)

**Implemented enhancements:**

- Support TLS for elasticsearch connection [\#126](https://github.com/appscode/kubed/pull/126) ([aerokite](https://github.com/aerokite))

**Fixed bugs:**

- Installing kubed fails due to missing service account [\#121](https://github.com/appscode/kubed/issues/121)
- Cleanup search index when a namespace is deleted. [\#109](https://github.com/appscode/kubed/issues/109)

**Closed issues:**

- Vault Integration [\#119](https://github.com/appscode/kubed/issues/119)
- Support auth for Elasticsearch janitor [\#64](https://github.com/appscode/kubed/issues/64)

**Merged pull requests:**

- Revendor for generator clients. [\#124](https://github.com/appscode/kubed/pull/124) ([tamalsaha](https://github.com/tamalsaha))
- Update chart to match recent convention [\#123](https://github.com/appscode/kubed/pull/123) ([tamalsaha](https://github.com/tamalsaha))
- Use correct service account for RBAC installer [\#122](https://github.com/appscode/kubed/pull/122) ([tamalsaha](https://github.com/tamalsaha))
- Fix command in Developer-guide doc [\#120](https://github.com/appscode/kubed/pull/120) ([the-redback](https://github.com/the-redback))

## [4.0.0-alpha.0](https://github.com/appscode/kubed/tree/4.0.0-alpha.0) (2017-09-05)
[Full Changelog](https://github.com/appscode/kubed/compare/0.2.0...4.0.0-alpha.0)

**Implemented enhancements:**

- Annotate replicated objects indicating they are a replica and the source [\#112](https://github.com/appscode/kubed/issues/112)

**Closed issues:**

- Notify about new CSR requests [\#73](https://github.com/appscode/kubed/issues/73)
- Support CRD [\#53](https://github.com/appscode/kubed/issues/53)

**Merged pull requests:**

- Forward CSR approved/denied events [\#117](https://github.com/appscode/kubed/pull/117) ([tamalsaha](https://github.com/tamalsaha))
- Use kutil package for utils [\#116](https://github.com/appscode/kubed/pull/116) ([tamalsaha](https://github.com/tamalsaha))
- Annotate copied configmaps & secrets with kubed.appscode.com/origin [\#115](https://github.com/appscode/kubed/pull/115) ([tamalsaha](https://github.com/tamalsaha))
- Use client-go 4.0.0 [\#114](https://github.com/appscode/kubed/pull/114) ([tamalsaha](https://github.com/tamalsaha))
- Fix config object. [\#105](https://github.com/appscode/kubed/pull/105) ([tamalsaha](https://github.com/tamalsaha))

## [0.2.0](https://github.com/appscode/kubed/tree/0.2.0) (2017-08-04)
[Full Changelog](https://github.com/appscode/kubed/compare/0.1.0...0.2.0)

**Closed issues:**

- Use Title case for notifier names [\#99](https://github.com/appscode/kubed/issues/99)
- Add pushover [\#98](https://github.com/appscode/kubed/issues/98)
- Indicate cluster name in the CHAT/SMS version [\#94](https://github.com/appscode/kubed/issues/94)

**Merged pull requests:**

- Prepare docs for 0.2.0 release. [\#104](https://github.com/appscode/kubed/pull/104) ([tamalsaha](https://github.com/tamalsaha))
- Add support for new ByPush interface. [\#103](https://github.com/appscode/kubed/pull/103) ([tamalsaha](https://github.com/tamalsaha))
- Add cluster name [\#102](https://github.com/appscode/kubed/pull/102) ([tamalsaha](https://github.com/tamalsaha))
- Support pushover.net [\#100](https://github.com/appscode/kubed/pull/100) ([tamalsaha](https://github.com/tamalsaha))
- Add DCO [\#97](https://github.com/appscode/kubed/pull/97) ([tamalsaha](https://github.com/tamalsaha))
- Use robfig/cron master since we don't need delete feature. [\#96](https://github.com/appscode/kubed/pull/96) ([tamalsaha](https://github.com/tamalsaha))
- Fix domains to namespaces. [\#95](https://github.com/appscode/kubed/pull/95) ([tamalsaha](https://github.com/tamalsaha))

## [0.1.0](https://github.com/appscode/kubed/tree/0.1.0) (2017-08-01)
**Implemented enhancements:**

- Fix API Response Types [\#77](https://github.com/appscode/kubed/pull/77) ([sadlil](https://github.com/sadlil))
- WIP: Enable Kubed Health API [\#66](https://github.com/appscode/kubed/pull/66) ([sadlil](https://github.com/sadlil))
- Add Service to ServiceMonitor Reverse index [\#46](https://github.com/appscode/kubed/pull/46) ([sadlil](https://github.com/sadlil))
- Remove voyager and searchlight controller from Kubed [\#10](https://github.com/appscode/kubed/pull/10) ([sadlil](https://github.com/sadlil))
- Reverse index for ServiceMonitor to Prometheus [\#47](https://github.com/appscode/kubed/pull/47) ([sadlil](https://github.com/sadlil))
- Pod to service Reverse index and Full Text search [\#21](https://github.com/appscode/kubed/pull/21) ([sadlil](https://github.com/sadlil))

**Fixed bugs:**

- Restarting kubed did not fix existing namespaces [\#13](https://github.com/appscode/kubed/issues/13)
- invalid memory address or nil pointer dereference [\#59](https://github.com/appscode/kubed/issues/59)
- Assign TypeMeta [\#40](https://github.com/appscode/kubed/issues/40)
- Make snapshotter storage inline properly [\#84](https://github.com/appscode/kubed/pull/84) ([tamalsaha](https://github.com/tamalsaha))
- Forwarding events only if recently added [\#67](https://github.com/appscode/kubed/pull/67) ([tamalsaha](https://github.com/tamalsaha))

**Closed issues:**

- Remove Voyager & Searchlight from Kubed [\#7](https://github.com/appscode/kubed/issues/7)
- Move prometheus YAML here? [\#5](https://github.com/appscode/kubed/issues/5)
- Add README.md for promwatcher [\#2](https://github.com/appscode/kubed/issues/2)
- Local volumes does not work for cluster snapshot [\#52](https://github.com/appscode/kubed/issues/52)
- Tutorial.md -\> 404 [\#49](https://github.com/appscode/kubed/issues/49)
- Install as critical addon [\#36](https://github.com/appscode/kubed/issues/36)
- Sync configmap/secret based on label [\#27](https://github.com/appscode/kubed/issues/27)
- Use Kubernetes style response objects [\#26](https://github.com/appscode/kubed/issues/26)
- Support RBAC [\#25](https://github.com/appscode/kubed/issues/25)
- Send email for Warning events [\#24](https://github.com/appscode/kubed/issues/24)
- Create full-text search index for Pharm [\#17](https://github.com/appscode/kubed/issues/17)
- Keep backup of deleted or updated objects [\#16](https://github.com/appscode/kubed/issues/16)
- Notify cluster admin about soon to be expired certs. [\#15](https://github.com/appscode/kubed/issues/15)
- Notify cluster admin when some resource is deleted [\#11](https://github.com/appscode/kubed/issues/11)
- Backup etcd [\#9](https://github.com/appscode/kubed/issues/9)
- Turn kubed into a reverse index [\#8](https://github.com/appscode/kubed/issues/8)
- Pass configurations in a secret [\#6](https://github.com/appscode/kubed/issues/6)

**Merged pull requests:**

- Upload snapshot file in .tar.gz form [\#92](https://github.com/appscode/kubed/pull/92) ([tamalsaha](https://github.com/tamalsaha))
- Take the first backup using go routine. [\#91](https://github.com/appscode/kubed/pull/91) ([tamalsaha](https://github.com/tamalsaha))
- Fix Hipchat notifications [\#90](https://github.com/appscode/kubed/pull/90) ([tamalsaha](https://github.com/tamalsaha))
- Only watch for warning events [\#89](https://github.com/appscode/kubed/pull/89) ([tamalsaha](https://github.com/tamalsaha))
- Support overwriting old snapshot files. [\#88](https://github.com/appscode/kubed/pull/88) ([tamalsaha](https://github.com/tamalsaha))
- Test osm credential using `osm lc` [\#87](https://github.com/appscode/kubed/pull/87) ([tamalsaha](https://github.com/tamalsaha))
- Support multiple receivers for each notification [\#85](https://github.com/appscode/kubed/pull/85) ([tamalsaha](https://github.com/tamalsaha))
- Fix panic: check reverse index enabled. [\#80](https://github.com/appscode/kubed/pull/80) ([tamalsaha](https://github.com/tamalsaha))
- Rename kubed-notifier to notifier-info [\#79](https://github.com/appscode/kubed/pull/79) ([tamalsaha](https://github.com/tamalsaha))
- Update local snapshotter installer scripts. [\#78](https://github.com/appscode/kubed/pull/78) ([tamalsaha](https://github.com/tamalsaha))
- Show how to use multiple notifiers [\#76](https://github.com/appscode/kubed/pull/76) ([tamalsaha](https://github.com/tamalsaha))
- Document config [\#74](https://github.com/appscode/kubed/pull/74) ([tamalsaha](https://github.com/tamalsaha))
- User docs - part 15 [\#71](https://github.com/appscode/kubed/pull/71) ([tamalsaha](https://github.com/tamalsaha))
- Obfuscate secrets in index and recycle bin [\#69](https://github.com/appscode/kubed/pull/69) ([tamalsaha](https://github.com/tamalsaha))
- Update apiServer config [\#65](https://github.com/appscode/kubed/pull/65) ([tamalsaha](https://github.com/tamalsaha))
- Document janitors [\#62](https://github.com/appscode/kubed/pull/62) ([tamalsaha](https://github.com/tamalsaha))
- Update event forwarder docs [\#61](https://github.com/appscode/kubed/pull/61) ([tamalsaha](https://github.com/tamalsaha))
- Document event forwarder [\#58](https://github.com/appscode/kubed/pull/58) ([tamalsaha](https://github.com/tamalsaha))
- User docs - recycle bin [\#57](https://github.com/appscode/kubed/pull/57) ([tamalsaha](https://github.com/tamalsaha))
- Use docs - part 2 [\#56](https://github.com/appscode/kubed/pull/56) ([tamalsaha](https://github.com/tamalsaha))
- User Docs - part 1 [\#50](https://github.com/appscode/kubed/pull/50) ([tamalsaha](https://github.com/tamalsaha))
- Require config to pass notification receiver address [\#48](https://github.com/appscode/kubed/pull/48) ([tamalsaha](https://github.com/tamalsaha))
- Cleanup Reverse Index [\#44](https://github.com/appscode/kubed/pull/44) ([tamalsaha](https://github.com/tamalsaha))
- Assign TypeKind [\#43](https://github.com/appscode/kubed/pull/43) ([tamalsaha](https://github.com/tamalsaha))
- Add event forwarder. [\#38](https://github.com/appscode/kubed/pull/38) ([tamalsaha](https://github.com/tamalsaha))
- Index resources for searching [\#37](https://github.com/appscode/kubed/pull/37) ([tamalsaha](https://github.com/tamalsaha))
- Refine cluster config [\#34](https://github.com/appscode/kubed/pull/34) ([tamalsaha](https://github.com/tamalsaha))
- Watch everything [\#31](https://github.com/appscode/kubed/pull/31) ([tamalsaha](https://github.com/tamalsaha))
- Add docs from stash [\#29](https://github.com/appscode/kubed/pull/29) ([tamalsaha](https://github.com/tamalsaha))
- Generate reference docs [\#28](https://github.com/appscode/kubed/pull/28) ([tamalsaha](https://github.com/tamalsaha))
- Notify admin exp certs [\#23](https://github.com/appscode/kubed/pull/23) ([ashiquzzaman33](https://github.com/ashiquzzaman33))
- Pass configurations in a secret [\#20](https://github.com/appscode/kubed/pull/20) ([ashiquzzaman33](https://github.com/ashiquzzaman33))
- Remove provider name flag [\#12](https://github.com/appscode/kubed/pull/12) ([tamalsaha](https://github.com/tamalsaha))
- Add documentation for Prometheus Watcher [\#4](https://github.com/appscode/kubed/pull/4) ([aerokite](https://github.com/aerokite))
- Change package to github.com/appscode/kubed [\#3](https://github.com/appscode/kubed/pull/3) ([tamalsaha](https://github.com/tamalsaha))
- Add Prometheus TPR watcher [\#1](https://github.com/appscode/kubed/pull/1) ([aerokite](https://github.com/aerokite))
- Use docs - part 13 [\#68](https://github.com/appscode/kubed/pull/68) ([tamalsaha](https://github.com/tamalsaha))
- Turn janitors into an array [\#55](https://github.com/appscode/kubed/pull/55) ([tamalsaha](https://github.com/tamalsaha))
- Add kubed check command to verify cluster config [\#54](https://github.com/appscode/kubed/pull/54) ([tamalsaha](https://github.com/tamalsaha))
- Update Elastic to Elasticsearch [\#45](https://github.com/appscode/kubed/pull/45) ([tamalsaha](https://github.com/tamalsaha))
- Sync configmap & secret with annotation kubernetes.appscode.com/sync [\#42](https://github.com/appscode/kubed/pull/42) ([tamalsaha](https://github.com/tamalsaha))
- Various bug fixes [\#39](https://github.com/appscode/kubed/pull/39) ([tamalsaha](https://github.com/tamalsaha))
- Organize operator [\#35](https://github.com/appscode/kubed/pull/35) ([tamalsaha](https://github.com/tamalsaha))
- Cleanup config format [\#33](https://github.com/appscode/kubed/pull/33) ([tamalsaha](https://github.com/tamalsaha))
- Allow recovering deleted Kube objects [\#32](https://github.com/appscode/kubed/pull/32) ([tamalsaha](https://github.com/tamalsaha))
- Add cluster backup command from appctl [\#30](https://github.com/appscode/kubed/pull/30) ([tamalsaha](https://github.com/tamalsaha))
- Use client-go [\#18](https://github.com/appscode/kubed/pull/18) ([tamalsaha](https://github.com/tamalsaha))



\* *This Change Log was automatically generated by [github_changelog_generator](https://github.com/skywinder/Github-Changelog-Generator)*