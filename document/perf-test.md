> zerodb性能测试



## Proxy单机性能测试
[nanxing@zerodb-proxy003 ~]$ lscpu
Architecture:          x86_64
CPU op-mode(s):        32-bit, 64-bit
Byte Order:            Little Endian
CPU(s):                4
On-line CPU(s) list:   0-3
Thread(s) per core:    2
Core(s) per socket:    2
Socket(s):             1
NUMA node(s):          1
Vendor ID:             GenuineIntel
CPU family:            6
Model:                 85
Model name:            Intel(R) Xeon(R) Platinum 8163 CPU @ 2.50GHz
Stepping:              4
CPU MHz:               2500.006
BogoMIPS:              5000.01
Hypervisor vendor:     KVM
Virtualization type:   full
L1d cache:             32K
L1i cache:             32K
L2 cache:              1024K
L3 cache:              33792K
NUMA node0 CPU(s):     0-3
Flags:                 fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush mmx fxsr sse sse2 ss ht syscall nx pdpe1gb rdtscp lm constant_tsc rep_good nopl eagerfpu pni pclmulqdq ssse3 fma cx16 pcid sse4_1 sse4_2 x2apic movbe popcnt tsc_deadline_timer aes xsave avx f16c rdrand hypervisor lahf_lm abm 3dnowprefetch fsgsbase tsc_adjust bmi1 hle avx2 smep bmi2 erms invpcid rtm mpx avx512f rdseed adx smap avx512cd xsaveopt xsavec xgetbv1
[nanxing@zerodb-proxy003 ~]$ sudo sysctl -a | grep tcp
net.ipv4.tcp_abort_on_overflow = 0
net.ipv4.tcp_adv_win_scale = 1
net.ipv4.tcp_allowed_congestion_control = cubic reno
net.ipv4.tcp_app_win = 31
net.ipv4.tcp_autocorking = 1
net.ipv4.tcp_available_congestion_control = cubic reno
net.ipv4.tcp_base_mss = 512
net.ipv4.tcp_challenge_ack_limit = 100
net.ipv4.tcp_congestion_control = cubic
net.ipv4.tcp_dsack = 1
net.ipv4.tcp_early_retrans = 3
net.ipv4.tcp_ecn = 2
net.ipv4.tcp_fack = 1
net.ipv4.tcp_fastopen = 0
net.ipv4.tcp_fastopen_key = 00000000-00000000-00000000-00000000
net.ipv4.tcp_fin_timeout = 60
net.ipv4.tcp_frto = 2
net.ipv4.tcp_invalid_ratelimit = 500
net.ipv4.tcp_keepalive_intvl = 75
net.ipv4.tcp_keepalive_probes = 9
net.ipv4.tcp_keepalive_time = 7200
net.ipv4.tcp_limit_output_bytes = 262144
net.ipv4.tcp_low_latency = 0
net.ipv4.tcp_max_orphans = 32768
net.ipv4.tcp_max_ssthresh = 0
net.ipv4.tcp_max_syn_backlog = 1024
net.ipv4.tcp_max_tw_buckets = 5000
net.ipv4.tcp_mem = 185340	247120	370680
net.ipv4.tcp_min_tso_segs = 2
net.ipv4.tcp_moderate_rcvbuf = 1
net.ipv4.tcp_mtu_probing = 0
net.ipv4.tcp_no_metrics_save = 0
net.ipv4.tcp_notsent_lowat = -1
net.ipv4.tcp_orphan_retries = 0
net.ipv4.tcp_reordering = 3
net.ipv4.tcp_retrans_collapse = 1
net.ipv4.tcp_retries1 = 3
net.ipv4.tcp_retries2 = 15
net.ipv4.tcp_rfc1337 = 0
net.ipv4.tcp_rmem = 4096	87380	6291456
net.ipv4.tcp_sack = 1
net.ipv4.tcp_slow_start_after_idle = 1
net.ipv4.tcp_stdurg = 0
net.ipv4.tcp_syn_retries = 6
net.ipv4.tcp_synack_retries = 2
net.ipv4.tcp_syncookies = 1
net.ipv4.tcp_thin_dupack = 0
net.ipv4.tcp_thin_linear_timeouts = 0
net.ipv4.tcp_timestamps = 1
net.ipv4.tcp_tso_win_divisor = 3
net.ipv4.tcp_tw_recycle = 0
net.ipv4.tcp_tw_reuse = 0
net.ipv4.tcp_window_scaling = 1
net.ipv4.tcp_wmem = 4096	16384	4194304
net.ipv4.tcp_workaround_signed_windows = 0


###
####connections:10
[nanxing@zerodb-proxy003 ~]$ free -h
              total        used        free      shared  buff/cache   available
Mem:           7.6G        276M        286M        376K        7.1G        7.1G
Swap:            0B          0B          0B
[nanxing@zerodb-proxy003 ~]$ pstree -p 26192 | wc -l
14
----total-cpu-usage---- ---load-avg--- -dsk/total- -net/total- ---system--
usr sys idl wai hiq siq| 1m   5m  15m | read  writ| recv  send| int   csw
 42  15  37   0   0   6|2.78 2.91 2.83|   0     0 |  12M   12M|  71k   80k
 42  15  38   0   0   5|2.78 2.91 2.83|   0    16k|  12M   12M|  70k   77k
 43  15  37   0   0   6|2.78 2.91 2.83|   0     0 |  12M   12M|  70k   77k
 42  15  37   0   0   6|2.71 2.90 2.82|   0     0 |  12M   12M|  71k   80k
 43  15  38   0   0   5|2.71 2.90 2.82|   0     0 |  12M   12M|  69k   76k
 43  16  36   0   0   6|2.71 2.90 2.82|   0     0 |  12M   12M|  74k   85k
 42  15  37   0   0   6|2.71 2.90 2.82|   0     0 |  12M   12M|  73k   83k
 41  15  38   0   0   6|2.71 2.90 2.82|   0    12k|  12M   12M|  71k   80k
 43  14  38   0   0   5|2.58 2.87 2.81|   0     0 |  12M   12M|  70k   77k
 43  14  37   0   0   5|2.58 2.87 2.81|   0     0 |  12M   12M|  71k   79k
 42  15  37   0   0   6|2.58 2.87 2.81|   0     0 |  12M   12M|  72k   80k
 42  15  37   0   0   5|2.58 2.87 2.81|   0     0 |  12M   12M|  71k   79k
 42  15  37   0   0   6|2.58 2.87 2.81|   0    16k|  12M   12M|  71k   80k
 40  15  39   0   0   6|2.61 2.87 2.81|   0     0 |  12M   12M|  69k   76k
 41  15  38   0   0   6|2.61 2.87 2.81|   0     0 |  12M   12M|  69k   75k
 41  15  38   0   0   5|2.61 2.87 2.81|   0     0 |  12M   12M|  69k   76k
 43  14  38   0   0   6|2.61 2.87 2.81|   0     0 |  12M   12M|  69k   76k
 44  16  36   0   0   5|2.61 2.87 2.81|   0   372k|  12M   12M|  75k   87k
 44  14  37   0   0   5|2.72 2.89 2.82|   0     0 |  12M   12M|  70k   78k
 43  14  37   0   0   6|2.72 2.89 2.82|   0     0 |  12M   12M|  70k   77k
 42  14  38   0   0   6|2.72 2.89 2.82|   0     0 |  12M   12M|  69k   75k
 42  15  38   0   0   6|2.72 2.89 2.82|   0     0 |  12M   12M|  68k   75k
 42  15  37   0   0   6|2.72 2.89 2.82|   0    28k|  12M   12M|  71k   79k
 43  15  37   0   0   6|2.82 2.91 2.83|   0     0 |  12M   12M|  72k   82k
 42  15  38   0   0   5|2.82 2.91 2.83|   0     0 |  12M   12M|  69k   76k
 42  14  39   0   0   5|2.82 2.91 2.83|   0     0 |  12M   12M|  69k   77k
 42  15  37   0   0   6|2.82 2.91 2.83|   0     0 |  12M   12M|  70k   77k
 41  15  38   0   0   6|2.82 2.91 2.83|   0     0 |  12M   12M|  69k   77k
 43  15  37   0   0   5|2.92 2.92 2.83|   0     0 |  12M   12M|  70k   76k
 42  14  38   0   0   6|2.92 2.92 2.83|   0     0 |  12M   12M|  70k   78k
 42  15  37   0   0   6|2.92 2.92 2.83|   0     0 |  12M   12M|  71k   80k
 43  15  37   0   0   5|2.92 2.92 2.83|   0    16k|  12M   12M|  70k   77k
 43  14  38   0   0   5|2.92 2.92 2.83|   0     0 |  12M   12M|  69k   76k
 42  15  37   0   0   6|3.01 2.94 2.84|   0     0 |  12M   12M|  72k   81k
 42  15  38   0   0   6|3.01 2.94 2.84|   0     0 |  12M   12M|  70k   78k
 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2624s]      [r:10,w:0,u:0,d:0]  28318    0       28318    NaN        0.35         74067613

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2625s]      [r:10,w:0,u:0,d:0]  27833    0       27833    NaN        0.36         74095446

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2626s]      [r:10,w:0,u:0,d:0]  28458    0       28458    NaN        0.35         74123904

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2627s]      [r:10,w:0,u:0,d:0]  28109    0       28109    NaN        0.35         74152013

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2628s]      [r:10,w:0,u:0,d:0]  28359    0       28359    NaN        0.35         74180372

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2629s]      [r:10,w:0,u:0,d:0]  28137    0       28137    NaN        0.35         74208509

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2630s]      [r:10,w:0,u:0,d:0]  28251    0       28251    NaN        0.35         74236760

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2631s]      [r:10,w:0,u:0,d:0]  28512    0       28512    NaN        0.35         74265272

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2632s]      [r:10,w:0,u:0,d:0]  28263    0       28263    NaN        0.35         74293535

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2633s]      [r:10,w:0,u:0,d:0]  28291    0       28291    NaN        0.35         74321826

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2634s]      [r:10,w:0,u:0,d:0]  28255    0       28255    NaN        0.35         74350081

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2635s]      [r:10,w:0,u:0,d:0]  28332    0       28332    NaN        0.35         74378413

####connections:20
[nanxing@zerodb-proxy003 ~]$ free -h
              total        used        free      shared  buff/cache   available
