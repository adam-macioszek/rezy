package api

import (
	"os"
	"testing"
	"time"

	"github.com/adam-macioszek/rezy/config"
	db "github.com/adam-macioszek/rezy/db/sqlc"
	"github.com/adam-macioszek/rezy/random"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := config.Config{
		ApiTokenKey:      random.RandomString(32),
		ApiTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
