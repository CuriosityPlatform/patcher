package repo

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/infrastructure/mysql"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"patcher/pkg/patcherservice/app"
)

func NewPatchRepository(client mysql.Client) app.PatchRepository {
	return &patchRepository{client: client}
}

type patchRepository struct {
	client mysql.Client
}

func (repo *patchRepository) Find(id app.PatchID) (app.Patch, error) {
	const selectSQL = `SELECT patch_id, project, applied, author, device, created_at FROM patch WHERE patch_id = ?`

	binaryID, err := uuid.UUID(id).MarshalBinary()
	if err != nil {
		return app.Patch{}, errors.WithStack(err)
	}

	var patch patchSqlx

	err = repo.client.Get(&patch, selectSQL, binaryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return app.Patch{}, app.ErrPatchNotFound
		}
		return app.Patch{}, errors.WithStack(err)
	}

	return app.Patch{
		ID:        app.PatchID(patch.ID),
		Project:   app.Project(patch.Project),
		Applied:   patch.Applied,
		Author:    app.PatchAuthor(patch.Author),
		Device:    app.Device(patch.Device),
		CreatedAt: &patch.CreatedAt,
	}, err
}

func (repo *patchRepository) Store(patch app.Patch) error {
	insertSQL := `
		INSERT INTO patch (patch_id, project, applied, content, author, device, created_at) VALUES(?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY 
		UPDATE patch_id=VALUES(patch_id), project=VALUES(project), applied=VALUES(applied), %s, author=VALUES(author), device=VALUES(device), created_at=VALUES(created_at)
	`

	if patch.CreatedAt == nil {
		now := time.Now()
		patch.CreatedAt = &now
	}

	binaryUUID, err := uuid.UUID(patch.ID).MarshalBinary()
	if err != nil {
		return errors.WithStack(err)
	}

	if patch.Content != nil {
		insertSQL = fmt.Sprintf(insertSQL, "content=VALUES(content)")

		_, err = repo.client.Exec(
			insertSQL,
			binaryUUID,
			patch.Project,
			patch.Applied,
			patch.Content,
			patch.Author,
			patch.Device,
			patch.CreatedAt,
		)
	} else {
		insertSQL = fmt.Sprintf(insertSQL, "content=content")

		_, err = repo.client.Exec(
			insertSQL,
			binaryUUID,
			patch.Project,
			patch.Applied,
			"",
			patch.Author,
			patch.Device,
			patch.CreatedAt,
		)
	}

	if err != nil {
		return err
	}

	return nil
}

type patchSqlx struct {
	ID        uuid.UUID `db:"patch_id"`
	Project   string    `db:"project"`
	Applied   bool      `db:"applied"`
	Author    string    `db:"author"`
	Device    string    `db:"device"`
	CreatedAt time.Time `db:"created_at"`
}
