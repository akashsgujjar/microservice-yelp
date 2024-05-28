# Lab 2: Simple Performance Analysis

In this lab, we will evaluate the performance of Welp from the perspective of a user and developer. The goal of this lab is to build our intuition about analyzing application performance and to see how ideas from queueing theory can give us insights about our application.

<!-- - mention we want distribution of latencies since we need 99% -->

## Lab 2 Overview
Three parties are involved, with different performance goals, in every
cloud application: developers, users, and cloud operators. This lab will serve as an introduction to analyzing application performance through the lens of users (Assignment 1) and developers (Assignment 2). Cloud application developers desire a simple development experience for their applications
that allows them to elastically increase application service when
there is demand. As discussed in Lab 1, our choice of a microservices development model provides a great degree of simplicity for developers that want to quickly scale up applications.

Users, on the other hand, care most about quality service. Since
most cloud applications are interactive, one very important quality
metric is the service latency of the cloud application that the user
connects to via the Internet.

Studies have shown that users will tolerate about 500ms
of service latency. For our application, service latency refers to the full time period from when the user submits a request to when the user receives the response back. 
Anything above this makes users increasingly likely to assume that there is 
something wrong with the provided service and go elsewhere. Hence, this number 
is typically used as an upper bound on user quality of service tolerance. 

In this assignment we will analyze the performance of our microservices application and look at the relationship between service latency and throughput. As part of our analysis, we will utilize insights from queueing theory.

## Load Generators
To run our performance analysis, we need a load generator. The purpose of a load generator is to generate user requests under different conditions in order to simulate the behavior of our system in the real world. 

<div style="text-align: center;">
  <img src="../images/lab2-arch.png" alt="Simple Architecture" width="30%" height="30%">
</div>



### Open-loop model
In an open-loop load generator model, new request arrivals occur independently of completions. Requests are generated and sent based on a pre-defined trace/schedule or a random process such as a Poisson point process.

**Pros:** Useful for simulating specific user behavior or test scenarios. For example, we can set a given rate and see how burstier arrival distributions affect latency. Response time (latency) is also a reflection of service time and potential queueing delay, which gives us insight into the system when it is overloaded.

**Cons:** Potentially less accurate in replicating real-world usage patterns.

<div style="text-align: center;">
  <img src="../images/open-loop.png" alt="open loop" width="40%" height="40%">
</div>


**`wrk2`:**
We will use an open-loop version of the load generator `wrk2` for HTTP benchmarking on Welp. 


## Queueing Theory
Recall Little's Law from lecture, which gives the relationship between the mean throughput $`X`$, response time $`R`$ (a.k.a. service latency plus queueing delay), and number of tasks $`N`$ in a stable system:
$$ N = X\cdot R $$
Importantly, this relationship applies regardless of the distribution of any of the random variables as (provided the system is stable i.e. $`0 \leq U < 1`$ where $`U`$ is the utilization).  

Ultimately, our goal is to understand how our system behaves under various load conditions.

<div align="center">

| Independent Variables | $`\lambda`$ |
|-----------------------|---------|
| Responding Variables  | $`X`$, $`R`$ |
</div>


With an open-loop model, we directly influence the arrival rate $`\lambda`$ and measure $`X`$ and $`R`$. Since jobs in an open-loop generator arrive independent of previous job completion, it is now possible for the arrival rate to exceed the service rate $\mu$, meaning the system can become unstable ($`U > 1 \longleftrightarrow \lambda > \mu`$). At this point, the queue can grow without bound, which pushes the response times towards infinity. By fixing the load $`\lambda`$ we can guarantee the mean number of requests in the system, but not the exact number of requests. Unfortunately, this means that even in cases where $`\lambda < \mu`$, a bursty set of arrivals in a small time period could overload our system and incur queueing delays!

<!-- TODO: describe arrivals -->

## Assessing Service Performance

We typically assess performance by evaluating service latency provided
under a variety of load conditions. System load has a major effect on
the provided service. If a service is overloaded (i.e., the load exceeds the service rate), it is going to provide terrible 
service---a large portion of users' requests is going to take a while,
as other answers will be processed first. Lower load conditions are going
to provide a spread of better service latencies. To see how a system
performs under a variety of load conditions, we draw latency-throughput
curves. For each throughput level (presented on the x-axis), we assess
service latency (presented on the y-axis).

