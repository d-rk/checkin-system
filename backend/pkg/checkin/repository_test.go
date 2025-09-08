//go:build integration

package checkin

import (
	"context"
	"testing"

	"github.com/d-rk/checkin-system/pkg/app"
	"github.com/d-rk/checkin-system/pkg/database"
	"github.com/stretchr/testify/assert"
)

func TestGetLatestCheckinDate_EmptyTable_ReturnsNotFoundErr(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}
	db := database.SetupTestDB(t)
	repo := NewRepo(db)
	ctx := context.Background()
	ts, err := repo.GetLatestCheckinDate(ctx)
	assert.Nil(t, ts)
	assert.Equal(t, app.ErrNotFound, err)
}
