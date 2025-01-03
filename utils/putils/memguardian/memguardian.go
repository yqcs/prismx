package memguardian

import (
	"context"
	"errors"
	"sync/atomic"
	"time"

	units "github.com/docker/go-units"
	"prismx_cli/utils/putils/env"
)

var (
	DefaultInterval        time.Duration
	DefaultMaxUsedRamRatio float64

	DefaultMemGuardian *MemGuardian
)

const (
	MemGuardianEnabled            = "MEMGUARDIAN"
	MemGuardianMaxUsedRamRatioENV = "MEMGUARDIAN_MAX_RAM_RATIO"
	MemGuardianMaxUsedMemoryENV   = "MEMGUARDIAN_MAX_RAM"
	MemGuardianIntervalENV        = "MEMGUARDIAN_INTERVAL"
)

func init() {
	DefaultInterval = env.GetEnvOrDefault(MemGuardianMaxUsedRamRatioENV, time.Duration(time.Second*30))
	DefaultMaxUsedRamRatio = env.GetEnvOrDefault(MemGuardianMaxUsedRamRatioENV, float64(75))
	maxRam := env.GetEnvOrDefault(MemGuardianMaxUsedRamRatioENV, "")

	options := []MemGuardianOption{
		WitInterval(DefaultInterval),
		WithMaxRamRatioWarning(DefaultMaxUsedRamRatio),
	}
	if maxRam != "" {
		options = append(options, WithMaxRamAmountWarning(maxRam))
	}

	var err error
	DefaultMemGuardian, err = New(options...)
	if err != nil {
		panic(err)
	}
}

type MemGuardianOption func(*MemGuardian) error

// WithInterval defines the ticker interval of the memory monitor
func WitInterval(d time.Duration) MemGuardianOption {
	return func(mg *MemGuardian) error {
		mg.t = time.NewTicker(d)
		return nil
	}
}

// WithCallback defines an optional callback if the warning ration is exceeded
func WithCallback(f func()) MemGuardianOption {
	return func(mg *MemGuardian) error {
		mg.f = f
		return nil
	}
}

// WithMaxRamRatioWarning defines the ratio (1-100) threshold of the warning state (and optional callback invocation)
func WithMaxRamRatioWarning(ratio float64) MemGuardianOption {
	return func(mg *MemGuardian) error {
		if ratio == 0 || ratio > 100 {
			return errors.New("ratio must be between 1 and 100")
		}
		mg.ratio = ratio
		return nil
	}
}

// WithMaxRamAmountWarning defines the max amount of used RAM in bytes threshold of the warning state (and optional callback invocation)
func WithMaxRamAmountWarning(maxRam string) MemGuardianOption {
	return func(mg *MemGuardian) error {
		size, err := units.FromHumanSize(maxRam)
		if err != nil {
			return err
		}
		mg.maxMemory = uint64(size)
		return nil
	}
}

type MemGuardian struct {
	t         *time.Ticker
	f         func()
	ctx       context.Context
	cancel    context.CancelFunc
	Warning   atomic.Bool
	ratio     float64
	maxMemory uint64
}

// New mem guadian instance with user defined options
func New(options ...MemGuardianOption) (*MemGuardian, error) {
	mg := &MemGuardian{}
	for _, option := range options {
		if err := option(mg); err != nil {
			return nil, err
		}
	}

	mg.ctx, mg.cancel = context.WithCancel(context.TODO())

	return mg, nil
}

// Run the instance monitor (cancel using the Stop method or context parameter)
func (mg *MemGuardian) Run(ctx context.Context) error {
	for {
		select {
		case <-mg.ctx.Done():
			mg.Close()
			return nil
		case <-ctx.Done():
			mg.Close()
			return nil
		case <-mg.t.C:
			usedRatio, used, err := UsedRam()
			if err != nil {
				return err
			}

			isRatioOverThreshold := mg.ratio > 0 && usedRatio >= mg.ratio
			isAmountOverThreshold := mg.maxMemory > 0 && used >= mg.maxMemory
			if isRatioOverThreshold || isAmountOverThreshold {
				mg.Warning.Store(true)
				if mg.f != nil {
					mg.f()
				}
			} else {
				mg.Warning.Store(false)
			}
		}
	}
}

// Close and stops the instance
func (mg *MemGuardian) Close() {
	mg.cancel()
	mg.t.Stop()
}

// Calculate the system absolute ratio of used RAM
func UsedRam() (ratio float64, used uint64, err error) {
	si, err := GetSysInfo()
	if err != nil {
		return 0, 0, err
	}

	return si.UsedPercent(), si.UsedRam(), nil
}
