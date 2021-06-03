`sqlhbenchmarks` is a sibling benchmarking package for `github.com/nofeaturesonlybugs/sqlh`.  I decided to break the benchmarks out of `sqlh` because:  

* The benchmark code is rather verbose and greatly increases the code base for `sqlh`.
* Keeping the benchmark code in `sqlh` increases the number of dependencies in `sqlh` even though it's not using them for anything practical.

My primary goals with my benchmark testing were to evaluate:

1. Is `sqlh` comparable in performance to currently accepted and widely used packages?
2. What are the weak areas where I can come back later for optimizations?
3. What is the **true** cost of `reflect` when pointed at a real database when other factors start to matter (network IO, disk IO at the database server)?
4. To act as something of a stress test by working with record sets of various sizes.

I'll summarize my results with the following:
1. I think `sqlh.Scanner` is roughly on par with existing packages and APIs; thus I feel confident using it in my applications anywhere I'd have considered using `sqlx` or `sqlscan`.
2. `model.Models` is probably sufficient for small to medium applications but is lacking in a few benchmarks (not to mention the somewhat limited current feature set).  I plan to use it in my small to medium sized applications if that's worth anything to you.
3. Once you factor in a real database the cost of `reflect` starts to dwindle.

## Repeating These Tests  
If you have any desire to repeat these benchmarks for your environment or as part of evaluating which SQL package to use in Go:  

* Create a `TEST_POSTGRES` environment variable with a correct DSN for `lib/pq` to run the Postgres tests.  The user in the DSN will need to be able to perform some `ALTER TABLE` statements; see `schema.go` for the exact statements.
* Create a `TEST_SQLITE` environment variable with a correct DSN for Sqlite.  Out of the box this package uses `modernc.org/sqlite`; you can make slight alterations to `functions_sqlite.go` to point it at `mattn` instead.

## Notes on Sqlite  
My `model.Models` type only supports grammars with a `RETURNING` clause; therefore to benchmark Sqlite I needed to use version 3.35.  Originally I was using the `github.com/mattn` package (and a specific commit for Sqlite 3.35) but was having trouble when switching back and forth between Windows and Debian for benchmarks.  Eventually I substituted `github.com/mattn`'s Sqlite for `modernc.org/sqlite`.  This satisfied my desire of having something other than Postgres to benchmark against however I do not use Sqlite professionally; I don't know if the `modernc` version of Sqlite is production ready (seems to be experimental).  Even though I do not present them here I will say the `mattn` Sqlite benchmarks were more performant when I did have it working.

## Notes on `gorm`  
Since `gorm` was relatively easy to point at Postgres I included it in the `lib/pq` driver benchmarks.  I did not feel like struggling to get `gorm` to behave with `modernc` Sqlite or `sqlmock` so it is not present in those benchmarks.

`gorm` has very good performance when inserting or updating slices of records.  Obviously it is using some type of bulk operation internally when asked to work with slices of records.  In such benchmarks where I've included `gorm` it is not really a fair comparison -- the other libraries could probably achieve similar results if written to do so.  However this improved performance with `gorm` "just worked" and required no extra effort on my part.  I felt that was noteworthy and included it in such benchmarks even if the comparison is unfair.  Kudos to the `gorm` team in this regard.

## Notes on `squirrel`  
`squirrel` was an interesting experience.  I'd had no experience with the package prior to including it in my benchmarks.  I do find it to be an improvement in placing query arguments next to where they're used in the query, especially when creating `UPDATE` statements.  I'm not overly fond of the introduction of new types to handle things like prepared statements although I understand the reasoning.  `squirrel` is generally more memory hungry than other packages.

## Hardware
```
goos: windows
goarch: amd64
pkg: github.com/nofeaturesonlybugs/sqlh/benchmarks
cpu: Intel(R) Core(TM) i7-7700K CPU @ 4.20GHz
```

## Scanning with `sqlmock`  
The following tests:  

