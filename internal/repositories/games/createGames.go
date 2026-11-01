package games

import (
	"context"
	models "steam_games/internal/models/games"
)

func (r *Repository) CreateGame(ctx context.Context, g models.Game) (int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// La colonne games.id n'a pas de séquence : verrou transactionnel
	// pour éviter que deux créations calculent le même identifiant.
	if _, err = tx.ExecContext(ctx, "SELECT pg_advisory_xact_lock(174829)"); err != nil {
		return 0, err
	}
	var id int
	err = tx.QueryRowContext(ctx, `INSERT INTO games
        (id,title,required_age,is_free,description,image,website,price,metacritic_score,recommendations,release_date,release_year)
        SELECT COALESCE(MAX(id),0)+1,$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11 FROM games
        RETURNING id`, g.Title, g.RequiredAge, g.IsFree, g.Description, g.Image, g.Website, g.Price, g.MetacriticScore, g.Recommendations, g.ReleaseDate, g.ReleaseYear).Scan(&id)
	if err != nil {
		return 0, err
	}
	if err = saveGameAssociations(ctx, tx, id, g); err != nil {
		return 0, err
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return id, nil
}
