---
title: Notifiers
description: Notifiers
menu:
  product_kubed_0.3.1:
    identifier: tutorials-notifiers
    name: Notifiers
    parent: tutorials
    weight: 20
product_name: kubed
menu_name: product_kubed_0.3.1
section_menu_id: tutorials
---

> New to Kubed? Please start [here](/docs/guides/README.md).

# Supported Notifiers
Kubed can send notifications via Email, SMS or Chat for various operations using [appscode/go-notify](https://github.com/appscode/go-notify) library. To connect to these services, you need to create a Secret with the appropriate keys. Then pass the secret name to Kubed by setting `notifierSecretName` field in Kubed cluster config. Also, set `clusterName` to a meaningful name for you. This cluster name will be prefixed to any notification sent via Email/SMS/Chat so that you can identify the source easily.

## Hipchat
To receive chat notifications in Hipchat, create a Secret with the following key:

| Name                         | Description                               |
|------------------------------|-------------------------------------------|
| HIPCHAT_AUTH_TOKEN           | `Required` Hipchat [api access token](https://developer.atlassian.com/hipchat/guide/hipchat-rest-api/api-access-tokens). You can use room notification tokens, if you are planning to send notifications to a single room.   |
| HIPCHAT_BASE_URL             | `Optional` Base url of Hipchat server                        |
| HIPCHAT_CA_CERT_DATA         | `Optional` PEM encoded CA certificate used by Hipchat server |
| HIPCHAT_INSECURE_SKIP_VERIFY | `Optional` If set to `true`, skips SSL verification          |

```console
$ echo -n 'your-hipchat-auth-token' > HIPCHAT_AUTH_TOKEN
$ kubectl create secret generic notifier-config -n kube-system \
    --from-file=./HIPCHAT_AUTH_TOKEN
secret "notifier-config" created
```
```yaml
apiVersion: v1
data:
  HIPCHAT_AUTH_TOKEN: eW91ci1oaXBjaGF0LWF1dGgtdG9rZW4=
kind: Secret
metadata:
  creationTimestamp: 2017-07-25T01:54:37Z
  name: notifier-config
  namespace: kube-system
  resourceVersion: "2244"
  selfLink: /api/v1/namespaces/kube-system/secrets/notifier-config
  uid: 372bc159-70dc-11e7-9b0b-080027503732
type: Opaque
```

Now, to receive notifications via Hipchat, configure receiver as below:
 - notifier: `Hipchat`
 - to: a list of chat room names

```yaml
clusterName: unicorn
recycleBin:
  handleUpdates: false
  path: /tmp/kubed/trash
  receivers:
  - notifier: Hipchat
    to:
    - ops-alerts
  ttl: 168h
notifierSecretName: notifier-config
```


## Mailgun
To receive email notifications via Mailgun, create a Secret with the following keys:

| Name                    | Description                                  |
|-------------------------|----------------------------------------------|
| MAILGUN_DOMAIN          | `Required` domain name for Mailgun account   |
| MAILGUN_API_KEY         | `Required` Mailgun API Key                   |
| MAILGUN_FROM            | `Required` sender email address              |
| MAILGUN_PUBLIC_API_KEY  | `Optional` Mailgun public API Key            |

```console
$ echo -n 'your-mailgun-domain' > MAILGUN_DOMAIN
$ echo -n 'no-reply@example.com' > MAILGUN_FROM
$ echo -n 'your-mailgun-api-key' > MAILGUN_API_KEY
$ echo -n 'your-mailgun-public-api-key' > MAILGUN_PUBLIC_API_KEY
$ kubectl create secret generic notifier-config -n kube-system \
    --from-file=./MAILGUN_DOMAIN \
    --from-file=./MAILGUN_FROM \
    --from-file=./MAILGUN_API_KEY \
    --from-file=./MAILGUN_PUBLIC_API_KEY
secret "notifier-config" created
```
```yaml
apiVersion: v1
data:
  MAILGUN_API_KEY: eW91ci1tYWlsZ3VuLWFwaS1rZXk=
  MAILGUN_DOMAIN: eW91ci1tYWlsZ3VuLWRvbWFpbg==
  MAILGUN_FROM: bm8tcmVwbHlAZXhhbXBsZS5jb20=
  MAILGUN_PUBLIC_API_KEY: bWFpbGd1bi1wdWJsaWMtYXBpLWtleQ==
kind: Secret
metadata:
  creationTimestamp: 2017-07-25T01:31:24Z
  name: notifier-config
  namespace: kube-system
  resourceVersion: "714"
  selfLink: /api/v1/namespaces/kube-system/secrets/notifier-config
  uid: f8e91037-70d8-11e7-9b0b-080027503732
type: Opaque
```

Now, to receive notifications via Mailgun, configure receiver as below:
 - notifier: `Mailgun`
 - to: a list of email addresses

```yaml
clusterName: unicorn
recycleBin:
  handleUpdates: false
  path: /tmp/kubed/trash
  receivers:
  - notifier: Mailgun
    to:
    - ops-alerts@example.com
  ttl: 168h
notifierSecretName: notifier-config
```


## SMTP
To configure any email provider, use SMTP notifier. To receive email notifications via SMTP services, create a Secret with the following keys:

| Name                      | Description                                           |
|---------------------------|-------------------------------------------------------|
| SMTP_HOST                 | `Required` Host address of smtp server                |
| SMTP_PORT                 | `Required` Port of smtp server                        |
| SMTP_INSECURE_SKIP_VERIFY | `Required` If set to `true`, skips SSL verification   |
| SMTP_USERNAME             | `Required` Username                                   |
| SMTP_PASSWORD             | `Required` Password                                   |
| SMTP_FROM                 | `Required` Sender email address                       |


```console
$ echo -n 'your-smtp-host' > SMTP_HOST
$ echo -n 'your-smtp-port' > SMTP_PORT
$ echo -n 'your-smtp-insecure-skip-verify' > SMTP_INSECURE_SKIP_VERIFY
$ echo -n 'your-smtp-username' > SMTP_USERNAME
$ echo -n 'your-smtp-password' > SMTP_PASSWORD
$ echo -n 'your-smtp-from' > SMTP_FROM
$ kubectl create secret generic notifier-config -n kube-system \
    --from-file=./SMTP_HOST \
    --from-file=./SMTP_PORT \
    --from-file=./SMTP_INSECURE_SKIP_VERIFY \
    --from-file=./SMTP_USERNAME \
    --from-file=./SMTP_PASSWORD \
    --from-file=./SMTP_FROM
secret "notifier-config" created
```

To configure Kubed to send email notifications using a GMail account, set the Secrets like below:
```
$ echo -n 'smtp.gmail.com' > SMTP_HOST
$ echo -n '587' > SMTP_PORT
$ echo -n 'your-gmail-adress' > SMTP_USERNAME
$ echo -n 'your-gmail-password' > SMTP_PASSWORD
$ echo -n 'your-gmail-address' > SMTP_FROM
```

Now, to receive notifications via SMTP, configure receiver as below:
 - notifier: `SMTP`
 - to: a list of email addresses

```yaml
clusterName: unicorn
recycleBin:
  handleUpdates: false
  path: /tmp/kubed/trash
  receivers:
  - notifier: SMTP
    to:
    - ops-alerts@example.com
  ttl: 168h
notifierSecretName: notifier-config
```


## Twilio
To receive SMS notifications via Twilio, create a Secret with the following keys:

| Name                | Description                                        |
|---------------------|----------------------------------------------------|
| TWILIO_ACCOUNT_SID  | `Required` Twilio account SID                      |
| TWILIO_AUTH_TOKEN   | `Required` Twilio authentication token             |
| TWILIO_FROM         | `Required` Sender mobile number                    |

```console
$ echo -n 'your-twilio-account-sid' > TWILIO_ACCOUNT_SID
$ echo -n 'your-twilio-auth-token' > TWILIO_AUTH_TOKEN
$ echo -n 'your-twilio-from' > TWILIO_FROM
$ kubectl create secret generic notifier-config -n kube-system \
    --from-file=./TWILIO_ACCOUNT_SID \
    --from-file=./TWILIO_AUTH_TOKEN \
    --from-file=./TWILIO_FROM
secret "notifier-config" created
```
```yaml
apiVersion: v1
data:
  TWILIO_ACCOUNT_SID: eW91ci10d2lsaW8tYWNjb3VudC1zaWQ=
  TWILIO_AUTH_TOKEN: eW91ci10d2lsaW8tYXV0aC10b2tlbg==
  TWILIO_FROM: eW91ci10d2lsaW8tZnJvbQ==
kind: Secret
metadata:
  creationTimestamp: 2017-07-26T17:38:38Z
  name: notifier-config
  namespace: kube-system
  resourceVersion: "27787"
  selfLink: /api/v1/namespaces/kube-system/secrets/notifier-config
  uid: 41f57a61-7229-11e7-af79-08002738e55e
type: Opaque
```

Now, to receive notifications via Twilio, configure receiver as below:
 - notifier: `Twilio`
 - to: a list of receiver mobile numbers

```yaml
clusterName: unicorn
recycleBin:
  handleUpdates: false
  path: /tmp/kubed/trash
  receivers:
  - notifier: Twilio
    to:
    - +1-999-888-1234
  ttl: 168h
notifierSecretName: notifier-config
```


## Slack
To receive chat notifications in Slack, create a Secret with the following keys:

| Name             | Description                      |
|------------------|----------------------------------|
| SLACK_AUTH_TOKEN | `Required` Slack [legacy auth token](https://api.slack.com/custom-integrations/legacy-tokens) |

```console
$ echo -n 'your-slack-auth-token' > SLACK_AUTH_TOKEN
$ kubectl create secret generic notifier-config -n kube-system \
    --from-file=./SLACK_AUTH_TOKEN
secret "notifier-config" created
```
```yaml
apiVersion: v1
data:
  SLACK_AUTH_TOKEN: eW91ci1zbGFjay1hdXRoLXRva2Vu
kind: Secret
metadata:
  creationTimestamp: 2017-07-25T01:58:58Z
  name: notifier-config
  namespace: kube-system
  resourceVersion: "2534"
  selfLink: /api/v1/namespaces/kube-system/secrets/notifier-config
  uid: d2571817-70dc-11e7-9b0b-080027503732
type: Opaque
```

Now, to receive notifications via Slack, configure receiver as below:
 - notifier: `Slack`
 - to: a list of chat room names

```yaml
clusterName: unicorn
recycleBin:
  handleUpdates: false
  path: /tmp/kubed/trash
  receivers:
  - notifier: Slack
    to:
    - '#ops-alerts'
  ttl: 168h
notifierSecretName: notifier-config
```


## Plivo
To receive SMS notifications via Plivo, create a Secret with the following keys:

| Name              | Description                             |
|-------------------|-----------------------------------------|
| PLIVO_AUTH_ID     | `Required` Plivo auth ID                |
| PLIVO_AUTH_TOKEN  | `Required` Plivo authentication token   |
| PLIVO_FROM        | `Required` Sender mobile number         |

```console
$ echo -n 'your-plivo-auth-id' > PLIVO_AUTH_ID
$ echo -n 'your-plivo-auth-token' > PLIVO_AUTH_TOKEN
$ echo -n 'your-plivo-from' > PLIVO_FROM
$ kubectl create secret generic notifier-config -n kube-system \
    --from-file=./PLIVO_AUTH_ID \
    --from-file=./PLIVO_AUTH_TOKEN \
    --from-file=./PLIVO_FROM
secret "notifier-config" created
```
```yaml
apiVersion: v1
data:
  PLIVO_AUTH_ID: eW91ci1wbGl2by1hdXRoLWlk
  PLIVO_AUTH_TOKEN: eW91ci1wbGl2by1hdXRoLXRva2Vu
  PLIVO_FROM: eW91ci1wbGl2by1mcm9t
kind: Secret
metadata:
  creationTimestamp: 2017-07-25T02:00:02Z
  name: notifier-config
  namespace: kube-system
  resourceVersion: "2606"
  selfLink: /api/v1/namespaces/kube-system/secrets/notifier-config
  uid: f8dade1c-70dc-11e7-9b0b-080027503732
type: Opaque
```

Now, to receive notifications via Plivo, configure receiver as below:
 - notifier: `Plivo`
 - to: a list of receiver mobile numbers

```yaml
clusterName: unicorn
recycleBin:
  handleUpdates: false
  path: /tmp/kubed/trash
  receivers:
  - notifier: Plivo
    to:
    - +1-999-888-1234
  ttl: 168h
notifierSecretName: notifier-config
```


## Pushover.net
To receive push notifications via Pushover.net, create a Secret with the following keys:

| Name               | Description                                                                    |
|--------------------|--------------------------------------------------------------------------------|
| PUSHOVER_TOKEN     | `Required` Pushover.net token.                                                 |
| PUSHOVER_USER_KEY  | `Required` User key or group key.                                              |
| PUSHOVER_TITLE     | `Optional` Message's title, otherwise your app's name is used.                 |
| PUSHOVER_URL       | `Optional` A supplementary URL to show with your message.                      |
| PUSHOVER_URL_TITLE | `Optional` A title for the supplementary URL, otherwise just the URL is shown. |
| PUSHOVER_PRIORITY  | `Optional` Send as -2 to generate no notification/alert, -1 to always send as a quiet notification, 1 to display as high-priority and bypass the user's quiet hours, or 2 to also require confirmation from the user.   |
| PUSHOVER_SOUND     | `Optional` The name of one of the sounds supported by device clients to override the user's default sound choice.   |


```console
$ echo -n 'your-pushover-token' > PUSHOVER_TOKEN
$ echo -n 'your-pushover-user-key' > PUSHOVER_USER_KEY
$ echo -n 'your-pushover-title' > PUSHOVER_TITLE
$ echo -n 'your-pushover-url' > PUSHOVER_URL
$ echo -n 'your-pushover-url-title' > PUSHOVER_URL_TITLE
$ echo -n 'your-pushover-priority' > PUSHOVER_PRIORITY
$ echo -n 'your-pushover-sound' > PUSHOVER_SOUND
$ kubectl create secret generic notifier-config -n kube-system \
    --from-file=./PUSHOVER_TOKEN \
    --from-file=./PUSHOVER_USER_KEY \
    --from-file=./PUSHOVER_TITLE \
    --from-file=./PUSHOVER_URL \
    --from-file=./PUSHOVER_URL_TITLE \
    --from-file=./PUSHOVER_PRIORITY \
    --from-file=./PUSHOVER_SOUND
secret "notifier-config" created
```
```yaml
apiVersion: v1
data:
  PUSHOVER_PRIORITY: eW91ci1wdXNob3Zlci1wcmlvcml0eQ==
  PUSHOVER_SOUND: eW91ci1wdXNob3Zlci1zb3VuZA==
  PUSHOVER_TITLE: eW91ci1wdXNob3Zlci10aXRsZQ==
  PUSHOVER_TOKEN: eW91ci1wdXNob3Zlci10b2tlbg==
  PUSHOVER_URL: eW91ci1wdXNob3Zlci11cmw=
  PUSHOVER_URL_TITLE: eW91ci1wdXNob3Zlci11cmwtdGl0bGU=
  PUSHOVER_USER_KEY: eW91ci1wdXNob3Zlci11c2VyLWtleQ==
kind: Secret
metadata:
  creationTimestamp: 2017-08-04T05:13:07Z
  name: notifier-config
  namespace: kube-system
  resourceVersion: "33711872"
  selfLink: /api/v1/namespaces/kube-system/secrets/notifier-config
  uid: 99df75a8-78d3-11e7-acfa-42010af00141
type: Opaque
```

Now, to receive notifications via Pushover.net, configure receiver as below:
 - notifier: `Pushover`
 - to: a list of devices where notifications will be sent. If list is empty, all devices will be notified.

```yaml
clusterName: unicorn
recycleBin:
  handleUpdates: false
  path: /tmp/kubed/trash
  receivers:
  - notifier: Pushover
    to:
    - my-phone
  ttl: 168h
notifierSecretName: notifier-config
```


## Using multiple notifiers
Kubed supports using different notifiers in different scenarios. First add the credentials for the different notifiers in the same Secret `notifier-config` and deploy that to Kubernetes. Then in the Kubed cluster config, specify the appropriate notifier for each feature.

```yaml
apiVersion: v1
data:
  MAILGUN_API_KEY: eW91ci1tYWlsZ3VuLWFwaS1rZXk=
  MAILGUN_DOMAIN: eW91ci1tYWlsZ3VuLWRvbWFpbg==
  MAILGUN_FROM: bm8tcmVwbHlAZXhhbXBsZS5jb20=
  SLACK_AUTH_TOKEN: eW91ci1zbGFjay1hdXRoLXRva2Vu
kind: Secret
metadata:
  creationTimestamp: 2017-07-25T01:58:58Z
  name: notifier-config
  namespace: kube-system
  resourceVersion: "2534"
  selfLink: /api/v1/namespaces/kube-system/secrets/notifier-config
  uid: d2571817-70dc-11e7-9b0b-080027503732
type: Opaque


clusterName: unicorn
eventForwarder:
  warningEvents:
    handle: true
    namespaces:
    - kube-system
  receivers:
  - notifier: Mailgun
    to:
    - ops@example.com
  - notifier: Slack
    to:
    - #ops-alerts
recycleBin:
  handleUpdates: false
  path: /tmp/kubed/trash
  receivers:
  - notifier: mailgun
    to:
    - ops@example.com
  - notifier: slack
    to:
    - '#ops-alerts'
  ttl: 168h
notifierSecretName: notifier-config
```

## Next Steps
 - Learn how to use Kubed to take periodic snapshots of a Kubernetes cluster [here](/docs/guides/cluster-snapshot.md).
 - To setup a recycle bin for deleted and/or updated Kubernetes objects, please visit [here](/docs/guides/recycle-bin.md).
 - Need to keep some configuration synchronized across namespaces? Try [Kubed config syncer](/docs/guides/config-syncer.md).
 - Want to keep an eye on your cluster with automated notifications? Setup Kubed [event forwarder](/docs/guides/event-forwarder.md).
 - Out of disk space because of too much logs in Elasticsearch or metrics in InfluxDB? Configure [janitors](/docs/guides/janitors.md) to delete old data.
 - Wondering what features are coming next? Please visit [here](/ROADMAP.md).
 - Want to hack on Kubed? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
