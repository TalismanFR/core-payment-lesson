package pg

import (
	"auth/internal/config"
	"auth/internal/model/principal"
	"auth/internal/repo"
	"auth/pkg/logger"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

var _ repo.PrincipalRepo = (*RepoPG)(nil)

type RepoPG struct {
	tableName string
	db        *sql.DB
	logger    logger.Logger
}

func NewRepoPG(cfg *config.PostgresConfig, logger logger.Logger) (*RepoPG, error) {
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBname, cfg.Sslmode)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &RepoPG{
		tableName: "principal",
		db:        db,
		logger:    logger,
	}, nil
}

func (r *RepoPG) CreatePrincipal(ctx context.Context, p principal.Principal) error {

	r.logger.Debug("CreatePrincipal")
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault, ReadOnly: false})
	if err != nil {
		r.logger.Error(err.Error())
		return err
	}

	defer func() {
		if err := tx.Rollback(); err != nil {
			r.logger.Error("Rollback: err: %s", err)
		}
	}()

	q := "INSERT INTO principal(email, hashed_password, role) VALUES ($1, $2, $3)"
	r.logger.Info("RepoPG.GetPrincipal: query: %s", q)

	if _, err = tx.ExecContext(ctx, q, p.Email, p.HashedPassword, p.Role); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		r.logger.Error("Tx.Commit: err: %s", err)
		return err
	}

	return nil
}

func (r *RepoPG) UpdatePrincipal(ctx context.Context, p principal.Principal) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault, ReadOnly: false})
	if err != nil {
		r.logger.Error("BeginTx: err: %s", err)
		return err
	}

	defer func() {
		if err := tx.Rollback(); err != nil {
			r.logger.Error("Rollback: err: %s", err)
		}
	}()

	q := "UPDATE public.principal SET hashed_password=$1, role=$2 WHERE email = $3"
	r.logger.Info("RepoPG.GetPrincipal: query: %s", q)

	if _, err := tx.ExecContext(ctx, q, p.HashedPassword, p.Role, p.Email); err != nil {
		r.logger.Error("Tx.ExecContext: update failed: %v", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		r.logger.Error("Tx.Commit: err: %s", err)
		return err
	}

	return nil
}

func (r *RepoPG) GetPrincipal(ctx context.Context, email string) (*principal.Principal, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault, ReadOnly: true})
	if err != nil {
		r.logger.Error("BeginTx: err: %s", err)
		return nil, err
	}

	defer func() {
		if err := tx.Rollback(); err != nil {
			r.logger.Error("Rollback: err: %s", err)
		}
	}()

	q := "SELECT email, hashed_password, role FROM public.principal WHERE email = $1"
	r.logger.Info("RepoPG.GetPrincipal: query: %s", q)

	row := tx.QueryRowContext(ctx, q, email)

	p := &principal.Principal{}

	if err := row.Scan(&p.Email, &p.HashedPassword, &p.Role); err != nil {
		r.logger.Error("Row.Scan: err: %s", err)
		return nil, err
	}
	if err := row.Err(); err != nil {
		r.logger.Error("Tx.ExecContext: select failed: %v", err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		r.logger.Error("Tx.Commit: err: %s", err)
		return nil, err
	}

	return p, nil
}

// TODO: remove
func (r *RepoPG) Migrate(ctx context.Context) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault, ReadOnly: false})
	if err != nil {
		r.logger.Error("BeginTx: err: %s", err)
		return err
	}

	defer func() {
		if err := tx.Rollback(); err != nil {
			r.logger.Error("Rollback: err: %s", err)
		}
	}()

	q := `CREATE TABLE principal(email text, hashed_password bytea, role text);
CREATE INDEX on principal (email);
`

	_, err = tx.ExecContext(ctx, q)

	if err != nil {
		r.logger.Error("Tx.Commit: err: %s", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		r.logger.Error("Tx.Commit: err: %s", err)
		return err
	}

	return nil
}
