package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)


func TestTransferTx(t *testing.T) {
	// Initialize the store
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before: account1 balance:", account1.Balance, ", account2 balance:", account2.Balance)

	errs := make(chan error)

	n := 10
	amount := int64(10)
	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func() {
			fmt.Println(">> tx:", fromAccountID, "->", toAccountID, "amount:", amount)
			ctx := context.Background()
			_, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <- errs
		require.NoError(t, err)
	}

	updateAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updateAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance, updateAccount1.Balance)
	require.Equal(t, account2.Balance, updateAccount2.Balance)
}