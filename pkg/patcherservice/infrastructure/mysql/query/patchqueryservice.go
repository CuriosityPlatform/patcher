package query

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/infrastructure/mysql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"patcher/pkg/patcherservice/app"
)

func NewPatchQueryService(client mysql.Client) app.PatchQueryService {
	return &patchQueryService{client: client}
}

type patchQueryService struct {
	client mysql.Client
}

func (service *patchQueryService) GetPatch(id app.PatchID) (app.Patch, error) {
	const selectSQL = `SELECT patch_id, project, applied, author, device, created_at FROM patch WHERE patch_id = ?`

	binaryID, err := uuid.UUID(id).MarshalBinary()
	if err != nil {
		return app.Patch{}, errors.WithStack(err)
	}

	var patch patchSqlx

	err = service.client.Get(&patch, selectSQL, binaryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return app.Patch{}, app.ErrPatchNotFound
		}
		return app.Patch{}, errors.WithStack(err)
	}

	return app.Patch{
		ID:        app.PatchID(patch.ID),
		Applied:   patch.Applied,
		Author:    app.PatchAuthor(patch.Author),
		Device:    app.Device(patch.Device),
		CreatedAt: &patch.CreatedAt,
	}, err
}

func (service *patchQueryService) GetPatchContent(id app.PatchID) (app.PatchContent, error) {
	const selectSQL = `SELECT content FROM patch WHERE patch_id = ?`

	binaryID, err := uuid.UUID(id).MarshalBinary()
	if err != nil {
		return app.PatchContent{}, errors.WithStack(err)
	}

	var patchContent []byte

	err = service.client.Get(&patchContent, selectSQL, binaryID)
	if err != nil {
		return nil, err
	}

	return patchContent, nil
}

//nolint:exportloopref
func (service *patchQueryService) GetPatches(spec app.PatchSpecification) ([]app.Patch, error) {
	selectSQL := `SELECT patch_id, project, applied, author, device, created_at FROM patch`

	conditions, args, err := getWhereConditionsBySpec(spec)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if conditions != "" {
		selectSQL += fmt.Sprintf(` WHERE %s`, conditions)
	}

	var patches []patchSqlx

	err = service.client.Select(&patches, selectSQL, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if len(patches) == 0 {
		return nil, errors.WithStack(err)
	}

	result := make([]app.Patch, 0, len(patches))

	for _, patch := range patches {
		result = append(result, app.Patch{
			ID:        app.PatchID(patch.ID),
			Project:   app.Project(patch.Project),
			Applied:   patch.Applied,
			Author:    app.PatchAuthor(patch.Author),
			Device:    app.Device(patch.Device),
			CreatedAt: &patch.CreatedAt,
		})
	}

	return result, nil
}

//nolint
func getWhereConditionsBySpec(spec app.PatchSpecification) (string, []interface{}, error) {
	var conditions []string
	var params []interface{}

	if len(spec.PatchIDS) != 0 {
		ids, err := patchIDsToBinaryUUIDs(spec.PatchIDS)
		if err != nil {
			return "", nil, errors.WithStack(err)
		}
		sqlQuery, args, err := sqlx.In(`patch_id IN (?)`, ids)
		if err != nil {
			return "", nil, errors.WithStack(err)
		}
		conditions = append(conditions, sqlQuery)
		for _, arg := range args {
			params = append(params, arg)
		}
	}

	if len(spec.Projects) != 0 {
		sqlQuery, args, err := sqlx.In(`project IN (?)`, spec.Projects)
		if err != nil {
			return "", nil, errors.WithStack(err)
		}
		conditions = append(conditions, sqlQuery)
		for _, arg := range args {
			params = append(params, arg)
		}
	}

	if len(spec.Authors) != 0 {
		sqlQuery, args, err := sqlx.In(`author IN (?)`, spec.Authors)
		if err != nil {
			return "", nil, errors.WithStack(err)
		}
		conditions = append(conditions, sqlQuery)
		for _, arg := range args {
			params = append(params, arg)
		}
	}

	if len(spec.Devices) != 0 {
		sqlQuery, args, err := sqlx.In(`device IN (?)`, spec.Devices)
		if err != nil {
			return "", nil, errors.WithStack(err)
		}
		conditions = append(conditions, sqlQuery)
		for _, arg := range args {
			params = append(params, arg)
		}
	}

	if spec.After != nil {
		conditions = append(conditions, "created_at > ?")
		params = append(params, spec.After)
	}

	if spec.Before != nil {
		conditions = append(conditions, "created_at < ?")
		params = append(params, spec.Before)
	}

	if spec.ShowApplied != nil {
		conditions = append(conditions, "applied = ?")
		intApplied := 0
		if *spec.ShowApplied {
			intApplied = 1
		}
		params = append(params, intApplied)
	}

	return strings.Join(conditions, " AND "), params, nil
}

func patchIDsToBinaryUUIDs(uuids []app.PatchID) ([][]byte, error) {
	res := make([][]byte, len(uuids))
	for i, id := range uuids {
		binaryUUID, err := uuid.UUID(id).MarshalBinary()
		if err != nil {
			return nil, err
		}
		res[i] = binaryUUID
	}
	return res, nil
}

type patchSqlx struct {
	ID        uuid.UUID `db:"patch_id"`
	Project   string    `db:"project"`
	Applied   bool      `db:"applied"`
	Author    string    `db:"author"`
	Device    string    `db:"device"`
	CreatedAt time.Time `db:"created_at"`
}