* Use `sqlmock` for setting up a mock database.
* Perform `SELECT * FROM table LIMIT %v` where the limits are: 5, 50, 100, 500, 1000, & 10000
* And scan into the following struct:
```go
type SaleReport struct {
	Id                 int    `json:"id" db:"pk"`
	CreatedTime        string `json:"created_time" db:"created_tmz"`
	ModifiedTime       string `json:"modified_time" db:"modified_tmz"`
	Price              int    `json:"price" db:"price"`
	Quantity           int    `json:"quantity" db:"quantity"`
	Total              int    `json:"total" db:"total"`
	CustomerId         int    `json:"customer_id" db:"customer_id"`
	CustomerFirst      string `json:"customer_first" db:"customer_first"`
	CustomerLast       string `json:"customer_last" db:"customer_last"`
	VendorId           int    `json:"vendor_id" db:"vendor_id"`
	VendorName         string `json:"vendor_name" db:"vendor_name"`
	VendorDescription  string `json:"vendor_description" db:"vendor_description"`
	VendorContactId    int    `json:"vendor_contact_id" db:"vendor_contact_id"`
	VendorContactFirst string `json:"vendor_contact_first" db:"vendor_contact_first"`
	VendorContactLast  string `json:"vendor_contact_last" db:"vendor_contact_last"`
}
```
My interpretation of the results:  
* Without a real database the cost of `reflect` is more pronounced for small record selections.
* Objectively `sqlh.Scanner` is the *worst* performing although not by much.  There's three reasons this may be the case:  
  * `set` and therefore `sqlh` does not use `FieldByName` methods in `reflect` due to a bug in the `reflect` package.
  * `set` instantiates deeply nested structs if they are pointers and `nil`
  * `set`'s struct traversal when requesting fields by mapped names *could* be improved; I have some thoughts on how but nothing concrete as of yet.
```bash
# records
database/sql_5_rows-8              10000	    103287 ns/op	    2917 B/op	      36 allocs/op
sqlx_5_rows-8                       5990	    247959 ns/op	    3607 B/op	      39 allocs/op
scany_5_rows-8                      4126	    327388 ns/op	    3086 B/op	      40 allocs/op
sqlh_5_rows-8                       3156	    403795 ns/op	    3832 B/op	      43 allocs/op
# 50 records
database/sql_50_rows-8         	    2790	    459121 ns/op	    2902 B/op	      36 allocs/op
sqlx_50_rows-8                 	    2626	    495046 ns/op	    3623 B/op	      39 allocs/op
scany_50_rows-8               	    2392	    543458 ns/op	    3128 B/op	      40 allocs/op
sqlh_50_rows-8                 	    2046	    604777 ns/op	    3864 B/op	      43 allocs/op
# 100 records
database/sql_100_rows-8             1966	    697026 ns/op	    2920 B/op	      36 allocs/op
sqlx_100_rows-8                	    1551	    766337 ns/op	    3663 B/op	      39 allocs/op
scany_100_rows-8               	    1596	    732904 ns/op	    3118 B/op	      40 allocs/op
sqlh_100_rows-8                	    1444	    761937 ns/op	    3911 B/op	      43 allocs/op
# 500 records
database/sql_500_rows-8        	    1528	    819504 ns/op	    3008 B/op	      36 allocs/op
sqlx_500_rows-8                	    1519	    889164 ns/op	    3677 B/op	      39 allocs/op
scany_500_rows-8               	    1333	   1041175 ns/op	    3247 B/op	      41 allocs/op
sqlh_500_rows-8                     1345	    996992 ns/op	    4036 B/op	      44 allocs/op
# 1,000 records
database/sql_1000_rows-8            1388	    944315 ns/op	    3037 B/op	      36 allocs/op
sqlx_1000_rows-8                    1232	    966885 ns/op	    3795 B/op	      40 allocs/op
scany_1000_rows-8                   1230	    975264 ns/op	    3457 B/op	      42 allocs/op
sqlh_1000_rows-8                    1134	    998770 ns/op	    4259 B/op	      46 allocs/op
# 10,000 records
database/sql_10000_rows-8           1081	   1174419 ns/op	    4684 B/op	      45 allocs/op
sqlx_10000_rows-8                    970	   1169252 ns/op	    6331 B/op	      60 allocs/op
scany_10000_rows-8                  1065	   1131304 ns/op	    7770 B/op	      68 allocs/op
sqlh_10000_rows-8                   1003	   1190989 ns/op	    8940 B/op	      83 allocs/op
```

## Scanning Postgres using `lib/pq` driver:  
The following tests:  

