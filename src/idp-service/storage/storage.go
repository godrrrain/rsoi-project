package storage

import (
	"context"
	"fmt"
	"time"

	"lab2/src/idp-service/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgStorage struct {
	pool *pgxpool.Pool
}

func NewPgStorage(ctx context.Context, postgresURL string) (*PgStorage, error) {
	config, err := pgxpool.ParseConfig(postgresURL)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return &PgStorage{pool: pool}, nil
}

func (ps *PgStorage) Close() {
	ps.pool.Close()
}

func (ps *PgStorage) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	user.UserUid = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	err := ps.pool.QueryRow(ctx,
		`INSERT INTO idp_users (user_uid, username, email, password_hash, full_name, role)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, user_uid, username, email, password_hash, full_name, role, created_at, updated_at`,
		user.UserUid, user.Username, user.Email, user.PasswordHash, user.FullName, user.Role).
		Scan(&user.ID, &user.UserUid, &user.Username, &user.Email, &user.PasswordHash, &user.FullName, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ps *PgStorage) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user := &models.User{}
	err := ps.pool.QueryRow(ctx,
		`SELECT id, user_uid, username, email, password_hash, full_name, role, created_at, updated_at
		 FROM idp_users WHERE username = $1`,
		username).
		Scan(&user.ID, &user.UserUid, &user.Username, &user.Email, &user.PasswordHash, &user.FullName, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (ps *PgStorage) GetUserByUid(ctx context.Context, uid uuid.UUID) (*models.User, error) {
	user := &models.User{}
	err := ps.pool.QueryRow(ctx,
		`SELECT id, user_uid, username, email, password_hash, full_name, role, created_at, updated_at
		 FROM idp_users WHERE user_uid = $1`,
		uid).
		Scan(&user.ID, &user.UserUid, &user.Username, &user.Email, &user.PasswordHash, &user.FullName, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (ps *PgStorage) SaveAuthorizationCode(ctx context.Context, authCode *models.AuthCode) error {
	_, err := ps.pool.Exec(ctx,
		`INSERT INTO auth_codes (code, user_uid, client_id, redirect_uri, scope, expires_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		authCode.Code, authCode.UserUid, authCode.ClientId, authCode.RedirectUri, authCode.Scope, authCode.ExpiresAt)

	return err
}

func (ps *PgStorage) GetAuthorizationCode(ctx context.Context, code string) (*models.AuthCode, error) {
	authCode := &models.AuthCode{}
	err := ps.pool.QueryRow(ctx,
		`SELECT id, code, user_uid, client_id, redirect_uri, scope, expires_at, created_at
		 FROM auth_codes WHERE code = $1`,
		code).
		Scan(&authCode.ID, &authCode.Code, &authCode.UserUid, &authCode.ClientId, &authCode.RedirectUri, &authCode.Scope, &authCode.ExpiresAt, &authCode.CreatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("authorization code not found or expired")
		}
		return nil, err
	}

	if time.Now().After(authCode.ExpiresAt) {
		ps.pool.Exec(ctx, "DELETE FROM auth_codes WHERE id = $1", authCode.ID)
		return nil, fmt.Errorf("authorization code expired")
	}

	return authCode, nil
}

func (ps *PgStorage) DeleteAuthorizationCode(ctx context.Context, code string) error {
	_, err := ps.pool.Exec(ctx, "DELETE FROM auth_codes WHERE code = $1", code)
	return err
}

func (ps *PgStorage) SaveRefreshToken(ctx context.Context, refreshToken *models.RefreshToken) error {
	_, err := ps.pool.Exec(ctx,
		`INSERT INTO refresh_tokens (token, user_uid, client_id, scope, expires_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		refreshToken.Token, refreshToken.UserUid, refreshToken.ClientId, refreshToken.Scope, refreshToken.ExpiresAt)

	return err
}

func (ps *PgStorage) GetRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	refreshToken := &models.RefreshToken{}
	err := ps.pool.QueryRow(ctx,
		`SELECT id, token, user_uid, client_id, scope, expires_at, created_at
		 FROM refresh_tokens WHERE token = $1`,
		token).
		Scan(&refreshToken.ID, &refreshToken.Token, &refreshToken.UserUid, &refreshToken.ClientId, &refreshToken.Scope, &refreshToken.ExpiresAt, &refreshToken.CreatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("refresh token not found or expired")
		}
		return nil, err
	}

	if time.Now().After(refreshToken.ExpiresAt) {
		ps.pool.Exec(ctx, "DELETE FROM refresh_tokens WHERE id = $1", refreshToken.ID)
		return nil, fmt.Errorf("refresh token expired")
	}

	return refreshToken, nil
}