Mem:           7.6G        277M        283M        376K        7.1G        7.1G
Swap:            0B          0B          0B
----total-cpu-usage---- ---load-avg--- -dsk/total- -net/total- ---system--
usr sys idl wai hiq siq| 1m   5m  15m | read  writ| recv  send| int   csw
 58  14  20   0   0   8|3.76 3.65 3.26|   0     0 |  19M   19M|  56k   15k
 57  16  19   0   0   8|3.76 3.65 3.26|   0     0 |  19M   20M|  57k   15k
 58  14  19   0   0   8|3.76 3.65 3.26|   0     0 |  19M   20M|  58k   16k
 57  15  19   0   0   9|3.76 3.65 3.26|   0     0 |  19M   20M|  58k   15k
 57  15  20   0   0   8|3.76 3.65 3.26|   0     0 |  19M   19M|  57k   15k
 57  15  20   0   0   8|3.78 3.65 3.26|   0    12k|  19M   20M|  57k   15k
 58  16  19   0   0   8|3.78 3.65 3.26|   0     0 |  19M   20M|  57k   15k
 56  15  20   0   0   9|3.78 3.65 3.26|   0     0 |  19M   20M|  57k   15k
 58  14  19   0   0   8|3.78 3.65 3.26|   0     0 |  19M   20M|  57k   14k
 58  15  19   0   0   8|3.78 3.65 3.26|   0    44k|  19M   20M|  58k   15k
 59  14  20   0   0   8|3.79 3.66 3.26|   0     0 |  19M   20M|  57k   15k
 57  15  20   0   0   8|3.79 3.66 3.26|   0     0 |  19M   19M|  58k   17k
 59  14  19   0   0   9|3.79 3.66 3.26|   0     0 |  19M   20M|  58k   15k
 58  15  19   0   0   8|3.79 3.66 3.26|   0     0 |  19M   20M|  58k   15k
 57  14  20   0   0   8|3.79 3.66 3.26|   0    16k|  19M   19M|  57k   16k
 58  14  19   0   0   8|3.73 3.65 3.26|   0     0 |  19M   20M|  57k   16k
 56  15  19   0   0   9|3.73 3.65 3.26|   0     0 |  19M   20M|  58k   15k
 58  14  19   0   0   8|3.73 3.65 3.26|   0     0 |  19M   19M|  57k   15k
 57  15  20   0   0   8|3.73 3.65 3.26|   0     0 |  19M   19M|  57k   15k
 59  14  20   0   0   7|3.73 3.65 3.26|   0    20k|  19M   19M|  56k   17k
 58  15  19   0   0   8|3.67 3.64 3.26|   0     0 |  19M   20M|  58k   15k
 58  15  19   0   0   8|3.67 3.64 3.26|   0     0 |  19M   20M|  57k   15k
 58  14  19   0   0   9|3.67 3.64 3.26|   0     0 |  19M   20M|  58k   16k
 57  15  20   0   0   8|3.67 3.64 3.26|   0     0 |  19M   19M|  57k   16k
 58  15  19   0   0   8|3.67 3.64 3.26|   0     0 |  19M   20M|  58k   15k
 57  16  20   0   0   8|3.70 3.64 3.26|   0    56k|  19M   20M|  57k   16k
 57  15  20   0   0   8|3.70 3.64 3.26|   0     0 |  19M   20M|  58k   14k
 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [1015s]      [r:20,w:0,u:0,d:0]  46172    0       46172    NaN        0.43         46613684

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [1016s]      [r:20,w:0,u:0,d:0]  44821    0       44821    NaN        0.44         46658505

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [1017s]      [r:20,w:0,u:0,d:0]  45642    0       45642    NaN        0.44         46704147

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [1018s]      [r:20,w:0,u:0,d:0]  46009    0       46009    NaN        0.43         46750156

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [1019s]      [r:20,w:0,u:0,d:0]  46283    0       46283    NaN        0.43         46796439

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [1020s]      [r:20,w:0,u:0,d:0]  45911    0       45911    NaN        0.43         46842350

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [1021s]      [r:20,w:0,u:0,d:0]  45923    0       45923    NaN        0.43         46888273

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [1022s]      [r:20,w:0,u:0,d:0]  45941    0       45941    NaN        0.43         46934214

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [1023s]      [r:20,w:0,u:0,d:0]  45589    0       45589    NaN        0.44         46979803

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [1024s]      [r:20,w:0,u:0,d:0]  45550    0       45550    NaN        0.44         47025353

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [1025s]      [r:20,w:0,u:0,d:0]  45731    0       45731    NaN        0.44         47071084

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [1026s]      [r:20,w:0,u:0,d:0]  46353    0       46353    NaN        0.43         47117437

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [1027s]      [r:20,w:0,u:0,d:0]  46116    0       46116    NaN        0.43         47163553

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [1028s]      [r:20,w:0,u:0,d:0]  46346    0       46346    NaN        0.43         47209899
 ####connections:40
 [nanxing@zerodb-proxy003 ~]$ free -h
               total        used        free      shared  buff/cache   available
 Mem:           7.6G        286M        274M        376K        7.1G        7.1G
 Swap:            0B          0B          0B
----total-cpu-usage---- ---load-avg--- -dsk/total- -net/total- ---system--
usr sys idl wai hiq siq| 1m   5m  15m | read  writ| recv  send| int   csw
 62  16  13   0   0   9|3.85 3.90 3.68|   0     0 |  21M   22M|  56k 5633
 61  16  13   0   0  10|3.85 3.90 3.68|   0     0 |  21M   22M|  56k 5598
 63  16  12   0   0   9|3.85 3.90 3.68|   0    12k|  21M   22M|  56k 5715
 62  15  14   0   0  10|3.85 3.90 3.68|   0     0 |  21M   21M|  56k 5892
 61  16  13   0   0  10|3.86 3.90 3.68|   0     0 |  21M   22M|  56k 5329
 63  16  12   0   0  10|3.86 3.90 3.68|   0     0 |  21M   22M|  56k 5662
 62  16  13   0   0   9|3.86 3.90 3.68|   0     0 |  21M   22M|  56k 5783
 62  15  14   0   0   9|3.86 3.90 3.68|   0    12k|  21M   22M|  56k 5622
 61  17  13   0   0   9|3.86 3.90 3.68|   0    40k|  21M   21M|  56k 5638
 62  15  12   0   0  10|3.95 3.92 3.69|   0     0 |  21M   22M|  55k 5659
 62  15  13   0   0  10|3.95 3.92 3.69|   0     0 |  21M   22M|  56k 5645
 61  16  13   0   0  10|3.95 3.92 3.69|   0     0 |  21M   22M|  56k 5648
 63  16  12   0   0   9|3.95 3.92 3.69|   0     0 |  21M   22M|  56k 6132
 62  16  12   0   0   9|3.95 3.92 3.69|   0    12k|  21M   21M|  56k 6001
 62  16  12   0   0  10|3.96 3.92 3.69|   0     0 |  21M   22M|  56k 5435
 64  15  12   0   0   9|3.96 3.92 3.69|   0     0 |  21M   22M|  56k 5831
 63  16  12   0   0   9|3.96 3.92 3.69|   0     0 |  21M   22M|  56k 5565
 62  15  13   0   0   9|3.96 3.92 3.69|   0    56k|  21M   21M|  55k 5365
 63  15  12   0   0   9|3.96 3.92 3.69|   0     0 |  21M   22M|  56k 5621
 63  15  13   0   0   9|3.96 3.92 3.69|   0     0 |  21M   22M|  56k 5419
 61  16  13   0   0  10|3.96 3.92 3.69|   0     0 |  21M   22M|  56k 5416
 63  15  13   0   0   9|3.96 3.92 3.69|   0     0 |  21M   21M|  55k 6019
 62  16  13   0   0   9|3.96 3.92 3.69|   0     0 |  21M   22M|  55k 5637
 61  16  13   0   0  10|3.96 3.92 3.69|   0     0 |  21M   22M|  55k 6009
 63  15  13   0   0   9|3.96 3.92 3.69|   0     0 |  21M   22M|  56k 5519
 61  16  13   0   0   9|3.96 3.92 3.69|   0     0 |  21M   22M|  56k 5826
 62  16  12   0   0  10|3.96 3.92 3.69|   0    16k|  21M   22M|  56k 5539
 62  16  13   0   0   9|3.96 3.92 3.69|   0     0 |  21M   22M|  56k 5597
 62  15  13   0   0  10|3.96 3.92 3.69|   0  4096B|  21M   22M|  56k 5611
 61  17  13   0   0  10|3.97 3.92 3.70|   0     0 |  21M   22M|  56k 5588'
 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [1538s]      [r:40,w:0,u:0,d:0]  51066    0       51066    NaN        0.78         77926130

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [1539s]      [r:40,w:0,u:0,d:0]  51064    0       51064    NaN        0.78         77977194

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [1540s]      [r:40,w:0,u:0,d:0]  50988    0       50988    NaN        0.78         78028182

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [1541s]      [r:40,w:0,u:0,d:0]  50952    0       50952    NaN        0.78         78079134

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [1542s]      [r:40,w:0,u:0,d:0]  50831    0       50831    NaN        0.79         78129965

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [1543s]      [r:40,w:0,u:0,d:0]  50846    0       50846    NaN        0.78         78180811

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [1544s]      [r:40,w:0,u:0,d:0]  50536    0       50536    NaN        0.79         78231347

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [1545s]      [r:40,w:0,u:0,d:0]  51269    0       51269    NaN        0.78         78282616

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [1546s]      [r:40,w:0,u:0,d:0]  50835    0       50835    NaN        0.79         78333451

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [1547s]      [r:40,w:0,u:0,d:0]  50642    0       50642    NaN        0.79         78384093

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [1548s]      [r:40,w:0,u:0,d:0]  50365    0       50365    NaN        0.79         78434458

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [1549s]      [r:40,w:0,u:0,d:0]  50516    0       50516    NaN        0.79         78484974

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [1550s]      [r:40,w:0,u:0,d:0]  50200    0       50200    NaN        0.80         78535174

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [1551s]      [r:40,w:0,u:0,d:0]  51189    0       51189    NaN        0.78         78586363

####connections:80
[nanxing@zerodb-proxy003 ~]$ free -h
              total        used        free      shared  buff/cache   available
Mem:           7.6G        300M        263M        376K        7.1G        7.1G
Swap:            0B          0B          0B
[nanxing@zerodb-proxy003 ~]$ pstree -p 26192 | wc -l
15
----total-cpu-usage---- ---load-avg--- -dsk/total- -net/total- ---system--
usr sys idl wai hiq siq| 1m   5m  15m | read  writ| recv  send| int   csw
 63  16   8   0   0  13|4.13 4.13 4.12|   0     0 |  22M   23M|  56k 3470
 65  15   7   0   0  14|4.13 4.13 4.12|   0     0 |  22M   23M|  56k 3186
 64  15   8   0   0  13|4.13 4.13 4.12|   0    28k|  22M   23M|  56k 3213
 63  16   8   0   0  13|4.13 4.13 4.12|   0     0 |  22M   23M|  56k 3497
 64  15   9   0   0  12|4.13 4.13 4.12|   0  8756k|  22M   22M|  56k 3546
 64  15   9   0   0  12|4.12 4.13 4.12|   0     0 |  22M   23M|  56k 3439
 65  16   7   0   0  12|4.12 4.13 4.12|   0     0 |  22M   23M|  56k 3055
 65  15   8   0   0  12|4.12 4.13 4.12|   0     0 |  22M   23M|  55k 3306
 65  16   7   0   0  12|4.12 4.13 4.12|   0    60k|  22M   23M|  57k 3427
 66  15   8   0   0  11|4.12 4.13 4.12|   0     0 |  22M   23M|  56k 3377
 65  17   7   0   0  11|4.11 4.13 4.12|   0     0 |  22M   23M|  56k 3310
 63  17   8   0   0  12|4.11 4.13 4.12|   0     0 |  22M   23M|  56k 3363
 62  17   8   0   0  12|4.11 4.13 4.12|   0     0 |  22M   23M|  56k 3219
 65  16   8   0   0  11|4.11 4.13 4.12|   0    12k|  22M   23M|  56k 3353
 65  17   7   0   0  12|4.11 4.13 4.12|   0     0 |  22M   23M|  57k 3060
 64  17   8   0   0  11|4.10 4.12 4.12|   0     0 |  22M   23M|  56k 3600
 62  17   9   0   0  12|4.10 4.12 4.12|   0     0 |  22M   23M|  57k 3120
 63  16   7   0   0  13|4.10 4.12 4.12|   0     0 |  22M   23M|  56k 3423
 63  17   7   0   0  13|4.10 4.12 4.12|   0     0 |  22M   23M|  57k 3084
 63  17   7   0   0  12|4.10 4.12 4.12|   0     0 |  22M   23M|  56k 4185
 64  15   8   0   0  13|4.09 4.12 4.12|   0     0 |  22M   23M|  56k 2999
 63  17   8   0   0  13|4.09 4.12 4.12|   0     0 |  22M   23M|  56k 3457
 64  16   8   0   0  13|4.09 4.12 4.12|   0     0 |  22M   23M|  56k 3318
 63  16   8   0   0  13|4.09 4.12 4.12|   0    28k|  22M   23M|  56k 3326
 62  15   9   0   0  14|4.09 4.12 4.12|   0     0 |  22M   23M|  56k 3002
 63  16   7   0   0  14|4.08 4.12 4.11|   0     0 |  22M   23M|  56k 3075
 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [4079s]      [r:80,w:0,u:0,d:0]  53680    0       53680    NaN        1.49         218074953

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [4080s]      [r:80,w:0,u:0,d:0]  53747    0       53747    NaN        1.49         218128700

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [4081s]      [r:80,w:0,u:0,d:0]  53606    0       53606    NaN        1.49         218182306

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [4082s]      [r:80,w:0,u:0,d:0]  54043    0       54043    NaN        1.48         218236349

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [4083s]      [r:80,w:0,u:0,d:0]  54077    0       54077    NaN        1.48         218290426

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [4084s]      [r:80,w:0,u:0,d:0]  53770    0       53770    NaN        1.49         218344196

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [4085s]      [r:80,w:0,u:0,d:0]  53497    0       53497    NaN        1.49         218397693

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [4086s]      [r:80,w:0,u:0,d:0]  53322    0       53322    NaN        1.50         218451015

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [4087s]      [r:80,w:0,u:0,d:0]  53685    0       53685    NaN        1.49         218504700

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [4088s]      [r:80,w:0,u:0,d:0]  53607    0       53607    NaN        1.49         218558307

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [4089s]      [r:80,w:0,u:0,d:0]  53702    0       53702    NaN        1.49         218612009

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [4090s]      [r:80,w:0,u:0,d:0]  53231    0       53231    NaN        1.50         218665240

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [4091s]      [r:80,w:0,u:0,d:0]  53129    0       53129    NaN        1.51         218718369
####connections:160
[nanxing@zerodb-proxy003 ~]$ free -h
              total        used        free      shared  buff/cache   available
