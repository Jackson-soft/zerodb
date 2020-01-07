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

全表负载 100 万

###
####connections:10
zeroProxy
----total-cpu-usage---- ---load-avg--- -dsk/total- -net/total- ---system--
usr sys idl wai hiq siq| 1m   5m  15m | read  writ| recv  send| int   csw
 40  15  40   0   0   5|2.49 2.51 2.59|   0     0 |  12M   12M|  73k   86k
 40  15  40   0   0   5|2.49 2.51 2.59|   0     0 |  12M   12M|  72k   84k
 40  16  41   0   0   4|2.49 2.51 2.59|   0     0 |  12M   12M|  73k   85k
 39  15  41   0   0   5|2.49 2.51 2.59|   0    12k|  11M   12M|  71k   83k
 41  15  40   0   0   5|2.61 2.53 2.60|   0     0 |  12M   12M|  71k   82k
 41  14  40   0   0   5|2.61 2.53 2.60|   0     0 |  12M   12M|  73k   86k
 41  15  39   0   0   5|2.61 2.53 2.60|   0     0 |  11M   12M|  74k   88k
 40  15  40   0   0   5|2.61 2.53 2.60|   0     0 |  12M   12M|  76k   92k
 40  16  39   0   0   5|2.61 2.53 2.60|   0     0 |  11M   12M|  75k   91k
 40  16  39   0   0   5|2.48 2.50 2.59|   0    40k|  11M   12M|  70k   80k
 40  16  39   0   0   6|2.48 2.50 2.59|   0  4096B|  11M   12M|  71k   83k
 39  16  40   0   0   5|2.48 2.50 2.59|   0     0 |  11M   12M|  67k   77k
 41  15  40   0   0   5|2.48 2.50 2.59|   0     0 |  11M   12M|  68k   78k
 39  16  40   0   0   5|2.48 2.50 2.59|   0    68k|  12M   12M|  72k   84k
 39  17  40   0   0   4|2.52 2.51 2.59|   0     0 |  12M   12M|  74k   88k
 39  15  40   0   0   5|2.52 2.51 2.59|   0     0 |  12M   12M|  72k   84k
 39  16  40   0   0   5|2.52 2.51 2.59|   0     0 |  11M   12M|  72k   85k
 39  16  41   0   0   4|2.52 2.51 2.59|   0     0 |  11M   12M|  72k   84k
 40  16  41   0   0   4|2.52 2.51 2.59|   0    12k|  12M   12M|  71k   82k
 40  15  40   0   0   5|2.56 2.52 2.59|   0     0 |  12M   12M|  72k   84k
 39  16  40   0   0   5|2.56 2.52 2.59|   0     0 |  11M   12M|  75k   91k
 40  16  40   0   0   5|2.56 2.52 2.59|   0     0 |  12M   12M|  74k   87k
 38  16  41   0   0   5|2.56 2.52 2.59|   0     0 |  11M   12M|  73k   86k
 40  16  39   0   0   5|2.56 2.52 2.59|   0     0 |  12M   12M|  75k   90k
 40  16  39   0   0   5|2.67 2.55 2.60|   0  8192B|  11M   12M|  78k   94k
 39  16  41   0   0   5|2.67 2.55 2.60|   0    28k|  11M   12M|  75k   90k
 40  15  40   0   0   5|2.67 2.55 2.60|   0     0 |  11M   12M|  74k   89k
 41  15  40   0   0   4|2.67 2.55 2.60|   0     0 |  11M   12M|  74k   90k
 39  16  40   0   0   5|2.67 2.55 2.60|   0     0 |  11M   12M|  75k   89k
mysql:
----total-cpu-usage---- ---load-avg--- -dsk/total- -net/total- ---system--
usr sys idl wai hiq siq| 1m   5m  15m | read  writ| recv  send| int   csw
 33   7  59   0   0   2|1.10 0.87 0.93|   0    16k| 770k 2208k|9773    12k
 33   6  59   1   0   2|1.10 0.87 0.93|   0    20k| 783k 2245k|9806    13k
 30   6  63   0   0   2|1.10 0.87 0.93|   0    20k| 771k 2210k|9663    12k
 33   5  61   0   0   2|1.09 0.87 0.93|   0    16k| 769k 2205k|9663    12k
 34   6  59   0   0   2|1.09 0.87 0.93|   0    16k| 769k 2204k|9792    12k
 33   5  61   0   0   1|1.09 0.87 0.93|   0    16k| 780k 2234k|9788    13k
 32   6  61   0   0   2|1.09 0.87 0.93|   0    16k| 768k 2201k|9675    13k
 32   6  61   0   0   2|1.09 0.87 0.93|   0    20k| 775k 2220k|9699    13k
 33   6  61   0   0   1|1.16 0.89 0.93|   0    20k| 778k 2227k|9872    13k
 32   6  60   0   0   2|1.16 0.89 0.93|   0    16k| 780k 2234k|9805    13k
 33   6  59   1   0   2|1.16 0.89 0.93|   0    16k| 778k 2231k|9742    12k
 33   5  60   0   0   2|1.16 0.89 0.93|   0    16k| 785k 2251k|9812    12k
 33   6  59   0   0   2|1.16 0.89 0.93|   0    20k| 771k 2225k|9893    13k
 31   5  62   0   0   2|1.07 0.88 0.93|   0    20k| 779k 2233k|9787    13k
 33   6  59   0   0   2|1.07 0.88 0.93|   0    16k| 784k 2245k|9805    12k
 32   5  61   0   0   2|1.07 0.88 0.93|   0    16k| 779k 2231k|9777    12k
 31   6  61   0   0   2|1.07 0.88 0.93|   0    20k| 758k 2176k|9566    12k
 32   7  59   0   0   2|1.07 0.88 0.93|   0    20k| 767k 2196k|9642    12k
 33   6  59   0   0   2|1.14 0.89 0.93|   0    44k| 764k 2191k|9594    13k
 31   6  62   0   0   1|1.14 0.89 0.93|   0    32k| 787k 2254k|9768    13k
 34   5  59   1   0   2|1.14 0.89 0.93|   0    16k| 767k 2200k|9680    12k
 33   6  59   0   0   2|1.14 0.89 0.93|   0    16k| 782k 2239k|9855    13k
 33   6  60   0   0   2|1.14 0.89 0.93|   0    20k| 767k 2200k|9797    13k
 36   6  57   0   0   2|1.13 0.90 0.93|   0    16k| 776k 2223k|9757    12k
 34   6  59   0   0   2|1.13 0.90 0.93|   0    16k| 781k 2237k|9805    12k
 35   5  58   0   0   2|1.13 0.90 0.93|   0    16k| 769k 2202k|9726    12k
time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[7888s]      [r:10,w:0,u:0,d:0]  26156    0       26156    NaN        0.38         206837514

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[7889s]      [r:10,w:0,u:0,d:0]  25822    0       25822    NaN        0.39         206863336

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[7890s]      [r:10,w:0,u:0,d:0]  25875    0       25875    NaN        0.39         206889211

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[7891s]      [r:10,w:0,u:0,d:0]  26237    0       26237    NaN        0.38         206915448

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[7892s]      [r:10,w:0,u:0,d:0]  25904    0       25904    NaN        0.38         206941352

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[7893s]      [r:10,w:0,u:0,d:0]  26147    0       26147    NaN        0.38         206967499

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[7894s]      [r:10,w:0,u:0,d:0]  26058    0       26058    NaN        0.38         206993557

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[7895s]      [r:10,w:0,u:0,d:0]  25954    0       25954    NaN        0.38         207019511

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[7896s]      [r:10,w:0,u:0,d:0]  26089    0       26089    NaN        0.38         207045600

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[7897s]      [r:10,w:0,u:0,d:0]  26316    0       26316    NaN        0.38         207071916

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[7898s]      [r:10,w:0,u:0,d:0]  25913    0       25913    NaN        0.38         207097829

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[7899s]      [r:10,w:0,u:0,d:0]  26065    0       26065    NaN        0.38         207123894

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[7900s]      [r:10,w:0,u:0,d:0]  26060    0       26060    NaN        0.38         207149954

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[7901s]      [r:10,w:0,u:0,d:0]  26194    0       26194    NaN        0.38         207176148

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[7902s]      [r:10,w:0,u:0,d:0]  25464    0       25464    NaN        0.39         207201612

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[7903s]      [r:10,w:0,u:0,d:0]  26061    0       26061    NaN        0.38         207227673

