package models

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTrie(t *testing.T) {
	root := NewTrie("")
	root.Add("1.2.3")
	root.Add("1.2.4")
	root.Add("1.4")

	require.Equal(t, false, root.Exist("1.3.3"))
	require.Equal(t, false, root.Exist("1.5"))
	require.Equal(t, false, root.Exist("1.2"))

	require.Equal(t, true, root.Exist("1.2.3"))
	require.Equal(t, true, root.Exist("1.2.4"))
	require.Equal(t, true, root.Exist("1.4"))

}
