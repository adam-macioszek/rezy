package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	db "github.com/adam-macioszek/rezy/db/sqlc"
	"github.com/gin-gonic/gin"
)

// add validator for better validations of fields
type MakeBookingParams struct {
	ID        int64 `json:"id" binding:"required,min=1"`
	TableSize int32 `json:"table_size" binding:"required,min=1"`
}

func (server *Server) createBooking(ctx *gin.Context) {
	var bookingRequest MakeBookingParams
	err := ctx.ShouldBindJSON(&bookingRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
	}
	arg := db.MakeReservationParams{
		ReservationID: bookingRequest.ID,
		TableSize:     bookingRequest.TableSize,
	}

	if !server.validReservation(ctx, arg.ReservationID, arg.TableSize) {
		return
	}
	result, err := server.store.MakeReservation(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validReservation(ctx *gin.Context, reservationID int64, partySize int32) bool {
	reservation, err := server.store.GetReservation(ctx, reservationID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}
	if reservation.Booked {
		alreadyBookedError := errors.New("this reservation is already booked, please choose another")
		ctx.JSON(http.StatusBadRequest, errorResponse(alreadyBookedError))
		return false
	}
	if reservation.TableSize < partySize {
		tableTooSmall := errors.New("this table size is too small for your requested party size")
		ctx.JSON(http.StatusBadRequest, errorResponse(tableTooSmall))
		return false
	}
	//Might be worth making this logic a little more user friendly,
	//i.e if within 5/10 min of start time can still make reservation?
	if reservation.StartTime.Before(time.Now()) {
		startTimeElapsed := errors.New("this reservations start time has already begun")
		ctx.JSON(http.StatusBadRequest, errorResponse(startTimeElapsed))
		return false
	}
	return true

}
