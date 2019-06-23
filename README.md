# GomiSig
Slack notification of garbage calendar

# Required
gcloud or server machine

# Deploy
Example of deploying function at your GCP project

## Architecture
`Cloud Scheduler -> Cloud Pub/Sub -> Cloud Functions`

## Getting Started
* deploy this function

You can get Webhook URL [here](https://slack.com/services/new/incoming-webhook).
```
gcloud functions deploy Gomidashi \
 --project <project-name> \
 --runtime go111 \
 --trigger-topic gomidashi-notify \
 --set-env-vars SLACK_WEBHOOK_URL=<your-webhook-url> \
 --region asia-northeast1
```

* deploy scheduler example
```
gcloud beta scheduler jobs create pubsub Gomidashi \
  --project <project-id> \
  --schedule '0 8 * * *' \
  --topic gomidashi-notify \
  --message-body '{"mention":"channel", "channel":"gomi","area":"north"}' \
  --time-zone 'Asia/Tokyo' \
  --description 'Slack notification scheduler of garbage calendar'
```

# Special Thanks
JSONデータの作成は@ysakashinさんの[tsukuba-gc-lite](https://github.com/ysakasin/tsukuba-gc-lite)を使わせていただきました。
