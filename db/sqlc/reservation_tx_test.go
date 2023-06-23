package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSimpleReservation(t *testing.T) {

	reservation := createRandomReservation(t)
	arg := MakeReservationParams{
		ReservationID: reservation.ID,
		TableSize:     reservation.TableSize,
	}
	reservation, err := store.MakeReservation(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, reservation)
	require.Equal(t, arg.ReservationID, reservation.ID)
	require.LessOrEqual(t, arg.TableSize, reservation.TableSize)
	require.Equal(t, true, reservation.Booked)

}

func TestBookedReservation(t *testing.T) {

	reservation := createCustomReservation(t, 5, time.Now(), time.Now(), true, 30)
	arg := MakeReservationParams{
		ReservationID: reservation.ID,
		TableSize:     reservation.TableSize,
	}
	reservation, err := store.MakeReservation(context.Background(), arg)
	require.Error(t, err)

}

func TestOptimizedReservation(t *testing.T) {
	reservations := createReservationList(t, 5, 8)
	arg := MakeReservationParams{
		ReservationID: reservations[2].ID,
		TableSize:     1,
	}
	for i := 0; i < 3; i++ {
		reservation, err := store.MakeReservation(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, reservation)
		require.Equal(t, reservation.Booked, true)
		require.Equal(t, int32(i+1), reservation.TableSize)
		require.Equal(t, reservations[i].ID, reservation.ID)
	}
	_, err := store.MakeReservation(context.Background(), arg)
	require.Error(t, err)
}