* Are pointed at a local network instance of Postgres.
* Perform `SELECT * FROM table LIMIT %v` where the limits are: 5, 50, 100, 500, & 1000
* And scan into the following struct:
```go
type Address struct {
	Id           int    `json:"id" db:"pk" model:"key,auto" gorm:"column:pk;primaryKey"`
	CreatedTime  Time   `json:"created_time" db:"created_tmz" model:"inserted" gorm:"-"`
	ModifiedTime Time   `json:"modified_time" db:"modified_tmz" model:"inserted,updated" gorm:"-"`
	Street       string `json:"street"`
	City         string `json:"city"`
	State        string `json:"state"`
	Zip          string `json:"zip"`
	//
	pushModified Time
}
```

These results are a little more difficult to draw conclusions from.  We're at the mercy of network IO and load on the database, which should have been fairly consistent throughout the benchmarks.  We're also at the mercy of how Postgres stores records in memory and decides to use them for consecuitive similar requests; in other words running
a `SELECT` statement could easily affect the performance of the next one.

I ran all of my benchmarks multiple times and in any result below where one is doing particularly bad (or good) I very likely had a result where it landed on the other end of the spectrum.

In general I conclude that all of them perform roughly on par with each other in my simple benchmark.  Certainly there are cases where one package is conclusively better than another but that is beyond the scope of my goals.
```bash
# 5 records
database/sql_5_rows-8         	    6807	    211957 ns/op	    2753 B/op	      84 allocs/op
GORM_5_rows-8                 	    6654	    234121 ns/op	    7934 B/op	     201 allocs/op
sqlx_5_rows-8                 	    7126	    195034 ns/op	    3448 B/op	      92 allocs/op
scany_5_rows-8                	    6697	    179223 ns/op	    5333 B/op	     188 allocs/op
sqlh_5_rows-8                 	    6957	    187767 ns/op	    4406 B/op	     104 allocs/op
# 50 records
database/sql_50_rows-8        	    3846	    281356 ns/op	   19434 B/op	     710 allocs/op
GORM_50_rows-8                	    3645	    389633 ns/op	   36408 B/op	    1418 allocs/op
sqlx_50_rows-8                	    2569	    439411 ns/op	   23174 B/op	     763 allocs/op
scany_50_rows-8               	    2108	    486062 ns/op	   28727 B/op	     904 allocs/op
sqlh_50_rows-8                	    2108	    490222 ns/op	   29477 B/op	     865 allocs/op
# 100 records
database/sql_100_rows-8       	    2242	    455634 ns/op	   37986 B/op	    1412 allocs/op
GORM_100_rows-8               	    2048	    592872 ns/op	   68158 B/op	    2773 allocs/op
sqlx_100_rows-8               	    2461	    529780 ns/op	   45153 B/op	    1515 allocs/op
scany_100_rows-8              	    2382	    563152 ns/op	   54625 B/op	    1706 allocs/op
sqlh_100_rows-8               	    2518	    491763 ns/op	   57407 B/op	    1717 allocs/op
# 500 records
database/sql_500_rows-8       	    1041	   1302823 ns/op	  191128 B/op	    7669 allocs/op
GORM_500_rows-8               	     534	   2023785 ns/op	  325790 B/op	   14227 allocs/op
sqlx_500_rows-8               	     918	   1264888 ns/op	  226455 B/op	    8168 allocs/op
scany_500_rows-8              	     750	   1783926 ns/op	  266028 B/op	    8758 allocs/op
sqlh_500_rows-8               	     840	   1644985 ns/op	  283867 B/op	    9168 allocs/op
# 1,000 records
database/sql_1000_rows-8      	     440	   2557591 ns/op	  385111 B/op	   15681 allocs/op
GORM_1000_rows-8              	     315	   3682919 ns/op	  650095 B/op	   28736 allocs/op
sqlx_1000_rows-8              	     243	   5090719 ns/op	  454025 B/op	   16677 allocs/op
scany_1000_rows-8             	     187	   5714156 ns/op	  534367 B/op	   17776 allocs/op
sqlh_1000_rows-8              	     222	   5068332 ns/op	  578752 B/op	   18684 allocs/op
```

## Scanning Postgres using `modernc.org/sqlite` (Sqlite 3.35) driver:  
This benchmark was performed with the same `SELECT` and scan destination as the Postgres `SELECT` benchmark.