Mem:           7.6G        327M        263M        376K        7.1G        7.1G
Swap:            0B          0B          0B
----total-cpu-usage---- ---load-avg--- -dsk/total- -net/total- ---system--
usr sys idl wai hiq siq| 1m   5m  15m | read  writ| recv  send| int   csw
 67  18   4   0   0  10|4.12 4.14 4.07|   0     0 |  23M   23M|  56k 1950
 67  18   4   0   0  10|4.11 4.14 4.07|   0     0 |  23M   23M|  56k 2031
 68  17   5   0   0  10|4.11 4.14 4.07|   0     0 |  23M   23M|  57k 2271
 67  18   5   0   0  11|4.11 4.14 4.07|   0     0 |  23M   23M|  56k 2130
 69  17   5   0   0   9|4.11 4.14 4.07|   0    48k|  23M   23M|  56k 2252
 67  17   5   0   0  11|4.11 4.14 4.07|   0     0 |  23M   23M|  56k 2013
 68  17   5   0   0  10|4.10 4.13 4.07|   0     0 |  23M   23M|  56k 2010
 67  19   4   0   0  10|4.10 4.13 4.07|   0     0 |  23M   23M|  56k 2235
 67  18   5   0   0  10|4.10 4.13 4.07|   0     0 |  23M   23M|  57k 1943
 69  17   5   0   0  10|4.10 4.13 4.07|   0     0 |  23M   23M|  56k 1961
 69  17   4   0   0   9|4.10 4.13 4.07|   0     0 |  23M   24M|  56k 1925
 69  16   4   0   0  10|4.10 4.13 4.07|   0     0 |  23M   23M|  56k 1954
 68  18   5   0   0  10|4.10 4.13 4.07|   0     0 |  23M   23M|  56k 1928
 67  18   5   0   0   9|4.10 4.13 4.07|   0    12k|  23M   23M|  56k 1891
 69  16   5   0   0   9|4.10 4.13 4.07|   0     0 |  23M   23M|  57k 2167
 68  18   4   0   0  10|4.10 4.13 4.07|   0     0 |  23M   24M|  56k 2011
 68  17   5   0   0   9|4.17 4.15 4.07|   0     0 |  23M   23M|  56k 2052
 69  17   4   0   0  10|4.17 4.15 4.07|   0     0 |  23M   23M|  57k 2477
 70  17   4   0   0   9|4.17 4.15 4.07|   0     0 |  23M   24M|  57k 1964
 69  17   5   0   0  10|4.17 4.15 4.07|   0     0 |  23M   23M|  57k 2282
 70  17   5   0   0   9|4.17 4.15 4.07|   0     0 |  23M   23M|  56k 2097
 68  17   5   0   0  10|4.16 4.14 4.07|   0     0 |  23M   23M|  56k 2326
 68  18   5   0   0   9|4.16 4.14 4.07|   0     0 |  23M   23M|  56k 1918
 68  18   5   0   0   9|4.16 4.14 4.07|   0    12k|  23M   23M|  57k 2051
 68  17   5   0   0  10|4.16 4.14 4.07|   0  8192B|  23M   23M|  56k 2173
 time            thds               tps     wtps    rtps    cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
 [1467s]      [r:160,w:0,u:0,d:0]  54863    0       54863   0         NaN        2.92         80597513

 time            thds               tps     wtps    rtps    cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
 [1468s]      [r:160,w:0,u:0,d:0]  55215    0       55215   0         NaN        2.90         80652728

 time            thds               tps     wtps    rtps    cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
 [1469s]      [r:160,w:0,u:0,d:0]  55660    0       55660   0         NaN        2.87         80708388

 time            thds               tps     wtps    rtps    cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
 [1470s]      [r:160,w:0,u:0,d:0]  55684    0       55684   0         NaN        2.87         80764072

 time            thds               tps     wtps    rtps    cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
 [1471s]      [r:160,w:0,u:0,d:0]  55023    0       55023   0         NaN        2.91         80819095

 time            thds               tps     wtps    rtps    cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
 [1472s]      [r:160,w:0,u:0,d:0]  54812    0       54812   0         NaN        2.92         80873907

 time            thds               tps     wtps    rtps    cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
 [1473s]      [r:160,w:0,u:0,d:0]  55408    0       55408   0         NaN        2.89         80929315

 time            thds               tps     wtps    rtps    cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
 [1474s]      [r:160,w:0,u:0,d:0]  54555    0       54555   0         NaN        2.93         80983870

 time            thds               tps     wtps    rtps    cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
 [1475s]      [r:160,w:0,u:0,d:0]  55387    0       55387   0         NaN        2.89         81039257

 time            thds               tps     wtps    rtps    cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
 [1476s]      [r:160,w:0,u:0,d:0]  55052    0       55052   0         NaN        2.90         81094309

 time            thds               tps     wtps    rtps    cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
 [1477s]      [r:160,w:0,u:0,d:0]  55119    0       55119   0         NaN        2.91         81149428

 time            thds               tps     wtps    rtps    cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
 [1478s]      [r:160,w:0,u:0,d:0]  55641    0       55641   0         NaN        2.87         81205069

 ####connections:200
 [nanxing@zerodb-proxy003 ~]$ free -h
               total        used        free      shared  buff/cache   available
 Mem:           7.6G        339M        250M        376K        7.1G        7.0G
 Swap:            0B          0B          0B
 ----total-cpu-usage---- ---load-avg--- -dsk/total- -net/total- ---system--
 usr sys idl wai hiq siq| 1m   5m  15m | read  writ| recv  send| int   csw
  69  17   2   0   0  12|4.08 4.14 4.07|   0     0 |  23M   24M|  56k 1432
  70  18   2   0   0  11|4.07 4.13 4.07|   0     0 |  23M   24M|  56k 1412
  70  18   2   0   0  10|4.07 4.13 4.07|   0     0 |  23M   24M|  56k 1491
  70  18   2   0   0  11|4.07 4.13 4.07|   0    12k|  24M   24M|  56k 1438
  69  19   1   0   0  11|4.07 4.13 4.07|   0     0 |  23M   24M|  56k 1388
  71  16   1   0   0  11|4.07 4.13 4.07|   0     0 |  23M   24M|  56k 1596
  70  18   1   0   0  11|4.06 4.13 4.07|   0     0 |  24M   24M|  57k 1454
  70  18   1   0   0  11|4.06 4.13 4.07|   0     0 |  23M   24M|  56k 1461
  70  17   2   0   0  12|4.06 4.13 4.07|   0     0 |  23M   24M|  56k 1364
  70  18   1   0   0  11|4.06 4.13 4.07|   0     0 |  23M   24M|  57k 1383
  68  19   2   0   0  11|4.06 4.13 4.07|   0     0 |  23M   24M|  56k 1387
  71  17   2   0   0  11|4.06 4.13 4.07|   0   316k|  23M   24M|  56k 1411
  69  18   2   0   0  12|4.06 4.13 4.07|   0     0 |  23M   24M|  56k 1514
  69  18   2   0   0  12|4.06 4.13 4.07|   0    72k|  23M   24M|  57k 2043
  68  19   1   0   0  12|4.06 4.13 4.07|   0  4096B|  22M   23M|  55k 2501
  68  18   2   0   0  12|4.06 4.13 4.07|   0     0 |  23M   24M|  56k 1564
  69  17   2   0   0  12|4.05 4.13 4.07|   0     0 |  23M   24M|  56k 1524
  68  18   3   0   0  11|4.05 4.13 4.07|   0     0 |  23M   23M|  56k 1790
  68  18   3   0   0  11|4.05 4.13 4.07|   0     0 |  23M   24M|  56k 1544
  69  17   3   0   0  12|4.05 4.13 4.07|   0   140k|  23M   24M|  56k 1646
  68  17   3   0   0  12|4.05 4.13 4.07|   0     0 |  23M   24M|  56k 1519
  70  17   2   0   0  11|4.05 4.13 4.07|   0   160k|  23M   24M|  56k 1596
  68  18   3   0   0  11|4.05 4.13 4.07|   0     0 |  23M   24M|  57k 1499
  69  17   3   0   0  11|4.05 4.13 4.07|   0     0 |  23M   24M|  56k 1614
  68  18   3   0   0  11|4.05 4.13 4.07|   0     0 |  23M   24M|  56k 1520
  time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
  [1945s]      [r:200,w:0,u:0,d:0]  56192    0       56192    NaN        3.56         109119636

  time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
  [1946s]      [r:200,w:0,u:0,d:0]  56127    0       56127    NaN        3.56         109175763

  time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
  [1947s]      [r:200,w:0,u:0,d:0]  56414    0       56414    NaN        3.54         109232177

  time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
  [1948s]      [r:200,w:0,u:0,d:0]  56337    0       56337    NaN        3.55         109288514

  time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
  [1949s]      [r:200,w:0,u:0,d:0]  56350    0       56350    NaN        3.55         109344864

  time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
  [1950s]      [r:200,w:0,u:0,d:0]  56225    0       56225    NaN        3.56         109401089

  time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
  [1951s]      [r:200,w:0,u:0,d:0]  56120    0       56120    NaN        3.56         109457209

  time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
  [1952s]      [r:200,w:0,u:0,d:0]  56052    0       56052    NaN        3.57         109513261

  time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
  [1953s]      [r:200,w:0,u:0,d:0]  55978    0       55978    NaN        3.57         109569239

  time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
  [1954s]      [r:200,w:0,u:0,d:0]  56276    0       56276    NaN        3.55         109625515

  time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
  [1955s]      [r:200,w:0,u:0,d:0]  56268    0       56268    NaN        3.55         109681783

  time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
  [1956s]      [r:200,w:0,u:0,d:0]  56523    0       56523    NaN        3.54         109738306

  time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
  [1957s]      [r:200,w:0,u:0,d:0]  56197    0       56197    NaN        3.56         109794503

  time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
  [1958s]      [r:200,w:0,u:0,d:0]  56435    0       56435    NaN        3.54         109850938
####connections:400
[nanxing@zerodb-proxy003 ~]$ free -h
              total        used        free      shared  buff/cache   available
