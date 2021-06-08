# psb - promscale benchmarking
psb is small tool to test query execution performance of Promscale

`psb` simply spawns multiple workers as instructed by user to send queries listed in input file to Promscale server. If there are multiple queries in input file, queries are benchmarked one after another. After finishing benchmark of one query, stats are reported and then `psb` moves to next query.

At any point in the execution, user can hit Ctrl-C to stop benchmark. It still displayes report of results collected in that time.

## Installation
Installation of `psb` is very easy. Just `go get` the package and you will have `psb` binary installed in go bin directory.

`go get github.com/siddharth178/psb`

## Usage

`psb` accepts multiple options to control how you want to run the benchmark. Queries to benchmark should be provided in a csv file with `|` as field separator. Each line should contain four fields - query, start time, end time, step.

`psb` assumes Promscale server is running at given server URL. Please see Misc section below to know how I setup Promscale locally.

Here is sample input file
```
cat obs-queries.csv 
demo_cpu_usage_seconds_total{mode="idle"}|1597056698698|1597059548699|15000
avg by(instance) (demo_cpu_usage_seconds_total)|1597057698698|1597058548699|60000
avg without(instance, mode) (demo_cpu_usage_seconds_total)|1597056698698|1597059548699|120000
topk by(instance) (2, demo_cpu_usage_seconds_total)|1597058698698|1597059548699|60000
avg(max by(mode) (demo_cpu_usage_seconds_total))|1597055698698|1597057548699|600000
```

### Sample run of psb
```
$ ~/go/bin/psb --help
Usage of /home/siddharth/go/bin/psb:
  -cpus int
    	number of CPUs to use for go runtime (default 12)
  -f string
    	queries will be read from this file (default "obs-queries.csv")
  -i int
    	number of query executions (default 10)
  -n int
    	number of workers (default 1)
  -t int
    	client timeout (in ms) (default 10000)
  -url string
    	Promscale server URL (default "http://localhost:9201")

$ ~/go/bin/psb -f temp.csv 
psb - benchmarking Promscale
Found 3 queries in input file: temp.csv

benchmark config
server url: http://localhost:9201
# concurrent workers/clients: 1
# iterations of each query: 10
# CPUs to use: 12
# client timeout: 10s

2021/06/08 11:10:56 Promscale health OK, moving to benchmark
2021/06/08 11:10:56 Ctrl-C handler registerd
2021/06/08 11:10:57 ready to benchmark
# 0 query: demo_cpu_usage_seconds_total{mode="idle"}|1597056698698|1597059548699|15000
2021/06/08 11:10:57 waiting for all results
2021/06/08 11:10:57 wait over, work is done, no more pending results
2021/06/08 11:10:57 report is ready: 10
--------
query str: demo_cpu_usage_seconds_total{mode="idle"}|1597056698698|1597059548699|15000
#queries: 10
#failed queries: 0
total time: 49.544202ms
min time: 3.612998ms
max time: 9.679498ms
avg time: 4.837793ms
med time: 4.601731ms
--------

# 1 query: avg by(instance) (demo_cpu_usage_seconds_total)|1597057698698|1597058548699|60000
2021/06/08 11:10:57 waiting for all results
2021/06/08 11:10:57 wait over, work is done, no more pending results
2021/06/08 11:10:57 report is ready: 10
--------
query str: avg by(instance) (demo_cpu_usage_seconds_total)|1597057698698|1597058548699|60000
#queries: 10
#failed queries: 0
total time: 46.540355ms
min time: 3.66021ms
max time: 6.269597ms
avg time: 4.533881ms
med time: 4.405724ms
--------

# 2 query: avg without(instance, mode) (demo_cpu_usage_seconds_total)|1597056698698|1597059548699|120000
2021/06/08 11:10:57 waiting for all results
2021/06/08 11:10:57 wait over, work is done, no more pending results
2021/06/08 11:10:57 report is ready: 10
--------
query str: avg without(instance, mode) (demo_cpu_usage_seconds_total)|1597056698698|1597059548699|120000
#queries: 10
#failed queries: 0
total time: 64.266125ms
min time: 5.59759ms
max time: 7.226903ms
avg time: 6.307687ms
med time: 6.253603ms
--------
```

