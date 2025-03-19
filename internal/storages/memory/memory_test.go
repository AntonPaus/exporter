package memory

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMemoryStorage_Get(t *testing.T) {
	// type fields struct {
	// 	mu sync.Mutex
	// 	g  map[string]gauge
	// 	c  map[string]counter
	// }
	// type args struct {
	// 	mName string
	// 	mType string
	// }
	// tests := []struct {
	// 	name    string
	// 	fields  fields
	// 	args    args
	// 	want    any
	// 	wantErr bool
	// }{
	// 	// TODO: Add test cases.
	// }
	// for _, tt := range tests {
	// }
	storage := NewMemoryStorage()
	t.Run("Gauge", func(t *testing.T) {
		got1, err := storage.Update("g1", "gauge", 3.1)
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
	// t.Run("Counter", func(t *testing.T) {
	// 	got1, err := storage.Update("c1", "counter", 3)
	// 	require.NoError(t, err)
	// 	require.Equal(t, 3, got1)
	// 	got2, err := storage.Update("c1", "counter", 3)
	// 	require.NoError(t, err)
	// 	require.Equal(t, 6, got2)
	// 	got3, err := storage.Update("c2", "counter", 2)
	// 	require.NoError(t, err)
	// 	require.Equal(t, 2, got3)
	// })
	t.Run("Wrong type", func(t *testing.T) {
		_, err := storage.Update("g1", "wrong", 3.2)
		require.Error(t, err)
	})
}