```bash
# 5 records
database/sql_5_rows-8          12458	     96676 ns/op	    7320 B/op	     250 allocs/op
sqlx_5_rows-8                  10000	    101621 ns/op	    8038 B/op	     258 allocs/op
scany_5_rows-8                 10000	    112007 ns/op	    9834 B/op	     354 allocs/op
sqlh_5_rows-8                  10000	    103653 ns/op	    8990 B/op	     270 allocs/op
# 50 records
database/sql_50_rows-8          3724	    305367 ns/op	   67448 B/op	    2275 allocs/op
sqlx_50_rows-8                  3951	    326165 ns/op	   71224 B/op	    2328 allocs/op
scany_50_rows-8                 3415	    347603 ns/op	   76601 B/op	    2469 allocs/op
sqlh_50_rows-8                  3453	    350259 ns/op	   77392 B/op	    2430 allocs/op
# 100 records
database/sql_100_rows-8         2318	    532861 ns/op	  134089 B/op	    4526 allocs/op
sqlx_100_rows-8                 2119	    564847 ns/op	  140973 B/op	    4629 allocs/op
scany_100_rows-8                2013	    603275 ns/op	  150671 B/op	    4820 allocs/op
sqlh_100_rows-8                 1963	    603979 ns/op	  153653 B/op	    4831 allocs/op
# 500 records
database/sql_500_rows-8          501	   2428738 ns/op	  673477 B/op	   23171 allocs/op
sqlx_500_rows-8                  456	   2662145 ns/op	  709732 B/op	   23674 allocs/op
scany_500_rows-8                 438	   2707386 ns/op	  746132 B/op	   24265 allocs/op
sqlh_500_rows-8                  435	   2776479 ns/op	  771244 B/op	   24676 allocs/op
# 1,000 records
database/sql_1000_rows-8         253	   4725117 ns/op	 1351002 B/op	   46671 allocs/op
sqlx_1000_rows-8                 228	   4974677 ns/op	 1423145 B/op	   47674 allocs/op
scany_1000_rows-8                224	   5307015 ns/op	 1494309 B/op	   48766 allocs/op
sqlh_1000_rows-8                 222	   5395924 ns/op	 1544802 B/op	   49676 allocs/op
```

## Postgres - Dumb Insert  
The following benchmarks show iterating a set of `X` records and inserting them without any database transactions or prepared statements.  Certainly you want to avoid this in your application if possible but there are times where it's what your application will need to do.

Results are fairly consistent across the board.  
```bash
# 5 records
database/sql_insert_5_row(s)-8         	     220	   5045438 ns/op	    6789 B/op	     180 allocs/op
GORM_insert_5_row(s)-8                 	     208	   6058386 ns/op	   33622 B/op	     471 allocs/op
squirrel_insert_5_row(s)-8             	     216	   5424851 ns/op	   28262 B/op	     665 allocs/op
sqlh/model_insert_5_row(s)-8         	     225	   6118051 ns/op	    8701 B/op	     205 allocs/op
# 50 records
database/sql_insert_50_row(s)-8        	      22	  57294218 ns/op	   67819 B/op	    1800 allocs/op
GORM_insert_50_row(s)-8                	      20	  62065835 ns/op	  334900 B/op	    4710 allocs/op
squirrel_insert_50_row(s)-8            	      18	  56739317 ns/op	  282610 B/op	    6653 allocs/op
sqlh/model_insert_50_row(s)-8        	      21	  60169900 ns/op	   86988 B/op	    2051 allocs/op
# 100 records
database/sql_insert_100_row(s)-8       	      10	 117970890 ns/op	  135688 B/op	    3602 allocs/op
GORM_insert_100_row(s)-8               	       9	 126497011 ns/op	  671555 B/op	    9427 allocs/op
squirrel_insert_100_row(s)-8           	      10	 118073370 ns/op	  564725 B/op	   13301 allocs/op
sqlh/model_insert_100_row(s)-8       	      10	 125292220 ns/op	  173843 B/op	    4100 allocs/op
# 500 records
database/sql_insert_500_row(s)-8       	       2	 565426900 ns/op	  677756 B/op	   18003 allocs/op
GORM_insert_500_row(s)-8               	       2	 650741850 ns/op	 3347560 B/op	   47080 allocs/op
squirrel_insert_500_row(s)-8           	       2	 573679700 ns/op	 2825428 B/op	   66529 allocs/op
sqlh/model_insert_500_row(s)-8       	       2	 615725850 ns/op	  869668 B/op	   20508 allocs/op
# 1,000 records
database/sql_insert_1000_row(s)-8      	       1	1163333900 ns/op	 1357776 B/op	   36028 allocs/op
GORM_insert_1000_row(s)-8              	       1	1314936200 ns/op	 6693088 B/op	   94130 allocs/op
squirrel_insert_1000_row(s)-8          	       1	1259743000 ns/op	 5651088 B/op	  133051 allocs/op
sqlh/model_insert_1000_row(s)-8      	       1	1208060600 ns/op	 1739528 B/op	   41015 allocs/op
```

