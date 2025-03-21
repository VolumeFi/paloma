package metrix_test

import (
	"testing"

	keepertest "github.com/palomachain/paloma/v2/testutil/keeper"
	"github.com/palomachain/paloma/v2/testutil/nullify"
	"github.com/palomachain/paloma/v2/x/metrix"
	"github.com/palomachain/paloma/v2/x/metrix/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
	}

	k, ctx := keepertest.MetrixKeeper(t)
	metrix.InitGenesis(ctx, *k, genesisState)
	got := metrix.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)
}
