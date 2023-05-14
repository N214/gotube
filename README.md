# Gotube

A Golang implementation of the PubSubHubbub notification subscriber to send notifications to slack when a new video of your favorite YouTuber is available.

## Architectures

PubSubHubbub -> Cloud functions -> Cloud Storage -> Slack

This implementation ignores the update video notification by leveraging a state file in the GCS

## TODO
- Add ngrock doc for local run
- Split checkDataHistory into CheckDataHistory and writeDataHistory


## Local run