package sap_segmentation

import (
	"fmt"

	"github.com/tdx/sap_segmentation/model"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PGStorer struct {
	db   *sqlx.DB
	stmt *sqlx.Stmt
}

// var errNoRowsAffected = errors.New("failed to store segment: no rows affected")

//
func NewPGStorer(cfg *DbConnection) (*PGStorer, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DbName)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("storer: failed to connect to postgres: %w", err)
	}

	if err := migration(db); err != nil {
		return nil, fmt.Errorf("storer: migration failed: %w", err)
	}

	stmt, err := db.Preparex(
		"INSERT INTO sap_segmentation(address_sap_id, adr_segment, segment_id) " +
			"VALUES($1,$2,$3) " +
			"ON CONFLICT (address_sap_id) " +
			"DO UPDATE SET " +
			"adr_segment=excluded.adr_segment, " +
			"segment_id=excluded.segment_id")
	if err != nil {
		return nil, fmt.Errorf("storer: failed to prepare statement: %w", err)
	}

	return &PGStorer{
		db:   db,
		stmt: stmt,
	}, nil
}

//
func (s *PGStorer) Store(segs []model.Segmentation) error {

	if len(segs) == 0 {
		return nil
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for i := range segs {
		seq := segs[i]
		_, err := s.stmt.Exec(seq.AddressSapID, seq.AdrSegment, seq.SegmentID)
		if err != nil {
			return err
		}

		// if res.RowsAffected() != 1 {
		// 	return errNoRowsAffected
		// }
	}

	return tx.Commit()
}

//
func migration(db *sqlx.DB) error {
	sql := `CREATE TABLE IF NOT EXISTS sap_segmentation(
    id              serial,
    address_sap_id  varchar(255) UNIQUE,
    adr_segment     varchar(16),
    segment_id      bigint
);`
	_, err := db.Exec(sql)
	return err
}
