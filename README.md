# messenger-app

A messaging app using Go standard library that implements Event-Driven Architecture (EDA).

## Technologies used

- Go channels for queueing
- MongoDB for persistence

## Planned additions / adjustments

- make events more lightweight --> store only `message_id` instead of the whole message
- move to a server set-up --> using Gin Gonic