####connections:20
proxy:
----total-cpu-usage---- ---load-avg--- -dsk/total- -net/total- ---system--
usr sys idl wai hiq siq| 1m   5m  15m | read  writ| recv  send| int   csw
 56  14  22   0   0   8|3.07 3.61 3.44|   0    28k|  19M   19M|  61k   26k
 54  16  23   0   0   8|3.07 3.61 3.44|   0  4096B|  19M   19M|  62k   26k
 56  15  21   0   0   8|3.07 3.61 3.44|   0     0 |  19M   19M|  62k   27k
 54  14  23   0   0   9|3.06 3.60 3.44|   0     0 |  19M   19M|  60k   24k
 55  15  22   0   0   8|3.06 3.60 3.44|   0     0 |  19M   19M|  60k   26k
 55  15  22   0   0   8|3.06 3.60 3.44|   0     0 |  19M   19M|  61k   25k
 54  15  23   0   0   8|3.06 3.60 3.44|   0    20k|  19M   19M|  61k   26k
 56  14  22   0   0   8|3.06 3.60 3.44|   0     0 |  19M   19M|  61k   25k
 55  13  23   0   0   9|3.14 3.61 3.44|   0     0 |  19M   19M|  61k   26k
 53  15  23   0   0   9|3.14 3.61 3.44|   0     0 |  19M   19M|  61k   25k
 53  15  22   0   0   9|3.14 3.61 3.44|   0     0 |  19M   19M|  62k   26k
 54  15  22   0   0   9|3.14 3.61 3.44|   0  8192B|  19M   19M|  62k   27k
 53  15  23   0   0   9|3.14 3.61 3.44|   0     0 |  19M   19M|  61k   27k
 54  15  22   0   0   9|3.21 3.62 3.44|   0     0 |  19M   20M|  61k   24k
 53  15  23   0   0   9|3.21 3.62 3.44|   0     0 |  19M   19M|  61k   25k
 53  15  22   0   0   9|3.21 3.62 3.44|   0    12k|  19M   19M|  61k   26k
 53  16  22   0   0   9|3.21 3.62 3.44|   0     0 |  19M   20M|  61k   24k
 53  16  22   0   0   9|3.21 3.62 3.44|   0     0 |  19M   19M|  61k   23k
 55  16  21   0   0   9|3.27 3.62 3.45|   0     0 |  19M   19M|  61k   25k
 55  15  22   0   0   7|3.27 3.62 3.45|   0     0 |  19M   19M|  61k   25k
 55  15  22   0   0   9|3.27 3.62 3.45|   0     0 |  19M   20M|  61k   24k
 53  16  23   0   0   8|3.27 3.62 3.45|   0     0 |  19M   19M|  61k   24k
 53  16  22   0   0   9|3.27 3.62 3.45|   0     0 |  19M   19M|  61k   25k
 54  15  23   0   0   9|3.33 3.63 3.45|   0     0 |  19M   19M|  61k   25k
 55  15  22   0   0   9|3.33 3.63 3.45|   0     0 |  19M   19M|  62k   26k
 55  14  23   0   0   8|3.33 3.63 3.45|   0    24k|  19M   19M|  61k   26k
 55  16  22   0   0   8|3.33 3.63 3.45|   0     0 |  19M   20M|  61k   24k
 55  15  22   0   0   8|3.33 3.63 3.45|   0     0 |  19M   19M|  62k   27k
 54  15  22   0   0   8|3.38 3.63 3.45|   0     0 |  19M   19M|  61k   25k
 55  15  22   0   0   8|3.38 3.63 3.45|   0     0 |  19M   19M|  60k   26k
 53  14  25   0   0   8|3.38 3.63 3.45|   0    12k|  18M   19M|  60k   28k
 55  16  22   0   0   7|3.38 3.63 3.45|   0    36k|  19M   19M|  61k   26k
mysql:
----total-cpu-usage---- ---load-avg--- -dsk/total- -net/total- ---system--
usr sys idl wai hiq siq| 1m   5m  15m | read  writ| recv  send| int   csw
 47   9  41   0   0   4|1.57 1.82 1.73|   0    16k|1210k 3466k|  14k   17k
 47   8  42   0   0   4|1.57 1.82 1.73|   0    20k|1213k 3491k|  14k   18k
 49   8  40   0   0   3|1.60 1.83 1.73|   0    20k|1233k 3535k|  15k   18k
 43   8  47   0   0   3|1.60 1.83 1.73|   0    16k|1129k 3234k|  14k   17k
 50   8  37   1   0   4|1.60 1.83 1.73|   0    16k|1256k 3600k|  15k   18k
 44  10  43   0   0   4|1.60 1.83 1.73|   0    56k|1143k 3276k|  14k   17k
 46   8  44   0   0   3|1.60 1.83 1.73|   0    24k|1084k 3105k|  13k   16k
 48   8  40   0   0   4|1.56 1.81 1.73|   0    32k|1213k 3475k|  14k   17k
 50   9  38   0   0   3|1.56 1.81 1.73|   0    16k|1284k 3677k|  15k   18k
 51   8  37   0   0   4|1.56 1.81 1.73|   0    16k|1304k 3740k|  15k   19k
 48   9  39   0   0   3|1.56 1.81 1.73|   0    96k|1289k 3692k|  15k   19k
 50  10  36   1   0   4|1.56 1.81 1.73|   0    20k|1300k 3728k|  15k   19k
 49   8  40   0   0   4|1.75 1.85 1.74|   0    16k|1193k 3420k|  14k   17k
 46   8  43   0   0   3|1.75 1.85 1.74|   0    16k|1121k 3215k|  13k   16k
 48   7  41   0   0   4|1.75 1.85 1.74|   0    20k|1109k 3176k|  13k   16k
 44   8  46   0   0   3|1.75 1.85 1.74|   0    16k|1104k 3164k|  13k   16k
 49   8  39   0   0   4|1.75 1.85 1.74|   0    24k|1153k 3306k|  14k   16k
 51   9  36   0   0   4|1.93 1.88 1.75|   0    16k|1164k 3336k|  14k   16k
 51   9  36   0   0   3|1.93 1.88 1.75|   0    16k|1274k 3651k|  15k   18k
 52   9  35   0   0   4|1.93 1.88 1.75|   0    16k|1297k 3715k|  15k   18k
 52   8  36   0   0   5|1.93 1.88 1.75|   0    16k|1295k 3712k|  15k   18k
 51   9  36   0   0   4|1.93 1.88 1.75|   0    20k|1294k 3708k|  15k   18k
 54   8  35   0   0   4|1.78 1.85 1.74|   0    16k|1272k 3647k|  15k   18k
 53   9  34   0   0   5|1.78 1.85 1.74|   0    16k|1273k 3647k|  15k   18k
 53   8  36   0   0   3|1.78 1.85 1.74|   0    16k|1274k 3653k|  15k   18k
 51   8  36   0   0   4|1.78 1.85 1.74|   0    24k|1278k 3663k|  15k   18k
 53   9  35   0   0   4|1.78 1.85 1.74|   0    20k|1274k 3651k|  15k   18k
 49   9  38   0   0   4|1.79 1.86 1.74|   0    16k|1270k 3637k|  15k   18k
 49   9  39   0   0   3|1.79 1.86 1.74|   0    16k|1241k 3558k|  15k   18k
 53   8  34   1   0   5|1.79 1.86 1.74|   0    16k|1291k 3700k|  15k   18k
 54   9  33   0   0   4|1.79 1.86 1.74|   0    16k|1268k 3634k|  15k   17k