## Postgres - Insert Slice w/ Begin(), Prepare(), Exec().
The following benchmarks show `Begin() -> Prepare() -> Exec()` to insert a slice of `X` records.
```bash
# 5 records
database/sql_begin+prepare+insert_5_row(s)-8            1490	    736679 ns/op	    6604 B/op	     165 allocs/op
GORM_slice+insert_5_row(s)-8                             950	   1067537 ns/op	    9815 B/op	     143 allocs/op
squirrel_begin+prepare+insert_5_row(s)-8                1084	   1092230 ns/op	   28024 B/op	     645 allocs/op
sqlh/model_begin+prepare+insert_5_row(s)-8               568	   2596587 ns/op	    8082 B/op	     196 allocs/op
# 50 records
database/sql_begin+prepare+insert_50_row(s)-8            138	   8354679 ns/op	   66298 B/op	    1652 allocs/op
GORM_slice+insert_50_row(s)-8                            501	   2473041 ns/op	   57312 B/op	     793 allocs/op
squirrel_begin+prepare+insert_50_row(s)-8                127	  11681451 ns/op	  280444 B/op	    6456 allocs/op
sqlh/model_begin+prepare+insert_50_row(s)-8              122	  12672856 ns/op	   67769 B/op	    1683 allocs/op
# 100 records
database/sql_begin+prepare+insert_100_row(s)-8            70	  24493961 ns/op	  132531 B/op	    3304 allocs/op
GORM_slice+insert_100_row(s)-8                           396	   2777564 ns/op	  107026 B/op	    1597 allocs/op
squirrel_begin+prepare+insert_100_row(s)-8                45	  26156067 ns/op	  560691 B/op	   12910 allocs/op
sqlh/model_begin+prepare+insert_100_row(s)-8              67	  23372439 ns/op	  134008 B/op	    3335 allocs/op
# 500 records
database/sql_begin+prepare+insert_500_row(s)-8             9	 111672778 ns/op	  664056 B/op	   16536 allocs/op
GORM_slice+insert_500_row(s)-8                           139	   8918647 ns/op	  583896 B/op	    8009 allocs/op
squirrel_begin+prepare+insert_500_row(s)-8                13	 113849262 ns/op	 2802340 B/op	   64551 allocs/op
sqlh/model_begin+prepare+insert_500_row(s)-8              14	  89883943 ns/op	  663012 B/op	   16543 allocs/op
# 1,000 records
database/sql_begin+prepare+insert_1000_row(s)-8            6	 176948517 ns/op	 1324336 B/op	   33033 allocs/op
GORM_slice+insert_1000_row(s)-8                           74	  16430857 ns/op	 1176609 B/op	   16018 allocs/op
squirrel_begin+prepare+insert_1000_row(s)-8                5	 227757900 ns/op	 5601816 B/op	  129079 allocs/op
sqlh/model_begin+prepare+insert_1000_row(s)-8              6	 195512150 ns/op	 1326716 B/op	   33074 allocs/op
```

## Postgres - Dumb Update  
The following benchmarks show iterating a set of `X` records and updating them without any database transactions or prepared statements.  Certainly you want to avoid this in your application if possible but there are times where it's what your application will need to do.

