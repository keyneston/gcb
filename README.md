# Generic Circuit Breakers

Heavily Inspired by [hystrix-go](https://github.com/afex/hystrix-go).

## Goals

* Golang Circuit Breakers with a clean and simple to use API
* Hierarchial circuits. 

## Hierarchical Circuits

__Currently Theoretical__

When calls are made to GCB a name is passed that is used to lookup which
circuit to use. Traditionally this lets circuits be configured and triggered on
a one by one basis. With a hierarchical setup we can create configurations that
where settings are set and triggers are triggered in a cascading manner.


For example we could have the key `www.example.com!80`, and
`www.example.com!443`. We can then configure it such that `www.example.com` and
all child circuits fail at 90%, but then set the `443` child circuit to fail at
95%. Now if the average of both the `443` and `80` circuits hits 90% both child
circuits will trigger.

## Future Thoughts

* Use a pool or similar for things like timers to stop memory churn
* API for gathering circuit data and values
