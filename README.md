# go-thanksgiving

Google's thanksgiving-doodle in 24h, patched to run locally without crappengine.

See also: https://blog.golang.org/from-zero-to-go-launching-on-google

![Turkey](https://blog.golang.org/from-zero-to-go-launching-on-google_image03.png)

## Build and run

```
go build turkey.go
./turkey
Server running at 8080
```

Now visit http://localhost:8080/thumb/12345678 and play around :)

## Benchmark

```
ab -n 10000 -c 10

This is ApacheBench, Version 2.3 <$Revision: 1604373 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking localhost (be patient)
Completed 1000 requests
^C

Server Software:        
Server Hostname:        localhost
Server Port:            8080

Document Path:          /thumb/12345678
Document Length:        25472 bytes

Concurrency Level:      10
Time taken for tests:   4.725 seconds
Complete requests:      1323
Failed requests:        0
Total transferred:      33929171 bytes
HTML transferred:       33768604 bytes
Requests per second:    279.99 [#/sec] (mean)
Time per request:       35.715 [ms] (mean)
Time per request:       3.572 [ms] (mean, across all concurrent requests)
Transfer rate:          7012.31 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.1      0       1
Processing:    11   35  17.7     32     124
Waiting:        4   24  17.1     21     114
Total:         11   35  17.7     33     124

Percentage of the requests served within a certain time (ms)
  50%     33
  66%     41
  75%     46
  80%     50
  90%     60
  95%     67
  98%     77
  99%     86
 100%    124 (longest request)
```