time            thds              tps     wtps    rtps      w-rsp(ms)  r-rsp(ms)    total-number
[1374s]      [r:20,w:0,u:0,d:0]  42139    0       42139     NaN        0.47         59290415

time            thds              tps     wtps    rtps      w-rsp(ms)  r-rsp(ms)    total-number
[1375s]      [r:20,w:0,u:0,d:0]  43506    0       43506     NaN        0.46         59333921

time            thds              tps     wtps    rtps      w-rsp(ms)  r-rsp(ms)    total-number
[1376s]      [r:20,w:0,u:0,d:0]  43522    0       43522     NaN        0.46         59377443

time            thds              tps     wtps    rtps      w-rsp(ms)  r-rsp(ms)    total-number
[1377s]      [r:20,w:0,u:0,d:0]  43753    0       43753     NaN        0.46         59421196

time            thds              tps     wtps    rtps      w-rsp(ms)  r-rsp(ms)    total-number
[1378s]      [r:20,w:0,u:0,d:0]  43563    0       43563     NaN        0.46         59464759

time            thds              tps     wtps    rtps      w-rsp(ms)  r-rsp(ms)    total-number
[1379s]      [r:20,w:0,u:0,d:0]  43192    0       43192     NaN        0.46         59507951

time            thds              tps     wtps    rtps      w-rsp(ms)  r-rsp(ms)    total-number
[1380s]      [r:20,w:0,u:0,d:0]  43090    0       43090     NaN        0.46         59551041

time            thds              tps     wtps    rtps      w-rsp(ms)  r-rsp(ms)    total-number
[1381s]      [r:20,w:0,u:0,d:0]  42992    0       42992     NaN        0.46         59594033

time            thds              tps     wtps    rtps      w-rsp(ms)  r-rsp(ms)    total-number
[1382s]      [r:20,w:0,u:0,d:0]  43354    0       43354     NaN        0.46         59637387

time            thds              tps     wtps    rtps      w-rsp(ms)  r-rsp(ms)    total-number
[1383s]      [r:20,w:0,u:0,d:0]  43458    0       43458     NaN        0.46         59680845

time            thds              tps     wtps    rtps      w-rsp(ms)  r-rsp(ms)    total-number
[1384s]      [r:20,w:0,u:0,d:0]  43140    0       43140     NaN        0.46         59723985

time            thds              tps     wtps    rtps      w-rsp(ms)  r-rsp(ms)    total-number
[1385s]      [r:20,w:0,u:0,d:0]  43508    0       43508     NaN        0.46         59767493

time            thds              tps     wtps    rtps      w-rsp(ms)  r-rsp(ms)    total-number
[1386s]      [r:20,w:0,u:0,d:0]  43387    0       43387     NaN        0.46         59810880

time            thds              tps     wtps    rtps      w-rsp(ms)  r-rsp(ms)    total-number
[1387s]      [r:20,w:0,u:0,d:0]  43468    0       43468     NaN        0.46         59854348

####connections:40
proxy:
----total-cpu-usage---- ---load-avg--- -dsk/total- -net/total- ---system--
usr sys idl wai hiq siq| 1m   5m  15m | read  writ| recv  send| int   csw
 63  15  13   0   0   9|3.78 3.83 3.91|   0     0 |  22M   23M|  57k 6077
 63  15  14   0   0   8|3.88 3.85 3.91|   0     0 |  22M   22M|  57k 5753
 62  15  14   0   0   9|3.88 3.85 3.91|   0    84k|  22M   22M|  57k 5888
 63  15  13   0   0   8|3.88 3.85 3.91|   0   112k|  22M   22M|  57k 5831
 62  15  14   0   0   9|3.88 3.85 3.91|   0     0 |  22M   23M|  57k 6097
 62  15  14   0   0   9|3.88 3.85 3.91|   0     0 |  22M   23M|  57k 5717
 61  16  14   0   0   8|3.89 3.85 3.91|   0     0 |  22M   22M|  57k 6400
 62  16  13   0   0   9|3.89 3.85 3.91|   0    20k|  22M   22M|  57k 6351
 62  16  13   0   0   9|3.89 3.85 3.91|   0  4096B|  22M   22M|  56k 6064
 61  16  15   0   0   9|3.89 3.85 3.91|   0     0 |  22M   22M|  56k 6644
 61  16  14   0   0   9|3.89 3.85 3.91|   0     0 |  21M   22M|  56k 7014
 62  17  12   0   0   9|3.98 3.87 3.92|   0     0 |  22M   23M|  56k 5740
 62  16  14   0   0   8|3.98 3.87 3.92|   0    24k|  22M   22M|  56k 6353
 63  15  14   0   0   8|3.98 3.87 3.92|   0     0 |  22M   23M|  57k 5866
 61  16  14   0   0   9|3.98 3.87 3.92|   0     0 |  22M   22M|  57k 6192
 62  15  14   0   0   9|3.98 3.87 3.92|   0     0 |  22M   22M|  56k 5740
 61  16  14   0   0   9|3.98 3.87 3.92|   0     0 |  22M   22M|  57k 6295
 63  15  14   0   0   9|3.98 3.87 3.92|   0     0 |  22M   23M|  57k 6507
 61  16  14   0   0   9|3.98 3.87 3.92|   0     0 |  22M   22M|  56k 6295
 61  16  14   0   0   9|3.98 3.87 3.92|   0  4096B|  22M   23M|  57k 6101
 61  16  14   0   0   9|3.98 3.87 3.92|   0     0 |  22M   22M|  57k 6119
 63  15  14   0   0   9|3.98 3.88 3.92|   0     0 |  22M   23M|  57k 5875
 62  16  14   0   0   8|3.98 3.88 3.92|   0     0 |  22M   23M|  57k 5956
 62  16  14   0   0   8|3.98 3.88 3.92|   0    16k|  22M   22M|  57k 6126
 61  16  14   0   0   9|3.98 3.88 3.92|   0     0 |  22M   22M|  57k 6164
 61  16  14   0   0   9|3.98 3.88 3.92|   0     0 |  22M   22M|  56k 5953
 62  15  14   0   0   8|3.98 3.88 3.92|   0     0 |  22M   23M|  57k 6138
 60  16  14   0   0   9|3.98 3.88 3.92|   0     0 |  22M   23M|  57k 6521
 61  15  14   0   0  10|3.98 3.88 3.92|   0     0 |  22M   22M|  57k 6605