Mem:           7.6G        388M        293M        376K        7.0G        7.0G
Swap:            0B          0B          0B
----total-cpu-usage---- ---load-avg--- -dsk/total- -net/total- ---system--
usr sys idl wai hiq siq| 1m   5m  15m | read  writ| recv  send| int   csw
 70  18   2   0   0  11|4.10 4.15 3.81|   0     0 |  24M   24M|  57k 1422
 69  19   2   0   0  11|4.09 4.15 3.81|   0     0 |  24M   24M|  56k 1277
 69  19   2   0   0  11|4.09 4.15 3.81|   0    16k|  23M   24M|  56k 1320
 70  18   2   0   0  11|4.09 4.15 3.81|   0     0 |  24M   24M|  57k 1243
 70  18   2   0   0  11|4.09 4.15 3.81|   0     0 |  23M   24M|  56k 1285
 69  18   3   0   0  11|4.09 4.15 3.81|   0     0 |  23M   24M|  56k 1236
 69  19   2   0   0  11|4.08 4.15 3.81|   0     0 |  24M   24M|  56k 1170
 70  17   2   0   0  11|4.08 4.15 3.81|   0     0 |  23M   24M|  56k 1297
 69  18   2   0   0  11|4.08 4.15 3.81|   0     0 |  24M   24M|  56k 1171
 69  17   3   0   0  11|4.08 4.15 3.81|   0     0 |  23M   24M|  56k 1291
 68  19   2   0   0  11|4.08 4.15 3.81|   0     0 |  23M   24M|  56k 1274
 69  18   2   0   0  11|4.08 4.14 3.81|   0    16k|  24M   24M|  56k 1205
 68  19   2   0   0  11|4.08 4.14 3.81|   0     0 |  23M   24M|  56k 1513
 69  18   2   0   0  11|4.08 4.14 3.81|   0     0 |  24M   24M|  56k 1368
 68  19   2   0   0  11|4.08 4.14 3.81|   0     0 |  23M   24M|  55k 1228
 69  19   1   0   0  11|4.08 4.14 3.81|   0     0 |  24M   24M|  56k 1187
 68  19   2   0   0  11|4.07 4.14 3.81|   0     0 |  23M   24M|  56k 1239
 69  19   2   0   0  11|4.07 4.14 3.81|   0    12k|  23M   24M|  54k 1354
 71  18   1   0   0  11|4.07 4.14 3.81|   0     0 |  24M   24M|  56k 1206
 69  19   2   0   0  11|4.07 4.14 3.81|   0     0 |  24M   24M|  56k 1090
 70  18   1   0   0  11|4.07 4.14 3.81|   0     0 |  24M   24M|  56k 1182
 69  18   2   0   0  11|4.06 4.14 3.81|   0     0 |  23M   24M|  56k 1199
 69  19   1   0   0  11|4.06 4.14 3.81|   0     0 |  24M   24M|  56k 1268
 69  18   2   0   0  11|4.06 4.14 3.81|   0     0 |  24M   24M|  56k 1206
 69  18   1   0   0  11|4.06 4.14 3.81|   0     0 |  24M   24M|  56k 1230
 69  18   2   0   0  11|4.06 4.14 3.81|   0     0 |  24M   24M|  56k 1170
 69  18   2   0   0  11|4.06 4.14 3.81|   0   316k|  23M   24M|  56k 1307
 time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2223s]      [r:400,w:0,u:0,d:0]  56509    0       56509    NaN        7.09         123476096

 time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2224s]      [r:400,w:0,u:0,d:0]  56988    0       56988    NaN        7.02         123533084

 time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2225s]      [r:400,w:0,u:0,d:0]  56027    0       56027    NaN        7.12         123589111

 time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2226s]      [r:400,w:0,u:0,d:0]  56337    0       56337    NaN        7.12         123645448

 time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2227s]      [r:400,w:0,u:0,d:0]  56214    0       56214    NaN        7.12         123701662

 time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2228s]      [r:400,w:0,u:0,d:0]  55882    0       55882    NaN        7.16         123757544

 time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2229s]      [r:400,w:0,u:0,d:0]  56654    0       56654    NaN        7.05         123814198

 time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2230s]      [r:400,w:0,u:0,d:0]  56803    0       56803    NaN        7.05         123871001

 time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2231s]      [r:400,w:0,u:0,d:0]  56293    0       56293    NaN        7.10         123927294

 time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2232s]      [r:400,w:0,u:0,d:0]  56707    0       56707    NaN        7.04         123984001

 time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2233s]      [r:400,w:0,u:0,d:0]  56845    0       56845    NaN        7.05         124040846

 time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2234s]      [r:400,w:0,u:0,d:0]  57249    0       57249    NaN        6.98         124098095

 time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2235s]      [r:400,w:0,u:0,d:0]  56528    0       56528    NaN        7.06         124154623
####connections:20000
[nanxing@zerodb-proxy003 ~]$ free -h
              total        used        free      shared  buff/cache   available
Mem:           7.6G        3.4G        295M        376K        3.9G        3.9G
Swap:            0B          0B          0B
----total-cpu-usage---- ---load-avg--- -dsk/total- -net/total- ---system--
usr sys idl wai hiq siq| 1m   5m  15m | read  writ| recv  send| int   csw
 69  17   3   0   0  12|4.14 4.14 4.01|   0     0 |  15M   20M|  37k  916
 67  18   2   0   0  13|4.14 4.14 4.01|   0     0 |  17M   23M|  41k 1068
 62  21   4   0   0  14|4.13 4.14 4.01|   0     0 |  18M   23M|  42k 1029
 68  16   3   0   0  12|4.13 4.14 4.01|   0     0 |  15M   20M|  36k  948
 63  19   4   0   0  14|4.13 4.14 4.01|   0     0 |  17M   23M|  42k 1019
 69  17   2   0   0  13|4.13 4.14 4.01|   0    12k|  15M   20M|  37k 1094
 62  21   4   0   0  14|4.13 4.14 4.01|   0     0 |  17M   23M|  43k  917
 60  20   5   0   0  14|4.12 4.14 4.01|   0     0 |  17M   23M|  42k 1104
 69  16   3   0   0  12|4.12 4.14 4.01|   0     0 |  15M   20M|  36k  951
 63  19   4   0   0  14|4.12 4.14 4.01|   0     0 |  17M   23M|  41k 1184
 65  19   3   0   0  13|4.12 4.14 4.01|   0    12k|  17M   22M|  39k 1126
 66  18   4   0   0  13|4.12 4.14 4.01|   0     0 |  16M   21M|  38k 1005
 63  19   5   0   0  13|4.19 4.15 4.02|   0     0 |  17M   23M|  41k 1161
 68  17   2   0   0  13|4.19 4.15 4.02|   0     0 |  15M   20M|  35k 1005
 62  20   6   0   0  13|4.19 4.15 4.02|   0     0 |  17M   23M|  42k  970
 60  20   6   0   0  14|4.19 4.15 4.02|   0    16k|  17M   22M|  42k 1125
 68  17   3   0   0  12|4.19 4.15 4.02|   0     0 |  15M   20M|  32k 1242
 62  21   4   0   0  13|4.17 4.15 4.02|   0     0 |  17M   23M|  41k 1119
 64  19   5   0   0  12|4.17 4.15 4.02|   0     0 |  16M   21M|  39k 1039
 64  19   5   0   0  12|4.17 4.15 4.02|   0     0 |  16M   21M|  39k  927
 61  20   6   0   0  13|4.17 4.15 4.02|   0    16k|  17M   23M|  42k 1148
 67  18   3   0   0  12|4.17 4.15 4.02|   0     0 |  15M   20M|  36k  973
 62  19   6   0   0  13|4.16 4.15 4.02|   0     0 |  17M   23M|  42k 1009
 61  19   6   0   0  13|4.16 4.15 4.02|   0     0 |  17M   23M|  42k 1225
 69  17   3   0   0  12|4.16 4.15 4.02|   0     0 |  15M   20M|  36k  925
 64  20   4   0   0  13|4.16 4.15 4.02|   0    12k|  17M   23M|  41k 1063
time            thds                 tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[1184s]      [r:20000,w:0,u:0,d:0]  42446    0       42446    NaN        494.63       46454910

time            thds                 tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[1185s]      [r:20000,w:0,u:0,d:0]  41532    0       41532    NaN        481.91       46496442

time            thds                 tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[1186s]      [r:20000,w:0,u:0,d:0]  35840    0       35840    NaN        554.87       46532282

time            thds                 tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[1187s]      [r:20000,w:0,u:0,d:0]  41443    0       41443    NaN        482.02       46573725

time            thds                 tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[1188s]      [r:20000,w:0,u:0,d:0]  36256    0       36256    NaN        483.06       46609981

time            thds                 tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[1189s]      [r:20000,w:0,u:0,d:0]  42141    0       42141    NaN        536.43       46652122

time            thds                 tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[1190s]      [r:20000,w:0,u:0,d:0]  42242    0       42242    NaN        476.72       46694364

time            thds                 tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[1191s]      [r:20000,w:0,u:0,d:0]  35493    0       35493    NaN        550.65       46729857

time            thds                 tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[1192s]      [r:20000,w:0,u:0,d:0]  42284    0       42284    NaN        480.04       46772141

time            thds                 tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[1193s]      [r:20000,w:0,u:0,d:0]  41781    0       41781    NaN        478.32       46813922

time            thds                 tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[1194s]      [r:20000,w:0,u:0,d:0]  35719    0       35719    NaN        559.17       46849641

time            thds                 tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[1195s]      [r:20000,w:0,u:0,d:0]  41674    0       41674    NaN        478.03       46891315

time            thds                 tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[1196s]      [r:20000,w:0,u:0,d:0]  36440    0       36440    NaN        504.66       46927755

time            thds                 tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[1197s]      [r:20000,w:0,u:0,d:0]  42097    0       42097    NaN        515.72       46969852

time            thds                 tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[1198s]      [r:20000,w:0,u:0,d:0]  42017    0       42017    NaN        479.08       47011869


## Cobar2.0单机性能测试
####connections:10
[nanxing@zerodb-proxy003 cobar-server-2.0.0]$ free -h
              total        used        free      shared  buff/cache   available
Mem:           7.6G        4.9G        262M        400K        2.5G        2.5G
Swap:            0B          0B          0B
[nanxing@zerodb-proxy003 cobar-server-2.0.0]$ pstree -p 5529 | wc -l
1580
----total-cpu-usage---- ---load-avg--- -dsk/total- -net/total- ---system--
usr sys idl wai hiq siq| 1m   5m  15m | read  writ| recv  send| int   csw
 39  27  29   0   0   5|2.46 2.68 2.09| 456k    0 |  10M   10M| 103k  204k
 35  28  32   0   0   5|2.35 2.65 2.08|   0   508k|  10M   10M| 104k  207k
 24  27  43   0   0   5|2.35 2.65 2.08|   0     0 |  11M   11M| 121k  236k
 24  27  44   0   0   5|2.35 2.65 2.08|   0   648k|  11M   11M| 120k  235k
 25  26  44   0   0   5|2.35 2.65 2.08|   0     0 |  11M   11M| 120k  234k
 24  26  44   0   0   5|2.35 2.65 2.08|   0     0 |  11M   11M| 121k  235k
 24  28  44   0   0   4|2.56 2.69 2.10|   0     0 |  11M   11M| 120k  234k
 24  26  45   0   0   5|2.56 2.69 2.10|   0    56k|  11M   11M| 121k  236k
 24  28  44   0   0   5|2.56 2.69 2.10|   0  9116k|  11M   11M| 120k  234k
 25  26  44   0   0   5|2.56 2.69 2.10|   0     0 |  11M   11M| 120k  235k
 24  27  44   0   0   5|2.56 2.69 2.10|   0     0 |  11M   11M| 121k  235k
 23  27  45   0   0   5|2.67 2.71 2.11|   0     0 |  11M   11M| 121k  235k
 24  27  44   0   0   4|2.67 2.71 2.11|   0     0 |  11M   11M| 121k  236k
 25  25  45   0   0   5|2.67 2.71 2.11|   0     0 |  11M   11M| 120k  234k
 23  28  44   0   0   4|2.67 2.71 2.11|   0    60k|  11M   11M| 120k  234k
 25  27  43   0   0   5|2.67 2.71 2.11|   0     0 |  11M   11M| 120k  235k
 23  27  44   0   0   6|2.54 2.68 2.10|   0     0 |  11M   11M| 121k  235k
 24  27  43   0   0   5|2.54 2.68 2.10|   0     0 |  11M   11M| 120k  234k
 24  26  44   0   0   6|2.54 2.68 2.10|   0  1472k|  11M   11M| 120k  234k
 25  26  44   0   0   5|2.54 2.68 2.10|   0     0 |  11M   11M| 121k  236k
 25  26  44   0   0   5|2.54 2.68 2.10|   0     0 |  11M   11M| 119k  233k
 24  26  44   0   0   6|2.66 2.71 2.11|   0     0 |  11M   11M| 121k  235k
 23  27  45   0   0   5|2.66 2.71 2.11|   0     0 |  11M   11M| 120k  234k
 24  26  44   0   0   6|2.66 2.71 2.11|   0    60k|  11M   11M| 120k  234k
 23  27  45   0   0   5|2.66 2.71 2.11|   0     0 |  11M   11M| 120k  234k
 24  26  44   0   0   5|2.66 2.71 2.11|   0     0 |  11M   11M| 120k  234k
 25  27  43   0   0   5|2.60 2.69 2.11|   0     0 |  11M   11M| 119k  232k
 24  27  44   0   0   6|2.60 2.69 2.11|   0     0 |  11M   11M| 120k  236k
 24  27  44   0   0   5|2.60 2.69 2.11|   0     0 |  11M   11M| 121k  236k
time            thds              tps     wtps    rtps    w-rsp(ms)  r-rsp(ms)    total-number
[1185s]      [r:10,w:0,u:0,d:0]  22870    0       22870   NaN        0.44         26797260

time            thds              tps     wtps    rtps    w-rsp(ms)  r-rsp(ms)    total-number
[1186s]      [r:10,w:0,u:0,d:0]  22857    0       22857   NaN        0.44         26820117

time            thds              tps     wtps    rtps    w-rsp(ms)  r-rsp(ms)    total-number
[1187s]      [r:10,w:0,u:0,d:0]  22738    0       22738   NaN        0.44         26842855

time            thds              tps     wtps    rtps    w-rsp(ms)  r-rsp(ms)    total-number
[1188s]      [r:10,w:0,u:0,d:0]  22884    0       22884   NaN        0.44         26865739

