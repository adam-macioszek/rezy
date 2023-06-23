package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Store interface {
	MakeReservation(ctx context.Context, arg MakeReservationParams) (Reservation, error)
	Querier
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error:%v rollback error: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()

}

type MakeReservationParams struct {
	ReservationID int64 `json:"id"`
	TableSize     int32 `json:"table_size"`
}
type MakeReservationResult struct {
	ReservationID int64     `json:"id"`
	TableSize     int32     `json:"table_size"`
	Booked        bool      `json:"booked"`
	StartTime     time.Time `json:"start_time"`
	Duration      int32     `json:"duration"`
}

func (store *SQLStore) MakeReservation(ctx context.Context, arg MakeReservationParams) (Reservation, error) {
	var result Reservation
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result, err = q.GetReservation(ctx, arg.ReservationID)
		if err != nil {
			return err
		}
		if result.Booked {
			return errors.New("reservation is already booked, sorry")
		}

		optimizedResult, err := q.GetOptimizedReservation(ctx, GetOptimizedReservationParams{
			TableSize: arg.TableSize,
			Booked:    false,
			StartTime: result.StartTime,
		})
		if err != nil {
			return err
		}

		result, err = q.UpdateReservation(ctx, UpdateReservationParams{
			ID:        optimizedResult.ID,
			TableSize: optimizedResult.TableSize,
			Booked:    true,
			StartTime: optimizedResult.StartTime,
			Duration:  optimizedResult.Duration,
		})
		if err != nil {
			return err
		}
		return nil
		//rollback
	})
	return result, err
}
