# delay-queue
**This is project is still working in progress. Don't use it in production.**

This project is a delay queue. It is based on Youzan 有赞 delay queue. Currently,
it is based on Redis for storage. It will support more types of storages in the 
future.

The design is inspire from Youzan Delay Queue.

https://tech.youzan.com/queuing_delay/

### Design Plans
1. This delay queue could be scaled and HA.
2. Supporting different clustered storage like redis, mysql, rabbitmq, memory and so on.
3. Support different protocol like http, tcp and grpc.
4. Be OOD to separate logic and different storage implementations.
5. Unit Tested.

### Usage
Currently, the client has not written yet. Therefore, I write http
request to give examples.

````http request
// Push job
POST 127.0.0.1:8080/topic/mytopic/job
body: {"id": "myid1","delay":10, "ttr":4, "body":"body"}

// response 
{
    "message": "ok",
    "success": true
}
````

````http request
// Pop job
GET 127.0.0.1:8080/topic/mytopic/job

// response
{
    "id": "myid1",
    "success": true,
    "value": "body"
}
````

````http request
// Delete job
DELETE 127.0.0.1:8080/topic/mytopic/job/myid1

// response
{
    "message": "ok",
    "success": true
}
````

````http request
// Delete job
PUT 127.0.0.1:8080/topic/mytopic/job/myid1

// response
{
    "message": "ok",
    "success": true
}
````

### Designs

##### - Component
<img src="/doc/pic/delay-queue.png" width="100%">



