package examples

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMakeConfig(t *testing.T) {
	err := MakeConfig()
	require.NoError(t, err)
}
