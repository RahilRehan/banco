package db

import (
	"context"
	"testing"
	"time"

	"github.com/RahilRehan/banco/db/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, acc1 Account, acc2 Account) Transfer {
	params := CreateTransferParams{
		FromAccountID: acc1.ID,
		ToAccountID:   acc2.ID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), params)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, params.Amount, transfer.Amount)
	require.Equal(t, params.ToAccountID, transfer.ToAccountID)
	require.Equal(t, params.FromAccountID, transfer.FromAccountID)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)
	createRandomTransfer(t, acc1, acc2)
}

func TestGetTransfer(t *testing.T) {
	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)
	transfer := createRandomTransfer(t, acc1, acc2)

	trasnfer2, err := testQueries.GetTransfer(context.Background(), transfer.ID)

	require.NoError(t, err)
	require.NotEmpty(t, trasnfer2)

	require.Equal(t, transfer.ID, trasnfer2.ID)
	require.Equal(t, transfer.FromAccountID, trasnfer2.FromAccountID)
	require.Equal(t, transfer.ToAccountID, trasnfer2.ToAccountID)
	require.WithinDuration(t, transfer.CreatedAt, trasnfer2.CreatedAt, time.Second)
}

func TestListTransfers(t *testing.T) {
	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomTransfer(t, acc1, acc2)
	}

	params := ListTransfersParams{
		FromAccountID: acc1.ID,
		ToAccountID:   acc2.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), params)

	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}
