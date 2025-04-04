package memory

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStorage_Get(t *testing.T) {
	storage, _ := NewStorage(300, "./storage", false)
	defer storage.Terminate()
	t.Run("Gauge", func(t *testing.T) {
		got1, err := storage.Update("g1", "gauge", 3.1)
		fmt.Println("got1", got1)
		require.NoError(t, err)
		require.Equal(t, 3.1, got1)
		got2, err := storage.Update("g1", "gauge", 3.2)
		require.NoError(t, err)
		require.Equal(t, 3.2, got2)
		got3, err := storage.Update("g2", "gauge", 3.3)
		require.NoError(t, err)
		require.Equal(t, 3.3, got3)
		got4, err := storage.Get("g2", "gauge")
		require.NoError(t, err)
		require.Equal(t, 3.3, got4)
	})
	t.Run("Wrong type", func(t *testing.T) {
		_, err := storage.Update("g1", "wrong", 3.2)
		require.Error(t, err)
	})
}
