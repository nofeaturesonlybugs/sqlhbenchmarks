`sqlx v1.3.4` vs local performance branches of `sqlh` and `set`.

Package `set` is included because it is the `sqlh` dependency that provides the `reflect` gateway into Go structs.

## Method

Benchmarks are run and compared with:

```bash
# windows/amd64 Intel(R) Core(TM) i7-7700K CPU @ 4.20GHz
# go version go1.16 windows/amd64

go test -run abcxyz -benchmem -bench "Select/sqlx" -count 5 > perf-sqlx.txt  # sqlx
go test -run abcxyz -benchmem -bench "Select/sqlh" -count 5 > perf-new.txt  # new local branches

# modify files step (see below)
benchstat perf-sqlx.txt perf-new.txt
```

`benchstat` can only compare benchmarks if they have the same name. Therefore in the files `perf-sqlx.txt` and `perf-new.txt` I replace all `sqlx_` with empty string and do the same for `sqlh_`. Then `benchstat` can compare the benchmarks.

## `ns/op`

CPU performance is nearly identical. Occasionally one slightly outperforms the other.

```
name                  old time/op    new time/op    delta
Libpq/5_rows-8           180µs ± 1%     167µs ± 1%   -7.08%  (p=0.016 n=5+4)
Libpq/50_rows-8          314µs ± 5%     318µs ± 9%     ~     (p=1.000 n=5+5)
Libpq/100_rows-8         478µs ±19%     424µs ± 5%     ~     (p=0.310 n=5+5)
Libpq/500_rows-8        2.19ms ±47%    1.44ms ±11%     ~     (p=0.690 n=5+5)
Libpq/1000_rows-8       2.72ms ±25%    2.86ms ± 3%     ~     (p=1.000 n=5+5)
Sqlite/5_rows-8          103µs ± 1%     104µs ± 1%   +1.86%  (p=0.016 n=5+5)
Sqlite/50_rows-8         322µs ± 1%     326µs ± 1%   +1.39%  (p=0.008 n=5+5)
Sqlite/100_rows-8        563µs ± 1%     574µs ± 0%   +1.89%  (p=0.008 n=5+5)
Sqlite/500_rows-8       2.50ms ± 0%    2.58ms ± 2%   +3.05%  (p=0.008 n=5+5)
Sqlite/1000_rows-8      4.95ms ± 2%    5.11ms ± 0%   +3.11%  (p=0.008 n=5+5)
Sqlmock/5_rows-8         306µs ±66%     321µs ±67%     ~     (p=0.690 n=5+5)
Sqlmock/50_rows-8        579µs ±18%     597µs ±15%     ~     (p=0.841 n=5+5)
Sqlmock/100_rows-8       773µs ± 9%     777µs ± 8%     ~     (p=1.000 n=5+5)
Sqlmock/500_rows-8       937µs ± 6%     944µs ± 7%     ~     (p=0.841 n=5+5)
Sqlmock/1000_rows-8     1.10ms ± 7%    1.08ms ± 6%     ~     (p=0.690 n=5+5)
Sqlmock/10000_rows-8    1.23ms ± 4%    1.23ms ± 5%     ~     (p=0.841 n=5+5)
```

## `allocated memory`

`sqlx` is marginally lower. Note however that `sqlh` has quite a number _more_ total allocations and I plan to lower these in a future release with pooling.

