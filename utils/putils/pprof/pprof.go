package pprof

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"
	"time"

	"prismx_cli/utils/putils/env"
)

const (
	PPROFSwitchENV = "PPROF"
	MemProfileENV  = "MEM_PROFILE_DIR"
	CPUProfileENV  = "CPU_PROFILE_DIR"
	PPROFTimeENV   = "PPROF_TIME"
	MemProfileRate = "MEM_PROFILE_RATE"
)

func init() {
	if env.GetEnvOrDefault(PPROFSwitchENV, 0) == 1 {
		log.Printf("[+] GOOS: %v\n", runtime.GOOS)
		log.Printf("[+] GOARCH: %v\n", runtime.GOARCH)
		log.Printf("[+] Command: %v\n", strings.Join(os.Args, " "))
		log.Println("Available PPROF Config Options:")
		log.Printf("%-16v - directory to write memory profiles to\n", MemProfileENV)
		log.Printf("%-16v - directory to write cpu profiles to\n", CPUProfileENV)
		log.Printf("%-16v - polling time for cpu and memory profiles (with unit ex: 10s)\n", PPROFTimeENV)
		log.Printf("%-16v - memory profiling rate (default 4096)\n", MemProfileRate)

		memProfilesDir := env.GetEnvOrDefault(MemProfileENV, "memdump")
		cpuProfilesDir := env.GetEnvOrDefault(CPUProfileENV, "cpuprofile")
		pprofTimeDuration := env.GetEnvOrDefault(PPROFTimeENV, time.Duration(3)*time.Second)
		pprofRate := env.GetEnvOrDefault(MemProfileRate, 4096)

		_ = os.MkdirAll(memProfilesDir, 0755)
		_ = os.MkdirAll(cpuProfilesDir, 0755)

		runtime.MemProfileRate = pprofRate
		log.Printf("profile: memory profiling enabled (rate %d), %s\n", runtime.MemProfileRate, memProfilesDir)
		log.Printf("profile: ticker enabled (rate %s)\n", pprofTimeDuration)

		// cpu ticker and profiler
		go func() {
			ticker := time.NewTicker(pprofTimeDuration)
			count := 0
			buff := bytes.Buffer{}
			log.Printf("profile: cpu profiling enabled (ticker %s)\n", pprofTimeDuration)
			for {
				err := pprof.StartCPUProfile(&buff)
				if err != nil {
					log.Fatalf("profile: could not start cpu profile: %s\n", err)
				}
				<-ticker.C
				pprof.StopCPUProfile()
				if err := os.WriteFile(filepath.Join(cpuProfilesDir, "cpuprofile-t"+strconv.Itoa(count)+".out"), buff.Bytes(), 0755); err != nil {
					log.Fatalf("profile: could not write cpu profile: %s\n", err)
				}
				buff.Reset()
				count++
			}
		}()

		// memory ticker and profiler
		go func() {
			ticker := time.NewTicker(pprofTimeDuration)
			count := 0
			log.Printf("profile: memory profiling enabled (ticker %s)\n", pprofTimeDuration)
			for {
				<-ticker.C
				var buff bytes.Buffer
				if err := pprof.WriteHeapProfile(&buff); err != nil {
					log.Printf("profile: could not write memory profile: %s\n", err)
				}
				err := os.WriteFile(filepath.ToSlash(filepath.Join(memProfilesDir, "memprofile-t"+strconv.Itoa(count)+".out")), buff.Bytes(), 0755)
				if err != nil {
					log.Printf("profile: could not write memory profile: %s\n", err)
				}
				count++
			}
		}()
	}
}
