# Files for presentation "Concurrent Optimization Methods Using Go"
Written and presented by Brodie Roberts on behalf of Index Excahnge.

Presented at the Hack The 6ix workshop at UWaterloo on June 5, 2018.

## Setup

```
# Generate 10M data sample
mkdir -p data/10M
time go run generator.go -n 100000 > data/10M/server_1.log
time for i in {2..100}; do cp data/10M/server_1.log data/10M/server_$i.log; done

# Generate 300M data sample (>20 minutes)
mkdir -p data/300M
time for i in {1..300}; do go run generator.go -r -n 1000000 > data/300M/server_$i.log; done
```


## Baseline: Bash
```
# Lower Bound
time cat data/10M/* | wc -l

# Upper Bound
time cat data/10M/* | cut -f2 | grep -v "^0$" | sort -n | uniq -c | wc -l
time cat data/10M/* | grep -P "^.*\t0\t" | cut -f3,4,5 | sort -n | uniq -c | wc -l
```


## Version 1: Go
```
# Run basic program
go build V1.go && ./V1 data/10M

# Profile program
go build V1.go && ./V1 -prof data/10M
go tool pprof V1 prof.dat
go tool pprof -png V1 prof.dat > V1.png

# Strip-down testing of V1: add channel buffering
```


## Version 2: Channel Buffering
```
# Run V2
go build V2.go && ./V2 data/10M

# Profile program
go build V2.go && ./V2 -prof data/10M
go tool pprof -png V2 prof.dat > V2.png

# Add channel batching to V1
```


## Version 3: Channel Batching
```
# Run V3
go build V3.go && ./V3 -b 64 data/10M

# Profile program
go build V3.go && ./V3 -prof -b 64 data/10M
go tool pprof -png V3 prof.dat > V3.png

# Strip-down testing of V3: add concurrency
```


## Version 4: Parse Concurrency
```
# Run V4
go build V4.go && ./V4 -b 64 -c 4 data/10M

# Profile program
go build V4.go && ./V4 -prof -b 64 -c 4 data/10M
go tool pprof -png V4 prof.dat > V4.png

# Strip-down testing of V4: add MakeKey
```


## Version 5: Optimise Sequential Section
```
## Run V5
go build V5.go && ./V5 -b 64 -c 16 data/10M

# Profile program
go build V5.go && ./V5 -prof -b 64 -c 16 data/300M
go tool pprof -png V5 prof.dat > V5.png
```


### Slow Commands
```
# Versions of the above commands on the large sample data (several >20 minutes)
go run V1.go data/300M
go run V2.go data/300M
go run V3.go -b 64 data/300M
go run V4.go -b 64 -c 4 data/300M
go run V5.go -b 64 -c 16 data/300M
```

