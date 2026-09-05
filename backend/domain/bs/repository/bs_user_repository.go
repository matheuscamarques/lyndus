package repository

import (
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/infra/config"
	"bitbucket.org/lyndus/backend/infra/db"
	"bitbucket.org/lyndus/backend/infra/logger"
	"database/sql"
	"fmt"
	"go.uber.org/zap"
)

type BSUserRepository struct {
	db.Connector
}



func NewBSUserRepository() (repo BSUserRepository) {
	repo.UsePostgres()
	return repo
}

func (cur BSUserRepository) Update(bsUser entity.BSUser) error {
	command := `UPDATE bs_user SET
						phone = $1,
						name  = $2,
						email = $3
				WHERE id= $4`

	stmt, err := cur.Conn.Prepare(command)
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {

		}
	}(stmt)
	_, err = stmt.Exec(bsUser.Phone, bsUser.Name, bsUser.Email, bsUser.ID)

	return err
}

func (cur BSUserRepository) Get(bsID, bsUserID int) (BSUser entity.BSUser, err error) {
	command := `SELECT id, 
       				   phone,
       				   username,
       				   name,
       				   email
					FROM bs_user
					WHERE bs_id = $1 AND  id = $2 `

	err = cur.Conn.Get(&BSUser, command, bsID, bsUserID)
	if err == sql.ErrNoRows {
		return BSUser, nil
	}
	if err != nil {
		return BSUser, err
	}

	command = `SELECT bs_permission_id as id,
       				  access_level as level_id 
			   		FROM bs_user_permission
					WHERE bs_user_id= $1`

	err = cur.Conn.Select(&BSUser.Permissions, command, bsUserID)

	if err == sql.ErrNoRows {
		err = nil
	}

	for i := range BSUser.Permissions {
		BSUser.Permissions[i].LevelName = config.LEVEL[BSUser.Permissions[i].LevelId]
		BSUser.Permissions[i].Desc = config.BSPERMISSIONS[BSUser.Permissions[i].Id]
	}
	return BSUser, err
}

//func (cur BSUserRepository) AddPermissions(bsUser entity.BSUser) (PermissionsID []utils.LevelPermission) {

func (cur BSUserRepository) AddPermissions(bsUser entity.BSUser) {

	fmtError := func(bsUserID, permissionID, levelID int, err error) []zap.Field {
		return []zap.Field{
			zap.Int("bs_user_id", bsUserID),
			zap.Int("bs_permission_id", permissionID),
			zap.Int("access_level", levelID),
			zap.String("erro", err.Error()),
		}
	}

	for i := range bsUser.Permissions {
		exist, err := cur.PermissionExist(bsUser.ID, bsUser.Permissions[i].Id)
		if err != nil {
			erros := fmtError(bsUser.ID, bsUser.Permissions[i].Id, bsUser.Permissions[i].LevelId, err)
			logger.Error("FALHA AO ADICIONAR PERMISSÃO AO REGISTRAR BS USUÁRIO", erros...)
		}
		if !exist {
			command := `INSERT INTO  bs_user_permission(bs_user_id,
                                						bs_permission_id,
                                						access_level)
						VALUES($1,$2,$3) RETURNING ID`
			stmt, err := cur.Conn.Prepare(command)

			if err != nil {
				erros := fmtError(bsUser.ID, bsUser.Permissions[i].Id, bsUser.Permissions[i].LevelId, err)
				logger.Error("FALHA AO ADICIONAR PERMISSÃO AO REGISTRAR bs USUÁRIO", erros...)
			}

			defer stmt.Close()
			var id int
			err = stmt.QueryRow(bsUser.ID, bsUser.Permissions[i].Id, bsUser.Permissions[i].LevelId).Scan(&id)

			if err != nil {
				erros := fmtError(bsUser.ID, bsUser.Permissions[i].Id, bsUser.Permissions[i].LevelId, err)
				logger.Error("FALHA AO ADICIONAR PERMISSÃO AO REGISTRAR BS USUÁRIO", erros...)
			}

			//PermissionsID = append(PermissionsID, utils.LevelPermission{
			//	Id:      bsUser.Permissions[i].Id,
			//	LevelId: bsUser.Permissions[i].LevelId,
			//})
		} else {
			command := `UPDATE bs_user_permission SET access_level = $1 
							WHERE bs_user_id = $2 AND bs_permission_id = $3
							RETURNING ID`
			stmt, err := cur.Conn.Prepare(command)

			if err != nil {
				erros := fmtError(bsUser.ID, bsUser.Permissions[i].Id, bsUser.Permissions[i].LevelId, err)
				logger.Error("FALHA AO ALTERAR PERMISSÃO AO REGISTRAR BS USUÁRIO", erros...)
			}

			defer stmt.Close()

			var id int
			err = stmt.QueryRow(bsUser.Permissions[i].LevelId, bsUser.ID, bsUser.Permissions[i].Id).Scan(&id)

			if err != nil {
				erros := fmtError(bsUser.ID, bsUser.Permissions[i].Id, bsUser.Permissions[i].LevelId, err)
				logger.Error("FALHA AO ALTERAR PERMISSÃO AO REGISTRAR BS USUÁRIO", erros...)
			}

			//PermissionsID = append(PermissionsID, utils.LevelPermission{
			//	Id:      bsUser.Permissions[i].Id,
			//	LevelId: bsUser.Permissions[i].LevelId,
			//})
		}
	}
	return
}

func (cur BSUserRepository) Create(bsUser entity.BSUser) (int, error) {
	command := `INSERT INTO  bs_user(phone,
                     				username,
                     				name,
                     				created_by,
                     				status,
                     				bs_id,
                     				master,
                     				email,
                     				authentication_id)  
						VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING ID`

	stmt, err := cur.Conn.Prepare(command)
	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	var id int
	err = stmt.QueryRow(
		bsUser.Phone,
		bsUser.Username,
		bsUser.Name,
		bsUser.CreatedBy,
		bsUser.Status,
		bsUser.BsID,
		bsUser.Master,
		bsUser.Email,
		bsUser.AuthenticationID,
	).Scan(&id)

	return id, err
}

func (cur BSUserRepository) GetAll(bsID, status int) (BSUsers []entity.BSUser, err error) {

	command := `SELECT id,
       				   phone,
       				   username,
       				   name,
       				   email 
					FROM bs_user
					WHERE bs_id = $1 AND status = $2`

	err = cur.Conn.Select(&BSUsers, command, bsID, status)
	if err == sql.ErrNoRows {
		err = nil
	}

	return BSUsers, err
}

func (cur BSUserRepository) GenerateUsername(bsUser entity.BSUser) (string, error) {
	command := `SELECT count(*) 
					FROM bs_user 
					WHERE bs_id=$1`
	stmt, err := cur.Conn.Prepare(command)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	var n int
	err = stmt.QueryRow(bsUser.BsID).Scan(&n)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%03d", n+1), err
}
func (cur BSUserRepository) PermissionExist(bsUserID int, permissionID int) (bool, error) {
	command := `SELECT count(*)
					FROM bs_user_permission 
					WHERE bs_user_id=$1 AND bs_permission_id=$2`
	stmt, err := cur.Conn.Prepare(command)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	var total int
	err = stmt.QueryRow(bsUserID, permissionID).Scan(&total)
	if err != nil {
		return false, err
	}
	return total > 0, err
}


func (b BSUserRepository) GetMasterUserAuthID(id int) (userID int,err error) {
	command := `SELECT authentication_id FROM bs_user WHERE master = true AND bs_id = $1`
	err = b.Conn.Get(&userID, command, id)
	return
}