package main

import (
	"context"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestApplyForcedConfigOptions(t *testing.T) {
	t.Run("Sets override config we expect", func(t *testing.T) {
		cmd := cobra.Command{}
		serverCtx := server.NewDefaultContext()
		ctx := context.WithValue(context.Background(), server.ServerContextKey, serverCtx)

		cmd.SetContext(ctx)
		err := applyForcedConfigOptions(&cmd)
		require.NoError(t, err)

		result := server.GetServerContextFromCmd(&cmd)

		// Assert each thing we tried to set
		require.Equal(t, 200*time.Millisecond, result.Config.Consensus.TimeoutCommit)
		require.Equal(t, 200*time.Millisecond, result.Config.Consensus.TimeoutPrecommit)
		require.Equal(t, 200*time.Millisecond, result.Config.Consensus.TimeoutPrevote)
		require.Equal(t, 200*time.Millisecond, result.Config.Consensus.TimeoutPropose)
		require.Equal(t, 50*time.Millisecond, result.Config.Consensus.TimeoutPrecommitDelta)
		require.Equal(t, 50*time.Millisecond, result.Config.Consensus.TimeoutPrevoteDelta)
		require.Equal(t, 50*time.Millisecond, result.Config.Consensus.TimeoutProposeDelta)
	})
}
