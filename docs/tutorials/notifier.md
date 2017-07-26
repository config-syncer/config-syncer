> New to KubeDB? Please start [here](/docs/tutorials/README.md).

```console
$ echo -n 'your-mailgun-domain' > MAILGUN_DOMAIN
$ echo -n 'no-reply@example.com' > MAILGUN_FROM
$ echo -n 'your-mailgun-api-key' > MAILGUN_API_KEY
$ echo -n 'your-mailgun-public-api-key' > MAILGUN_PUBLIC_API_KEY
$ kubectl create secret generic kubed-notifier \
    --from-file=./MAILGUN_DOMAIN \
    --from-file=./MAILGUN_FROM \
    --from-file=./MAILGUN_API_KEY \
    --from-file=./MAILGUN_PUBLIC_API_KEY
secret "kubed-notifier" created
```



```console
$ echo -n 'your-hipchat-auth-token' > HIPCHAT_AUTH_TOKEN
$ kubectl create secret generic kubed-notifier \
    --from-file=./HIPCHAT_AUTH_TOKEN
secret "kubed-notifier" created
```

```console
$ echo -n 'your-smtp-host' > SMTP_HOST
$ echo -n 'your-smtp-port' > SMTP_PORT
$ echo -n 'your-smtp-insecure-skip-verify' > SMTP_INSECURE_SKIP_VERIFY
$ echo -n 'your-smtp-username' > SMTP_USERNAME
$ echo -n 'your-smtp-password' > SMTP_PASSWORD
$ echo -n 'your-smtp-from' > SMTP_FROM
$ kubectl create secret generic kubed-notifier \
    --from-file=./SMTP_HOST \
    --from-file=./SMTP_PORT \
    --from-file=./SMTP_INSECURE_SKIP_VERIFY \
    --from-file=./SMTP_USERNAME \
    --from-file=./SMTP_PASSWORD \
    --from-file=./SMTP_FROM
secret "kubed-notifier" created
```

```console
$ echo -n 'your-twilio-account-sid' > TWILIO_ACCOUNT_SID
$ echo -n 'your-twilio-auth-token' > TWILIO_AUTH_TOKEN
$ echo -n 'your-twilio-from' > TWILIO_FROM
$ kubectl create secret generic kubed-notifier \
    --from-file=./TWILIO_ACCOUNT_SID \
    --from-file=./TWILIO_AUTH_TOKEN \
    --from-file=./TWILIO_FROM
secret "kubed-notifier" created
```


```console
$ echo -n 'your-slack-auth-token' > SLACK_AUTH_TOKEN
$ echo -n 'your-slack-channel' > SLACK_CHANNEL
$ kubectl create secret generic kubed-notifier \
    --from-file=./SLACK_AUTH_TOKEN \
    --from-file=./SLACK_CHANNEL
secret "kubed-notifier" created
```


```console
$ echo -n 'your-plivo-auth-id' > PLIVO_AUTH_ID
$ echo -n 'your-plivo-auth-token' > PLIVO_AUTH_TOKEN
$ echo -n 'your-plivo-from' > PLIVO_FROM
$ kubectl create secret generic kubed-notifier \
    --from-file=./PLIVO_AUTH_ID \
    --from-file=./PLIVO_AUTH_TOKEN \
    --from-file=./PLIVO_FROM
secret "kubed-notifier" created
```
