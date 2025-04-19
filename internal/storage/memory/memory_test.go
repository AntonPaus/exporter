package memory

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStorage_Get(t *testing.T) {
	storage, _ := NewStorage(300, "./storage", false)
	defer storage.Terminate()
	t.Run("Gauge", func(t *testing.T) {
		got1, err := storage.UpdateGauge("g1", 3.1)
		require.NoError(t, err)
		require.Equal(t, 3.1, float64(got1))
		got2, err := storage.UpdateGauge("g1", 3.2)
		require.NoError(t, err)
		require.Equal(t, 3.2, float64(got2))
		got3, err := storage.UpdateGauge("g2", 3.3)
		require.NoError(t, err)
		require.Equal(t, 3.3, float64(got3))
		got4, err := storage.GetGauge("g2")
		require.NoError(t, err)
		require.Equal(t, 3.3, float64(got4))
	})
}