Results are fairly consistent with `gorm` being something of an outlier:  
```bash
# 5 records
database/sql_update_5_row(s)-8              482	   2249577 ns/op	    6537 B/op	     165 allocs/op
GORM_update_5_row(s)-8                      723	   1735263 ns/op	   29224 B/op	     335 allocs/op
squirrel_update_5_row(s)-8                  738	   2463613 ns/op	   38720 B/op	     875 allocs/op
sqlh/model_update_5_row(s)-8                580	   2557353 ns/op	    8413 B/op	     195 allocs/op
# 50 records
database/sql_update_50_row(s)-8              74	  17060749 ns/op	   65468 B/op	    1650 allocs/op
GORM_update_50_row(s)-8                     100	  12066319 ns/op	  292367 B/op	    3352 allocs/op
squirrel_update_50_row(s)-8                  62	  26438821 ns/op	  387376 B/op	    8753 allocs/op
sqlh/model_update_50_row(s)-8                43	  23284305 ns/op	   84257 B/op	    1950 allocs/op
# 100 records
database/sql_update_100_row(s)-8             27	  46719315 ns/op	  130933 B/op	    3301 allocs/op
GORM_update_100_row(s)-8                     61	  24044793 ns/op	  584659 B/op	    6704 allocs/op
squirrel_update_100_row(s)-8                 24	  59805871 ns/op	  774675 B/op	   17505 allocs/op
sqlh/model_update_100_row(s)-8               24	  54183550 ns/op	  168450 B/op	    3901 allocs/op
# 500 records
database/sql_update_500_row(s)-8              6	 213189917 ns/op	  658580 B/op	   16997 allocs/op
GORM_update_500_row(s)-8                     14	  80748021 ns/op	 2923198 B/op	   33519 allocs/op
squirrel_update_500_row(s)-8                  4	 439392025 ns/op	 3877772 B/op	   88019 allocs/op
sqlh/model_update_500_row(s)-8                4	 267054850 ns/op	  844580 B/op	   19752 allocs/op
# 1,000 records
database/sql_update_1000_row(s)-8             3	 466049733 ns/op	 1321160 B/op	   34498 allocs/op
GORM_update_1000_row(s)-8                     7	 197242614 ns/op	 5848568 B/op	   67040 allocs/op
squirrel_update_1000_row(s)-8                 2	 779891850 ns/op	 7760052 B/op	  176554 allocs/op
sqlh/model_update_1000_row(s)-8               3	 382446033 ns/op	 1690760 B/op	   39752 allocs/op
```

## Postgres - Update Slice w/ Begin(), Prepare(), Exec().
The following benchmarks show `Begin() -> Prepare() -> Exec()` to insert a slice of `X` records.
```bash
# 5 records
database/sql_begin+prepare+update_5_row(s)-8         	    1110	   1336426 ns/op	    6895 B/op	     150 allocs/op
GORM_update_5_row(s)-8                               	    3848	    520668 ns/op	   14241 B/op	     162 allocs/op
squirrel_begin+prepare+update_5_row(s)-8             	     841	   1529204 ns/op	   38957 B/op	     855 allocs/op
sqlh/model_begin+prepare+update_5_row(s)-8                   649	   2008840 ns/op	    7885 B/op	     173 allocs/op
# 50 records
database/sql_begin+prepare+update_50_row(s)-8        	     100	  10373155 ns/op	   69153 B/op	    1501 allocs/op
GORM_update_50_row(s)-8                              	     535	   7558777 ns/op	  110557 B/op	     820 allocs/op
squirrel_begin+prepare+update_50_row(s)-8            	      81	  17105143 ns/op	  389943 B/op	    8553 allocs/op
sqlh/model_begin+prepare+update_50_row(s)-8                   97	  17899106 ns/op	   70424 B/op	    1568 allocs/op
# 100 records
database/sql_begin+prepare+update_100_row(s)-8       	      63	  19669398 ns/op	  138172 B/op	    3001 allocs/op
GORM_update_100_row(s)-8                             	     280	  10390849 ns/op	  205307 B/op	    1625 allocs/op
squirrel_begin+prepare+update_100_row(s)-8           	      36	  28659156 ns/op	  779430 B/op	   17104 allocs/op
sqlh/model_begin+prepare+update_100_row(s)-8                  64	  21718161 ns/op	  139811 B/op	    3119 allocs/op
# 500 records
database/sql_begin+prepare+update_500_row(s)-8       	      12	  97196517 ns/op	  695215 B/op	   15498 allocs/op
GORM_update_500_row(s)-8                             	      74	  21223735 ns/op	 1006385 B/op	    8681 allocs/op
squirrel_begin+prepare+update_500_row(s)-8           	      10	 123098130 ns/op	 3903352 B/op	   86021 allocs/op
sqlh/model_begin+prepare+update_500_row(s)-8                  10	 135324480 ns/op	  697669 B/op	   15767 allocs/op
# 1,000 records
database/sql_begin+prepare+update_1000_row(s)-8      	       7	 206369114 ns/op	 1396044 B/op	   31509 allocs/op
GORM_update_1000_row(s)-8                            	      39	  38828905 ns/op	 2207693 B/op	   17690 allocs/op
squirrel_begin+prepare+update_1000_row(s)-8          	       6	 269227417 ns/op	 7809578 B/op	  172543 allocs/op
sqlh/model_begin+prepare+update_1000_row(s)-8                  6	 229407333 ns/op	 1397428 B/op	   31773 allocs/op
```

