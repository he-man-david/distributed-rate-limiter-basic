# Design and build rate limiter service in Go

## Context:

There are many ways to RL users, many policies, and each user may have different policies. E.G - the way you RL a Netflix user by USER_ID is very different from RL a developer’s API_KEY as an API provider. A Netflix user is a person, and will not change regions frequently (even if they share accounts), as there is a physical limitation. API_KEY is used in another dev’s backend, and they may make calls from N regions. For API_KEY you would need local and global (distributed) RL via synchronization, whereas for Netflix users only local RL suffice. This is just one example. A real life production level RL can get very complex, but we will just focus on API_KEY RL.

## Functional Requirements:

1. Rate limit requests based on API-KEY
2. 10,000 requests per minute ~ roughly 160 RPS
3. Linear backoff of 30% for 3 minutes
   - If user hit RL, for next 3 minutes they are limited to 7000 RPM
   - If they hit RL again, each time will linear backoff (Last RPM \* 70%)
4. Distributed rate limiter service that needs synchronization
5. BONUS - deploy via minikube
6. BONUS - testing

## Technical Requirements:

1. Building demo on personal machine, loads are not real life
2. Eventually consistent
3. Fault tolerant, able to handle RL failure
4. Crash recovery
5. Highly available + extreme perf low latency

## Networking protocol:

Since RL is a feature of API gateways, and we can run RL as a separate proxy within AG, we can use RPC protocol. Reasons are:

1. AG already takes in HTTP request from users, validates request, sessions…etc
2. AG will convert HTTP request to RPC, gather, and then dispatch for downstream services
3. RPC is lightweight, and message size is 30-40% of HTTP, which saves cost. E.G - 1 user HTTP request can result in 10+ downstream calls, so minimizing package size can save lot of money.
4. We can use gRPC for all our micro-services, which offer IDL and code generation

With the above assumptions, since only our AG will be talking to RL (as it is an embedded but separate proxy service), we can use gRPC.

## RL strategy:

We will use the sliding-window RL strategy, implemented using linked list for each API_KEY. Reason is we want to prevent burst traffic from causing damage. The downside is higher memory costs, but that may be tolerable.

- Each API_KEY will have it’s own sliding-window LL
- Each LL node will just be timestamp, 32 bit integer, which is 4 bytes

## Calculations:

If we assume that our rate limit is 10K request per minute per API_KEY, each rate limiter service will store 40KB of LL data per API_KEY. Which comes out to 40GB for 1mm users. Since our RL lives inside AG, and AG traffics will be partitioned via API_KEY, if the partition is 100K keys per AG, then our RL will only need 4GB of memory at the most.

## High level design:

![High Level Design](/static/design.png)

## Point of consideration

Our logic for RL is quite simple:

    T  = timeframe for RL in epoch ms
    R  = the # of requests per T
    Tn = Time now epoch ms
    T1 = earliest time epoch ms in our LL
    TR = total requests received thus far

    if (TR >= R && Tn - T1 < T) return false
    else return true

Everything is accessible locally, but TR and T1 is difficult to determine. Since our system is distributed, TR and T1 needs to be global. Since we are synching, we can determine TR with a small tolerable rate of error (10,000 vs 10,050).

But "what is global T1?". To answer this question, not only do we track the LL for each RL Node replicate separately, we also need to have a min heap for tracking the T1 for N # of LL.

## Rate Module

This module will expose gRPC functions:

    allowRequest()
    - Request: API_KEY
    - Reponse: True/False

## Policy Config Module

This module is incharged of storing metadata for each API_KEY. Idea is that different users will have different limits, depending on what service they purchased. Some can pay extra to allow burst request limits..etc When rate limits are violated, Rate Module will check Policy Config Module to figure out what to do.

## Rate Tracker Instance

For each API_KEY, it is not as simple as a single LL, if we want to have a distributed rate limiter. To allow concurrent global rate limiting, synchronization is required. To facilitate this process, we would need this Rate Tracker Instance for each key. This is basically a collection of LL, each RL node would have it's own LL stored locally. We would also have a min heap to store t1.
This way, when syncs are coming in from N nodes, we append to their own LL in our RTI (rate tracker instance). Our min heap tracks the global t1, which is just the smallest of each LL's head.

This essentially solves the "What is t1?" problem.

## Synchronization Module

This module will be responsible for handling incoming sync requests, and broadcasting the node's own requests. To minimize bandwidth cost, this message can just be a small JSON of ratelimit_ID, API_KEY, timestamp. Incoming sync request will contain ratelimit_ID, and timestamp also, which we will append to LL of the respective NODE, for the given API_KEY.

## Scope

1. Rate limiter proxy
2. Rate Module
3. Policy Config Module
4. Synchronization Module
5. Rate Tracker Instance
6. Support distributed RL & sync
7. Deploy via minikube to test
8. E2E testing

## Not in scope

1. Crash recovery, as that would need coordinator
2. We are not building a AG
3. Not dealing with partioning for now
