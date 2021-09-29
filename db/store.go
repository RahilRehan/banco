package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v and rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"fromAccountID"`
	ToAccountID   int64 `json:"toAccountID"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"fromAccount"`
	ToAccount   Account  `json:"toAccount"`
	FromEntry   Entry    `json:"fromEntry"`
	ToEntry     Entry    `json:"toEntry"`
}

func (store *Store) TransferTx(ctx context.Context, args TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(args))
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.FromAccountID,
			Amount:    -args.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.ToAccountID,
			Amount:    args.Amount,
		})
		if err != nil {
			return err
		}

		// not so better way
		// account1, err := q.GetAccountForUpdate(context.Background(), args.FromAccountID)
		// if err != nil {
		// 	return err
		// }

		// result.FromAccount, err = q.UpdateAccount(context.Background(), UpdateAccountParams{args.FromAccountID, account1.Balance - args.Amount})
		// if err != nil {
		// 	return err
		// }

		// account2, err := q.GetAccountForUpdate(context.Background(), args.ToAccountID)
		// if err != nil {
		// 	return err
		// }

		// result.ToAccount, err = q.UpdateAccount(context.Background(), UpdateAccountParams{args.ToAccountID, account2.Balance + args.Amount})
		// if err != nil {
		// 	return err
		// }

		// better way
		// result.FromAccount, err = q.AddAccountBalance(context.Background(), AddAccountBalanceParams{-args.Amount, args.FromAccountID})
		// if err != nil {
		// 	return err
		// }

		// result.ToAccount, err = q.AddAccountBalance(context.Background(), AddAccountBalanceParams{args.Amount, args.ToAccountID})
		// if err != nil {
		// 	return err
		// }

		// avoid deadlock while updating
		if args.FromAccountID < args.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, args.FromAccountID, -args.Amount, args.ToAccountID, args.Amount)
			if err != nil {
				return err
			}
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, args.ToAccountID, args.Amount, args.FromAccountID, -args.Amount)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return TransferTxResult{}, err
	}
	return result, nil
}

func addMoney(ctx context.Context, q *Queries, accountID1, amount1, accountID2, amount2 int64) (account1, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{amount1, accountID1})
	if err != nil {
		return
	}
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{amount2, accountID2})
	if err != nil {
		return
	}
	return
}
