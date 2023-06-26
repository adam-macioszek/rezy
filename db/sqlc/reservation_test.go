package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/adam-macioszek/rezy/random"
	"github.com/stretchr/testify/require"
)

func createReservationList(t *testing.T, n int, h int) []Reservation {
	var result []Reservation
	for i := 1; i <= n; i++ {
		arg := CreateReservationParams{
			TableSize: int32(i),
			StartTime: time.Date(2023, 6, 23, h, 30, 0, 0, time.UTC),
			Booked:    false,
			Duration:  30,
		}
		reservation, err := testQueries.CreateReservation(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, reservation)
		require.Equal(t, arg.TableSize, reservation.TableSize)
		require.Equal(t, arg.Booked, reservation.Booked)
		require.Equal(t, arg.Duration, reservation.Duration)
		require.Equal(t, arg.StartTime.Year(), reservation.StartTime.Year())
		require.Equal(t, arg.StartTime.Month(), reservation.StartTime.Month())
		require.Equal(t, arg.StartTime.Day(), reservation.StartTime.Day())
		require.Equal(t, arg.StartTime.Hour(), reservation.StartTime.Hour())
		require.Equal(t, arg.StartTime.Minute(), reservation.StartTime.Minute())
		require.Equal(t, arg.StartTime.Second(), reservation.StartTime.Second())
		result = append(result, reservation)
	}
	return result
}

func createCustomReservation(t *testing.T, tableSize int32, start time.Time, booked bool, duration int) Reservation {
	arg := CreateReservationParams{
		TableSize: tableSize,
		StartTime: start,
		Booked:    booked,
		Duration:  int32(duration),
	}
	reservation, err := testQueries.CreateReservation(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, reservation)
	require.Equal(t, arg.TableSize, reservation.TableSize)
	require.Equal(t, arg.Booked, reservation.Booked)
	require.Equal(t, arg.Duration, reservation.Duration)
	require.Equal(t, arg.StartTime.Year(), reservation.StartTime.Year())
	require.Equal(t, arg.StartTime.Month(), reservation.StartTime.Month())
	require.Equal(t, arg.StartTime.Day(), reservation.StartTime.Day())
	require.Equal(t, arg.StartTime.Hour(), reservation.StartTime.Hour())
	require.Equal(t, arg.StartTime.Minute(), reservation.StartTime.Minute())
	require.Equal(t, arg.StartTime.Second(), reservation.StartTime.Second())
	return reservation
}

// TODO: generate random values to test
func createRandomReservation(t *testing.T) Reservation {
	arg := CreateReservationParams{
		TableSize: random.RandomInt(1, 10),
		StartTime: time.Now(),
		Booked:    false,
		Duration:  30,
	}
	reservation, err := testQueries.CreateReservation(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, reservation)
	require.Equal(t, arg.TableSize, reservation.TableSize)
	require.Equal(t, arg.Booked, reservation.Booked)
	require.Equal(t, arg.Duration, reservation.Duration)
	require.Equal(t, arg.StartTime.Year(), reservation.StartTime.Year())
	require.Equal(t, arg.StartTime.Month(), reservation.StartTime.Month())
	require.Equal(t, arg.StartTime.Day(), reservation.StartTime.Day())
	require.Equal(t, arg.StartTime.Hour(), reservation.StartTime.Hour())
	require.Equal(t, arg.StartTime.Minute(), reservation.StartTime.Minute())
	require.Equal(t, arg.StartTime.Second(), reservation.StartTime.Second())
	return reservation
}