## Misc
### Instructions followed to setup local Promscale and upload sample data
I created small scripts to setup, start and stop timescale-db with promscale. It its lot easier to work with these than using docker commands.

Call `setup.sh` to download and run timescaledb with promscale extension. It also creates necessary network config in docker.
```
siddharth@siddharth-ubuntu:~/source/psb$ ./setup.sh 
Importing the environment variables
Destroying existing database and docker containers
timescaledb
Creating docker network bridge (ignore error, if run the second time)
Error response from daemon: network with name promscale-timescaledb already exists
Running docker image
cccf0450bf7f2097afa33e6237e2dc229c2812011d802e93f78002fe6f75800d

siddharth@siddharth-ubuntu:~/source/psb$ docker ps
CONTAINER ID   IMAGE                                          COMMAND                  CREATED         STATUS         PORTS                                       NAMES
cccf0450bf7f   timescaledev/promscale-extension:latest-pg12   "docker-entrypoint.s…"   3 seconds ago   Up 2 seconds   0.0.0.0:5433->5432/tcp, :::5433->5432/tcp   timescaledb

siddharth@siddharth-ubuntu:~/source/psb$ docker network ls
NETWORK ID     NAME                    DRIVER    SCOPE
41eea213c57c   promscale-timescaledb   bridge    local
```

Once done, one can login to this instance using `psql` like following. Password is available in `config.sh`
```
iddharth@siddharth-ubuntu:~/source/psb$ psql -U postgres -h localhost -p 5433
Password for user postgres: 
psql (12.7 (Ubuntu 12.7-0ubuntu0.20.04.1), server 12.4)
Type "help" for help.

postgres=# 
```

Download latest promscale binary from project release page. One can now run promscale server with following command
```
siddharth@siddharth-ubuntu:~/tools/bin$ promscale-0.4.1 --db-name postgres --db-password password --db-ssl-mode allow --db-host localhost --db-port 5433
level=info ts=2021-06-05T07:25:22.797Z caller=runner.go:29 msg="Version:0.4.1; Commit Hash: "
level=info ts=2021-06-05T07:25:22.798Z caller=runner.go:30 config="&{ListenAddr::9201 PgmodelCfg:{CacheConfig:{SeriesCacheInitialSize:250000 seriesCacheMemoryMaxFlag:{kind:0 value:50} SeriesCacheMemoryMaxBytes:6658123366 MetricsCacheSize:10000 LabelsCacheSize:10000} AppName:promscale@0.4.1 Host:localhost Port:5433 User:postgres password:**** Database:postgres SslMode:allow DbConnectRetries:0 DbConnectionTimeout:1m0s IgnoreCompressedChunks:false AsyncAcks:false ReportInterval:0 WriteConnectionsPerProc:4 MaxConnections:-1 UsesHA:false DbUri: EnableStatementsCache:true} LogCfg:{Level:info Format:logfmt} APICfg:{AllowedOrigin:^(?:.*)$ ReadOnly:false HighAvailability:false AdminAPIEnabled:false TelemetryPath:/metrics Auth:0xc000284d70 MultiTenancy:<nil> EnableFeatures: EnabledFeaturesList:[] MaxQueryTimeout:2m0s SubQueryStepInterval:1m0s LookBackDelta:5m0s MaxSamples:50000000 MaxPointsPerTs:11000} LimitsCfg:{targetMemoryFlag:{kind:0 value:80} TargetMemoryBytes:13316246732} TenancyCfg:{SkipTenantValidation:false EnableMultiTenancy:false AllowNonMTWrites:false ValidTenantsStr:allow-all ValidTenantsList:[]} ConfigFile:config.yml TLSCertFile: TLSKeyFile: HaGroupLockID:0 PrometheusTimeout:-1ns ElectionInterval:5s Migrate:true StopAfterMigrate:false UseVersionLease:true InstallExtensions:true UpgradeExtensions:true UpgradePrereleaseExtensions:false}"
level=info ts=2021-06-05T07:25:23.039Z caller=extension.go:225 msg="successfully created promscale extension at version 0.1.1"
level=warn ts=2021-06-05T07:25:23.062Z caller=client.go:148 msg="No adapter leader election. Group lock id is not set. Possible duplicate write load if running multiple connectors"
level=warn ts=2021-06-05T07:25:23.070Z caller=config.go:169 msg="had to reduce the number of copiers due to connection limits: wanted 48, reduced to 25"
level=info ts=2021-06-05T07:25:23.072Z caller=client.go:116 msg="application_name=promscale@0.4.1 host=localhost port=5433 user=postgres dbname=postgres password='****' sslmode=allow connect_timeout=60" numCopiers=25 pool_max_conns=50 pool_min_conns=12 statement_cache="512 statements"
level=info ts=2021-06-05T07:25:23.108Z caller=runner.go:56 msg="Starting up..."
level=info ts=2021-06-05T07:25:23.108Z caller=runner.go:57 msg=Listening addr=:9201
```