time            thds              tps     wtps    rtps    w-rsp(ms)  r-rsp(ms)    total-number
[1189s]      [r:10,w:0,u:0,d:0]  22870    0       22870   NaN        0.44         26888609

time            thds              tps     wtps    rtps    w-rsp(ms)  r-rsp(ms)    total-number
[1190s]      [r:10,w:0,u:0,d:0]  22819    0       22819   NaN        0.44         26911428

time            thds              tps     wtps    rtps    w-rsp(ms)  r-rsp(ms)    total-number
[1191s]      [r:10,w:0,u:0,d:0]  22868    0       22868   NaN        0.44         26934296

time            thds              tps     wtps    rtps    w-rsp(ms)  r-rsp(ms)    total-number
[1192s]      [r:10,w:0,u:0,d:0]  22737    0       22737   NaN        0.44         26957033

time            thds              tps     wtps    rtps    w-rsp(ms)  r-rsp(ms)    total-number
[1193s]      [r:10,w:0,u:0,d:0]  22756    0       22756   NaN        0.44         26979789

time            thds              tps     wtps    rtps    w-rsp(ms)  r-rsp(ms)    total-number
[1194s]      [r:10,w:0,u:0,d:0]  22514    0       22514   NaN        0.44         27002303

time            thds              tps     wtps    rtps    w-rsp(ms)  r-rsp(ms)    total-number
[1195s]      [r:10,w:0,u:0,d:0]  22920    0       22920   NaN        0.44         27025223

time            thds              tps     wtps    rtps    w-rsp(ms)  r-rsp(ms)    total-number
[1196s]      [r:10,w:0,u:0,d:0]  22810    0       22810   NaN        0.44         27048033

time            thds              tps     wtps    rtps    w-rsp(ms)  r-rsp(ms)    total-number
[1197s]      [r:10,w:0,u:0,d:0]  22790    0       22790   NaN        0.44         27070823
####connections:11
----total-cpu-usage---- ---load-avg--- -dsk/total- -net/total- ---system--
usr sys idl wai hiq siq| 1m   5m  15m | read  writ| recv  send| int   csw
 27  29  37   0   0   7|3.24 3.27 3.63|   0     0 |  12M   12M| 121k  243k
 27  29  37   0   0   7|3.24 3.27 3.63|   0     0 |  12M   12M| 121k  243k
 27  28  38   0   0   7|3.24 3.27 3.63|   0     0 |  12M   12M| 120k  241k
 27  28  37   0   0   7|3.24 3.27 3.63|   0     0 |  12M   12M| 120k  242k
 26  29  37   0   0   7|3.24 3.27 3.63|   0     0 |  12M   12M| 119k  240k
 26  29  38   0   0   7|3.06 3.23 3.61|   0     0 |  12M   12M| 122k  245k
 26  30  37   0   0   7|3.06 3.23 3.61|   0     0 |  12M   12M| 121k  244k
 27  29  37   0   0   7|3.06 3.23 3.61|   0    12k|  12M   12M| 121k  244k
 26  28  38   0   0   8|3.06 3.23 3.61|   0    12k|  12M   12M| 121k  243k
 27  28  37   0   0   7|3.06 3.23 3.61|   0     0 |  12M   12M| 122k  244k
 26  30  38   0   0   7|2.97 3.21 3.60|   0     0 |  12M   12M| 121k  245k
 28  28  38   0   0   6|2.97 3.21 3.60|   0     0 |  12M   12M| 120k  242k
 26  29  38   0   0   7|2.97 3.21 3.60|   0     0 |  12M   12M| 119k  241k
 27  28  38   0   0   7|2.97 3.21 3.60|   0   416k|  12M   12M| 120k  241k
 28  27  38   0   0   7|2.97 3.21 3.60|   0     0 |  12M   12M| 120k  243k
 27  29  37   0   0   6|2.98 3.21 3.60|   0     0 |  12M   12M| 121k  244k
 26  30  37   0   0   7|2.98 3.21 3.60|   0     0 |  12M   12M| 122k  245k
 28  29  37   0   0   6|2.98 3.21 3.60|   0     0 |  12M   12M| 121k  243k
 28  28  38   0   0   6|2.98 3.21 3.60|   0   116k|  12M   12M| 121k  243k
 27  29  38   0   0   6|2.98 3.21 3.60|   0     0 |  12M   12M| 121k  243k
 25  30  37   0   0   7|2.90 3.19 3.59|   0     0 |  12M   12M| 122k  244k
 27  29  37   0   0   7|2.90 3.19 3.59|   0     0 |  12M   12M| 120k  243k
 27  29  37   0   0   6|2.90 3.19 3.59|   0     0 |  12M   12M| 121k  243k
time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[470s]       [r:11,w:0,u:0,d:0]  24355    0       24355    NaN        0.45         11403849

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[471s]       [r:11,w:0,u:0,d:0]  24319    0       24319    NaN        0.45         11428168

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[472s]       [r:11,w:0,u:0,d:0]  24385    0       24385    NaN        0.45         11452553

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[473s]       [r:11,w:0,u:0,d:0]  24381    0       24381    NaN        0.45         11476934

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[474s]       [r:11,w:0,u:0,d:0]  24377    0       24377    NaN        0.45         11501311

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[475s]       [r:11,w:0,u:0,d:0]  24334    0       24334    NaN        0.45         11525645

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[476s]       [r:11,w:0,u:0,d:0]  24233    0       24233    NaN        0.45         11549878

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[477s]       [r:11,w:0,u:0,d:0]  24356    0       24356    NaN        0.45         11574234

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[478s]       [r:11,w:0,u:0,d:0]  24325    0       24325    NaN        0.45         11598559

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[479s]       [r:11,w:0,u:0,d:0]  24365    0       24365    NaN        0.45         11622924

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[480s]       [r:11,w:0,u:0,d:0]  23982    0       23982    NaN        0.46         11646906

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[481s]       [r:11,w:0,u:0,d:0]  24455    0       24455    NaN        0.45         11671361
####connections:12
----total-cpu-usage---- ---load-avg--- -dsk/total- -net/total- ---system--
usr sys idl wai hiq siq| 1m   5m  15m | read  writ| recv  send| int   csw
 27  32  36   0   0   6|3.92 3.95 4.12|   0     0 |  12M   12M| 122k  250k
 28  31  34   0   0   8|3.92 3.95 4.12|   0    16k|  13M   13M| 124k  255k
 29  31  33   0   0   7|3.92 3.95 4.12|   0     0 |  13M   13M| 122k  253k
 30  29  34   0   0   7|3.92 3.95 4.12|   0     0 |  13M   13M| 124k  256k
 30  30  33   0   0   7|3.92 3.95 4.12|   0     0 |  13M   13M| 123k  255k
 32  30  31   0   0   7|4.08 3.98 4.13|   0    40k|  12M   13M| 119k  246k
 30  29  34   0   0   7|4.08 3.98 4.13|   0     0 |  13M   13M| 124k  254k
 30  30  33   0   0   7|4.08 3.98 4.13|   0  8824k|  13M   13M| 123k  253k
 30  30  34   0   0   7|4.08 3.98 4.13|   0     0 |  13M   13M| 124k  255k
 28  32  33   0   0   6|4.08 3.98 4.13|   0     0 |  13M   13M| 122k  253k
 29  31  34   0   0   7|4.00 3.97 4.12|   0   136k|  13M   13M| 123k  254k
 29  30  34   0   0   7|4.00 3.97 4.12|   0     0 |  13M   13M| 124k  257k
 30  30  33   0   0   6|4.00 3.97 4.12|   0     0 |  13M   13M| 123k  254k
 30  30  34   0   0   7|4.00 3.97 4.12|   0     0 |  13M   13M| 122k  252k
 31  29  34   0   0   7|4.00 3.97 4.12|   0     0 |  13M   13M| 121k  252k
 30  29  33   0   0   7|3.92 3.95 4.12|   0     0 |  13M   13M| 123k  255k
 29  31  34   0   0   7|3.92 3.95 4.12|   0     0 |  13M   13M| 123k  255k
 29  30  33   0   0   8|3.92 3.95 4.12|   0    24k|  13M   13M| 123k  254k
 28  30  33   0   0   9|3.92 3.95 4.12|   0     0 |  13M   13M| 123k  254k
 28  30  34   0   0   8|3.92 3.95 4.12|   0     0 |  13M   13M| 124k  257k
 29  28  34   0   0   9|4.00 3.97 4.12|   0     0 |  13M   13M| 124k  256k
 29  28  33   0   0   9|4.00 3.97 4.12|   0     0 |  13M   13M| 124k  257k
 28  29  34   0   0   9|4.00 3.97 4.12|   0     0 |  13M   13M| 124k  256k
 28  29  33   0   0  10|4.00 3.97 4.12|   0     0 |  13M   13M| 124k  256k
 29  29  33   0   0  10|4.00 3.97 4.12|   0     0 |  13M   13M| 124k  255k
 28  29  34   0   0   9|3.76 3.92 4.11|   0     0 |  13M   13M| 123k  255k
time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[337s]       [r:12,w:0,u:0,d:0]  26053    0       26053    NaN        0.46         8765152

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[338s]       [r:12,w:0,u:0,d:0]  26163    0       26163    NaN        0.46         8791315

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[339s]       [r:12,w:0,u:0,d:0]  26219    0       26219    NaN        0.46         8817534

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[340s]       [r:12,w:0,u:0,d:0]  26260    0       26260    NaN        0.46         8843794

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[341s]       [r:12,w:0,u:0,d:0]  26104    0       26104    NaN        0.46         8869898

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[342s]       [r:12,w:0,u:0,d:0]  26090    0       26090    NaN        0.46         8895988

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[343s]       [r:12,w:0,u:0,d:0]  26296    0       26296    NaN        0.46         8922284

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[344s]       [r:12,w:0,u:0,d:0]  26329    0       26329    NaN        0.45         8948613

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[345s]       [r:12,w:0,u:0,d:0]  26225    0       26225    NaN        0.46         8974838

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[346s]       [r:12,w:0,u:0,d:0]  25977    0       25977    NaN        0.46         9000815

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[347s]       [r:12,w:0,u:0,d:0]  26227    0       26227    NaN        0.46         9027042

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[348s]       [r:12,w:0,u:0,d:0]  26072    0       26072    NaN        0.46         9053114

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[349s]       [r:12,w:0,u:0,d:0]  26158    0       26158    NaN        0.46         9079272

####connections:13
----total-cpu-usage---- ---load-avg--- -dsk/total- -net/total- ---system--
usr sys idl wai hiq siq| 1m   5m  15m | read  writ| recv  send| int   csw
 30  30  31   0   0   9|3.97 4.11 4.32|   0     0 |  14M   14M| 126k  266k
 30  30  31   0   0  10|3.97 4.11 4.32|   0    28k|  14M   14M| 125k  266k
 29  30  31   0   0  10|3.97 4.11 4.32|   0     0 |  13M   13M| 123k  262k
 29  30  30   0   0  10|3.97 4.11 4.32|   0     0 |  14M   14M| 125k  267k
 30  30  31   0   0  10|4.14 4.14 4.33|   0  4096B|  14M   14M| 126k  268k
 29  30  30   0   0  11|4.14 4.14 4.33|   0     0 |  14M   14M| 125k  265k
 30  29  30   0   0  10|4.14 4.14 4.33|   0     0 |  14M   14M| 125k  265k
 30  31  30   0   0   9|4.14 4.14 4.33|   0     0 |  13M   14M| 124k  264k
 30  30  31   0   0   9|4.14 4.14 4.33|   0     0 |  13M   14M| 124k  265k
 31  30  30   0   0   9|4.13 4.14 4.33|   0    20k|  13M   14M| 124k  265k
 29  30  30   0   0  11|4.13 4.14 4.33|   0     0 |  14M   14M| 124k  265k
 30  30  30   0   0  10|4.13 4.14 4.33|   0     0 |  14M   14M| 125k  266k
 30  30  31   0   0   9|4.13 4.14 4.33|   0     0 |  14M   14M| 125k  267k
 30  30  30   0   0  10|4.13 4.14 4.33|   0     0 |  14M   14M| 126k  268k
 29  31  31   0   0   9|4.20 4.15 4.33|   0     0 |  14M   14M| 126k  268k
 30  30  30   0   0   9|4.20 4.15 4.33|   0    12k|  13M   13M| 124k  263k
 29  32  30   0   0   9|4.20 4.15 4.33|   0     0 |  14M   14M| 125k  265k
 31  30  30   0   0   9|4.20 4.15 4.33|   0     0 |  13M   13M| 124k  262k
 30  30  31   0   0   9|4.20 4.15 4.33|   0     0 |  14M   14M| 125k  268k
 30  30  31   0   0   9|4.18 4.15 4.33|   0     0 |  14M   14M| 126k  268k
 30  31  30   0   0   8|4.18 4.15 4.33|   0    44k|  14M   14M| 125k  267k
 30  32  30   0   0   8|4.18 4.15 4.33|   0    32k|  14M   14M| 126k  268k
 30  32  31   0   0   8|4.18 4.15 4.33|   0     0 |  14M   14M| 126k  268k
 30  31  31   0   0   8|4.18 4.15 4.33|   0     0 |  14M   14M| 126k  267k
 30  32  31   0   0   8|4.01 4.11 4.31|   0     0 |  14M   14M| 125k  267k
 31  31  30   0   0   8|4.01 4.11 4.31|   0     0 |  14M   14M| 125k  266k
 30  32  31   0   0   8|4.01 4.11 4.31|   0     0 |  14M   14M| 126k  267k
