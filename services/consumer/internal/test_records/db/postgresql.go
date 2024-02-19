package test_records

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"

	"consumer/internal/pkg/client/postgresql"
	"consumer/internal/test_records"
)

type repository struct {
	client postgresql.Client
}

func (r *repository) Create(ctx context.Context, record *test_records.TestRecord) error {
	q := `
		INSERT INTO test_records
			(text, created, stored)
		VALUES 
		    ($1, $2, $3)
		RETURNING id
	`
	if err := r.client.QueryRow(ctx, q, record.Text, record.Created, record.Stored).Scan(&record.ID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)

			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			fmt.Println(newErr)
			return newErr
		}
		return err
	}

	return nil
}

func (r *repository) Update(ctx context.Context, record *test_records.TestRecord) error {
	q := `
		UPDATE test_records
		SET 
		    text = $2,
		    created = $3,
		    stored = $4
		WHERE id = $1
	`
	return r.client.QueryRow(ctx, q, record.ID, record.Text, record.Created, record.Stored).Scan()
}

func (r *repository) FindAll(ctx context.Context) (u []test_records.TestRecord, err error) {
	q := `
    	SELECT id, text, created, stored from public.test_records;
	`

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	records := make([]test_records.TestRecord, 0)

	for rows.Next() {
		var rec test_records.TestRecord

		err = rows.Scan(&rec.ID, &rec.Text, &rec.Created, &rec.Stored)
		if err != nil {
			return nil, err
		}

		records = append(records, rec)

		if err = rows.Err(); err != nil {
			return nil, err
		}

	}
	return records, nil
}

func (r *repository) FindOne(ctx context.Context, id string) (test_records.TestRecord, error) {
	q := `
		SELECT 
		    id, 
		    text, 
		    created, 
		    stored 
		FROM public.test_records 
		WHERE id = $1;
	`

	var rec test_records.TestRecord

	err := r.client.QueryRow(ctx, q, id).Scan(&rec.ID, &rec.Text, &rec.Created, &rec.Stored)

	if err != nil {
		return test_records.TestRecord{}, err
	}

	return rec, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	q := `
		DELETE FROM public.test_records 
	    WHERE id = $1
	`
	return r.client.QueryRow(ctx, q, id).Scan()
}

func NewRepository(c postgresql.Client) test_records.Repository {
	return &repository{
		client: c,
	}
}
