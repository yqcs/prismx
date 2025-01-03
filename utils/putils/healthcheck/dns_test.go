package healthcheck

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDnsResolve(t *testing.T) {
	t.Run("Successful resolution", func(t *testing.T) {
		info := DnsResolve("scanme.sh", "1.1.1.1")
		assert.NoError(t, info.Error)
		assert.True(t, info.Successful)
		assert.Equal(t, "scanme.sh", info.Host)
		assert.Equal(t, "1.1.1.1", info.Resolver)
		assert.NotEmpty(t, info.IPAddresses)
	})

	t.Run("Unsuccessful resolution due to invalid host", func(t *testing.T) {
		info := DnsResolve("invalid.website", "1.1.1.1")
		assert.Error(t, info.Error)
	})

	t.Run("Unsuccessful resolution due to invalid resolver", func(t *testing.T) {
		info := DnsResolve("google.com", "invalid.resolver")
		assert.Error(t, info.Error)
	})
}