```
name                  old alloc/op   new alloc/op   delta
Libpq/5_rows-8          3.46kB ± 0%    3.60kB ± 0%   +3.86%  (p=0.008 n=5+5)
Libpq/50_rows-8         23.2kB ± 1%    23.7kB ± 0%   +2.06%  (p=0.008 n=5+5)
Libpq/100_rows-8        45.5kB ± 0%    46.1kB ± 0%   +1.32%  (p=0.016 n=5+4)
Libpq/500_rows-8         225kB ± 1%     228kB ± 0%   +1.03%  (p=0.008 n=5+5)
Libpq/1000_rows-8        452kB ± 1%     458kB ± 0%   +1.33%  (p=0.008 n=5+5)
Sqlite/5_rows-8         6.79kB ± 0%    6.91kB ± 0%   +1.81%  (p=0.008 n=5+5)
Sqlite/50_rows-8        58.7kB ± 0%    59.0kB ± 0%   +0.54%  (p=0.008 n=5+5)
Sqlite/100_rows-8        116kB ± 0%     117kB ± 0%   +0.94%  (p=0.008 n=5+5)
Sqlite/500_rows-8        582kB ± 0%     583kB ± 0%   +0.20%  (p=0.008 n=5+5)
Sqlite/1000_rows-8      1.17MB ± 0%    1.17MB ± 0%   +0.26%  (p=0.008 n=5+5)
Sqlmock/5_rows-8        3.63kB ± 1%    4.08kB ± 1%  +12.52%  (p=0.008 n=5+5)
Sqlmock/50_rows-8       3.67kB ± 0%    4.10kB ± 1%  +11.90%  (p=0.008 n=5+5)
Sqlmock/100_rows-8      3.65kB ± 2%    4.09kB ± 1%  +12.02%  (p=0.008 n=5+5)
Sqlmock/500_rows-8      3.72kB ± 1%    4.19kB ± 1%  +12.70%  (p=0.008 n=5+5)
Sqlmock/1000_rows-8     3.84kB ± 1%    4.36kB ± 1%  +13.48%  (p=0.008 n=5+5)
Sqlmock/10000_rows-8    6.63kB ± 3%    7.53kB ± 2%  +13.58%  (p=0.008 n=5+5)
```

## `number of allocations`

`sqlx` makes fewer allocations. Despite this the total allocated memory between the two is roughly comparable and I have plans to decrease allocations via pooling.

```
name                  old allocs/op  new allocs/op  delta
Libpq/5_rows-8            92.0 ± 0%     107.0 ± 0%  +16.30%  (p=0.008 n=5+5)
Libpq/50_rows-8            763 ± 0%       871 ± 0%  +14.15%  (p=0.008 n=5+5)
Libpq/100_rows-8         1.51k ± 0%     1.72k ± 0%  +13.80%  (p=0.008 n=5+5)
Libpq/500_rows-8         8.17k ± 0%     9.18k ± 0%  +12.39%  (p=0.016 n=4+5)
Libpq/1000_rows-8        16.7k ± 0%     18.7k ± 0%  +12.08%  (p=0.008 n=5+5)
Sqlite/5_rows-8            258 ± 0%       273 ± 0%   +5.81%  (p=0.008 n=5+5)
Sqlite/50_rows-8         2.33k ± 0%     2.44k ± 0%   +4.64%  (p=0.008 n=5+5)
Sqlite/100_rows-8        4.63k ± 0%     4.84k ± 0%   +4.52%  (p=0.008 n=5+5)
Sqlite/500_rows-8        23.7k ± 0%     24.7k ± 0%   +4.27%  (p=0.008 n=5+5)
Sqlite/1000_rows-8       47.7k ± 0%     49.7k ± 0%   +4.22%  (p=0.008 n=5+5)
Sqlmock/5_rows-8          39.0 ± 0%      43.0 ± 0%  +10.26%  (p=0.008 n=5+5)
Sqlmock/50_rows-8         39.0 ± 0%      43.0 ± 0%  +10.26%  (p=0.008 n=5+5)
Sqlmock/100_rows-8        39.0 ± 0%      43.0 ± 0%  +10.26%  (p=0.008 n=5+5)
Sqlmock/500_rows-8        39.0 ± 0%      44.0 ± 0%  +12.82%  (p=0.008 n=5+5)
Sqlmock/1000_rows-8       40.4 ± 1%      46.4 ± 1%  +14.85%  (p=0.008 n=5+5)
Sqlmock/10000_rows-8      62.0 ± 2%      89.4 ± 2%  +44.19%  (p=0.008 n=5+5)
```
