package fractal_test

import (
	"testing"

	"github.com/ibllex/go-fractal"
	"github.com/stretchr/testify/assert"
)

func TestParseIncludes(t *testing.T) {
	manager := fractal.NewManager(nil)

	t.Run("default", func(t *testing.T) {
		manager.ParseIncludes([]string{"author"})

		actual := manager.GetRequestedIncludes()
		expected := []string{"author"}

		assert.Equal(t, expected, actual)
	})

	t.Run("sub relations", func(t *testing.T) {
		manager.ParseIncludes([]string{"author.address"})

		actual := manager.GetRequestedIncludes()
		expected := []string{"author", "author.address"}

		assert.Equal(t, expected, actual)
	})

	t.Run("recursion limit", func(t *testing.T) {
		manager.SetRecursionLimit(3)
		manager.ParseIncludes([]string{"one.two.three.four.five"})

		actual := manager.GetRequestedIncludes()
		expected := []string{"one", "one.two", "one.two.three"}

		assert.Equal(t, expected, actual)
	})
}
