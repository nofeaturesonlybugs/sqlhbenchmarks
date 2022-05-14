`database/sql` (stdlib) vs local performance branches of `sqlh` and `set`.

Package `set` is included because it is the `sqlh` dependency that provides the `reflect` gateway into Go structs.

## Method

Benchmarks are run and compared with:

```bash
# windows/amd64 Intel(R) Core(TM) i7-7700K CPU @ 4.20GHz
# go version go1.16 windows/amd64

go test -run abcxyz -benchmem -bench "Select/database/sql" -count 5 > perf-stdlib.txt  # stdlib
go test -run abcxyz -benchmem -bench "Select/sqlh" -count 5 > perf-new.txt  # new local branches

# modify files step (see below)
benchstat perf-stdlib.txt perf-new.txt
```

`benchstat` can only compare benchmarks if they have the same name. Therefore in the files `perf-stdlib.txt` and `perf-new.txt` I replace all `database/sql_` with empty string and do the same for `sqlh_`. Then `benchstat` can compare the benchmarks.

## `ns/op`

`sqlh` performs roughly on par with standard library `database/sql`.

```
name                  old time/op    new time/op    delta
Libpq/5_rows-8           167µs ± 2%     167µs ± 1%     ~     (p=0.905 n=5+4)
Libpq/50_rows-8          367µs ± 8%     318µs ± 9%  -13.22%  (p=0.032 n=5+5)
Libpq/100_rows-8         541µs ± 8%     424µs ± 5%  -21.63%  (p=0.008 n=5+5)
Libpq/500_rows-8        1.98ms ±19%    1.44ms ±11%  -26.89%  (p=0.008 n=5+5)
Libpq/1000_rows-8       3.13ms ±10%    2.86ms ± 3%     ~     (p=0.095 n=5+5)
Sqlite/5_rows-8         99.4µs ± 1%   104.4µs ± 1%   +5.03%  (p=0.008 n=5+5)
Sqlite/50_rows-8         305µs ± 0%     326µs ± 1%   +6.89%  (p=0.008 n=5+5)
Sqlite/100_rows-8        534µs ± 1%     574µs ± 0%   +7.40%  (p=0.008 n=5+5)
Sqlite/500_rows-8       2.38ms ± 0%    2.58ms ± 2%   +8.26%  (p=0.008 n=5+5)
Sqlite/1000_rows-8      4.71ms ± 0%    5.11ms ± 0%   +8.52%  (p=0.008 n=5+5)
Sqlmock/5_rows-8         461µs ±37%     321µs ±67%     ~     (p=0.095 n=5+5)
Sqlmock/50_rows-8        698µs ±11%     597µs ±15%     ~     (p=0.056 n=5+5)
Sqlmock/100_rows-8       884µs ± 9%     777µs ± 8%  -12.08%  (p=0.032 n=5+5)
Sqlmock/500_rows-8      1.03ms ± 6%    0.94ms ± 7%   -8.12%  (p=0.016 n=5+5)
Sqlmock/1000_rows-8     1.19ms ± 4%    1.08ms ± 6%   -8.76%  (p=0.016 n=5+5)
Sqlmock/10000_rows-8    1.31ms ± 3%    1.23ms ± 5%   -6.00%  (p=0.016 n=5+5)
```

## `allocated memory`

`sqlh` and its reflect based approach do consume more memory than standard library. There are a few areas in `sqlh` and `set` where memory pooling could potentially bring `sqlh` closer to `database/sql`.