The previous Postgres tests are now repeated with Sqlite sans `gorm`.

## Sqlite - Dumb Insert  
```bash
# 5 records
database/sql_insert_5_row(s)-8             45	  29689611 ns/op	  266493 B/op	    6923 allocs/op
squirrel_insert_5_row(s)-8                 50	  25831794 ns/op	  288350 B/op	    7454 allocs/op
sqlh/model_insert_5_row(s)-8               51	  31900388 ns/op	  268402 B/op	    6955 allocs/op
# 50 records
database/sql_insert_50_row(s)-8             4	 257990950 ns/op	 2665130 B/op	   69301 allocs/op
squirrel_insert_50_row(s)-8                 4	 264287750 ns/op	 2883536 B/op	   74554 allocs/op
sqlh/model_insert_50_row(s)-8               4	 381673525 ns/op	 2684044 B/op	   69550 allocs/op
# 100 records
database/sql_insert_100_row(s)-8            2	 515282400 ns/op	 5330480 B/op	  138605 allocs/op
squirrel_insert_100_row(s)-8                2	 514017900 ns/op	 5766700 B/op	  149104 allocs/op
sqlh/model_insert_100_row(s)-8              2	 779631350 ns/op	 5368056 B/op	  139100 allocs/op
# 500 records
database/sql_insert_500_row(s)-8            1	2612650600 ns/op	26649912 B/op	  693013 allocs/op
squirrel_insert_500_row(s)-8                1	3168402800 ns/op	28837584 B/op	  745552 allocs/op
sqlh/model_insert_500_row(s)-8              1	2637685200 ns/op	26840328 B/op	  695505 allocs/op
# 1,000 records
database/sql_insert_1000_row(s)-8           1	5907092800 ns/op	53297856 B/op	 1386017 allocs/op
squirrel_insert_1000_row(s)-8               1	5839871100 ns/op	57670072 B/op	 1491095 allocs/op
sqlh/model_insert_1000_row(s)-8             1	5444616300 ns/op	53682080 B/op	 1391011 allocs/op
```

## Sqlite - Insert Slice w/ Begin(), Prepare(), Exec().
```bash
# 5 records
database/sql_begin+prepare+insert_5_row(s)-8        	1992	    595825 ns/op	  269003 B/op	    6975 allocs/op
squirrel_begin+prepare+insert_5_row(s)-8            	1833	    623820 ns/op	  290786 B/op	    7495 allocs/op
sqlh/model_begin+prepare+insert_5_row(s)-8               208	   5671543 ns/op	  269702 B/op	    6987 allocs/op
# 50 records
database/sql_begin+prepare+insert_50_row(s)-8       	 214	   6060468 ns/op	 2690042 B/op	   69754 allocs/op
squirrel_begin+prepare+insert_50_row(s)-8           	 192	   6271936 ns/op	 2907824 B/op	   74958 allocs/op
sqlh/model_begin+prepare+insert_50_row(s)-8              100	  10821204 ns/op	 2690489 B/op	   69763 allocs/op
# 100 records
database/sql_begin+prepare+insert_100_row(s)-8      	 106	  11184752 ns/op	 5379689 B/op	  139505 allocs/op
squirrel_begin+prepare+insert_100_row(s)-8          	  96	  12460612 ns/op	 5815093 B/op	  149910 allocs/op
sqlh/model_begin+prepare+insert_100_row(s)-8              76	  16489736 ns/op	 5380870 B/op	  139521 allocs/op
# 500 records
database/sql_begin+prepare+insert_500_row(s)-8      	  19	  57751174 ns/op	26901126 B/op	  697553 allocs/op
squirrel_begin+prepare+insert_500_row(s)-8          	  18	  62421122 ns/op	29076555 B/op	  749562 allocs/op
sqlh/model_begin+prepare+insert_500_row(s)-8              18	  61375433 ns/op	26900149 B/op	  697547 allocs/op
# 1,000 records
database/sql_begin+prepare+insert_1000_row(s)-8     	   9	 111906267 ns/op	53799579 B/op	 1395079 allocs/op
squirrel_begin+prepare+insert_1000_row(s)-8         	   8	 125255625 ns/op	58156813 B/op	 1499166 allocs/op
sqlh/model_begin+prepare+insert_1000_row(s)-8              9	 120388589 ns/op	53796854 B/op	 1395054 allocs/op
```