For cloud operators, it is typically not enough to keep average
service latencies below 500ms. Averages *ignore* the distribution of
service experiences. Thus, they could preclude a potentially large
fraction of users from acceptable service. Instead, we care about the
largest fraction of users, knowing that some users are bound to have
bad performance for reasons outside of our control (e.g., due to a bad
Internet connection). Hence, we assess a sufficiently large percentile
of users. This is typically 99% of users or more. The actual number
depends on our assessment of what is in our control and what is
not. For this reason, we will look at 99%-ile latencies in our analysis of users.

<div style="text-align: center;">
  <img src="../images/throughput-vs-latency-sla.png" alt="open loop" width="40%" height="40%">
</div>
Note: max point refers to the maximum throughput you can get under a specific tail latency target (i.e. 500ms)

## Lab 2 Prep Work:
We did not enforce synchronization of the in-memory data structure in lab 1. However, since our workloads might run into concurrent access/writes, you should make sure that your in-memory data structures is thread-safe (Note that the detail service we provide is not thread-safe.). Feel free to use a synchronization mechanism of your choice (some options for doing this include using a mutex, concurrent hash map, or channels). Once this is working and you are still passing the tests from lab 1 run the sample workload below. (Hint: you will *also* probably need to synchronize any additional data structures you used to implement the `MostPopular` and `SearchReviews` RPCs  for reservation and review respectively). 

Do a sample run of `wrk2` using the sample workload configuration file in `workloads/` to learn how `wrk2` works. Specify a load rate and record the corresponding 99%-ile latency and throughput. For example, we can run the below command (in this example, we record 9047.41 req/sec for the load, 12004.74 req/sec for the throughput and 2.88ms as the 99% response time).

```
(base) user@user-vm:~/cse453-welp$ ./wrk2/wrk -D const -t 10 -c 100 -d1m -s ./workloads/lab2/sample-workload.lua http://10.96.88.88:8080 -R 12000 -L
Running 1m test @ http://10.96.88.88:8080
  10 threads and 100 connections
...
  Thread Stats   Avg      Stdev     99%   +/- Stdev
    Latency     0.96ms  610.97us   2.88ms   87.84%
    Req/Sec     1.27k   104.00     1.55k    76.83%
  Latency Distribution (HdrHistogram - Recorded Latency)
 50.000%    0.86ms
...
100.000%   13.22ms
----------------------------------------------------------
  720247 completed requests in 1.00m, 91.36MB read
Completed Requests/sec (Throughput):  12004.74
Sent Requests/sec (Load):   9047.41
Transfer/sec:      1.52MB
```
Parameter `R` controls the rate of requests. Change `R` to increase/decrease load. Setting `t` and `c` equal to 10 will provide reasonable performance for your experiments. `D` controls the arrival distribution (`const` represents uniform distribution and `exp` represents exponential distribution). 

## Assignment 1: Performance Analysis (7 pts)
We will run a performance analsyis on our application and look at performance under different arrival distributions.  

Our goal is to generate throughput vs latency graphs for our simple workload under two different setups: (1) uniform arrival rate of requests and (2) exponential arrival rate of requests. Use the workload `./workloads/lab2/simple-mixed.lua` for this assignment. We generate the throughput vs latency graphs by collecting data for load vs latency and load vs throughput and then directly compare these sets of data pairs in latency vs throughput charts (you will have a total of 6 graphs, 3 for each setup).
<div style="display: flex;">
    <img src="../images/load-vs-latency.png" alt="Image 1" width="30%" height="30%" style="margin: 5px;">
    <img src="../images/load-vs-throughput.png" alt="Image 2" width="30%" height="30%" style="margin: 5px;">
    <img src="../images/throughput-vs-latency-sla.png" alt="Image 3" width="30%" height="30%" style="margin: 5px;">
</div>

No experiment setup is perfect. In the real world, there are always
factors outside of our control. From the Internet connection down to the temperature on a given day affecting server
performance, many factors can impact experiment outcome. Hence, re-running the
experiment is going to produce slightly different results. The
difference in results is generally described as noise or error. We
strive to create experimental environments where the noise is
minimized, leaving mostly the true result, or signal.

One way to control for noise is to run the experiment multiple times
and to calculate average data points from each run. However, while
this smoothes results, it does not give us an understanding as to the
extent of the noise. As experimenters, we should quantify the noise to
show that our experiment is indeed measuring the signal and not
presenting mostly noise. To do so, we attach error bars to each data
point. Hence, for each data point you generate make sure to run `wrk2` for at least 5 trials. 
Use the average values and include error bars in your graphs (use standard deviation for the error bars).