mysql:
----total-cpu-usage---- ---load-avg--- -dsk/total- -net/total- ---system--
usr sys idl wai hiq siq| 1m   5m  15m | read  writ| recv  send| int   csw
 59   9  28   0   0   5|3.21 3.17 3.12|   0    16k|1492k 4271k|  17k   20k
 58  10  27   0   0   5|3.21 3.17 3.12|   0    16k|1470k 4216k|  16k   19k
 61  10  24   0   0   5|3.21 3.17 3.12|   0    20k|1461k 4188k|  16k   19k
 58  11  26   0   0   5|3.21 3.17 3.12|   0    32k|1469k 4208k|  16k   19k
 57   9  30   0   0   4|3.21 3.17 3.12|   0    32k|1496k 4289k|  17k   20k
 62  11  23   0   0   5|3.28 3.18 3.13|   0    16k|1484k 4251k|  16k   19k
 58  10  27   0   0   6|3.28 3.18 3.13|   0    16k|1491k 4274k|  17k   20k
 56  10  28   0   0   5|3.28 3.18 3.13|   0    20k|1463k 4194k|  17k   20k
 59   9  26   0   0   5|3.28 3.18 3.13|   0    20k|1455k 4170k|  16k   19k
 54   9  31   1   0   5|3.28 3.18 3.13|   0    16k|1468k 4205k|  17k   20k
 59  11  23   1   0   6|3.89 3.31 3.17|   0    16k|1504k 4311k|  16k   18k
 60  10  26   0   0   5|3.89 3.31 3.17|   0    20k|1495k 4287k|  16k   19k
 59   8  28   0   0   5|3.89 3.31 3.17|   0    24k|1495k 4282k|  17k   19k
 58   9  27   0   0   6|3.89 3.31 3.17|   0    16k|1488k 4264k|  17k   20k
 54  11  30   0   0   5|3.89 3.31 3.17|   0    16k|1473k 4221k|  17k   20k
 60  10  26   0   0   5|3.74 3.29 3.16|   0    16k|1484k 4251k|  17k   20k
 57  11  27   0   0   5|3.74 3.29 3.16|   0    16k|1491k 4271k|  17k   19k
 55  10  30   0   0   5|3.74 3.29 3.16|   0    20k|1476k 4232k|  17k   20k
 59  10  27   1   0   5|3.74 3.29 3.16|   0    16k|1464k 4196k|  16k   20k
 57  11  28   0   0   5|3.74 3.29 3.16|   0    20k|1490k 4268k|  17k   20k
 57  11  27   0   0   5|3.76 3.30 3.17|   0    16k|1493k 4281k|  17k   20k
 57  10  27   0   0   6|3.76 3.30 3.17|   0    16k|1469k 4205k|  16k   19k
 59  10  27   0   0   4|3.76 3.30 3.17|   0    20k|1504k 4310k|  17k   20k
 57  10  28   0   0   5|3.76 3.30 3.17|   0    16k|1494k 4283k|  17k   20k
 58   9  28   0   0   5|3.76 3.30 3.17|   0    16k|1466k 4197k|  16k   19k
 56  10  30   0   0   4|3.78 3.31 3.17|   0    16k|1498k 4294k|  17k   20k
 55  10  30   0   0   4|3.78 3.31 3.17|   0    16k|1443k 4138k|  16k   19k
 58  11  26   0   0   6|3.78 3.31 3.17|   0    24k|1505k 4328k|  17k   20k
 54  10  30   0   0   6|3.78 3.31 3.17|   0    24k|1471k 4215k|  16k   18k
 61  11  22   0   0   6|3.78 3.31 3.17|   0    16k|1483k 4250k|  16k   19k
 59  10  26   0   0   6|4.04 3.37 3.19|   0    16k|1453k 4164k|  16k   19k
 56   9  30   0   0   5|4.04 3.37 3.19|   0    16k|1482k 4246k|  17k   20k
time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[6674s]      [r:40,w:0,u:0,d:0]  50687    0       50687    NaN        0.79         336403355

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[6675s]      [r:40,w:0,u:0,d:0]  51199    0       51199    NaN        0.78         336454554

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[6676s]      [r:40,w:0,u:0,d:0]  50595    0       50595    NaN        0.79         336505149

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[6677s]      [r:40,w:0,u:0,d:0]  50183    0       50183    NaN        0.80         336555332

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[6678s]      [r:40,w:0,u:0,d:0]  50893    0       50893    NaN        0.78         336606225

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[6679s]      [r:40,w:0,u:0,d:0]  51214    0       51214    NaN        0.78         336657439

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[6680s]      [r:40,w:0,u:0,d:0]  50279    0       50279    NaN        0.80         336707718

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[6681s]      [r:40,w:0,u:0,d:0]  50355    0       50355    NaN        0.79         336758073

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[6682s]      [r:40,w:0,u:0,d:0]  49355    0       49355    NaN        0.81         336807428

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[6683s]      [r:40,w:0,u:0,d:0]  49752    0       49752    NaN        0.80         336857180

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[6684s]      [r:40,w:0,u:0,d:0]  50602    0       50602    NaN        0.79         336907782

time            thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[6685s]      [r:40,w:0,u:0,d:0]  50090    0       50090    NaN        0.80         336957872

####connections:80
proxy:
----total-cpu-usage---- ---load-avg--- -dsk/total- -net/total- ---system--
usr sys idl wai hiq siq| 1m   5m  15m | read  writ| recv  send| int   csw
 65  16   9   0   0  10|4.24 4.19 4.11|   0     0 |  23M   24M|  57k 3407
 62  17   9   0   0  12|4.22 4.19 4.11|   0     0 |  23M   24M|  57k 3260
 65  16   9   0   0  11|4.22 4.19 4.11|   0     0 |  23M   24M|  57k 3549
 64  17   8   0   0  11|4.22 4.19 4.11|   0     0 |  23M   24M|  57k 3805
 64  16   9   0   0  12|4.22 4.19 4.11|   0     0 |  23M   24M|  57k 3367
 64  15   9   0   0  12|4.22 4.19 4.11|   0     0 |  23M   24M|  57k 3286
 62  17  10   0   0  11|4.20 4.19 4.11|   0     0 |  23M   24M|  57k 3397
 62  17   9   0   0  12|4.20 4.19 4.11|   0    16k|  23M   24M|  57k 3636
 63  17   9   0   0  12|4.20 4.19 4.11|   0  8192B|  23M   24M|  57k 3515
 63  17   8   0   0  12|4.20 4.19 4.11|   0     0 |  23M   24M|  57k 3305
 62  17   8   0   0  13|4.20 4.19 4.11|   0     0 |  23M   24M|  57k 3385
 64  17   8   0   0  11|4.19 4.18 4.11|   0     0 |  23M   24M|  57k 3277
 63  17   9   0   0  12|4.19 4.18 4.11|   0     0 |  23M   24M|  57k 3777
 63  16   9   0   0  12|4.19 4.18 4.11|   0    28k|  23M   24M|  56k 3954
 63  16   8   0   0  12|4.19 4.18 4.11|   0     0 |  23M   24M|  57k 3758
 62  16   9   0   0  12|4.19 4.18 4.11|   0     0 |  23M   24M|  57k 3244
 63  15   9   0   0  13|4.17 4.18 4.11|   0     0 |  23M   24M|  57k 3665
 63  17   9   0   0  11|4.17 4.18 4.11|   0     0 |  23M   24M|  57k 3357
 63  17   9   0   0  11|4.17 4.18 4.11|   0     0 |  23M   24M|  56k 3240
 64  15  10   0   0  12|4.17 4.18 4.11|   0     0 |  23M   24M|  57k 3446
 62  16  10   0   0  12|4.17 4.18 4.11|   0     0 |  23M   24M|  56k 3257
 63  17   8   0   0  12|4.16 4.18 4.11|   0     0 |  23M   24M|  57k 3376
 63  17   9   0   0  11|4.16 4.18 4.11|   0    12k|  23M   24M|  57k 3940
 64  17   9   0   0  11|4.16 4.18 4.11|   0     0 |  23M   24M|  57k 4032
 63  16  10   0   0  11|4.16 4.18 4.11|   0     0 |  23M   24M|  57k 3806
 63  16  10   0   0  11|4.16 4.18 4.11|   0     0 |  23M   23M|  56k 3415
 63  16  10   0   0  11|4.15 4.17 4.11|   0     0 |  23M   23M|  56k 3415
 64  17   9   0   0  11|4.15 4.17 4.11|   0     0 |  23M   24M|  57k 3649
 mysql:
----total-cpu-usage---- ---load-avg--- -dsk/total- -net/total- ---system--
usr sys idl wai hiq siq| 1m   5m  15m | read  writ| recv  send| int   csw
 52   9  34   0   0   5|5.03 4.12 4.19|   0    16k|1310k 3693k|  15k   16k
 61  10  24   0   0   5|4.95 4.12 4.19|   0    20k|1504k 4308k|  16k   19k
 63  11  21   0   0   6|4.95 4.12 4.19|   0    44k|1539k 4410k|  17k   20k
 61  12  21   0   0   6|4.95 4.12 4.19|   0    16k|1546k 4430k|  17k   20k
 61  11  22   0   0   7|4.95 4.12 4.19|   0    16k|1584k 4539k|  17k   21k
 54  11  29   0   0   6|4.95 4.12 4.19|   0    16k|1543k 4423k|  17k   21k
 60  10  25   0   0   6|4.95 4.13 4.19|   0    16k|1559k 4466k|  17k   20k
 58  10  26   0   0   6|4.95 4.13 4.19|   0    48k|1526k 4374k|  17k   20k
 48   9  38   0   0   5|4.95 4.13 4.19|   0    20k|1357k 3848k|  15k   18k
 52   7  37   0   0   5|4.95 4.13 4.19|   0    16k|1286k 3598k|  14k   16k
 57  11  26   0   0   6|4.95 4.13 4.19|   0    16k|1546k 4435k|  17k   20k
 58  10  27   0   0   6|4.56 4.06 4.17|   0    16k|1534k 4400k|  17k   20k
 56  10  28   0   0   6|4.56 4.06 4.17|   0    20k|1476k 4187k|  16k   19k
 57   9  29   0   0   5|4.56 4.06 4.17|   0    16k|1531k 4388k|  17k   21k
 59  11  23   1   0   6|4.56 4.06 4.17|   0    16k|1523k 4366k|  17k   20k
 59  10  26   0   0   5|4.56 4.06 4.17|   0    16k|1539k 4413k|  17k   20k
 60  10  25   0   0   6|4.19 4.00 4.15|   0    16k|1555k 4456k|  17k   20k
 60  11  23   0   0   6|4.19 4.00 4.15|   0    24k|1581k 4527k|  17k   21k
 57  10  28   0   0   5|4.19 4.00 4.15|   0    16k|1529k 4383k|  17k   21k
 56  10  27   0   0   6|4.19 4.00 4.15|   0    16k|1535k 4400k|  17k   20k
 57  11  26   0   0   6|4.19 4.00 4.15|   0    16k|1579k 4527k|  17k   21k
 62   9  22   0   0   6|4.02 3.96 4.14|   0    16k|1587k 4548k|  17k   20k
 61  11  21   0   0   6|4.02 3.96 4.14|   0    20k|1584k 4534k|  17k   21k
 58  11  26   0   0   5|4.02 3.96 4.14|   0    16k|1575k 4522k|  17k   21k
 60  11  22   0   0   7|4.02 3.96 4.14|   0    20k|1587k 4543k|  17k   21k
 60  11  24   0   0   5|4.02 3.96 4.14|   0    16k|1571k 4506k|  17k   20k
 59  11  25   0   0   6|4.09 3.98 4.14|   0    16k|1567k 4492k|  17k   21k
 60  10  24   0   0   6|4.09 3.98 4.14|   0    60k|1565k 4498k|  17k   21k
 time            thds              tps     wtps    rtps    rio    rio/op   wio    wio/op    rMB     rKB/op    wMB     wKB/op   cpu/op  freeMB  cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
 [12s]        [r:80,w:0,u:0,d:0]  53284    0       53284   0      0.00     0      0.00      0.00    0.00      0.00    0.00     0.00    0       0         NaN        1.50         633678

 time            thds              tps     wtps    rtps    rio    rio/op   wio    wio/op    rMB     rKB/op    wMB     wKB/op   cpu/op  freeMB  cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
 [13s]        [r:80,w:0,u:0,d:0]  52677    0       52677   0      0.00     0      0.00      0.00    0.00      0.00    0.00     0.00    0       0         NaN        1.52         686355

 time            thds              tps     wtps    rtps    rio    rio/op   wio    wio/op    rMB     rKB/op    wMB     wKB/op   cpu/op  freeMB  cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
 [14s]        [r:80,w:0,u:0,d:0]  53341    0       53341   0      0.00     0      0.00      0.00    0.00      0.00    0.00     0.00    0       0         NaN        1.50         739696

 time            thds              tps     wtps    rtps    rio    rio/op   wio    wio/op    rMB     rKB/op    wMB     wKB/op   cpu/op  freeMB  cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
 [15s]        [r:80,w:0,u:0,d:0]  52654    0       52654   0      0.00     0      0.00      0.00    0.00      0.00    0.00     0.00    0       0         NaN        1.52         792350

 time            thds              tps     wtps    rtps    rio    rio/op   wio    wio/op    rMB     rKB/op    wMB     wKB/op   cpu/op  freeMB  cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
 [16s]        [r:80,w:0,u:0,d:0]  53514    0       53514   0      0.00     0      0.00      0.00    0.00      0.00    0.00     0.00    0       0         NaN        1.49         845864

 time            thds              tps     wtps    rtps    rio    rio/op   wio    wio/op    rMB     rKB/op    wMB     wKB/op   cpu/op  freeMB  cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
 [17s]        [r:80,w:0,u:0,d:0]  53353    0       53353   0      0.00     0      0.00      0.00    0.00      0.00    0.00     0.00    0       0         NaN        1.50         899217

 time            thds              tps     wtps    rtps    rio    rio/op   wio    wio/op    rMB     rKB/op    wMB     wKB/op   cpu/op  freeMB  cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
 [18s]        [r:80,w:0,u:0,d:0]  53095    0       53095   0      0.00     0      0.00      0.00    0.00      0.00    0.00     0.00    0       0         NaN        1.51         952312

 time            thds              tps     wtps    rtps    rio    rio/op   wio    wio/op    rMB     rKB/op    wMB     wKB/op   cpu/op  freeMB  cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
 [19s]        [r:80,w:0,u:0,d:0]  53220    0       53220   0      0.00     0      0.00      0.00    0.00      0.00    0.00     0.00    0       0         NaN        1.50         1005532

 time            thds              tps     wtps    rtps    rio    rio/op   wio    wio/op    rMB     rKB/op    wMB     wKB/op   cpu/op  freeMB  cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
 [20s]        [r:80,w:0,u:0,d:0]  53088    0       53088   0      0.00     0      0.00      0.00    0.00      0.00    0.00     0.00    0       0         NaN        1.51         1058620

 time            thds              tps     wtps    rtps    rio    rio/op   wio    wio/op    rMB     rKB/op    wMB     wKB/op   cpu/op  freeMB  cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
 [21s]        [r:80,w:0,u:0,d:0]  52611    0       52611   0      0.00     0      0.00      0.00    0.00      0.00    0.00     0.00    0       0         NaN        1.52         1111231

 time            thds              tps     wtps    rtps    rio    rio/op   wio    wio/op    rMB     rKB/op    wMB     wKB/op   cpu/op  freeMB  cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
 [22s]        [r:80,w:0,u:0,d:0]  53520    0       53520   0      0.00     0      0.00      0.00    0.00      0.00    0.00     0.00    0       0         NaN        1.49         1164751

 time            thds              tps     wtps    rtps    rio    rio/op   wio    wio/op    rMB     rKB/op    wMB     wKB/op   cpu/op  freeMB  cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
 [23s]        [r:80,w:0,u:0,d:0]  52593    0       52593   0      0.00     0      0.00      0.00    0.00      0.00    0.00     0.00    0       0         NaN        1.52         1217344

 time            thds              tps     wtps    rtps    rio    rio/op   wio    wio/op    rMB     rKB/op    wMB     wKB/op   cpu/op  freeMB  cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
 [24s]        [r:80,w:0,u:0,d:0]  52996    0       52996   0      0.00     0      0.00      0.00    0.00      0.00    0.00     0.00    0       0         NaN        1.51         1270340