We can now see local tsdb instance getting utilized
```
siddharth@siddharth-ubuntu:~/source/psb$ psql -U postgres -h localhost -p 5433
Password for user postgres: 
psql (12.7 (Ubuntu 12.7-0ubuntu0.20.04.1), server 12.4)
Type "help" for help.

postgres=# \d
                           List of relations
    Schema     |          Name           |       Type        |  Owner   
---------------+-------------------------+-------------------+----------
 _prom_catalog | default                 | table             | postgres
 _prom_catalog | ha_leases               | table             | postgres
 _prom_catalog | ha_leases_logs          | table             | postgres
 _prom_catalog | ids_epoch               | table             | postgres
 _prom_catalog | label                   | table             | postgres
 _prom_catalog | label_id_seq            | sequence          | postgres
 _prom_catalog | label_key               | table             | postgres
 _prom_catalog | label_key_id_seq        | sequence          | postgres
 _prom_catalog | label_key_position      | table             | postgres
 _prom_catalog | metric                  | table             | postgres
 _prom_catalog | metric_id_seq           | sequence          | postgres
 _prom_catalog | remote_commands         | table             | postgres
 _prom_catalog | remote_commands_seq_seq | sequence          | postgres
 _prom_catalog | series                  | partitioned table | postgres
 _prom_catalog | series_id               | sequence          | postgres
 public        | prom_installation_info  | table             | postgres
 public        | prom_schema_migrations  | table             | postgres
(17 rows)

postgres=# 
```

Adding sample data to promscale
```
curl -v \
-H "Content-Type: application/x-protobuf" \
-H "Content-Encoding: snappy" \
-H "X-Prometheus-Remote-Write-Version: 0.1.0" \
--data-binary "@real-dataset.sz" \
"http://localhost:9201/write"

*   Trying 127.0.0.1:9201...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 9201 (#0)
> POST /write HTTP/1.1
> Host: localhost:9201
> User-Agent: curl/7.68.0
> Accept: */*
> Content-Type: application/x-protobuf
> Content-Encoding: snappy
> X-Prometheus-Remote-Write-Version: 0.1.0
> Content-Length: 4715088
> Expect: 100-continue
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 100 Continue
* We are completely uploaded and fine
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Date: Sat, 05 Jun 2021 08:09:16 GMT
< Content-Length: 0
< 
* Connection #0 to host localhost left intact
```