time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[644s]       [r:13,w:0,u:0,d:0]  27943    0       27943    NaN        0.46         17872447

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[645s]       [r:13,w:0,u:0,d:0]  28176    0       28176    NaN        0.46         17900623

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[646s]       [r:13,w:0,u:0,d:0]  27925    0       27925    NaN        0.46         17928548

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[647s]       [r:13,w:0,u:0,d:0]  27924    0       27924    NaN        0.46         17956472

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[648s]       [r:13,w:0,u:0,d:0]  27957    0       27957    NaN        0.46         17984429

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[649s]       [r:13,w:0,u:0,d:0]  27897    0       27897    NaN        0.46         18012326

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[650s]       [r:13,w:0,u:0,d:0]  27626    0       27626    NaN        0.47         18039952

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[651s]       [r:13,w:0,u:0,d:0]  27898    0       27898    NaN        0.46         18067850

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[652s]       [r:13,w:0,u:0,d:0]  27964    0       27964    NaN        0.46         18095814

####connections:15
----total-cpu-usage---- ---load-avg--- -dsk/total- -net/total- ---system--
usr sys idl wai hiq siq| 1m   5m  15m | read  writ| recv  send| int   csw
 34  33  25   0   0   8|5.32 4.93 4.75|   0     0 |  15M   15M| 123k  275k
 33  35  25   0   0   8|5.32 4.93 4.75|   0     0 |  15M   15M| 125k  280k
 33  34  25   0   0   8|5.32 4.93 4.75|   0     0 |  15M   15M| 125k  278k
 32  35  24   0   0   8|5.32 4.93 4.75|   0     0 |  15M   15M| 124k  278k
 33  35  25   0   0   8|5.32 4.93 4.75|   0    16k|  15M   15M| 123k  274k
 32  35  24   0   0   9|5.22 4.92 4.74|   0     0 |  15M   15M| 124k  279k
 32  36  24   0   0   8|5.22 4.92 4.74|   0     0 |  15M   15M| 123k  276k
 34  35  25   0   0   7|5.22 4.92 4.74|   0     0 |  15M   15M| 125k  279k
 34  35  25   0   0   6|5.22 4.92 4.74|   0     0 |  15M   15M| 123k  273k
 33  35  25   0   0   8|5.22 4.92 4.74|   0    52k|  15M   15M| 124k  277k
 33  35  25   0   0   8|5.04 4.88 4.73|   0     0 |  15M   15M| 124k  279k
 33  37  25   0   0   6|5.04 4.88 4.73|   0     0 |  15M   15M| 124k  279k
 35  34  24   0   0   6|5.04 4.88 4.73|   0     0 |  15M   15M| 124k  277k
 34  36  25   0   0   7|5.04 4.88 4.73|   0     0 |  15M   15M| 124k  278k
 33  35  24   0   0   7|5.04 4.88 4.73|   0    80k|  15M   15M| 125k  279k
 34  34  25   0   0   7|5.12 4.90 4.74|   0     0 |  15M   15M| 124k  278k
 33  36  25   0   0   6|5.12 4.90 4.74|   0     0 |  15M   15M| 124k  276k
 34  34  25   0   0   8|5.12 4.90 4.74|   0     0 |  15M   15M| 124k  278k
 34  35  25   0   0   7|5.12 4.90 4.74|   0     0 |  15M   15M| 125k  280k
 33  34  25   0   0   7|5.12 4.90 4.74|   0     0 |  15M   15M| 123k  273k
 34  34  24   0   0   8|4.87 4.85 4.72|   0    32k|  15M   15M| 125k  279k
 33  35  24   0   0   8|4.87 4.85 4.72|   0     0 |  15M   15M| 124k  278k
 32  35  25   0   0   8|4.87 4.85 4.72|   0     0 |  15M   15M| 124k  278k
 33  34  24   0   0   8|4.87 4.85 4.72|   0     0 |  15M   15M| 124k  278k
 34  34  25   0   0   8|4.87 4.85 4.72|   0  4096B|  15M   15M| 124k  278k
 33  35  24   0   0   8|4.56 4.79 4.70|   0     0 |  15M   15M| 125k  280k
 35  32  25   0   0   8|4.56 4.79 4.70|   0     0 |  15M   15M| 125k  279k
 34  34  25   0   0   8|4.56 4.79 4.70|   0     0 |  15M   15M| 125k  278k
 33  33  25   0   0   9|4.56 4.79 4.70|   0     0 |  15M   15M| 125k  280k
 34  33  24   0   0   9|4.56 4.79 4.70|   0     0 |  15M   15M| 124k  278k
 32  33  25   0   0  10|4.67 4.81 4.71|   0  4096B|  15M   15M| 123k  276k
time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[5849s]      [r:15,w:0,u:0,d:0]  31618    0       31618    NaN        0.47         184116616

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[5850s]      [r:15,w:0,u:0,d:0]  31667    0       31667    NaN        0.47         184148283

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[5851s]      [r:15,w:0,u:0,d:0]  31710    0       31710    NaN        0.47         184179993

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[5852s]      [r:15,w:0,u:0,d:0]  31322    0       31322    NaN        0.48         184211315

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[5853s]      [r:15,w:0,u:0,d:0]  31703    0       31703    NaN        0.47         184243018

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[5854s]      [r:15,w:0,u:0,d:0]  31629    0       31629    NaN        0.47         184274647

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[5855s]      [r:15,w:0,u:0,d:0]  31764    0       31764    NaN        0.47         184306411

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[5856s]      [r:15,w:0,u:0,d:0]  31720    0       31720    NaN        0.47         184338131

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[5857s]      [r:15,w:0,u:0,d:0]  31711    0       31711    NaN        0.47         184369842

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[5858s]      [r:15,w:0,u:0,d:0]  31633    0       31633    NaN        0.47         184401475

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[5859s]      [r:15,w:0,u:0,d:0]  31676    0       31676    NaN        0.47         184433151

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[5860s]      [r:15,w:0,u:0,d:0]  31821    0       31821    NaN        0.47         184464972

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[5861s]      [r:15,w:0,u:0,d:0]  31734    0       31734    NaN        0.47         184496706

####connections:20
[nanxing@zerodb-proxy003 cobar-server-2.0.0]$ free -h
              total        used        free      shared  buff/cache   available
Mem:           7.6G        4.9G        287M        400K        2.5G        2.5G
Swap:            0B          0B          0B
----total-cpu-usage---- ---load-avg--- -dsk/total- -net/total- ---system--
usr sys idl wai hiq siq| 1m   5m  15m | read  writ| recv  send| int   csw
 33  37  17   0   0  12|5.82 6.13 6.02|   0     0 |  18M   18M| 122k  302k
 33  37  18   0   0  12|5.82 6.13 6.02|   0    76k|  18M   18M| 124k  301k
 35  36  18   0   0  11|5.82 6.13 6.02|   0     0 |  18M   18M| 124k  301k
 33  38  18   0   0  12|5.82 6.13 6.02|   0     0 |  18M   18M| 126k  304k
 33  36  18   0   0  13|5.82 6.13 6.02|   0     0 |  18M   18M| 126k  305k
 33  38  18   0   0  12|6.23 6.21 6.07|   0    56k|  18M   18M| 124k  304k
 33  37  18   0   0  12|6.23 6.21 6.07|   0     0 |  18M   18M| 124k  304k
 34  37  18   0   0  12|6.23 6.21 6.07|   0     0 |  18M   18M| 127k  304k
 32  37  18   0   0  13|6.23 6.21 6.07|   0     0 |  18M   18M| 128k  305k
 34  36  18   0   0  13|6.23 6.21 6.07|   0     0 |  18M   18M| 127k  304k
 32  38  18   0   0  12|6.13 6.19 6.07|   0     0 |  18M   18M| 127k  303k
 32  38  18   0   0  12|6.13 6.19 6.07|   0  8192B|  18M   18M| 124k  303k
 34  36  17   0   0  12|6.13 6.19 6.07|   0     0 |  18M   18M| 123k  303k
 32  37  18   0   0  13|6.13 6.19 6.07|   0    28k|  18M   19M| 123k  305k
 34  36  18   0   0  12|6.13 6.19 6.07|   0     0 |  18M   18M| 123k  301k
 33  36  18   0   0  13|6.12 6.19 6.07|   0     0 |  18M   19M| 124k  305k
 34  35  18   0   0  13|6.12 6.19 6.07|   0   400k|  18M   18M| 123k  303k
 32  37  18   0   0  13|6.12 6.19 6.07|   0     0 |  18M   18M| 125k  304k
 33  36  18   0   0  13|6.12 6.19 6.07|   0     0 |  18M   18M| 124k  304k
 34  37  17   0   0  12|6.12 6.19 6.07|   0    24k|  18M   18M| 123k  302k
 35  37  18   0   0  11|6.11 6.19 6.07|   0     0 |  18M   18M| 122k  304k
 33  37  18   0   0  12|6.11 6.19 6.07|   0     0 |  18M   18M| 123k  305k
 34  37  17   0   0  12|6.11 6.19 6.07|   0     0 |  18M   18M| 123k  304k
 32  38  18   0   0  12|6.11 6.19 6.07|   0     0 |  18M   18M| 125k  304k
 34  37  18   0   0  12|6.11 6.19 6.07|   0     0 |  18M   18M| 124k  304k
time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[2145s]      [r:20,w:0,u:0,d:0]  36462    0       36462    NaN        0.55         80528606

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[2146s]      [r:20,w:0,u:0,d:0]  36468    0       36468    NaN        0.55         80565074

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[2147s]      [r:20,w:0,u:0,d:0]  36946    0       36946    NaN        0.54         80602020

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[2148s]      [r:20,w:0,u:0,d:0]  36630    0       36630    NaN        0.54         80638650

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[2149s]      [r:20,w:0,u:0,d:0]  37942    0       37942    NaN        0.53         80676592

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[2150s]      [r:20,w:0,u:0,d:0]  36743    0       36743    NaN        0.54         80713335

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[2151s]      [r:20,w:0,u:0,d:0]  36577    0       36577    NaN        0.55         80749912

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[2152s]      [r:20,w:0,u:0,d:0]  37730    0       37730    NaN        0.53         80787642

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[2153s]      [r:20,w:0,u:0,d:0]  37589    0       37589    NaN        0.53         80825231

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[2154s]      [r:20,w:0,u:0,d:0]  37596    0       37596    NaN        0.53         80862827

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[2155s]      [r:20,w:0,u:0,d:0]  37917    0       37917    NaN        0.53         80900744

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[2156s]      [r:20,w:0,u:0,d:0]  37888    0       37888    NaN        0.53         80938632

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[2157s]      [r:20,w:0,u:0,d:0]  37484    0       37484    NaN        0.53         80976116