## Cobar2.0单机性能测试
####connections: 10
cobar:
----total-cpu-usage---- ---load-avg--- -dsk/total- -net/total- ---system--
usr sys idl wai hiq siq| 1m   5m  15m | read  writ| recv  send| int   csw
 16  26  54   0   0   5|1.68 1.91 2.09|   0     0 |  11M   11M| 122k  235k
 16  25  54   0   0   5|1.68 1.91 2.09|   0     0 |  11M   11M| 122k  235k
 16  25  55   0   0   5|1.68 1.91 2.09|   0     0 |  11M   11M| 122k  235k
 15  24  55   0   0   5|1.68 1.91 2.09|   0     0 |  11M   11M| 120k  230k
 15  25  54   0   0   5|1.71 1.91 2.09|   0    12k|  11M   11M| 122k  235k
 15  25  55   0   0   5|1.71 1.91 2.09|   0     0 |  11M   11M| 121k  234k
 16  24  55   0   0   5|1.71 1.91 2.09|   0    44k|  11M   11M| 117k  226k
 16  25  54   0   0   5|1.71 1.91 2.09|   0     0 |  11M   11M| 123k  235k
 15  25  55   0   0   5|1.71 1.91 2.09|   0     0 |  11M   11M| 122k  235k
 15  26  55   0   0   5|1.81 1.93 2.10|   0    28k|  11M   11M| 122k  236k
 16  24  55   0   0   5|1.81 1.93 2.10|   0     0 |  11M   11M| 121k  234k
 17  24  55   0   0   4|1.81 1.93 2.10|   0    32k|  11M   11M| 120k  231k
 17  24  54   0   0   5|1.81 1.93 2.10|   0     0 |  11M   11M| 121k  232k
 15  26  54   0   0   5|1.81 1.93 2.10|   0     0 |  11M   11M| 123k  236k
 15  25  55   0   0   5|1.83 1.93 2.09|   0    72k|  11M   11M| 122k  235k
 16  25  54   0   0   5|1.83 1.93 2.09|   0     0 |  11M   11M| 123k  236k
 16  26  54   0   0   5|1.83 1.93 2.09|   0     0 |  11M   11M| 122k  234k
 15  25  55   0   0   5|1.83 1.93 2.09|   0     0 |  11M   11M| 120k  231k
 16  25  54   0   0   5|1.83 1.93 2.09|   0     0 |  11M   11M| 123k  237k
 16  25  54   0   0   5|2.00 1.97 2.10|   0    16k|  11M   11M| 122k  235k
 17  24  53   0   0   6|2.00 1.97 2.10|   0     0 |  11M   11M| 119k  230k
 14  25  54   0   0   6|2.00 1.97 2.10|   0   400k|  11M   11M| 122k  235k
 17  25  54   0   0   5|2.00 1.97 2.10|   0     0 |  11M   11M| 121k  232k
 17  25  53   0   0   5|2.00 1.97 2.10|   0     0 |  11M   11M| 121k  233k
mysql:
----total-cpu-usage---- ---load-avg--- -dsk/total- -net/total- ---system--
usr sys idl wai hiq siq| 1m   5m  15m | read  writ| recv  send| int   csw
 26   8  64   1   0   3|0.95 0.90 0.89|   0    20k| 994k 2242k|  16k   21k
 25   7  65   0   0   2|1.03 0.92 0.90|   0    16k| 982k 2217k|  16k   22k
 28   7  63   0   0   2|1.03 0.92 0.90|   0    20k|1018k 2298k|  16k   22k
 27   8  63   0   0   2|1.03 0.92 0.90|   0    16k|1019k 2302k|  16k   22k
 27   9  63   0   0   2|1.03 0.92 0.90|   0    16k|1019k 2295k|  16k   22k
 28   7  62   0   0   3|1.03 0.92 0.90|   0    16k|1014k 2289k|  16k   22k
 29   7  63   0   0   2|0.95 0.90 0.90|   0    16k|1015k 2291k|  16k   22k
 27   7  64   1   0   2|0.95 0.90 0.90|   0    84k| 996k 2249k|  16k   22k
 27   9  62   0   0   3|0.95 0.90 0.90|   0    20k| 982k 2214k|  16k   22k
 27   8  63   0   0   2|0.95 0.90 0.90|   0    60k|1009k 2276k|  16k   22k
 29   7  62   0   0   2|0.95 0.90 0.90|   0    16k| 996k 2245k|  16k   22k
 26   8  63   0   0   3|0.95 0.90 0.90|   0    16k| 993k 2239k|  16k   22k
 25   7  65   0   0   2|0.95 0.90 0.90|   0    24k|1008k 2291k|  16k   22k
 26   8  63   0   0   2|0.95 0.90 0.90|   0    20k|1005k 2269k|  16k   22k
 25   7  64   1   0   3|0.95 0.90 0.90|   0    16k|1005k 2264k|  16k   22k
 28   7  63   0   0   2|0.95 0.90 0.90|   0    16k| 987k 2232k|  16k   21k
 26   8  64   0   0   2|1.04 0.92 0.90|   0    16k|1014k 2288k|  16k   22k
 28   6  63   0   0   3|1.04 0.92 0.90|   0    20k| 987k 2225k|  16k   21k
 26   8  64   0   0   3|1.04 0.92 0.90|   0    16k| 998k 2255k|  16k   22k
 25   9  63   0   0   2|1.04 0.92 0.90|   0    16k|1012k 2283k|  16k   22k
 27   8  64   0   0   2|1.04 0.92 0.90|   0    16k| 990k 2236k|  16k   21k
 25   8  64   1   0   3|0.96 0.91 0.90|   0    16k| 992k 2243k|  16k   22k
 26   9  64   0   0   2|0.96 0.91 0.90|   0    20k| 987k 2225k|  16k   21k
 26   6  65   0   0   3|0.96 0.91 0.90|   0    16k| 963k 2173k|  15k   21k
 27   8  62   0   0   3|0.96 0.91 0.90|   0    16k|1011k 2283k|  16k   22k
 28   8  62   0   0   2|0.96 0.91 0.90|   0    16k|1008k 2275k|  16k   22k

time             thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[83973s]      [r:10,w:0,u:0,d:0]  21934    0       21934    NaN        0.45         1827409711

time             thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[83974s]      [r:10,w:0,u:0,d:0]  21864    0       21864    NaN        0.46         1827431575

time             thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[83975s]      [r:10,w:0,u:0,d:0]  21761    0       21761    NaN        0.46         1827453336

time             thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[83976s]      [r:10,w:0,u:0,d:0]  22019    0       22019    NaN        0.45         1827475355

time             thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[83977s]      [r:10,w:0,u:0,d:0]  21522    0       21522    NaN        0.46         1827496877

time             thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[83978s]      [r:10,w:0,u:0,d:0]  21937    0       21937    NaN        0.45         1827518814

time             thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[83979s]      [r:10,w:0,u:0,d:0]  21966    0       21966    NaN        0.45         1827540780

time             thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[83980s]      [r:10,w:0,u:0,d:0]  21763    0       21763    NaN        0.46         1827562543

time             thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[83981s]      [r:10,w:0,u:0,d:0]  21895    0       21895    NaN        0.46         1827584438

time             thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[83982s]      [r:10,w:0,u:0,d:0]  21804    0       21804    NaN        0.46         1827606242

time             thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[83983s]      [r:10,w:0,u:0,d:0]  21708    0       21708    NaN        0.46         1827627950

time             thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[83984s]      [r:10,w:0,u:0,d:0]  21910    0       21910    NaN        0.46         1827649860

time             thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[83985s]      [r:10,w:0,u:0,d:0]  21789    0       21789    NaN        0.46         1827671649

time             thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[83986s]      [r:10,w:0,u:0,d:0]  21634    0       21634    NaN        0.46         1827693283

