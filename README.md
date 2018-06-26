# glue

> glue two programs in a client/server manner, over stdin/stdout

There is an example app which uses two other cli programs as services, over stdin/stdout.

The basic idea is about being able to change a program parts without the need for recompiling the whole app. The initial ambitiion was to implement _HTTP over commandline_. But it turned out as not that useful.

* Is this really useful? 
  
  > On servers we use usually a microservices architecture with service discovery, monitoring, etc, etc with specific deployment strategies. So this might not be of much good use on our servers.

* What about clients?
  
  > It depends. I had this app that was doing some image processing that would eventually succeed. But in rare ocasions, it was bringing down whole application. I used NATS Server, embedded in my app to communicate with external (GO) process that was doing the image processing. Maybe there were some easier ways (it had to run on Windows too - though so far ran only on Linux). This is a try to have something sompler and more cross-platform.