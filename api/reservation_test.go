package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/adam-macioszek/rezy/config"
	mockdb "github.com/adam-macioszek/rezy/db/mock"
	db "github.com/adam-macioszek/rezy/db/sqlc"
	"github.com/adam-macioszek/rezy/random"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetReservationAPI(t *testing.T) {
	config, err := config.LoadConfig("./..")
	if err != nil {
		log.Println("cannot load config: ", err)
	}
	reservation := randomReservation()

	testCases := []struct {
		name          string
		reservationID int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{{
		name:          "OK",
		reservationID: reservation.ID,
		buildStubs: func(store *mockdb.MockStore) {
			store.EXPECT().
				GetReservation(gomock.Any(), gomock.Eq(reservation.ID)).
				Times(1).
				Return(reservation, nil)
		},
		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			require.Equal(t, http.StatusOK, recorder.Code)
			requireBodyMatchAccount(t, recorder.Body, reservation)
		},
	},
		{
			name:          "NotFound",
			reservationID: reservation.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetReservation(gomock.Any(), gomock.Eq(reservation.ID)).
					Times(1).
					Return(db.Reservation{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:          "InternalError",
			reservationID: reservation.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetReservation(gomock.Any(), gomock.Eq(reservation.ID)).
					Times(1).
					Return(db.Reservation{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		//add more cases
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)
			server, err := NewServer(config, store)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/reservation/%d", tc.reservationID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}

}

func randomReservation() db.Reservation {
	return db.Reservation{
		ID:        int64(random.RandomInt(0, 100)),
		StartTime: time.Now().Add(time.Hour),
		Booked:    false,
		Duration:  int32(time.Second),
		TableSize: (random.RandomInt(0, 100)),
	}
}
func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, reservation db.Reservation) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)
	var recievedReservation db.Reservation
	err = json.Unmarshal(data, &recievedReservation)
	require.NoError(t, err)
	require.Equal(t, reservation.ID, recievedReservation.ID)
	require.Equal(t, reservation.Booked, recievedReservation.Booked)
	require.Equal(t, reservation.TableSize, recievedReservation.TableSize)
}