####connections: 20
cobar:
----total-cpu-usage---- ---load-avg--- -dsk/total- -net/total- ---system--
usr sys idl wai hiq siq| 1m   5m  15m | read  writ| recv  send| int   csw
 24  37  28   0   0  11|5.17 5.11 5.09|   0     0 |  19M   19M| 135k  326k
 24  38  28   0   0  11|5.17 5.11 5.09|   0     0 |  18M   18M| 132k  321k
 23  38  28   0   0  11|5.17 5.11 5.09|   0     0 |  19M   19M| 135k  328k
 23  40  27   0   0  11|4.99 5.07 5.08|   0    20k|  19M   19M| 134k  326k
 25  38  27   0   0  10|4.99 5.07 5.08|   0     0 |  19M   19M| 135k  328k
 24  38  27   0   0  11|4.99 5.07 5.08|   0     0 |  19M   19M| 133k  324k
 24  39  26   0   0  11|4.99 5.07 5.08|   0     0 |  19M   19M| 133k  322k
 25  36  27   0   0  12|4.99 5.07 5.08|   0  4096B|  18M   19M| 133k  322k
 24  37  28   0   0  12|4.91 5.05 5.08|   0    60k|  18M   19M| 133k  322k
 25  36  28   0   0  12|4.91 5.05 5.08|   0    16k|  19M   19M| 133k  323k
 23  37  29   0   0  11|4.91 5.05 5.08|   0     0 |  18M   18M| 133k  322k
 24  36  29   0   0  11|4.91 5.05 5.08|   0     0 |  18M   18M| 134k  315k
 23  35  32   0   0  11|4.91 5.05 5.08|   0     0 |  17M   18M| 134k  313k
 23  36  30   0   0  11|4.84 5.04 5.07|   0    16k|  18M   18M| 134k  315k
 22  37  29   0   0  11|4.84 5.04 5.07|   0   132k|  18M   18M| 131k  312k
 23  38  28   0   0  11|4.84 5.04 5.07|   0     0 |  18M   18M| 135k  323k
 22  35  33   0   0  10|4.84 5.04 5.07|   0     0 |  17M   17M| 134k  308k
 23  36  31   0   0  11|4.84 5.04 5.07|   0     0 |  18M   18M| 136k  317k
 24  36  28   0   0  12|5.09 5.08 5.08|   0    36k|  19M   19M| 135k  325k
 25  37  28   0   0  11|5.09 5.08 5.08|   0   372k|  19M   19M| 133k  323k
 28  37  24   0   0  11|5.09 5.08 5.08|   0     0 |  18M   18M| 127k  308k
 24  37  28   0   0  11|5.09 5.08 5.08|   0     0 |  18M   19M| 134k  323k
mysql:
----total-cpu-usage---- ---load-avg--- -dsk/total- -net/total- ---system--
usr sys idl wai hiq siq| 1m   5m  15m | read  writ| recv  send| int   csw
 39  11  45   0   0   4|2.16 1.91 1.79|   0    16k|1515k 3415k|  22k   30k
 39  11  46   1   0   4|2.16 1.91 1.79|   0    20k|1472k 3322k|  21k   29k
 38  12  45   1   0   5|2.30 1.95 1.80|   0    16k|1492k 3361k|  22k   29k
 41  11  43   0   0   5|2.30 1.95 1.80|   0    16k|1572k 3549k|  22k   30k
 39  11  45   0   0   5|2.30 1.95 1.80|   0    16k|1511k 3410k|  22k   30k
 42  12  40   0   0   6|2.30 1.95 1.80|   0    20k|1637k 3691k|  23k   31k
 42  12  39   0   0   6|2.30 1.95 1.80|   0    76k|1685k 3800k|  24k   32k
 43  12  40   0   0   6|2.20 1.93 1.80|   0    16k|1692k 3820k|  24k   32k
 44  12  39   0   0   5|2.20 1.93 1.80|   0    20k|1677k 3784k|  24k   32k
 43  11  41   0   0   5|2.20 1.93 1.80|   0    16k|1663k 3749k|  24k   32k
 43  11  41   0   0   6|2.20 1.93 1.80|   0    20k|1656k 3736k|  24k   32k
 43  12  40   0   0   6|2.20 1.93 1.80|   0    20k|1689k 3810k|  24k   32k
 44  12  39   1   0   4|2.10 1.92 1.79|   0    16k|1695k 3822k|  24k   32k
 42  13  40   0   0   5|2.10 1.92 1.79|   0    16k|1661k 3750k|  24k   32k
 42  12  41   0   0   6|2.10 1.92 1.79|   0    16k|1654k 3735k|  23k   31k
 45  11  40   0   0   5|2.10 1.92 1.79|   0    20k|1676k 3782k|  24k   32k
 43  12  39   1   0   5|2.10 1.92 1.79|   0   384k|1669k 3779k|  24k   32k
 43  11  40   0   0   6|2.01 1.90 1.79|   0    16k|1676k 3780k|  24k   32k
 42  12  40   0   0   6|2.01 1.90 1.79|   0    16k|1682k 3795k|  24k   32k
 42  13  40   0   0   5|2.01 1.90 1.79|   0    20k|1657k 3739k|  23k   31k
 45  12  39   0   0   4|2.01 1.90 1.79|   0    20k|1682k 3795k|  24k   32k
 44  13  39   0   0   4|2.01 1.90 1.79|   0    44k|1685k 3803k|  24k   32k
time             thds              tps     wtps    rtps    rio    rio/op   wio    wio/op    rMB     rKB/op    wMB     wKB/op   cpu/op  freeMB  cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
[11081s]      [r:20,w:0,u:0,d:0]  35846    0       35846   0      0.00     0      0.00      0.00    0.00      0.00    0.00     0.00    0       0         NaN        0.56         365022232

time             thds              tps     wtps    rtps    rio    rio/op   wio    wio/op    rMB     rKB/op    wMB     wKB/op   cpu/op  freeMB  cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
[11082s]      [r:20,w:0,u:0,d:0]  36487    0       36487   0      0.00     0      0.00      0.00    0.00      0.00    0.00     0.00    0       0         NaN        0.55         365058719

time             thds              tps     wtps    rtps    rio    rio/op   wio    wio/op    rMB     rKB/op    wMB     wKB/op   cpu/op  freeMB  cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
[11083s]      [r:20,w:0,u:0,d:0]  36593    0       36593   0      0.00     0      0.00      0.00    0.00      0.00    0.00     0.00    0       0         NaN        0.55         365095312

time             thds              tps     wtps    rtps    rio    rio/op   wio    wio/op    rMB     rKB/op    wMB     wKB/op   cpu/op  freeMB  cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
[11084s]      [r:20,w:0,u:0,d:0]  36794    0       36794   0      0.00     0      0.00      0.00    0.00      0.00    0.00     0.00    0       0         NaN        0.54         365132106

time             thds              tps     wtps    rtps    rio    rio/op   wio    wio/op    rMB     rKB/op    wMB     wKB/op   cpu/op  freeMB  cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
[11085s]      [r:20,w:0,u:0,d:0]  36864    0       36864   0      0.00     0      0.00      0.00    0.00      0.00    0.00     0.00    0       0         NaN        0.54         365168970

time             thds              tps     wtps    rtps    rio    rio/op   wio    wio/op    rMB     rKB/op    wMB     wKB/op   cpu/op  freeMB  cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
[11086s]      [r:20,w:0,u:0,d:0]  34052    0       34052   0      0.00     0      0.00      0.00    0.00      0.00    0.00     0.00    0       0         NaN        0.59         365203022

time             thds              tps     wtps    rtps    rio    rio/op   wio    wio/op    rMB     rKB/op    wMB     wKB/op   cpu/op  freeMB  cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
[11087s]      [r:20,w:0,u:0,d:0]  33069    0       33069   0      0.00     0      0.00      0.00    0.00      0.00    0.00     0.00    0       0         NaN        0.60         365236091

time             thds              tps     wtps    rtps    rio    rio/op   wio    wio/op    rMB     rKB/op    wMB     wKB/op   cpu/op  freeMB  cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
[11088s]      [r:20,w:0,u:0,d:0]  32394    0       32394   0      0.00     0      0.00      0.00    0.00      0.00    0.00     0.00    0       0         NaN        0.62         365268485

time             thds              tps     wtps    rtps    rio    rio/op   wio    wio/op    rMB     rKB/op    wMB     wKB/op   cpu/op  freeMB  cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
[11089s]      [r:20,w:0,u:0,d:0]  32882    0       32882   0      0.00     0      0.00      0.00    0.00      0.00    0.00     0.00    0       0         NaN        0.61         365301367

time             thds              tps     wtps    rtps    rio    rio/op   wio    wio/op    rMB     rKB/op    wMB     wKB/op   cpu/op  freeMB  cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
[11090s]      [r:20,w:0,u:0,d:0]  34335    0       34335   0      0.00     0      0.00      0.00    0.00      0.00    0.00     0.00    0       0         NaN        0.58         365335702

time             thds              tps     wtps    rtps    rio    rio/op   wio    wio/op    rMB     rKB/op    wMB     wKB/op   cpu/op  freeMB  cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
[11091s]      [r:20,w:0,u:0,d:0]  33605    0       33605   0      0.00     0      0.00      0.00    0.00      0.00    0.00     0.00    0       0         NaN        0.59         365369307

