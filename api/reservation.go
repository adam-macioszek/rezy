package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/adam-macioszek/rezy/db/sqlc"
	"github.com/gin-gonic/gin"
)

// min for time and start dat, time.Now()
// min time for duration
type createReservationRequest struct {
	TableSize int32     `json:"table_size" binding:"required,min=1"`
	StartTime time.Time `json:"start_time" binding:"required"`
	Booked    bool      `json:"booked"`
	Duration  int32     `json:"Duration" binding:"required,min=1"`
}
type getReservationRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}
type listReservationRequest struct {
	PageID   int32 `form:"page_id,default=1" binding:"min=1"`
	PageSize int32 `form:"page_size,default=5" binding:"min=5,max=15"`
}

func (server *Server) createReservation(ctx *gin.Context) {
	var reservationRequest createReservationRequest
	if err := ctx.ShouldBindJSON(&reservationRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.CreateReservationParams{
		TableSize: reservationRequest.TableSize,
		StartTime: reservationRequest.StartTime,
		Booked:    reservationRequest.Booked,
		Duration:  reservationRequest.Duration,
	}
	reservation, err := server.store.CreateReservation(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, reservation)
}

func (server *Server) getReservation(ctx *gin.Context) {
	var reservationRequest getReservationRequest
	if err := ctx.ShouldBindUri(&reservationRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	reservation, err := server.store.GetReservation(ctx, int64(reservationRequest.ID))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, reservation)
}

func (server *Server) listReservation(ctx *gin.Context) {
	var listRequest listReservationRequest

	if err := ctx.ShouldBindQuery(&listRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListAvailableReservationsParams{
		Limit:  listRequest.PageSize,
		Offset: (listRequest.PageID - 1) * listRequest.PageSize,
	}
	reservations, err := server.store.ListAvailableReservations(ctx, arg)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, reservations)
}
