`scany v1.0.0` vs local performance branches of `sqlh` and `set`.

Package `set` is included because it is the `sqlh` dependency that provides the `reflect` gateway into Go structs.

## Method

Benchmarks are run and compared with:

```bash
# windows/amd64 Intel(R) Core(TM) i7-7700K CPU @ 4.20GHz
# go version go1.16 windows/amd64

go test -run abcxyz -benchmem -bench "Select/scany" -count 5 > perf-scany.txt  # scany
go test -run abcxyz -benchmem -bench "Select/sqlh" -count 5 > perf-new.txt  # new local branches

# modify files step (see below)
benchstat perf-scany.txt perf-new.txt
```

`benchstat` can only compare benchmarks if they have the same name. Therefore in the files `perf-scany.txt` and `perf-new.txt` I replace all `scany_` with empty string and do the same for `sqlh_`. Then `benchstat` can compare the benchmarks.

## `ns/op`

CPU performance is roughly comparable with `sqlh` winning out slightly with SQLite and a bit more against Postgres. I'm not sure why `sqlh` seems to win out consistently; it could be that both Postgres and SQLite were behaving "a bit better" when I ran benchmarks for `sqlh`. It could also be that `scany` has inefficiencies based on the underlying database; I really don't see how that would be the case.

However I ran the `scany` benchmarks under the same conditions as `sqlh` and `sqlx` and -- for example -- comparing `sqlx` and `sqlh` is much closer with neither winning out against the other consistently like we see below.

I would be curious if someone other than me produces similar results with `scany` vs `sqlh`.

```
name                  old time/op    new time/op    delta
Libpq/5_rows-8           194µs ± 6%     167µs ± 1%  -13.65%  (p=0.016 n=5+4)
Libpq/50_rows-8          358µs ± 9%     318µs ± 9%  -11.13%  (p=0.016 n=5+5)
Libpq/100_rows-8         589µs ±14%     424µs ± 5%  -27.99%  (p=0.008 n=5+5)
Libpq/500_rows-8        3.30ms ±48%    1.44ms ±11%  -56.24%  (p=0.008 n=5+5)
Libpq/1000_rows-8       4.13ms ±42%    2.86ms ± 3%  -30.88%  (p=0.008 n=5+5)
Sqlite/5_rows-8          114µs ± 1%     104µs ± 1%   -8.13%  (p=0.008 n=5+5)
Sqlite/50_rows-8         349µs ± 1%     326µs ± 1%   -6.52%  (p=0.008 n=5+5)
Sqlite/100_rows-8        607µs ± 0%     574µs ± 0%   -5.58%  (p=0.008 n=5+5)
Sqlite/500_rows-8       2.71ms ± 1%    2.58ms ± 2%   -4.92%  (p=0.008 n=5+5)
Sqlite/1000_rows-8      5.33ms ± 0%    5.11ms ± 0%   -4.14%  (p=0.008 n=5+5)
Sqlmock/5_rows-8         299µs ±65%     321µs ±67%     ~     (p=0.841 n=5+5)
Sqlmock/50_rows-8        588µs ±18%     597µs ±15%     ~     (p=1.000 n=5+5)
Sqlmock/100_rows-8       791µs ± 9%     777µs ± 8%     ~     (p=0.841 n=5+5)
Sqlmock/500_rows-8       961µs ± 5%     944µs ± 7%     ~     (p=0.548 n=5+5)
Sqlmock/1000_rows-8     1.10ms ± 6%    1.08ms ± 6%     ~     (p=0.548 n=5+5)
Sqlmock/10000_rows-8    1.30ms ± 5%    1.23ms ± 5%     ~     (p=0.056 n=5+5)
```

## `allocated memory`

Against Postgres and SQLite `sqlh` once again is slightly better. `scany` has the advantage against Sqlmock. This is interesting but somewhat irrelevant as its not a real database engine; I include it because Sqlmock does not have the same IO requirements as a real database and I think it provides a better comparison of performance regarding the `reflect` mapping of data to Go structs. So in that regard `sqlh` can be improved.

