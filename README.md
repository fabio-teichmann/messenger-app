# messenger-app

A messaging app using Go standard library that implements Event-Driven Architecture (EDA).

## Technologies used

- Go channels for queueing
- MongoDB for persistence

## Planned additions / adjustments

- mechanism to notify all users that chat with sender (`USER_ONLINE` event)
- user persistence
- make events more lightweight --> store only `message_id` instead of the whole message
- containerisation / clusterization
  - k8s package Mosquito (?)
- move to a server set-up --> using Gin Gonic
  - user login/auth
- basic FE interface for demonstration
  - login-screen
  - overview chats
  - user profile
  - chat interactions
