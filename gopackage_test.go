package alfred_go_pacakge_search

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSearchGoPackages(t *testing.T) {
	results, err := SearchGoPackages("elasticsearch", 10)
	require.NoError(t, err)

	require.Greater(t, len(results), 0)

	for _, v := range results {
		fmt.Printf("package: %+v\n", v)
	}
}