time             thds              tps     wtps    rtps    rio    rio/op   wio    wio/op    rMB     rKB/op    wMB     wKB/op   cpu/op  freeMB  cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
[11092s]      [r:20,w:0,u:0,d:0]  35754    0       35754   0      0.00     0      0.00      0.00    0.00      0.00    0.00     0.00    0       0         NaN        0.56         365405061

time             thds              tps     wtps    rtps    rio    rio/op   wio    wio/op    rMB     rKB/op    wMB     wKB/op   cpu/op  freeMB  cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
[11093s]      [r:20,w:0,u:0,d:0]  36744    0       36744   0      0.00     0      0.00      0.00    0.00      0.00    0.00     0.00    0       0         NaN        0.54         365441805

time             thds              tps     wtps    rtps    rio    rio/op   wio    wio/op    rMB     rKB/op    wMB     wKB/op   cpu/op  freeMB  cacheMB   w-rsp(ms)  r-rsp(ms)    total-number
[11094s]      [r:20,w:0,u:0,d:0]  36799    0       36799   0      0.00     0      0.00      0.00    0.00      0.00    0.00     0.00    0       0         NaN        0.54         365478604
####connections: 40
cobar:
----total-cpu-usage---- ---load-avg--- -dsk/total- -net/total- ---system--
usr sys idl wai hiq siq| 1m   5m  15m | read  writ| recv  send| int   csw
 28  45  13   0   0  14|8.33 8.33 8.46|   0     0 |  25M   25M| 115k  338k
 28  44  13   0   0  15|8.22 8.31 8.45|   0    16k|  25M   25M| 113k  334k
 28  44  13   0   0  14|8.22 8.31 8.45|   0     0 |  25M   25M| 115k  341k
 29  44  13   0   0  14|8.22 8.31 8.45|   0    24k|  25M   25M| 113k  335k
 29  44  13   0   0  15|8.22 8.31 8.45|   0     0 |  25M   25M| 114k  339k
 28  45  13   0   0  14|8.22 8.31 8.45|   0     0 |  25M   25M| 115k  340k
 29  44  13   0   0  15|8.36 8.34 8.46|   0     0 |  25M   25M| 114k  336k
 28  44  13   0   0  15|8.36 8.34 8.46|   0    36k|  25M   25M| 115k  341k
 27  46  13   0   0  14|8.36 8.34 8.46|   0  9044k|  25M   25M| 114k  337k
 28  44  13   0   0  15|8.36 8.34 8.46|   0     0 |  25M   25M| 114k  339k
 28  44  13   0   0  15|8.36 8.34 8.46|   0     0 |  25M   25M| 115k  339k
 29  42  13   0   0  15|8.42 8.35 8.46|   0     0 |  25M   25M| 113k  330k
 27  45  13   0   0  15|8.42 8.35 8.46|   0   100k|  25M   25M| 114k  340k
 29  43  13   0   0  14|8.42 8.35 8.46|   0   212k|  25M   25M| 115k  340k
 28  45  13   0   0  14|8.42 8.35 8.46|   0     0 |  25M   25M| 114k  336k
 29  44  13   0   0  15|8.42 8.35 8.46|   0     0 |  25M   25M| 114k  339k
 29  44  13   0   0  14|8.70 8.41 8.48|   0     0 |  25M   25M| 114k  336k
 28  44  13   0   0  14|8.70 8.41 8.48|   0    20k|  25M   25M| 115k  338k
 30  42  13   0   0  15|8.70 8.41 8.48|   0  8192B|  25M   25M| 112k  328k
 30  42  13   0   0  15|8.70 8.41 8.48|   0     0 |  25M   25M| 115k  335k
 30  43  13   0   0  14|8.70 8.41 8.48|   0     0 |  25M   25M| 113k  330k
mysql:
----total-cpu-usage---- ---load-avg--- -dsk/total- -net/total- ---system--
usr sys idl wai hiq siq| 1m   5m  15m | read  writ| recv  send| int   csw
 57  15  21   0   0   7|3.71 3.93 4.09|   0    20k|2262k 5088k|  28k   37k
 55  16  22   0   0   7|3.71 3.93 4.09|   0    16k|2241k 5056k|  28k   37k
 56  16  20   0   0   8|3.71 3.93 4.09|   0    16k|2280k 5142k|  29k   37k
 57  15  21   0   0   7|3.71 3.93 4.09|   0    20k|2295k 5169k|  29k   37k
 55  16  20   0   0   9|3.71 3.93 4.09|   0    16k|2286k 5149k|  28k   37k
 58  15  20   0   0   7|3.57 3.90 4.08|   0    20k|2289k 5160k|  29k   38k
 55  15  22   0   0   8|3.57 3.90 4.08|   0    16k|2254k 5081k|  28k   37k
 56  15  21   0   0   8|3.57 3.90 4.08|   0    16k|2261k 5102k|  28k   37k
 55  16  22   0   0   7|3.57 3.90 4.08|   0    16k|2214k 4988k|  28k   37k
 56  16  22   0   0   7|3.57 3.90 4.08|   0    16k|2259k 5092k|  29k   37k
 56  16  20   0   0   8|3.68 3.92 4.08|   0    20k|2272k 5120k|  29k   37k
 56  16  21   0   0   7|3.68 3.92 4.08|   0    20k|2268k 5108k|  28k   38k
 56  15  21   0   0   8|3.68 3.92 4.08|   0    16k|2263k 5098k|  29k   37k
 53  16  21   1   0   9|3.68 3.92 4.08|   0    20k|2246k 5064k|  28k   37k
 58  14  20   0   0   8|3.68 3.92 4.08|   0    16k|2262k 5094k|  28k   37k
 56  15  22   0   0   7|3.95 3.97 4.10|   0   416k|2237k 5042k|  28k   37k
 55  16  20   0   0   9|3.95 3.97 4.10|   0    20k|2279k 5150k|  28k   37k
 57  15  21   0   0   7|3.95 3.97 4.10|   0    16k|2268k 5113k|  28k   37k
 56  15  21   0   0   8|3.95 3.97 4.10|   0    16k|2267k 5111k|  29k   37k
 55  17  21   0   0   7|3.95 3.97 4.10|   0    20k|2259k 5090k|  28k   37k
 56  16  21   0   0   7|3.87 3.95 4.09|   0    24k|2263k 5100k|  28k   37k
 53  18  21   0   0   9|3.87 3.95 4.09|   0    16k|2238k 5046k|  28k   37k
time             thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[60024s]      [r:40,w:0,u:0,d:0]  47809    0       47809    NaN        0.84         2952567816

time             thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[60025s]      [r:40,w:0,u:0,d:0]  48415    0       48415    NaN        0.82         2952616231

time             thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[60026s]      [r:40,w:0,u:0,d:0]  48579    0       48579    NaN        0.82         2952664810

time             thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[60027s]      [r:40,w:0,u:0,d:0]  48803    0       48803    NaN        0.82         2952713613

time             thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[60028s]      [r:40,w:0,u:0,d:0]  48651    0       48651    NaN        0.82         2952762264

time             thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[60029s]      [r:40,w:0,u:0,d:0]  48983    0       48983    NaN        0.82         2952811247

time             thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[60030s]      [r:40,w:0,u:0,d:0]  48938    0       48938    NaN        0.82         2952860185

time             thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[60031s]      [r:40,w:0,u:0,d:0]  48855    0       48855    NaN        0.82         2952909040

time             thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[60032s]      [r:40,w:0,u:0,d:0]  48997    0       48997    NaN        0.81         2952958037

time             thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[60033s]      [r:40,w:0,u:0,d:0]  49143    0       49143    NaN        0.81         2953007180

time             thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[60034s]      [r:40,w:0,u:0,d:0]  48101    0       48101    NaN        0.83         2953055281

time             thds              tps     wtps    rtps     w-rsp(ms)  r-rsp(ms)    total-number
[60035s]      [r:40,w:0,u:0,d:0]  48998    0       48998    NaN        0.82         2953104279
####connections: 80