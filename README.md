# delay-queue
This is project is still working in progress. It will be released 
a prototype in Feb.

A redis delay queue. It is based on Youzan 有赞 delay queue.

https://tech.youzan.com/queuing_delay/

### Design Plans
1. This delay queue could be scaled and HA.
2. Supporting different clustered storage like redis, mysql, rabbitmq, memory and so on.
3. Support different protocol like http, tcp and grpc.
4. Be OOD to separate logic and different storage implementations.
5. Unit Tested.
