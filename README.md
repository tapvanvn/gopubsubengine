# gopubsubengine
## wspubsub
- type
    -  pick_one : pick only one randomly subscriber to send message
- tier
    - one : make sure the receiver got the message. After 5 second if the receiver not reponse, resend or send to another subscriber.
    - two : make sure the receiver got the message and keep alive. if the current receiver die, send to another one and so on.

# This library is on development stage.
