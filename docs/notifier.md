```sh
$ echo -n 'your-mailgun-domain' > MAILGUN_DOMAIN
$ echo -n 'mailgun-from' > MAILGUN_FROM
$ echo -n 'mailgun-to' > MAILGUN_TO
$ echo -n 'your-mailgun-api-key' > MAILGUN_API_KEY
$ echo -n 'your-mailgun-public-api-key' > MAILGUN_PUBLIC_API_KEY
$ kubectl create secret generic kubed-notifier \
    --from-file=./MAILGUN_DOMAIN \
    --from-file=./MAILGUN_FROM \
    --from-file=./MAILGUN_TO \
    --from-file=./MAILGUN_API_KEY \
    --from-file=./MAILGUN_PUBLIC_API_KEY
secret "kubed-notifier" created
```