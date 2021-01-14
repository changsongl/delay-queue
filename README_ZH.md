# delay-queue

### 介绍
**这个项目还在持续开发当中，功能还不完善到生产使用。当基本功能完成和单元测试覆盖后，
将发布Beta版本，提供大家使用。**

这个项目是仿照有赞的延迟队列进行设计的。现在这个队列是通过Redis来进行存储和实现分布式高可用的，
这个项目会进行存储的拓展，如MYSQL，RABBITMQ，内存等实现方式的集成。

### 设计计划
1. 这个延迟队列是支持分布式高可用的。
2. 其支持多种存储引擎支持，如Redis, Mysql, Rabbit Mq, Memory等等。
3. 支持多种协议，如http, tcp, grpc等等。
4. 设计对象编程，可以轻松替换组件和实现多态。
5. 单元测试覆盖到100%。

### 使用用例
现在delay-queue的客户端还未提供，因此现在的用例为次项目已经支持的http请求。

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

### 设计

#### Terms
1. Job：需要异步处理的任务，是延迟队列里的基本单元。与具体的Topic关联在一起。
2. Topic：一组相同类型Job的集合（队列）。供消费者来订阅。

#### 任务
1. Topic：Job类型。可以理解成具体的业务名称。
2. Id：Job的唯一标识。用来检索和删除指定的Job信息。Topic和Id的组合应该是业务中唯一的。
3. Delay：Job需要延迟的时间。单位：秒。（服务端会将其转换为绝对时间）
4. TTR（time-to-run)：Job执行超时时间，超过此事件后，会将此Job再次发给消费者消费。单位：秒。
5. Body：Job的内容，供消费者做具体的业务处理，可以为json格式。


#### 组件

>有四个组件
>1. Job Pool: 用来存放所有Job的元信息。
>2. Delay Bucket: 是一组以时间为维度的有序队列，用来存放所有需要延迟的／已经被reserve的Job（这里只存放Job Id）。
>3. Timer: 负责实时扫描各个Bucket，并将delay时间大于等于当前时间的Job放入到对应的Ready Queue。
>4. Ready Queue: 存放处于Ready状态的Job（这里只存放Job Id），以供消费程序消费。

<img alt="delay-queue" src="https://tech.youzan.com/content/images/2016/03/delay-queue.png" width="80%">

#### 状态
>Job的状态一共有4种，同一时间下只能有一种状态。
>1. ready：可执行状态，等待消费。
>2. delay：不可执行状态，等待时钟周期。
>3. reserved：已被消费者读取，但还未得到消费者的响应（delete、finish）。
>4. deleted：已被消费完成或者已被删除。

<img alt="job-state" src="/doc/pic/job-state.png" width="80%">

### 项目计划
我将持续打磨这个项目，并且加入更多的功能和修复问题。我将会让这个项目可以投入到生产环境使用。
如果喜欢的话，欢迎给个星或者Fork参与进来，这里欢迎你的贡献！
 
### 如何贡献?
1. 在Issue里发布自己的问题或评论。
2. 我们会在问题中进行讨论，并进行设计如何开发。
3. Fork项目进行开发，并以develop分支，创建自己的分支进行开发。如fix-xxx, feature-xxx等等。
4. 开发完成后，发起PR合入develop。
5. Code Review后将会把你的代码合进分支。
 
### Reference

Youzan Design Concept [Youzan Link](https://tech.youzan.com/queuing_delay/)


### 如果看不到文档的图片
绑定github图片hosts: 199.232.96.133 raw.githubusercontent.com