```
name                  old alloc/op   new alloc/op   delta
Libpq/5_rows-8          5.39kB ± 1%    3.60kB ± 0%  -33.32%  (p=0.008 n=5+5)
Libpq/50_rows-8         28.6kB ± 0%    23.7kB ± 0%  -17.29%  (p=0.008 n=5+5)
Libpq/100_rows-8        54.4kB ± 1%    46.1kB ± 0%  -15.41%  (p=0.016 n=5+4)
Libpq/500_rows-8         266kB ± 1%     228kB ± 0%  -14.49%  (p=0.008 n=5+5)
Libpq/1000_rows-8        532kB ± 0%     458kB ± 0%  -13.89%  (p=0.008 n=5+5)
Sqlite/5_rows-8         8.72kB ± 0%    6.91kB ± 0%  -20.75%  (p=0.008 n=5+5)
Sqlite/50_rows-8        64.0kB ± 0%    59.0kB ± 0%   -7.81%  (p=0.008 n=5+5)
Sqlite/100_rows-8        125kB ± 0%     117kB ± 0%   -6.66%  (p=0.008 n=5+5)
Sqlite/500_rows-8        620kB ± 0%     583kB ± 0%   -5.94%  (p=0.008 n=5+5)
Sqlite/1000_rows-8      1.24MB ± 0%    1.17MB ± 0%   -5.82%  (p=0.008 n=5+5)
Sqlmock/5_rows-8        3.11kB ± 1%    4.08kB ± 1%  +31.19%  (p=0.008 n=5+5)
Sqlmock/50_rows-8       3.12kB ± 1%    4.10kB ± 1%  +31.56%  (p=0.008 n=5+5)
Sqlmock/100_rows-8      3.14kB ± 1%    4.09kB ± 1%  +30.31%  (p=0.008 n=5+5)
Sqlmock/500_rows-8      3.32kB ± 2%    4.19kB ± 1%  +26.14%  (p=0.008 n=5+5)
Sqlmock/1000_rows-8     3.52kB ± 2%    4.36kB ± 1%  +23.76%  (p=0.008 n=5+5)
Sqlmock/10000_rows-8    8.94kB ±10%    7.53kB ± 2%  -15.79%  (p=0.008 n=5+5)
```

## `number of allocations`

In general `scany` makes fewer total allocations.

```
name                  old allocs/op  new allocs/op  delta
Libpq/5_rows-8             192 ± 0%       107 ± 0%  -44.27%  (p=0.008 n=5+5)
Libpq/50_rows-8            908 ± 0%       871 ± 0%   -4.07%  (p=0.008 n=5+5)
Libpq/100_rows-8         1.71k ± 0%     1.72k ± 0%   +0.82%  (p=0.008 n=5+5)
Libpq/500_rows-8         8.76k ± 0%     9.18k ± 0%   +4.74%  (p=0.008 n=5+5)
Libpq/1000_rows-8        17.8k ± 0%     18.7k ± 0%   +5.15%  (p=0.008 n=5+5)
Sqlite/5_rows-8            358 ± 0%       273 ± 0%  -23.74%  (p=0.008 n=5+5)
Sqlite/50_rows-8         2.47k ± 0%     2.44k ± 0%   -1.50%  (p=0.008 n=5+5)
Sqlite/100_rows-8        4.82k ± 0%     4.84k ± 0%   +0.29%  (p=0.008 n=5+5)
Sqlite/500_rows-8        24.3k ± 0%     24.7k ± 0%   +1.71%  (p=0.008 n=5+5)
Sqlite/1000_rows-8       48.8k ± 0%     49.7k ± 0%   +1.88%  (p=0.008 n=5+5)
Sqlmock/5_rows-8          40.0 ± 0%      43.0 ± 0%   +7.50%  (p=0.008 n=5+5)
Sqlmock/50_rows-8         40.0 ± 0%      43.0 ± 0%   +7.50%  (p=0.008 n=5+5)
Sqlmock/100_rows-8        40.0 ± 0%      43.0 ± 0%   +7.50%  (p=0.008 n=5+5)
Sqlmock/500_rows-8        41.0 ± 0%      44.0 ± 0%   +7.32%  (p=0.008 n=5+5)
Sqlmock/1000_rows-8       42.4 ± 1%      46.4 ± 1%   +9.43%  (p=0.008 n=5+5)
Sqlmock/10000_rows-8      74.8 ± 7%      89.4 ± 2%  +19.52%  (p=0.008 n=5+5)
```
