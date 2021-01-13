# delay-queue
Translations:

- [中文文档](./README_ZH.md)

### 介绍
**This is project is still working in progress. Don't use it in production. I will release BETA when it is ready.**

This project is a delay queue. It is based on Youzan 有赞 delay queue. Currently,
it is based on Redis for storage. It will support more types of storages in the 
future.

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

#### Terms
1. Job: It is a task to be processed, and it is related to only one topic.
2. Topic: It is a set of jobs, it is implemented by a time-sorted queue.
 All consumers need to choose at least one topic to consume jobs.

#### Job
Jobs contain many properties like:
1. Topic: It could be a service name, users can define it depending on their
 business.
2. ID: it is unique key for inside of a topic. It's used to search job information
 in a topic. The combination of a topic and an ID should be unique in your
 business.
3. Delay: It defines how many second to be delay for the job. Unit: Second
4. TTR(time to run): It is job processing timeout. If consumer process this
 job more than TTR seconds, it might be sent to other consumer, if a consumer
 pop the topic.
5. Body: It is content of job. It is a string. You can put your json data to it.
 When you consume the job, you can decode it and run your logic.


#### Component

>There are 4 components in the delay queue.
>1. Job Pool: It saves all metadata of jobs.
>2. Delay Bucket: It is a time-sorted queue. It saves jobs that is waiting
 for being ready. There are more than one Bucket in the delay queue for
 higher throughput.
>3. Timer: It is a core component to scan the Delay Bucket. It pops out 
 ready jobs from Buckets and put then inside ready queue.
>4. Ready Queue: It is a queue for storing all ready jobs, which can be
 popped now. It is also only store the job id for the consumers.

<img alt="delay-queue" src="/doc/pic/delay-queue.png" width="80%">

#### States
>There are four states for jobs in the delay queue. The job can be only
> in one state at the time.
>1. Ready: It is ready to be consumed.
>2. Delay: It is waiting for the delay time, and it can't be consumed.
>3. reserved: It means the job has consumed by a consumer, but consumer
> hasn't ack the job. (Call delete、finish).
>4. Deleted: The job has finished or deleted.

<img alt="job-state" src="/doc/pic/job-state.png" width="80%">

### What's the plan of this project?
I will work on this project all the time! I will add more features and 
 fix bugs, and I will make this project ready to use in production. Star
 Or Fork it if you like it. I'm very welcome to you for contribution.
 
### How to contribute?
1. Level a message in the unsigned issue.
2. We will discuss how to do it, and I will assign the issue to you.
3. Fork the project, and checkout your branch from "develop" branch.
4. Submit the PR to "develop" branch.
5. It will be merged after code review.
 
### Reference

Youzan Design Concept [Youzan Link](https://tech.youzan.com/queuing_delay/)


### Can't See Image In China?
bind host: 199.232.96.133 raw.githubusercontent.com