Once you have generated your graphs, comment on the following questions in your writeup:
<!-- - What throughput can each setup handle before 99%-ile latencies exceed our SLA of 500ms? -->
- Comment on the shape of your load vs latency graphs. Do you see any differences?
- Under what load do you observe the knee points in the load-latency and load-throughput graphs?
- Where are the max points and what do they signify? 
- What level of noise do you observe in your throughput vs. latency and why might it be there?
- What differences do you see in latency and throughput between the
  setup with a uniform arrival rate and an exponential (burstier)
  arrival rate? Make a hypothesis about what might be happening in relation to the
  queues in the system (we will look at the queues of different services in the next section)
<!-- Once you have identified the max point, create graphs of the latency CDF for the application (run `wrk2` at the largest load the application can tolerate before latencies exceed the max point). Do this for the exponential and uniform arrival distribution. In your CDF, make sure to include the following percentiles [50.00%, 75.00%, 90.00%, 99.00%, 99.99%, 100.000%]. Along with the CDF, plot a vertical line indicating the average latency. Comment on the distribution of latencies for uniform and exponential arrival rates and the pros / cons of designing applications around averages. -->

## Assignment 2: Bottleneck Analysis (4 pts)
In the previous assignment we looked at the overall performance of the application by looking at our system as a [black box](https://en.wikipedia.org/wiki/Black-box_testing).

We identified the max point for the system by varying the load
conditions. However, if we want to improve the performance of the
application it's not entirely clear what components of the system we
should change. In particular, we would want to know the bottleneck of
the system––the component that limits the processing rate for the
entire system. We will run a simple bottleneck analysis to identify
the bottleneck that is limiting performance.

We can think about bottlenecks from the lens of queueing theory. In
any complex system, there are several interacting components, each
with their own queues of incoming and departing tasks. In assignment
1, we looked at the entire system as a single queue. We could also go
a level further and look at the queues of each pod/server in our
microservices architecture (frontend, detail, review,
reservation). The bottleneck in any multiple-queue system is the
component with an unboundedly growing queue length. This is because
tasks accumulate in front of the bottleneck service faster than it can
process it, due to its limited service rate. For example, if our
frontend was the bottleneck, when applying load just above the max
point, we would see very large (and continually growing) queue lengths
at the frontend and much smaller queue lengths at the detail,
reservation, and review servers. As queue lengths can be difficult to
access as a user of the cloud, we will intercept and measure service
latencies and use them as proxies for queue length.

Run a bottleneck analysis using the mixed workload from Assignment 1
and an exponential arrival rate. For that particular workload,
identify the bottleneck of our architecture and argue your case with
empirical evidence (e.g. graphs, tables). Once you identify the
bottlenck, suggest at least two detailed solutions for remedying the
bottleneck (e.g. adding more replicas, vertical scaling, etc.). Make
sure to discuss the pros/cons of each potential solution and how you
would design the solution so that it continues to correctly implement
the Welp application.

### Tips
#### Parsing RPC Data
You can use our provided script `parse_logs.py` to export the kubernetes log files in a cleaned CSV format. For example, let's post a restaurant review and look at the output.
```bash
curl "http://10.96.88.88:8080/post-review?user_name=dubs&restaurant_name=Hub+U+District+Seattle&review=a%20good%20place%20for%20food?&rating=2"
```
Running `kubectl logs <frontend-pod-name>` gives us the raw data for the duration of each RPC through the system on a per-request basis. You should see a similar output to the example below.
```bash
$ kubectl logs <frontend-pod-name>
2023/09/1 19:49:43 frontend server running at port: 8080
2023/09/12 19:50:21 grpc;/review.ReviewService/PostReview;{"user_name":"foo","restaurant_name":"Hub U District Seattle","review":"a good place for food?","rating":2};{"status":true};<nil>;10251
2023/09/12 19:50:21 frontend.postReviewHandler;{"user_name":"foo","restaurant_name":"Hub U District Seattle","review":"a good place for food?","rating":2};{};<nil>;10705
2023/09/1 19:53:10 grpc;/review.ReviewService/GetReview;{"restaurant_name":"Hub U District Seattle","user_name":"foo"};{"user_name":"foo","restaurant_name":"Hub U District Seattle","review":"a good place for food?","rating":2};<nil>;1454
2023/09/1 19:53:10 frontend.getReviewHandler;{"restaurant_name":"Hub U District Seattle","user_name":"foo"};{};<nil>;1537
```
This format is not particularly readable or easy to do analysis on. We can run pipe the output to our parser, which gives us a much cleaner CSV format.
```bash
$ kubectl logs <frontend-pod-name> | python scripts/parse_logs.py <frontend-pod-name>
Parsed data has been written to <frontend-pod-name>-logs.csv
```
Opening the log file gives the following output. You can copy-paste or export this CSV file to an application like Pandas, Excel, or Google to do your data analysis.
|     date    |    time    |             rpc              |                                               input                                               |            output            |                  error                  | duration(µs) |
|:-----------:|:----------:|:---------------------------:|:--------------------------------------------------------------------------------------------------:|:----------------------------:|:---------------------------------------:|:------------:|
| 09/1/2023  |  19:50:21  | review.ReviewService/PostReview | {"user_name":"foo","restaurant_name":"Hub U District Seattle","review":"a good place for food?","rating":2} | {"status":true}             | \<nil\>                                  |     10251    |
| 09/1/2023  |  19:50:21  | frontend.postReviewHandler    | {"user_name":"foo","restaurant_name":"Hub U District Seattle","review":"a good place for food?","rating":2} | {}                         | \<nil\>                                   |     10705    |
| 09/1/2023  |  19:53:10  | review.ReviewService/GetReview  | {"restaurant_name":"Hub U District Seattle","user_name":"foo"}                                     | {"user_name":"foo","restaurant_name":"Hub U District Seattle","review":"a good place for food?","rating":2} | \<nil\> |     1454     |
| 09/1/2023  |  19:53:10  | frontend.getReviewHandler     | {"restaurant_name":"Hub U District Seattle","user_name":"foo"}                                     | {}                         | \<nil\>                                   |     1537     |

#### Interpreting RPC Data

As mentioned earlier, our goal is to track latencies of RPCs on
different servers and use the measured latencies on a server as
proxies for the queue length on that server. The measured latencies
above, however, represent the total elapsed time between a request
arriving on a server up till the time its response departs the
server. What we want is the approximate duration spent on a particular
server. In the case of review, detail, reservation, the listed
duration *is* the time spent on those servers--these services do not
generate further requests as part of servicing an incoming request and
thus their response is just their local processing time. The frontend
is different--it generates further requests to the backend servers,
before responding. Hence, the actual service time for a request on the
frontend server is the tracked duration at the frontend minus the
tracked duration of any generated RPCs. For example, in the above
sample, the PostReview RPC "averages" 10,251 microseconds on the
review server and 10,705 - 10,251 = 454 microseconds on the frontend
server.

<!-- In reality, there are other sources of latency you would need to -->
<!-- account for, as well (inter process communication, network latency for -->
<!-- networked applications, etc.), but for the purposes of this analysis, -->
<!-- general heuristics of latencies spent at each service is sufficient. -->

## Assignment 3: Little's Law Analysis (4 pts)

Can you use Little's Law to determine the applications's utilization
at different operating points? We'll use the workload from assignment
2 for this, but to keep things simple, we use the uniform arrival
distribution (`-D const`).

Determine the average **service time** at the frontend, using
technique from assignment 2, over multiple stable (non-overloaded)
arrival rates. You can now use Little's Law to compute the frontend's
expected utilization given an arrival rate. Draw a graph that shows
the frontend's utilization on the y axis, for increasing arrival rates
on the x axis, until you compute a utilization of 1.

Conduct a new experiment, measuring average **response time** at the
client over multiple stable (non-overloaded) and unstable (overloaded)
arrival rates. Add another line to the graph, plotting the measured
average response time on the y axis, for each arrival rate. At which x
coordinate does your computed utilization reach 1? Do measured
response times around and to the right of this point start growing
towards infinity (indicating overload)? If not, what could explain the
difference?



<!-- the -->
<!-- latency over throughput graph. Then, take Little's Law and compute the -->
<!-- same graph using just one of the measured quantities. For example, -->
<!-- take average measured latencies and use Little's Law to compute -->
<!-- expected average throughput for a particular load point (expressed as -->
<!-- total number of tasks submitted to the system over the -->
<!-- experiment). Conversely, you can take the average measured througput -->
<!-- and compute the expected average latency. -->

<!-- Draw the computed and measured lines in the same graph and compare. Do -->
<!-- they match closely?  -->
