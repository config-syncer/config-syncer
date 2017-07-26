> New to Kubed? Please start [here](/docs/tutorials/README.md).

# Supported Notifiers
Kubed can send notifications via Email, SMS or Chat for various operations using [appscode/go-notify](https://github.com/appscode/go-notify) library. To connect to these services, you need to create a Secret with the appropriate keys. Then pass the secret name to Kubed by setting `notifierSecretName` field in Kubed cluster config.

## Hipchat
To configure Hipchat, create a Secret with the following key:

| Name                | Description                               |
|---------------------|-------------------------------------------|
| HIPCHAT_AUTH_TOKEN  | `Required` Hipchat authentication token   |

```console
$ echo -n 'your-hipchat-auth-token' > HIPCHAT_AUTH_TOKEN
$ kubectl create secret generic kubed-notifier -n kube-system \
    --from-file=./HIPCHAT_AUTH_TOKEN
secret "kubed-notifier" created
```
```yaml
apiVersion: v1
data:
  HIPCHAT_AUTH_TOKEN: eW91ci1oaXBjaGF0LWF1dGgtdG9rZW4=
kind: Secret
metadata:
  creationTimestamp: 2017-07-25T01:54:37Z
  name: kubed-notifier
  namespace: kube-system
  resourceVersion: "2244"
  selfLink: /api/v1/namespaces/kube-system/secrets/kubed-notifier
  uid: 372bc159-70dc-11e7-9b0b-080027503732
type: Opaque
```

Now, to receiver notifications via Hipchat, configure receiver as below:
 - notifier: `hipchat`
 - to: a list of chat room names

```yaml
recycleBin:
  handle_update: false
  path: /tmp/kubed
  receiver:
    notifier: hipchat
    to:
    - ops-alerts
  ttl: 168h
```

```console
$ echo -n 'your-mailgun-domain' > MAILGUN_DOMAIN
$ echo -n 'no-reply@example.com' > MAILGUN_FROM
$ echo -n 'your-mailgun-api-key' > MAILGUN_API_KEY
$ echo -n 'your-mailgun-public-api-key' > MAILGUN_PUBLIC_API_KEY
$ kubectl create secret generic kubed-notifier -n kube-system \
    --from-file=./MAILGUN_DOMAIN \
    --from-file=./MAILGUN_FROM \
    --from-file=./MAILGUN_API_KEY \
    --from-file=./MAILGUN_PUBLIC_API_KEY
secret "kubed-notifier" created
```


| Name                    | Description                                                                    |
| :---                    | :---                                                                           |
| MAILGUN_DOMAIN          | Set domain name for mailgun configuration                                      |
| MAILGUN_API_KEY         | Set mailgun API Key                                                            |
| MAILGUN_PUBLIC_API_KEY  | Set mailgun public API Key                                                     |
| MAILGUN_FROM            | Set sender address for notification                                            |


These environment variables will be set using `searchlight-icinga` Secret.

> Set `NOTIFY_VIA` to `mailgun`







```console
$ echo -n 'your-smtp-host' > SMTP_HOST
$ echo -n 'your-smtp-port' > SMTP_PORT
$ echo -n 'your-smtp-insecure-skip-verify' > SMTP_INSECURE_SKIP_VERIFY
$ echo -n 'your-smtp-username' > SMTP_USERNAME
$ echo -n 'your-smtp-password' > SMTP_PASSWORD
$ echo -n 'your-smtp-from' > SMTP_FROM
$ kubectl create secret generic kubed-notifier -n kube-system \
    --from-file=./SMTP_HOST \
    --from-file=./SMTP_PORT \
    --from-file=./SMTP_INSECURE_SKIP_VERIFY \
    --from-file=./SMTP_USERNAME \
    --from-file=./SMTP_PASSWORD \
    --from-file=./SMTP_FROM
secret "kubed-notifier" created
```

| Name                      | Description                                                                    |
| :---                      | :---                                                                           |
| SMTP_HOST                 | Set host address of smtp server                                                |
| SMTP_PORT                 | Set port of smtp server                                                        |
| SMTP_INSECURE_SKIP_VERIFY | Set `true` to skip ssl verification                                            |
| SMTP_USERNAME             | Set username                                                                   |
| SMTP_PASSWORD             | Set password                                                                   |
| SMTP_FROM                 | Set sender address for notification                                            |


These environment variables will be set using `searchlight-icinga` Secret.

> Set `NOTIFY_VIA` to `smtp`







```console
$ echo -n 'your-twilio-account-sid' > TWILIO_ACCOUNT_SID
$ echo -n 'your-twilio-auth-token' > TWILIO_AUTH_TOKEN
$ echo -n 'your-twilio-from' > TWILIO_FROM
$ kubectl create secret generic kubed-notifier -n kube-system \
    --from-file=./TWILIO_ACCOUNT_SID \
    --from-file=./TWILIO_AUTH_TOKEN \
    --from-file=./TWILIO_FROM
secret "kubed-notifier" created
```

| Name                | Description                                                                        |
| :---                | :---                                                                               |
| TWILIO_ACCOUNT_SID  | Set twilio account SID                                                             |
| TWILIO_AUTH_TOKEN   | Set twilio authentication token                                                    |
| TWILIO_FROM         | Set sender mobile number for notification                                          |



These environment variables will be set using `searchlight-icinga` Secret.

> Set `NOTIFY_VIA` to `twilio`






```console
$ echo -n 'your-slack-auth-token' > SLACK_AUTH_TOKEN
$ echo -n 'your-slack-channel' > SLACK_CHANNEL
$ kubectl create secret generic kubed-notifier -n kube-system \
    --from-file=./SLACK_AUTH_TOKEN \
    --from-file=./SLACK_CHANNEL
secret "kubed-notifier" created
```

##### envconfig for `slack`

| Name             | Description                                                               |
| :---             | :---                                                                      |
| SLACK_AUTH_TOKEN | Set slack access authentication token                                     |
| SLACK_CHANNEL    | Set slack channel name. For multiple channels, set comma separated names. |


#### Add Searchlight app
Add Searchlight app in your slack channel and use provided `bot_access_token`.

<a href="https://slack.com/oauth/authorize?scope=bot&client_id=31843174386.143405120770"><img alt="Add to Slack" height="40" width="139" src="https://platform.slack-edge.com/img/add_to_slack.png" srcset="https://platform.slack-edge.com/img/add_to_slack.png 1x, https://platform.slack-edge.com/img/add_to_slack@2x.png 2x" /></a>

#### Set Environment Variables

These environment variables will be set using `searchlight-icinga` Secret.

> Set `NOTIFY_VIA` to `slack`






```console
$ echo -n 'your-plivo-auth-id' > PLIVO_AUTH_ID
$ echo -n 'your-plivo-auth-token' > PLIVO_AUTH_TOKEN
$ echo -n 'your-plivo-from' > PLIVO_FROM
$ kubectl create secret generic kubed-notifier -n kube-system \
    --from-file=./PLIVO_AUTH_ID \
    --from-file=./PLIVO_AUTH_TOKEN \
    --from-file=./PLIVO_FROM
secret "kubed-notifier" created
```

| Name              | Description                                                                        |
| :---              | :---                                                                               |
| PLIVO_AUTH_ID     | Set plivo auth ID                                                                  |
| PLIVO_AUTH_TOKEN  | Set plivo authentication token                                                     |
| PLIVO_FROM        | Set sender mobile number for notification                                          |



These environment variables will be set using `searchlight-icinga` Secret.

> Set `NOTIFY_VIA` to `plivo`
