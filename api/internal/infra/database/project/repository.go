package project

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/analopesdev/project-manager-clicksign/api/internal/domain/project"
)

// orderings maps each supported sort to a fixed ORDER BY. The value from the request is only ever
// used as a key, never concatenated into the SQL. Every ordering ends with id so it is stable.
var orderings = map[string]string{
	project.SortName:      "lower(name), id",
	project.SortStartDate: "start_date DESC, id",
	// Upcoming deadlines first (soonest first), finished projects last (most recently finished first).
	project.SortEndDate: `(end_date < CURRENT_DATE),
		CASE WHEN end_date >= CURRENT_DATE THEN end_date END ASC,
		CASE WHEN end_date < CURRENT_DATE THEN end_date END DESC, id`,
}

type ProjectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

type scanner interface {
	Scan(dest ...any) error
}

func scanProject(s scanner) (*project.Project, error) {
	var (
		p          project.Project
		start, end time.Time
		cover      sql.NullString
	)
	err := s.Scan(&p.ID, &p.Name, &p.Client, &start, &end, &cover, &p.Favorite, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	p.StartDate = start.Format(project.DateLayout)
	p.EndDate = end.Format(project.DateLayout)
	p.CoverURL = cover.String
	return &p, nil
}

func (r *ProjectRepository) Create(ctx context.Context, p *project.Project) (*project.Project, error) {
	row := r.db.QueryRowContext(ctx,
		`INSERT INTO projects (name, client, start_date, end_date, cover_url)
		VALUES ($1, $2, $3, $4, NULLIF($5, ''))
		RETURNING id, name, client, start_date, end_date, cover_url, favorite, created_at, updated_at`,
		p.Name, p.Client, p.StartDate, p.EndDate, p.CoverURL)

	created, err := scanProject(row)
	if err != nil {
		return nil, fmt.Errorf("creating project: %w", err)
	}
	return created, nil
}

func (r *ProjectRepository) GetById(ctx context.Context, id string) (*project.Project, error) {
	p, err := scanProject(r.db.QueryRowContext(ctx, `SELECT id, name, client, start_date, end_date, cover_url, favorite, created_at, updated_at FROM projects WHERE id = $1`, id))

	if errors.Is(err, sql.ErrNoRows) {
		return nil, project.ErrProjectNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("getting project: %w", err)
	}

	return p, nil
}

func (r *ProjectRepository) List(ctx context.Context, search string, sort string, favoriteOnly bool) ([]*project.Project, int, error) {
	var (
		where []string
		args  []any
	)

	if search != "" {
		args = append(args, "%"+escapeLike(search)+"%")
		where = append(where, fmt.Sprintf("(name ILIKE $%d OR client ILIKE $%d)", len(args), len(args)))
	}
	if favoriteOnly {
		where = append(where, "favorite")
	}

	query := `SELECT id, name, client, start_date, end_date, cover_url, favorite, created_at, updated_at FROM projects`
	if len(where) > 0 {
		query += ` WHERE ` + strings.Join(where, " AND ")
	}

	if sort == "" {
		sort = project.SortName
	}
	orderBy, ok := orderings[sort]
	if !ok {
		return nil, 0, fmt.Errorf("unsupported sort %q", sort)
	}
	query += ` ORDER BY ` + orderBy

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("querying projects: %w", err)
	}
	defer rows.Close()

	projects := make([]*project.Project, 0)
	for rows.Next() {
		p, err := scanProject(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("scanning project: %w", err)
		}
		projects = append(projects, p)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterating projects: %w", err)
	}

	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT count(*) FROM projects`).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("counting projects: %w", err)
	}

	return projects, total, nil
}

func (r *ProjectRepository) Update(ctx context.Context, p *project.Project) (*project.Project, error) {
	row := r.db.QueryRowContext(ctx,
		`UPDATE projects
		SET name = $1, client = $2, start_date = $3, end_date = $4,
		    cover_url = NULLIF($5, ''), updated_at = now()
		WHERE id = $6
		RETURNING id, name, client, start_date, end_date, cover_url, favorite, created_at, updated_at`,
		p.Name, p.Client, p.StartDate, p.EndDate, p.CoverURL, p.ID)

	updated, err := scanProject(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, project.ErrProjectNotFound
	}
	return updated, err
}

func (r *ProjectRepository) SetFavorite(ctx context.Context, id string, favorite bool) (*project.Project, error) {
	row := r.db.QueryRowContext(ctx,
		`UPDATE projects SET favorite = $1 WHERE id = $2 RETURNING id, name, client, start_date, end_date, cover_url, favorite, created_at, updated_at`,
		favorite, id)

	p, err := scanProject(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, project.ErrProjectNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("setting favorite: %w", err)
	}
	return p, nil
}

func (r *ProjectRepository) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM projects WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("deleting project: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("deleting project: %w", err)
	}
	if n == 0 {
		return project.ErrProjectNotFound
	}
	return nil
}

func escapeLike(s string) string {
	return strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(s)
}