####connections:40
----total-cpu-usage---- ---load-avg--- -dsk/total- -net/total- ---system--
usr sys idl wai hiq siq| 1m   5m  15m | read  writ| recv  send| int   csw
 35  42   9   0   0  14|8.58 8.62 8.16|   0     0 |  23M   23M| 111k  309k
 36  41  10   0   0  14|8.58 8.62 8.16|   0     0 |  23M   23M| 110k  307k
 36  41   9   0   0  14|8.58 8.62 8.16|   0    60k|  23M   23M| 109k  307k
 35  42  10   0   0  14|8.58 8.62 8.16|   0     0 |  23M   23M| 110k  308k
 36  41   9   0   0  14|8.58 8.62 8.16|   0     0 |  23M   23M| 111k  311k
 36  41  10   0   0  14|9.09 8.72 8.19|   0     0 |  23M   23M| 110k  308k
 35  42  10   0   0  14|9.09 8.72 8.19|   0     0 |  23M   23M| 109k  305k
 36  42   9   0   0  13|9.09 8.72 8.19|   0    92k|  23M   23M| 111k  310k
 36  41   9   0   0  13|9.09 8.72 8.19|   0     0 |  23M   23M| 111k  309k
 35  42  10   0   0  13|9.09 8.72 8.19|   0     0 |  23M   23M| 109k  304k
 35  42   9   0   0  14|9.08 8.73 8.20|   0     0 |  23M   23M| 110k  308k
 36  42   9   0   0  13|9.08 8.73 8.20|   0    16k|  23M   23M| 111k  310k
 35  43   9   0   0  13|9.08 8.73 8.20|   0     0 |  23M   23M| 110k  310k
 35  42  10   0   0  13|9.08 8.73 8.20|   0     0 |  23M   23M| 111k  308k
 35  43   9   0   0  13|9.08 8.73 8.20|   0     0 |  23M   23M| 112k  311k
 35  42   9   0   0  14|8.76 8.66 8.18|   0     0 |  23M   23M| 111k  312k
 35  43   9   0   0  13|8.76 8.66 8.18|   0     0 |  23M   23M| 111k  309k
 37  40   9   0   0  14|8.76 8.66 8.18|   0  4096B|  23M   23M| 108k  308k
 35  43  10   0   0  12|8.76 8.66 8.18|   0     0 |  23M   23M| 110k  307k
 36  41  10   0   0  14|8.76 8.66 8.18|   0     0 |  23M   23M| 111k  308k
 35  42   9   0   0  14|8.94 8.70 8.20|   0     0 |  23M   23M| 111k  311k
 35  42   9   0   0  13|8.94 8.70 8.20|   0     0 |  23M   23M| 111k  308k
 37  40  10   0   0  13|8.94 8.70 8.20|   0    56k|  23M   23M| 111k  310k
 37  41   9   0   0  13|8.94 8.70 8.20|   0     0 |  23M   23M| 112k  311k
 34  43  10   0   0  14|8.94 8.70 8.20|   0     0 |  23M   23M| 110k  309k
 36  42   9   0   0  13|8.54 8.62 8.17|   0     0 |  23M   23M| 111k  310k
 36  41  10   0   0  13|8.54 8.62 8.17|   0     0 |  23M   23M| 111k  307k
 36  42   9   0   0  13|8.54 8.62 8.17|   0    12k|  23M   23M| 111k  312k
 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2306s]      [r:40,w:0,u:0,d:0]  47867    0       47867    NaN        0.83         109863583

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2307s]      [r:40,w:0,u:0,d:0]  47800    0       47800    NaN        0.84         109911383

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2308s]      [r:40,w:0,u:0,d:0]  48018    0       48018    NaN        0.83         109959401

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2309s]      [r:40,w:0,u:0,d:0]  47855    0       47855    NaN        0.83         110007256

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2310s]      [r:40,w:0,u:0,d:0]  48098    0       48098    NaN        0.83         110055354

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2311s]      [r:40,w:0,u:0,d:0]  48230    0       48230    NaN        0.83         110103584

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2312s]      [r:40,w:0,u:0,d:0]  48222    0       48222    NaN        0.83         110151806

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2313s]      [r:40,w:0,u:0,d:0]  47865    0       47865    NaN        0.83         110199671

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2314s]      [r:40,w:0,u:0,d:0]  47390    0       47390    NaN        0.84         110247061

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2315s]      [r:40,w:0,u:0,d:0]  46528    0       46528    NaN        0.86         110293589

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2316s]      [r:40,w:0,u:0,d:0]  48153    0       48153    NaN        0.83         110341742

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2317s]      [r:40,w:0,u:0,d:0]  48076    0       48076    NaN        0.83         110389818

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2318s]      [r:40,w:0,u:0,d:0]  47851    0       47851    NaN        0.83         110437669

 time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2319s]      [r:40,w:0,u:0,d:0]  47795    0       47795    NaN        0.84         110485464

####connections:80
 ----total-cpu-usage---- ---load-avg--- -dsk/total- -net/total- ---system--
 usr sys idl wai hiq siq| 1m   5m  15m | read  writ| recv  send| int   csw
  36  43   5   0   0  15|11.0 11.1 10.9|   0     0 |  26M   26M|  96k  306k
  36  43   5   0   0  16|11.0 11.1 10.9|   0     0 |  26M   27M|  98k  315k
  36  43   6   0   0  15|11.0 11.1 10.9|   0     0 |  26M   26M|  98k  309k
  36  43   6   0   0  16|11.0 11.1 10.9|   0    20k|  26M   27M|  98k  311k
  35  45   5   0   0  15|10.9 11.1 10.9|   0     0 |  27M   27M|  99k  318k
  35  45   6   0   0  15|10.9 11.1 10.9|   0     0 |  27M   27M| 100k  316k
  36  42   5   0   0  17|10.9 11.1 10.9|   0     0 |  27M   27M|  99k  318k
  35  43   5   0   0  16|10.9 11.1 10.9|   0     0 |  27M   27M|  99k  314k
  35  44   6   0   0  15|10.9 11.1 10.9|   0     0 |  26M   27M|  99k  312k
  34  44   5   0   0  17|12.1 11.3 11.0|   0     0 |  27M   27M| 100k  318k
  36  43   5   0   0  16|12.1 11.3 11.0|   0     0 |  27M   27M|  98k  317k
  36  43   6   0   0  15|12.1 11.3 11.0|   0     0 |  26M   27M|  98k  312k
  35  44   5   0   0  16|12.1 11.3 11.0|   0     0 |  27M   27M|  98k  314k
  36  43   5   0   0  16|12.1 11.3 11.0|   0     0 |  26M   27M|  98k  310k
  36  43   6   0   0  15|12.2 11.4 11.0|   0    16k|  26M   26M|  98k  309k
  35  43   6   0   0  15|12.2 11.4 11.0|   0    16k|  26M   27M|  96k  311k
  35  44   6   0   0  16|12.2 11.4 11.0|   0     0 |  26M   27M|  98k  311k
  35  45   5   0   0  16|12.2 11.4 11.0|   0     0 |  27M   27M| 100k  314k
  36  43   6   0   0  15|12.2 11.4 11.0|   0     0 |  26M   26M|  97k  305k
  37  43   5   0   0  16|12.3 11.4 11.0|   0   188k|  27M   27M|  98k  319k
  35  44   6   0   0  16|12.3 11.4 11.0|   0     0 |  27M   27M|  99k  314k
  36  43   5   0   0  16|12.3 11.4 11.0|   0     0 |  27M   27M|  98k  314k
time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[8833s]      [r:80,w:0,u:0,d:0]  54719    0       54719    NaN        1.46         482087521

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[8834s]      [r:80,w:0,u:0,d:0]  55065    0       55065    NaN        1.45         482142586

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[8835s]      [r:80,w:0,u:0,d:0]  55086    0       55086    NaN        1.45         482197672

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[8836s]      [r:80,w:0,u:0,d:0]  54909    0       54909    NaN        1.45         482252581

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[8837s]      [r:80,w:0,u:0,d:0]  55065    0       55065    NaN        1.45         482307646

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[8838s]      [r:80,w:0,u:0,d:0]  54170    0       54170    NaN        1.48         482361816

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[8839s]      [r:80,w:0,u:0,d:0]  54854    0       54854    NaN        1.46         482416670

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[8840s]      [r:80,w:0,u:0,d:0]  54991    0       54991    NaN        1.45         482471661

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[8841s]      [r:80,w:0,u:0,d:0]  54719    0       54719    NaN        1.46         482526380

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[8842s]      [r:80,w:0,u:0,d:0]  55155    0       55155    NaN        1.45         482581535
####connections:160
----total-cpu-usage---- ---load-avg--- -dsk/total- -net/total- ---system--
usr sys idl wai hiq siq| 1m   5m  15m | read  writ| recv  send| int   csw
 35  44   3   0   0  18|13.4 14.1 14.4|   0     0 |  30M   30M|  89k  323k
 36  44   3   0   0  18|13.4 14.1 14.4|   0     0 |  30M   30M|  90k  328k
 35  45   2   0   0  17|13.4 14.1 14.4|   0     0 |  29M   30M|  88k  331k
 36  44   3   0   0  18|13.4 14.1 14.4|   0     0 |  30M   30M|  90k  324k
 33  44   4   0   0  18|13.4 14.1 14.4|   0     0 |  29M   29M|  88k  312k
 34  44   3   0   0  20|13.4 14.1 14.4|   0     0 |  30M   30M|  88k  321k
 34  43   3   0   0  20|13.4 14.1 14.4|   0     0 |  30M   30M|  91k  323k
 35  43   4   0   0  18|13.4 14.1 14.4|   0     0 |  29M   29M|  92k  314k
 34  44   2   0   0  20|13.4 14.1 14.4|   0     0 |  30M   30M|  88k  329k
 34  42   3   0   0  21|13.4 14.1 14.4|   0    12k|  30M   30M|  89k  321k
 34  42   3   0   0  21|13.4 14.1 14.4|   0     0 |  30M   30M|  87k  322k
 33  44   3   0   0  20|13.4 14.1 14.4|   0     0 |  29M   30M|  87k  321k
 34  43   3   0   0  20|13.4 14.1 14.4|   0   328k|  29M   30M|  89k  326k
 34  44   3   0   0  20|13.8 14.1 14.4|   0     0 |  29M   29M|  87k  320k
 35  43   3   0   0  19|13.8 14.1 14.4|   0     0 |  29M   30M|  87k  320k
 36  43   3   0   0  18|13.8 14.1 14.4|   0     0 |  30M   30M|  90k  317k
 37  44   3   0   0  16|13.8 14.1 14.4|   0    32k|  30M   30M|  89k  320k
 37  45   2   0   0  16|13.8 14.1 14.4|   0     0 |  30M   30M|  87k  325k
 36  45   2   0   0  16|13.8 14.1 14.4|   0     0 |  30M   30M|  87k  327k
 37  45   3   0   0  15|13.8 14.1 14.4|   0     0 |  29M   29M|  87k  309k
 35  46   3   0   0  16|13.8 14.1 14.4|   0     0 |  30M   30M|  89k  327k
 36  45   4   0   0  16|13.8 14.1 14.4|   0     0 |  29M   29M|  91k  315k
 37  45   3   0   0  15|13.8 14.1 14.4|   0    32k|  29M   29M|  91k  318k
 36  46   2   0   0  16|13.6 14.1 14.4|   0     0 |  30M   30M|  89k  324k
 37  45   3   0   0  15|13.6 14.1 14.4|   0     0 |  30M   30M|  90k  327k
 38  43   3   0   0  16|13.6 14.1 14.4|   0     0 |  30M   30M|  89k  326k
 37  44   3   0   0  15|13.6 14.1 14.4|   0     0 |  29M   29M|  89k  312k
time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[2066s]      [r:160,w:0,u:0,d:0]  59696    0       59696    NaN        2.68         125160529

time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[2067s]      [r:160,w:0,u:0,d:0]  59867    0       59867    NaN        2.66         125220396

time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[2068s]      [r:160,w:0,u:0,d:0]  60983    0       60983    NaN        2.62         125281379

time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[2069s]      [r:160,w:0,u:0,d:0]  58737    0       58737    NaN        2.73         125340116

time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[2070s]      [r:160,w:0,u:0,d:0]  60945    0       60945    NaN        2.60         125401061

time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[2071s]      [r:160,w:0,u:0,d:0]  60399    0       60399    NaN        2.67         125461460

time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[2072s]      [r:160,w:0,u:0,d:0]  60188    0       60188    NaN        2.65         125521648

time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[2073s]      [r:160,w:0,u:0,d:0]  60806    0       60806    NaN        2.63         125582454

time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[2074s]      [r:160,w:0,u:0,d:0]  59853    0       59853    NaN        2.68         125642307

time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[2075s]      [r:160,w:0,u:0,d:0]  60002    0       60002    NaN        2.67         125702309