```
name                  old alloc/op   new alloc/op   delta
Libpq/5_rows-8          2.87kB ± 0%    3.60kB ± 0%  +25.13%  (p=0.008 n=5+5)
Libpq/50_rows-8         20.4kB ± 0%    23.7kB ± 0%  +15.86%  (p=0.008 n=5+5)
Libpq/100_rows-8        40.0kB ± 0%    46.1kB ± 0%  +15.10%  (p=0.016 n=5+4)
Libpq/500_rows-8         199kB ± 0%     228kB ± 0%  +14.27%  (p=0.008 n=5+5)
Libpq/1000_rows-8        401kB ± 0%     458kB ± 0%  +14.09%  (p=0.008 n=5+5)
Sqlite/5_rows-8         6.19kB ± 0%    6.91kB ± 0%  +11.64%  (p=0.008 n=5+5)
Sqlite/50_rows-8        55.8kB ± 0%    59.0kB ± 0%   +5.81%  (p=0.008 n=5+5)
Sqlite/100_rows-8        111kB ± 0%     117kB ± 0%   +5.46%  (p=0.008 n=5+5)
Sqlite/500_rows-8        554kB ± 0%     583kB ± 0%   +5.13%  (p=0.008 n=5+5)
Sqlite/1000_rows-8      1.11MB ± 0%    1.17MB ± 0%   +5.08%  (p=0.008 n=5+5)
Sqlmock/5_rows-8        2.92kB ± 1%    4.08kB ± 1%  +39.74%  (p=0.008 n=5+5)
Sqlmock/50_rows-8       2.92kB ± 1%    4.10kB ± 1%  +40.37%  (p=0.008 n=5+5)
Sqlmock/100_rows-8      2.93kB ± 0%    4.09kB ± 1%  +39.60%  (p=0.008 n=5+5)
Sqlmock/500_rows-8      2.99kB ± 1%    4.19kB ± 1%  +40.26%  (p=0.008 n=5+5)
Sqlmock/1000_rows-8     3.09kB ± 2%    4.36kB ± 1%  +40.79%  (p=0.008 n=5+5)
Sqlmock/10000_rows-8    5.16kB ± 3%    7.53kB ± 2%  +45.95%  (p=0.008 n=5+5)
```

## `number of allocations`

```
name                  old allocs/op  new allocs/op  delta
Libpq/5_rows-8            88.0 ± 0%     107.0 ± 0%  +21.59%  (p=0.008 n=5+5)
Libpq/50_rows-8            717 ± 0%       871 ± 0%  +21.48%  (p=0.008 n=5+5)
Libpq/100_rows-8         1.42k ± 0%     1.72k ± 0%  +21.41%  (p=0.008 n=5+5)
Libpq/500_rows-8         7.68k ± 0%     9.18k ± 0%  +19.59%  (p=0.008 n=5+5)
Libpq/1000_rows-8        15.7k ± 0%     18.7k ± 0%  +19.16%  (p=0.008 n=5+5)
Sqlite/5_rows-8            254 ± 0%       273 ± 0%   +7.48%  (p=0.008 n=5+5)
Sqlite/50_rows-8         2.28k ± 0%     2.44k ± 0%   +6.75%  (p=0.008 n=5+5)
Sqlite/100_rows-8        4.53k ± 0%     4.84k ± 0%   +6.70%  (p=0.008 n=5+5)
Sqlite/500_rows-8        23.2k ± 0%     24.7k ± 0%   +6.49%  (p=0.008 n=5+5)
Sqlite/1000_rows-8       46.7k ± 0%     49.7k ± 0%   +6.44%  (p=0.008 n=5+5)
Sqlmock/5_rows-8          36.0 ± 0%      43.0 ± 0%  +19.44%  (p=0.008 n=5+5)
Sqlmock/50_rows-8         36.0 ± 0%      43.0 ± 0%  +19.44%  (p=0.008 n=5+5)
Sqlmock/100_rows-8        36.0 ± 0%      43.0 ± 0%  +19.44%  (p=0.008 n=5+5)
Sqlmock/500_rows-8        36.0 ± 0%      44.0 ± 0%  +22.22%  (p=0.008 n=5+5)
Sqlmock/1000_rows-8       37.0 ± 0%      46.4 ± 1%  +25.41%  (p=0.016 n=4+5)
Sqlmock/10000_rows-8      48.0 ± 0%      89.4 ± 2%  +86.25%  (p=0.016 n=4+5)
```