## Sqlite - Dumb Update  
```bash
# 5 records
database/sql_update_5_row(s)-8               10000	    123116 ns/op	    4610 B/op	     175 allocs/op
squirrel_update_5_row(s)-8                    5836	    213105 ns/op	   37144 B/op	     935 allocs/op
sqlh/model_update_5_row(s)-8                  8781	    131921 ns/op	    6524 B/op	     205 allocs/op
# 50 records
database/sql_update_50_row(s)-8                988	   1238107 ns/op	   46097 B/op	    1751 allocs/op
squirrel_update_50_row(s)-8                    570	   2128517 ns/op	  371435 B/op	    9352 allocs/op
sqlh/model_update_50_row(s)-8                  906	   1331060 ns/op	   65286 B/op	    2050 allocs/op
# 100 records
database/sql_update_100_row(s)-8               483	   2446387 ns/op	   92150 B/op	    3501 allocs/op
squirrel_update_100_row(s)-8                   279	   4267277 ns/op	  742911 B/op	   18704 allocs/op
sqlh/model_update_100_row(s)-8                 457	   2667099 ns/op	  130623 B/op	    4102 allocs/op
# 500 records
database/sql_update_500_row(s)-8                97	  12766361 ns/op	  466980 B/op	   18001 allocs/op
squirrel_update_500_row(s)-8                    55	  21336253 ns/op	 3722212 B/op	   94011 allocs/op
sqlh/model_update_500_row(s)-8                  90	  13257092 ns/op	  654745 B/op	   20753 allocs/op
# 1,000 records
database/sql_update_1000_row(s)-8               49	  24498727 ns/op	  939114 B/op	   36502 allocs/op
squirrel_update_1000_row(s)-8                   27	  42796267 ns/op	 7453046 B/op	  188537 allocs/op
sqlh/model_update_1000_row(s)-8                 44	  26776673 ns/op	 1312121 B/op	   41767 allocs/op
```

## Sqlite - Update Slice w/ Begin(), Prepare(), Exec().
```bash
# 5 records
database/sql_begin+prepare+update_5_row(s)-8            9703	    130129 ns/op	    5847 B/op	     190 allocs/op
squirrel_begin+prepare+update_5_row(s)-8                5458	    218285 ns/op	   38304 B/op	     945 allocs/op
sqlh/model_begin+prepare+update_5_row(s)-8              8781	    136517 ns/op	    6591 B/op	     205 allocs/op
# 50 records
database/sql_begin+prepare+update_50_row(s)-8            916	   1297656 ns/op	   58457 B/op	    1900 allocs/op
squirrel_begin+prepare+update_50_row(s)-8                550	   2339813 ns/op	  383147 B/op	    9453 allocs/op
sqlh/model_begin+prepare+update_50_row(s)-8              859	   1350686 ns/op	   59912 B/op	    1960 allocs/op
# 100 records
database/sql_begin+prepare+update_100_row(s)-8           462	   2594040 ns/op	  116932 B/op	    3801 allocs/op
squirrel_begin+prepare+update_100_row(s)-8               277	   4367517 ns/op	  766107 B/op	   18904 allocs/op
sqlh/model_begin+prepare+update_100_row(s)-8             440	   2738170 ns/op	  119264 B/op	    3912 allocs/op
# 500 records
database/sql_begin+prepare+update_500_row(s)-8            91	  12915989 ns/op	  590369 B/op	   19494 allocs/op
squirrel_begin+prepare+update_500_row(s)-8                55	  21805969 ns/op	 3837766 B/op	   95006 allocs/op
sqlh/model_begin+prepare+update_500_row(s)-8              87	  13552601 ns/op	  595252 B/op	   19761 allocs/op
# 1,000 records
database/sql_begin+prepare+update_1000_row(s)-8           45	  25837742 ns/op	 1186918 B/op	   39500 allocs/op
squirrel_begin+prepare+update_1000_row(s)-8               26	  43744408 ns/op	 7684063 B/op	  190526 allocs/op
sqlh/model_begin+prepare+update_1000_row(s)-8             43	  27197272 ns/op	 1192095 B/op	   39770 allocs/op
```
