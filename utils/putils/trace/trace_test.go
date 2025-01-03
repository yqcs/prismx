package trace

import (
	"testing"
	"time"
)

func TestFunctionWithBeforeFunction(t *testing.T) {
	var beforeCalled bool
	_, _ = Trace(func() {
		if !beforeCalled {
			t.Errorf("Before function was not called before the main function")
		}
	}, WithBefore(func() {
		beforeCalled = true
	}))

	if !beforeCalled {
		t.Errorf("Before function was not called")
	}
}

func TestFunctionWithAfterFunction(t *testing.T) {
	var afterCalled bool
	_, _ = Trace(func() {
		if afterCalled {
			t.Errorf("After function was called before the main function finished")
		}
	}, WithAfter(func() {
		afterCalled = true
	}))

	if !afterCalled {
		t.Errorf("After function was not called")
	}
}

func TestFunctionTracing(t *testing.T) {
	metrics, _ := Trace(func() {
		time.Sleep(2 * time.Second)
	})

	if metrics.ExecutionDuration.Seconds() < 2 {
		t.Errorf("ExecutionDuration is less than expected: %v", metrics.ExecutionDuration)
	}

	if len(metrics.Snapshots) == 0 {
		t.Errorf("Memory snapshots are not captured")
	}

	if metrics.MinAllocMemory == 0 {
		t.Errorf("MinMemory not computed")
	}

	if metrics.MaxAllocMemory == 0 {
		t.Errorf("MaxMemory not computed")
	}

	if metrics.AvgAllocMemory == 0 {
		t.Errorf("AvgMemory not computed")
	}
}

func TestFunctionWithCustomStrategy(t *testing.T) {
	var customLogs []string
	metrics, _ := Trace(func() {
		time.Sleep(1 * time.Second)
	}, WithStrategy(&CustomStrategy{metrics: &Metrics{}, logs: &customLogs}))

	if len(customLogs) != 2 {
		t.Errorf("Custom logs not captured as expected")
	}

	if customLogs[0] != "Custom Before method started." {
		t.Errorf("Expected custom log for Before method not found")
	}

	if customLogs[1] != "Custom After method executed." {
		t.Errorf("Expected custom log for After method not found")
	}

	if metrics.ExecutionDuration.Seconds() < 1 {
		t.Errorf("ExecutionDuration is less than expected: %v", metrics.ExecutionDuration)
	}

	if len(metrics.Snapshots) != 0 {
		t.Errorf("Custom strategy should not capture snapshots")
	}
}

type CustomStrategy struct {
	metrics *Metrics
	logs    *[]string
}

func (c *CustomStrategy) Before() {
	*c.logs = append(*c.logs, "Custom Before method started.")
	c.metrics.StartTime = time.Now()
}

func (c *CustomStrategy) After() {
	*c.logs = append(*c.logs, "Custom After method executed.")
	c.metrics.FinishTime = time.Now()
	c.metrics.ExecutionDuration = c.metrics.FinishTime.Sub(c.metrics.StartTime)
}

func (c *CustomStrategy) GetMetrics() *Metrics {
	return c.metrics
}
