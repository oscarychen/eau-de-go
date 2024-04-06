package keys_test

import (
	"testing"

	"eau-de-go/pkg/keys"
	"github.com/stretchr/testify/assert"
)

func TestGetInMemoryAesKey(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		key := keys.GetInMemoryAesKey()
		assert.NotNil(t, key)
	})
}