func TestCreateReservation(t *testing.T) {
	reservation := createCustomReservation(t, 5, time.Now(), false, 30)
	err := testQueries.DeleteReservation(context.Background(), reservation.ID)
	require.NoError(t, err)

}
func TestGetReservation(t *testing.T) {
	reservation1 := createCustomReservation(t, 5, time.Now(), false, 30)
	response, err := testQueries.GetReservation(context.Background(), reservation1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, response)
	require.Equal(t, reservation1.StartTime.Year(), response.StartTime.Year())
	require.Equal(t, reservation1.StartTime.Month(), response.StartTime.Month())
	require.Equal(t, reservation1.StartTime.Day(), response.StartTime.Day())
	require.Equal(t, reservation1.StartTime.Hour(), response.StartTime.Hour())
	require.Equal(t, reservation1.StartTime.Minute(), response.StartTime.Minute())
	require.Equal(t, reservation1.StartTime.Second(), response.StartTime.Second())
	require.Equal(t, reservation1.Booked, response.Booked)
	require.Equal(t, reservation1.Duration, response.Duration)
	require.Equal(t, reservation1.TableSize, response.TableSize)
	require.Equal(t, reservation1.Booked, response.Booked)
	require.Equal(t, reservation1.Duration, response.Duration)

	err = testQueries.DeleteReservation(context.Background(), reservation1.ID)
	require.NoError(t, err)
}
func TestGetReservationAtTime(t *testing.T) {
	reservations := createReservationList(t, 5, 5)
	arg := GetOptimizedReservationParams{
		StartTime: reservations[2].StartTime,
		Booked:    reservations[2].Booked,
		TableSize: 3,
	}
	response, err := testQueries.GetOptimizedReservation(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, response)
	require.Equal(t, reservations[2].TableSize, response.TableSize)
	require.Equal(t, reservations[2].Booked, response.Booked)
	require.Equal(t, reservations[2].StartTime, response.StartTime)

	for _, res := range reservations {
		err := testQueries.DeleteReservation(context.Background(), res.ID)
		require.NoError(t, err)
	}
}

func TestUpdateReservation(t *testing.T) {
	reservation := createCustomReservation(t, 5, time.Now(), false, 30)
	arg := UpdateReservationParams{
		ID:        reservation.ID,
		TableSize: 10,
		StartTime: time.Date(2024, 6, 24, 21, 0, 0, 651387237, time.UTC),
		Booked:    true,
		Duration:  30,
	}
	response, err := testQueries.UpdateReservation(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, response)
	require.Equal(t, arg.StartTime.Year(), response.StartTime.Year())
	require.Equal(t, arg.StartTime.Month(), response.StartTime.Month())
	require.Equal(t, arg.StartTime.Day(), response.StartTime.Day())
	require.Equal(t, arg.StartTime.Hour(), response.StartTime.Hour())
	require.Equal(t, arg.StartTime.Minute(), response.StartTime.Minute())
	require.Equal(t, arg.StartTime.Second(), response.StartTime.Second())
	require.Equal(t, arg.Booked, response.Booked)
	require.Equal(t, arg.Duration, response.Duration)
	require.Equal(t, arg.TableSize, response.TableSize)

	err = testQueries.DeleteReservation(context.Background(), response.ID)
	require.NoError(t, err)
}

func TestDeleteReservation(t *testing.T) {
	reservation := createCustomReservation(t, 5, time.Now(), false, 30)
	err := testQueries.DeleteReservation(context.Background(), reservation.ID)
	require.NoError(t, err)

	response, err := testQueries.GetReservation(context.Background(), reservation.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, response)

}
func TestListReservation(t *testing.T) {
	for i := 0; i < 10; i++ {
		createCustomReservation(t, 5, time.Now(), false, 30)

	}
	arg := ListReservationsParams{
		Limit:  5,
		Offset: 0,
	}
	response, err := testQueries.ListReservations(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, response, 5)
	for _, reservation := range response {
		require.NotEmpty(t, reservation)
		err = testQueries.DeleteReservation(context.Background(), reservation.ID)
		require.NoError(t, err)
	}
}

func TestListAvailableReservation(t *testing.T) {
	for i := 0; i < 5; i++ {
		createCustomReservation(t, int32(2+i), time.Now(), false, 30)
	}
	arg := ListAvailableReservationsParams{
		Limit:  5,
		Offset: 0,
	}
	response, err := testQueries.ListAvailableReservations(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, response, 5)
	for _, reservation := range response {
		require.NotEmpty(t, reservation)
		require.Equal(t, false, reservation.Booked)
		err = testQueries.DeleteReservation(context.Background(), reservation.ID)
		require.NoError(t, err)
	}
}