time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[2076s]      [r:160,w:0,u:0,d:0]  58868    0       58868    NaN        2.71         125761177
####connections:200
----total-cpu-usage---- ---load-avg--- -dsk/total- -net/total- ---system--
usr sys idl wai hiq siq| 1m   5m  15m | read  writ| recv  send| int   csw
 36  45   3   0   0  16|16.8 16.3 16.1|   0     0 |  30M   30M|  86k  321k
 37  45   3   0   0  16|16.8 16.3 16.1|   0     0 |  30M   31M|  88k  323k
 36  47   2   0   0  16|16.8 16.3 16.1|   0     0 |  31M   31M|  85k  333k
 37  45   2   0   0  16|16.8 16.3 16.1|   0     0 |  30M   31M|  87k  326k
 37  45   2   0   0  16|16.8 16.3 16.1|   0     0 |  30M   30M|  86k  328k
 36  46   3   0   0  16|18.1 16.6 16.2|   0    20k|  30M   30M|  86k  320k
 36  46   3   0   0  15|18.1 16.6 16.2|   0     0 |  30M   31M|  87k  326k
 37  45   2   0   0  16|18.1 16.6 16.2|   0     0 |  30M   30M|  84k  330k
 37  45   2   0   0  16|18.1 16.6 16.2|   0   396k|  30M   30M|  88k  321k
 38  45   2   0   0  16|18.1 16.6 16.2|   0  4096B|  30M   30M|  87k  324k
 36  45   3   0   0  16|18.2 16.6 16.2|   0     0 |  30M   31M|  88k  323k
 37  45   2   0   0  16|18.2 16.6 16.2|   0     0 |  30M   30M|  86k  325k
 36  46   3   0   0  16|18.2 16.6 16.2|   0     0 |  30M   30M|  88k  327k
 37  45   2   0   0  16|18.2 16.6 16.2|   0   132k|  30M   31M|  86k  324k
 37  45   3   0   0  16|18.2 16.6 16.2|   0   184k|  30M   30M|  86k  321k
 37  46   2   0   0  15|17.6 16.6 16.2|   0     0 |  30M   30M|  86k  325k
 35  46   3   0   0  16|17.6 16.6 16.2|   0     0 |  30M   30M|  86k  320k
 36  46   2   0   0  16|17.6 16.6 16.2|   0     0 |  30M   31M|  87k  325k
 38  44   2   0   0  16|17.6 16.6 16.2|   0     0 |  30M   30M|  85k  321k
 37  45   2   0   0  16|17.6 16.6 16.2|   0    36k|  30M   31M|  85k  325k
 36  46   2   0   0  16|17.1 16.5 16.1|   0     0 |  30M   31M|  85k  324k
 36  44   3   0   0  17|17.1 16.5 16.1|   0     0 |  30M   30M|  87k  323k
 35  46   3   0   0  16|17.1 16.5 16.1|   0     0 |  30M   30M|  86k  319k
 34  47   2   0   0  17|17.1 16.5 16.1|   0     0 |  30M   30M|  85k  327k
 37  44   2   0   0  17|17.1 16.5 16.1|   0  4096B|  30M   31M|  88k  323k
 37  45   3   0   0  16|17.0 16.4 16.1|   0     0 |  30M   31M|  85k  323k
 37  44   2   0   0  17|17.0 16.4 16.1|   0     0 |  30M   31M|  87k  323k
 36  45   3   0   0  16|17.0 16.4 16.1|   0    16k|  30M   30M|  88k  320k
 34  48   2   0   0  16|17.0 16.4 16.1|   0     0 |  30M   30M|  86k  329k
 38  44   2   0   0  16|17.0 16.4 16.1|   0     0 |  30M   30M|  85k  329k
 time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2525s]      [r:200,w:0,u:0,d:0]  62893    0       62893    NaN        3.18         156776510

 time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2526s]      [r:200,w:0,u:0,d:0]  62818    0       62818    NaN        3.18         156839328

 time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2527s]      [r:200,w:0,u:0,d:0]  61195    0       61195    NaN        3.25         156900523

 time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2528s]      [r:200,w:0,u:0,d:0]  61689    0       61689    NaN        3.24         156962212

 time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2529s]      [r:200,w:0,u:0,d:0]  62761    0       62761    NaN        3.20         157024973

 time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2530s]      [r:200,w:0,u:0,d:0]  61982    0       61982    NaN        3.21         157086955

 time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2531s]      [r:200,w:0,u:0,d:0]  62440    0       62440    NaN        3.22         157149395

 time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2532s]      [r:200,w:0,u:0,d:0]  62882    0       62882    NaN        3.18         157212277

 time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2533s]      [r:200,w:0,u:0,d:0]  61662    0       61662    NaN        3.22         157273939

 time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2534s]      [r:200,w:0,u:0,d:0]  62177    0       62177    NaN        3.21         157336116

 time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2535s]      [r:200,w:0,u:0,d:0]  60173    0       60173    NaN        3.33         157396289

 time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2536s]      [r:200,w:0,u:0,d:0]  63255    0       63255    NaN        3.17         157459544

 time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2537s]      [r:200,w:0,u:0,d:0]  62386    0       62386    NaN        3.21         157521930

 time            thds               tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
 [2538s]      [r:200,w:0,u:0,d:0]  61624    0       61624    NaN        3.23         157583554

####connections:20000
----total-cpu-usage---- ---load-avg--- -dsk/total- -net/total- ---system--
usr sys idl wai hiq siq| 1m   5m  15m | read  writ| recv  send| int   csw
 36  45   0   0   0  19|22.1 22.1 22.5|   0  4096B|  27M   31M|  64k  291k
 36  44   1   0   0  19|22.0 22.1 22.5|   0     0 |  28M   32M|  63k  299k
 36  45   0   0   0  19|22.0 22.1 22.5|   0     0 |  28M   33M|  62k  302k
 37  45   0   0   0  18|22.0 22.1 22.5|   0    64k|  27M   29M|  65k  362k
 36  44   1   0   0  19|22.0 22.1 22.5|   0     0 |  28M   34M|  61k  313k
 37  45   1   0   0  18|22.0 22.1 22.5|   0     0 |  28M   31M|  65k  297k
 38  43   1   0   0  19|23.9 22.5 22.6|   0     0 |  27M   31M|  62k  300k
 37  45   1   0   0  18|23.9 22.5 22.6|   0     0 |  26M   29M|  60k  278k
 36  44   1   0   0  20|23.9 22.5 22.6|   0     0 |  28M   35M|  65k  293k
 36  45   1   0   0  19|23.9 22.5 22.6|   0     0 |  28M   32M|  63k  305k
 37  45   0   0   0  19|23.9 22.5 22.6|   0    20k|  28M   31M|  63k  299k
 37  45   1   0   0  18|25.6 22.9 22.7|   0     0 |  29M   33M|  57k  296k
 36  45   1   0   0  19|25.6 22.9 22.7|   0     0 |  28M   33M|  60k  298k
 36  43   3   0   0  18|25.6 22.9 22.7|   0     0 |  26M   29M|  63k  259k
 38  43   1   0   0  18|25.6 22.9 22.7|   0     0 |  26M   30M|  60k  257k
 37  46   0   0   0  17|25.6 22.9 22.7|   0    24k|  28M   29M|  67k  282k
 36  45   0   0   0  19|25.2 22.9 22.7|   0     0 |  29M   35M|  61k  302k
 36  45   0   0   0  19|25.2 22.9 22.7|   0     0 |  29M   34M|  59k  291k
 37  44   1   0   0  18|25.2 22.9 22.7|   0     0 |  27M   29M|  66k  306k
 36  45   1   0   0  18|25.2 22.9 22.7|   0     0 |  27M   31M|  65k  326k
 36  45   0   0   0  19|25.2 22.9 22.7|   0     0 |  27M   31M|  64k  348k
 36  46   0   0   0  18|24.3 22.7 22.7|   0    12k|  27M   33M|  60k  333k
 40  41   2   0   0  18|24.3 22.7 22.7|   0     0 |  27M   30M|  63k  281k
 36  46   1   0   0  18|24.3 22.7 22.7|   0     0 |  27M   29M|  70k  318k
 35  46   1   0   0  19|24.3 22.7 22.7|   0     0 |  28M   34M|  63k  280k
 36  46   1   0   0  18|24.3 22.7 22.7|   0    16k|  27M   29M|  71k  331k
 35  45   2   0   0  19|23.9 22.7 22.6|   0     0 |  29M   36M|  67k  293k
 37  44   1   0   0  18|23.9 22.7 22.6|   0    20k|  27M   29M|  71k  293k
time            thds                 tps     wtps    rtps    rio    rio/op   wio    wio/op    rMB     rKB/op    wMB     wKB/op   cpu/op  freeMB  cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
[2399s]      [r:20000,w:0,u:0,d:0]  51141    0       51141   0      0.00     0      0.00      0.00    0.00      0.00    0.00     0.00    0       0         NaN        335.80       134511302

time            thds                 tps     wtps    rtps    rio    rio/op   wio    wio/op    rMB     rKB/op    wMB     wKB/op   cpu/op  freeMB  cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
[2400s]      [r:20000,w:0,u:0,d:0]  57689    0       57689   0      0.00     0      0.00      0.00    0.00      0.00    0.00     0.00    0       0         NaN        383.77       134568991

time            thds                 tps     wtps    rtps    rio    rio/op   wio    wio/op    rMB     rKB/op    wMB     wKB/op   cpu/op  freeMB  cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
[2401s]      [r:20000,w:0,u:0,d:0]  57432    0       57432   0      0.00     0      0.00      0.00    0.00      0.00    0.00     0.00    0       0         NaN        323.58       134626423

time            thds                 tps     wtps    rtps    rio    rio/op   wio    wio/op    rMB     rKB/op    wMB     wKB/op   cpu/op  freeMB  cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
[2402s]      [r:20000,w:0,u:0,d:0]  57125    0       57125   0      0.00     0      0.00      0.00    0.00      0.00    0.00     0.00    0       0         NaN        382.47       134683548

time            thds                 tps     wtps    rtps    rio    rio/op   wio    wio/op    rMB     rKB/op    wMB     wKB/op   cpu/op  freeMB  cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
[2403s]      [r:20000,w:0,u:0,d:0]  50733    0       50733   0      0.00     0      0.00      0.00    0.00      0.00    0.00     0.00    0       0         NaN        340.72       134734281

time            thds                 tps     wtps    rtps    rio    rio/op   wio    wio/op    rMB     rKB/op    wMB     wKB/op   cpu/op  freeMB  cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
[2404s]      [r:20000,w:0,u:0,d:0]  56795    0       56795   0      0.00     0      0.00      0.00    0.00      0.00    0.00     0.00    0       0         NaN        386.29       134791076

time            thds                 tps     wtps    rtps    rio    rio/op   wio    wio/op    rMB     rKB/op    wMB     wKB/op   cpu/op  freeMB  cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
[2405s]      [r:20000,w:0,u:0,d:0]  55593    0       55593   0      0.00     0      0.00      0.00    0.00      0.00    0.00     0.00    0       0         NaN        344.32       134846669

time            thds                 tps     wtps    rtps    rio    rio/op   wio    wio/op    rMB     rKB/op    wMB     wKB/op   cpu/op  freeMB  cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
[2406s]      [r:20000,w:0,u:0,d:0]  55936    0       55936   0      0.00     0      0.00      0.00    0.00      0.00    0.00     0.00    0       0         NaN        370.38       134902605

time            thds                 tps     wtps    rtps    rio    rio/op   wio    wio/op    rMB     rKB/op    wMB     wKB/op   cpu/op  freeMB  cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
[2407s]      [r:20000,w:0,u:0,d:0]  57261    0       57261   0      0.00     0      0.00      0.00    0.00      0.00    0.00     0.00    0       0         NaN        293.96       134959866

time            thds                 tps     wtps    rtps    rio    rio/op   wio    wio/op    rMB     rKB/op    wMB     wKB/op   cpu/op  freeMB  cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
[2408s]      [r:20000,w:0,u:0,d:0]  61975    0       61975   0      0.00     0      0.00      0.00    0.00      0.00    0.00     0.00    0       0         NaN        370.41       135021841

time            thds                 tps     wtps    rtps    rio    rio/op   wio    wio/op    rMB     rKB/op    wMB     wKB/op   cpu/op  freeMB  cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
[2409s]      [r:20000,w:0,u:0,d:0]  53399    0       53399   0      0.00     0      0.00      0.00    0.00      0.00    0.00     0.00    0       0         NaN        359.95       135075240

time            thds                 tps     wtps    rtps    rio    rio/op   wio    wio/op    rMB     rKB/op    wMB     wKB/op   cpu/op  freeMB  cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
[2410s]      [r:20000,w:0,u:0,d:0]  59638    0       59638   0      0.00     0      0.00      0.00    0.00      0.00    0.00     0.00    0       0         NaN        281.49       135134878
