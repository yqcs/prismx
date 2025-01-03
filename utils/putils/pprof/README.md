## PProfiling Usage Guide

### Environment Variables

- `PPROF`: Enable or disable profiling. Set to 1 to enable.
- `MEM_PROFILE_DIR`: Directory to write memory profiles to.
- `CPU_PROFILE_DIR`: Directory to write CPU profiles to.
- `PPROF_TIME`: Polling time for CPU and memory profiles (with unit ex: 10s).
- `MEM_PROFILE_RATE`: Memory profiling rate (default 4096).


## How to Use

1. Set the environment variables as per your requirements.

```bash
export PPROF=1
export MEM_PROFILE_DIR=/path/to/memprofile
export CPU_PROFILE_DIR=/path/to/cpuprofile
export PPROF_TIME=10s
export MEM_PROFILE_RATE=4096
```

2. Run your Go application. The profiler will start automatically if PPROF is set to 1.

**Output**

- Memory profiles will be written to the directory specified by MEM_PROFILE_DIR.
- CPU profiles will be written to the directory specified by CPU_PROFILE_DIR.
- Profiles will be written at intervals specified by PPROF_TIME.
- Memory profiling rate is controlled by MEM_PROFILE_RATE.

### Example

```bash
[+] GOOS: linux
[+] GOARCH: amd64
[+] Command: /path/to/your/app
Available PPROF Config Options:
MEM_PROFILE_DIR   - directory to write memory profiles to
CPU_PROFILE_DIR   - directory to write cpu profiles to
PPROF_TIME        - polling time for cpu and memory profiles (with unit ex: 10s)
MEM_PROFILE_RATE  - memory profiling rate (default 4096)
profile: memory profiling enabled (rate 4096), /path/to/memprofile
profile: ticker enabled (rate 10s)
profile: cpu profiling enabled (ticker 10s)
```

### Note

- The polling time (PPROF_TIME) should be set according to your application's performance and profiling needs.
- The memory profiling rate (MEM_PROFILE_RATE) controls the granularity of the memory profiling. Higher values provide more detail but consume more resources.