Seeing data in database
```
postgres=# \dn
          List of schemas
          Name           |  Owner   
-------------------------+----------
 _prom_catalog           | postgres
 _prom_ext               | postgres
 _timescaledb_cache      | postgres
 _timescaledb_catalog    | postgres
 _timescaledb_config     | postgres
 _timescaledb_internal   | postgres
 prom_api                | postgres
 prom_data               | postgres
 prom_data_series        | postgres
 prom_info               | postgres
 prom_metric             | postgres
 prom_series             | postgres
 public                  | postgres
 timescaledb_information | postgres
(14 rows)


postgres=# select * from prom_metric.go_threads limit 5;
            time            | value | series_id |   labels   | instance_id | job_id 
----------------------------+-------+-----------+------------+-------------+--------
 2020-08-10 10:34:58.698+00 |    13 |        55 | {45,37,48} |          37 |     48
 2020-08-10 10:35:03.699+00 |    13 |        55 | {45,37,48} |          37 |     48
 2020-08-10 10:35:08.699+00 |    13 |        55 | {45,37,48} |          37 |     48
 2020-08-10 10:35:13.699+00 |    13 |        55 | {45,37,48} |          37 |     48
 2020-08-10 10:35:18.699+00 |    13 |        55 | {45,37,48} |          37 |     48
(5 rows)
```

### Querying Promscale
Promscale supports multiple types of Prometheus queries. But all the queries in `obs-queries.csv` are of range type. May be we should support range queries first.

Here is sample run of instant query on local setup
```
siddharth@siddharth-ubuntu:~/source/psb$ curl "http://localhost:9201/api/v1/query?query=up&time=2020-08-10T10:34:58.698Z" | jq .
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   279  100   279    0     0   136k      0 --:--:-- --:--:-- --:--:--  136k
{
  "status": "success",
  "data": {
    "resultType": "vector",
    "result": [
      {
        "metric": {
          "__name__": "up",
          "instance": "demo.promlabs.com:10000",
          "job": "demo"
        },
        "value": [
          1597055698.698,
          "1"
        ]
      },
      {
        "metric": {
          "__name__": "up",
          "instance": "demo.promlabs.com:10002",
          "job": "demo"
        },
        "value": [
          1597055698.698,
          "1"
        ]
      }
    ]
  }
}
```

Here is sample run of range query (1st from obs-queries.csv)
```
siddharth@siddharth-ubuntu:~/source/psb$ curl -g "http://localhost:9201/api/v1/query_range?start=1597056698.698&end=1597059548.699&step=60000" --data-urlencode 'query=demo_cpu_usage_seconds_total{mode="idle"}' | jq .
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   624  100   567  100    57   110k  11400 --:--:-- --:--:-- --:--:--  121k
{
  "status": "success",
  "data": {
    "resultType": "matrix",
    "result": [
      {
        "metric": {
          "__name__": "demo_cpu_usage_seconds_total",
          "instance": "demo.promlabs.com:10000",
          "job": "demo",
          "mode": "idle"
        },
        "values": [
          [
            1597056698.698,
            "16496977.984534835"
          ]
        ]
      },
      {
        "metric": {
          "__name__": "demo_cpu_usage_seconds_total",
          "instance": "demo.promlabs.com:10001",
          "job": "demo",
          "mode": "idle"
        },
        "values": [
          [
            1597056698.698,
            "16497259.987334812"
          ]
        ]
      },
      {
        "metric": {
          "__name__": "demo_cpu_usage_seconds_total",
          "instance": "demo.promlabs.com:10002",
          "job": "demo",
          "mode": "idle"
        },
        "values": [
          [
            1597056698.698,
            "16497074.938270388"
          ]
        ]
      }
    ]
  }
}
```

### References
* How to query Promscale using PromQL
  - https://prometheus.io/docs/prometheus/latest/querying/api/
  - There are 2 types of queries. Instant queries and Range queries
* Prometheus HTTP API supported by Promscale
  - https://github.com/timescale/promscale/blob/master/docs/prometheus_api.md#implemented-endpoints



## Questions
