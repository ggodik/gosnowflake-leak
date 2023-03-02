# gosnowflake-leak

Memory leak in Arrow handling for the gosnowflake driver as of v1.6.18

Code just reads ArrowBatches and releases

## Reproduce

Run the following command by connecting to snowflake via `DSN` env var

```bash
DSN=user:pass@account/db?warehouse=warehouse ./gosnowflake-leak.exe -rows 1000000
```

## Parameters

 * `DSN` - env var
 * -rows number of rows to generate for the `TABLE(GENERATOR)` query

## Sample

```
$ DSN=<HIDDEN>./gosnowflake-leak.exe -rows 10000000
2023/03/02 00:18:16 query=SELECT seq4() A, uniform(1, 10, RANDOM(12)) B, uniform(1, 100, RANDOM(121)) C, randstr(15, random()), randstr(25, random()), randstr(35, random())FROM TABLE(GENERATOR(ROWCOUNT=>10000000))
2023/03/02 00:18:21 memory usage alloc=2.9 MB totalAlloc=5.9 MB system=13.8 MB numGC=2
2023/03/02 00:18:26 memory usage alloc=2.9 MB totalAlloc=5.9 MB system=13.8 MB numGC=2
2023/03/02 00:18:31 memory usage alloc=2.9 MB totalAlloc=5.9 MB system=13.8 MB numGC=2
2023/03/02 00:18:36 memory usage alloc=2.9 MB totalAlloc=5.9 MB system=13.8 MB numGC=2
2023/03/02 00:18:41 memory usage alloc=2.9 MB totalAlloc=5.9 MB system=13.8 MB numGC=2
2023/03/02 00:18:46 memory usage alloc=2.9 MB totalAlloc=5.9 MB system=13.8 MB numGC=2
2023/03/02 00:18:51 memory usage alloc=8.3 MB totalAlloc=14.4 MB system=22.7 MB numGC=4
2023/03/02 00:18:56 memory usage alloc=41.1 MB totalAlloc=69.4 MB system=53.0 MB numGC=8
2023/03/02 00:19:01 memory usage alloc=67.2 MB totalAlloc=129.2 MB system=96.6 MB numGC=10
2023/03/02 00:19:06 memory usage alloc=100.6 MB totalAlloc=194.1 MB system=139.9 MB numGC=11
2023/03/02 00:19:11 memory usage alloc=172.8 MB totalAlloc=266.3 MB system=187.5 MB numGC=11
2023/03/02 00:19:16 memory usage alloc=186.6 MB totalAlloc=326.3 MB system=204.8 MB numGC=12
2023/03/02 00:19:21 memory usage alloc=246.4 MB totalAlloc=386.1 MB system=265.4 MB numGC=12
2023/03/02 00:19:26 memory usage alloc=249.7 MB totalAlloc=458.3 MB system=295.7 MB numGC=13
2023/03/02 00:19:31 memory usage alloc=317.1 MB totalAlloc=525.6 MB system=339.3 MB numGC=13
2023/03/02 00:19:36 memory usage alloc=379.3 MB totalAlloc=587.9 MB system=404.2 MB numGC=13
2023/03/02 00:19:41 memory usage alloc=341.3 MB totalAlloc=652.6 MB system=434.6 MB numGC=14
2023/03/02 00:19:46 memory usage alloc=398.8 MB totalAlloc=710.1 MB system=434.6 MB numGC=14
2023/03/02 00:19:51 memory usage alloc=463.5 MB totalAlloc=774.8 MB system=490.8 MB numGC=14
2023/03/02 00:19:56 memory usage alloc=424.3 MB totalAlloc=844.6 MB system=556.0 MB numGC=15
2023/03/02 00:20:01 memory usage alloc=491.6 MB totalAlloc=911.9 MB system=556.0 MB numGC=15
2023/03/02 00:20:06 memory usage alloc=556.3 MB totalAlloc=976.6 MB system=586.4 MB numGC=15
2023/03/02 00:20:11 memory usage alloc=618.7 MB totalAlloc=1.0 GB system=651.3 MB numGC=15
2023/03/02 00:20:16 memory usage alloc=676.1 MB totalAlloc=1.1 GB system=712.0 MB numGC=15
2023/03/02 00:20:21 memory usage alloc=743.3 MB totalAlloc=1.2 GB system=781.3 MB numGC=15
2023/03/02 00:20:26 memory usage alloc=803.1 MB totalAlloc=1.2 GB system=841.9 MB numGC=15
2023/03/02 00:20:31 memory usage alloc=650.2 MB totalAlloc=1.3 GB system=867.8 MB numGC=16
2023/03/02 00:20:36 memory usage alloc=722.4 MB totalAlloc=1.4 GB system=867.8 MB numGC=16
2023/03/02 00:20:41 memory usage alloc=784.7 MB totalAlloc=1.4 GB system=867.9 MB numGC=16
2023/03/02 00:20:46 memory usage alloc=854.5 MB totalAlloc=1.5 GB system=898.5 MB numGC=16
2023/03/02 00:20:51 memory usage alloc=921.8 MB totalAlloc=1.5 GB system=963.4 MB numGC=16
2023/03/02 00:20:56 memory usage alloc=979.2 MB totalAlloc=1.6 GB system=1.0 GB numGC=16
2023/03/02 00:21:01 memory usage alloc=1.0 GB totalAlloc=1.7 GB system=1.1 GB numGC=16
2023/03/02 00:21:06 memory usage alloc=879.7 MB totalAlloc=1.7 GB system=1.1 GB numGC=17
2023/03/02 00:21:11 memory usage alloc=939.5 MB totalAlloc=1.8 GB system=1.1 GB numGC=17
2023/03/02 00:21:16 memory usage alloc=994.5 MB totalAlloc=1.8 GB system=1.1 GB numGC=17
2023/03/02 00:21:21 memory usage alloc=1.1 GB totalAlloc=1.9 GB system=1.1 GB numGC=17
2023/03/02 00:21:26 memory usage alloc=1.1 GB totalAlloc=2.0 GB system=1.2 GB numGC=17
```

