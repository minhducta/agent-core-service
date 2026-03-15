package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/minhducta/agent-core-service/internal/domain"
	"github.com/minhducta/agent-core-service/pkg/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestDB(t *testing.T) (*database.DB, sqlmock.Sqlmock) {
	t.Helper()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	t.Cleanup(func() { sqlxDB.Close() })

	return &database.DB{DB: sqlxDB}, mock
}

func TestBotRepository_GetByID(t *testing.T) {
	tests := []struct {
		name      string
		botID     uuid.UUID
		wantNil   bool
		setupMock func(mock sqlmock.Sqlmock, id uuid.UUID)
	}{
		{
			name:    "found",
			botID:   uuid.New(),
			wantNil: false,
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				now := time.Now()
				rows := sqlmock.NewRows([]string{
					"id", "name", "role", "vibe", "emoji", "avatar_url",
					"api_key_hash", "last_seen_at", "status", "created_at", "updated_at",
				}).AddRow(id, "TestBot", "assistant", "helpful", "🤖", "https://example.com/avatar.png",
					"hashvalue", now, domain.BotStatusActive, now, now)
				mock.ExpectQuery("SELECT .+ FROM bots WHERE id").
					WithArgs(id).
					WillReturnRows(rows)
			},
		},
		{
			name:    "not found",
			botID:   uuid.New(),
			wantNil: true,
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				mock.ExpectQuery("SELECT .+ FROM bots WHERE id").
					WithArgs(id).
					WillReturnRows(sqlmock.NewRows([]string{}))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := newTestDB(t)
			repo := NewBotRepository(db)
			tc.setupMock(mock, tc.botID)

			bot, err := repo.GetByID(context.Background(), tc.botID)

			assert.NoError(t, err)
			if tc.wantNil {
				assert.Nil(t, bot)
			} else {
				assert.NotNil(t, bot)
				assert.Equal(t, tc.botID, bot.ID)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestBotRepository_GetByAPIKeyHash(t *testing.T) {
	db, mock := newTestDB(t)
	repo := NewBotRepository(db)

	hash := "sha256hashvalue"
	id := uuid.New()
	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "name", "role", "vibe", "emoji", "avatar_url",
		"api_key_hash", "last_seen_at", "status", "created_at", "updated_at",
	}).AddRow(id, "BotAlpha", "agent", "curious", "🦾", "", hash, now, domain.BotStatusActive, now, now)

	mock.ExpectQuery("SELECT .+ FROM bots WHERE api_key_hash").
		WithArgs(hash).
		WillReturnRows(rows)

	bot, err := repo.GetByAPIKeyHash(context.Background(), hash)
	assert.NoError(t, err)
	assert.NotNil(t, bot)
	assert.Equal(t, hash, bot.APIKeyHash)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBotRepository_UpdateLastSeen(t *testing.T) {
	db, mock := newTestDB(t)
	repo := NewBotRepository(db)

	id := uuid.New()

	mock.ExpectExec("UPDATE bots SET last_seen_at").
		WithArgs(id, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UpdateLastSeen(context.Background(), id)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBotMemoryRepository_Create(t *testing.T) {
	db, mock := newTestDB(t)
	repo := NewBotMemoryRepository(db)

	memory := &domain.BotMemory{
		ID:         uuid.New(),
		BotID:      uuid.New(),
		Type:       domain.MemoryTypeFact,
		Content:    "I prefer concise answers",
		Tags:       []string{"preference", "style"},
		Importance: 5,
	}

	mock.ExpectExec("INSERT INTO bot_memories").
		WithArgs(
			memory.ID, memory.BotID, memory.Type, memory.Content,
			sqlmock.AnyArg(), memory.Importance, memory.ExpiresAt,
			sqlmock.AnyArg(), sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Create(context.Background(), memory)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBotMemoryRepository_Delete(t *testing.T) {
	db, mock := newTestDB(t)
	repo := NewBotMemoryRepository(db)

	id := uuid.New()
	botID := uuid.New()

	mock.ExpectExec("DELETE FROM bot_memories WHERE id").
		WithArgs(id, botID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Delete(context.Background(), id, botID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
