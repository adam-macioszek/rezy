// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package db

import (
	"context"
)

type Querier interface {
	CreateReservation(ctx context.Context, arg CreateReservationParams) (Reservation, error)
	DeleteReservation(ctx context.Context, id int64) error
	GetOptimizedReservation(ctx context.Context, arg GetOptimizedReservationParams) (Reservation, error)
	GetReservation(ctx context.Context, id int64) (Reservation, error)
	ListAvailableReservations(ctx context.Context, arg ListAvailableReservationsParams) ([]Reservation, error)
	ListReservations(ctx context.Context, arg ListReservationsParams) ([]Reservation, error)
	UpdateReservation(ctx context.Context, arg UpdateReservationParams) (Reservation, error)
}

var _ Querier = (*Queries)(nil)