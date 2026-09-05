package repository

import (
	"bitbucket.org/lyndus/backend/domain/client/entity"
	"bitbucket.org/lyndus/backend/infra/config"
	"bitbucket.org/lyndus/backend/infra/criteria"
	"bitbucket.org/lyndus/backend/infra/db"
	"bitbucket.org/lyndus/backend/infra/logger"
	"bitbucket.org/lyndus/backend/infra/utils"
	"database/sql"
	"fmt"
	"go.uber.org/zap"
)

type ClientUserRepository struct {
	db.Connector
}

func NewClientUserRepository() ClientUserRepository {
	repo := ClientUserRepository{}
	repo.UsePostgres()
	repo.Table = "client_user"
	return repo
}

func (cur ClientUserRepository) Update(cliente entity.ClientUser) error {
	command := `UPDATE client_user SET
						phone=$1,
						name=$2,
						email=$3
					WHERE id= $4`

	stmt, err := cur.Conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(cliente.Phone, cliente.Name, cliente.Email, cliente.ID)

	return err
}

func (cur ClientUserRepository) Get(clientID, clientUserID int) (clientUser entity.ClientUser, err error) {
	command := `SELECT  id,
       					phone,
       					username,
       					name,
       					email,
       					status
					FROM client_user
					WHERE client_id = $1 AND  id = $2 `

	err = cur.Conn.Get(&clientUser, command, clientID, clientUserID)
	if err == sql.ErrNoRows {
		return clientUser, nil
	}
	if err != nil {
		return clientUser, err
	}

	command = `SELECT 	client_permission_id as id,
       						access_level as level_id 
			   			FROM client_user_permission
						WHERE client_user_id= $1`

	err = cur.Conn.Select(&clientUser.Permissions, command, clientUserID)
	if err == sql.ErrNoRows {
		err = nil
	}

	for i := range clientUser.Permissions {
		clientUser.Permissions[i].LevelName = config.LEVEL[clientUser.Permissions[i].LevelId]
		clientUser.Permissions[i].Desc = config.CLIENTPERMISSIONS[clientUser.Permissions[i].Id]
	}
	return clientUser, err
}

//func (cur ClientUserRepository) AddPermissions(cliente entity.ClientUser) (PermissionsID []utils.LevelPermission) {

func (cur ClientUserRepository) AddPermissions(cliente entity.ClientUser) {
	fmt.Printf("%+v\n", cliente)
	fmtError := func(cliID, permissionID, levelID int, err error) []zap.Field {
		return []zap.Field{
			zap.Int("client_user_id", cliID),
			zap.Int("client_permission_id", permissionID),
			zap.Int("access_level", levelID),
			zap.String("erro", err.Error()),
		}
	}

	for i := range cliente.Permissions {
		exist, err := cur.PermissionExist(cliente.ID, cliente.Permissions[i].Id)
		if err != nil {
			erros := fmtError(cliente.ID, cliente.Permissions[i].Id, cliente.Permissions[i].LevelId, err)
			logger.Error("FALHA AO ADICIONAR PERMISSÃO AO REGISTRAR CLIENTE USUÁRIO", erros...)
		}
		if !exist {
			command := `INSERT INTO  client_user_permission(client_user_id,
                                    						client_permission_id,
                                    						access_level)
							VALUES($1,$2,$3) RETURNING ID`
			stmt, err := cur.Conn.Prepare(command)

			if err != nil {
				erros := fmtError(cliente.ID, cliente.Permissions[i].Id, cliente.Permissions[i].LevelId, err)
				logger.Error("FALHA AO ADICIONAR PERMISSÃO AO REGISTRAR CLIENTE USUÁRO", erros...)
			}

			defer stmt.Close()
			var id int
			err = stmt.QueryRow(cliente.ID, cliente.Permissions[i].Id, cliente.Permissions[i].LevelId).Scan(&id)

			if err != nil {
				erros := fmtError(cliente.ID, cliente.Permissions[i].Id, cliente.Permissions[i].LevelId, err)
				logger.Error("FALHA AO ADICIONAR PERMISSÃO AO REGISTRAR CLIENTE USUÁRIO", erros...)
			}

			//PermissionsID = append(PermissionsID, utils.LevelPermission{
			//	Id:      cliente.Permissions[i].Id,
			//	LevelId: cliente.Permissions[i].LevelId,
			//})
		} else {
			command := `UPDATE client_user_permission SET access_level = $1 
							WHERE client_user_id = $2 AND client_permission_id = $3
							RETURNING ID`
			stmt, err := cur.Conn.Prepare(command)

			if err != nil {
				erros := fmtError(cliente.ID, cliente.Permissions[i].Id, cliente.Permissions[i].LevelId, err)
				logger.Error("FALHA AO ALTERAR PERMISSÃO AO REGISTRAR CLIENTE USUÁRO", erros...)
			}

			defer stmt.Close()

			var id int
			err = stmt.QueryRow(cliente.Permissions[i].LevelId, cliente.ID, cliente.Permissions[i].Id).Scan(&id)

			if err != nil {
				erros := fmtError(cliente.ID, cliente.Permissions[i].Id, cliente.Permissions[i].LevelId, err)
				logger.Error("FALHA AO ALTERAR PERMISSÃO AO REGISTRAR CLIENTE USUÁRIO", erros...)
			}

			//PermissionsID = append(PermissionsID, utils.LevelPermission{
			//	Id:      cliente.Permissions[i].Id,
			//	LevelId: cliente.Permissions[i].LevelId,
			//})
		}

	}

	return
}

func (cur *ClientUserRepository) Create(cliente entity.ClientUser) (int, error) {
	command := `INSERT INTO  client_user(phone,
                         				username,
                         				name,
                         				created_by,
                         				status,
                         				client_id,
                         				email,
                         				authentication_id)  
						VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING ID`

	stmt, err := cur.Conn.Prepare(command)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRow(
		cliente.Phone,
		cliente.Username,
		cliente.Name,
		cliente.CreatedBy,
		cliente.StatusID,
		cliente.ClientID,
		cliente.Email,
		cliente.AuthenticationID,
	).Scan(&id)

	return id, err
}

func (cur ClientUserRepository) GetAll(activePage, itemsPerPage, clientID, status int, search, orderBy string, sortDesc bool) (clientsUser []entity.ClientUser, totalItems, totalPages int, err error) {
	command := `SELECT  count(*)
					FROM client_user
					WHERE client_id = $1 
					  AND status = $2`

	if search != "" {
		command = command + ` AND unaccent(name) ILIKE $3 `
		err = cur.Conn.Get(&totalItems, command, clientID, status, "%"+search+"%")
	} else {
		err = cur.Conn.Get(&totalItems, command, clientID, status)
	}

	if err != nil {
		return clientsUser, totalItems, totalPages, err
	}
	rest := totalItems % itemsPerPage
	totalPages = totalItems / itemsPerPage
	if rest > 0 {
		totalPages += 1
	}

	command = `SELECT id,
       				  phone,
       				  username,
       				  name,
       				  email
					FROM client_user
					WHERE client_id = $3
					  AND status = $4 `

	if search != "" {
		command = command + ` AND (unaccent(name) ILIKE $5 `
	}

	if orderBy == "name" {
		command = command + ` ORDER BY name `
	} else if orderBy == "phone" {
		command = command + ` ORDER BY phone `
	} else if orderBy == "email" {
		command = command + ` ORDER BY email `
	} else if orderBy == "username" {
		command = command + ` ORDER BY username `
	} else {
		command = command + ` ORDER BY name `
	}

	if sortDesc {
		command = command + " DESC "
	}
	command = command + ` OFFSET $1 LIMIT $2`

	if search != "" {
		err = cur.Conn.Select(&clientsUser, command, (activePage-1)*itemsPerPage, itemsPerPage, clientID, status, "%"+search+"%")
	} else {
		err = cur.Conn.Select(&clientsUser, command, (activePage-1)*itemsPerPage, itemsPerPage, clientID, status)
	}
	if err == sql.ErrNoRows {
		err = nil
	}

	return clientsUser, totalItems, totalPages, err
}

func (cur ClientUserRepository) GenerateUsername(cliente entity.ClientUser) (string, error) {
	command := `SELECT count(*) 
					FROM client_user 
					WHERE client_id=$1`
	stmt, err := cur.Conn.Prepare(command)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	var n int
	err = stmt.QueryRow(cliente.ClientID).Scan(&n)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%03d", n+1), err
}
func (cur ClientUserRepository) PermissionExist(cliUserID int, levelId int) (bool, error) {
	command := `SELECT count(*) 
					FROM client_user_permission 
					WHERE client_user_id=$1 AND client_permission_id=$2`
	stmt, err := cur.Conn.Prepare(command)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	var total int
	err = stmt.QueryRow(cliUserID, levelId).Scan(&total)
	if err != nil {
		return false, err
	}
	return total > 0, err
}

func (cur ClientUserRepository) UserPermissions(clientUserID int) (permissions []utils.LevelPermission, err error) {
	command := `SELECT 	client_permission_id as id,
       						access_level as level_id 
			   			FROM client_user_permission
						WHERE client_user_id= $1`

	err = cur.Conn.Select(&permissions, command, clientUserID)
	if err == sql.ErrNoRows {
		err = nil
	}

	for i := range permissions {
		permissions[i].LevelName = config.LEVEL[permissions[i].LevelId]
		permissions[i].Desc = config.CLIENTPERMISSIONS[permissions[i].Id]
	}
	return permissions, err
}

func (cur ClientUserRepository) GetMasterUserAuthID(id int) (userID int, err error) {
	command := `SELECT authentication_id FROM client_user WHERE master = true AND client_id = $1`
	err = cur.Conn.Get(&userID, command, id)
	return
}

func (cur ClientUserRepository) GetAllV2(ctr criteria.Criteria) (response criteria.ClientClientUserResponse, err error) {
	command := `SELECT  id,
       					username,
       					name,
       					email
					FROM client_user
					WHERE client_id = $1 AND status=$2 AND lower(client_user::text) LIKE '%' || lower($3) || '%'`
	command = ctr.Query(command)
	err = cur.Conn.Select(&response.Items, command, ctr.Args...)
	if err == sql.ErrNoRows {
		return response, nil
	}
	err = ctr.ExecWithQuery(cur, `WHERE client_id = $1 AND status=$2 AND lower(client_user::text) LIKE '%' || lower($3) || '%'`, ctr.Args...)
	if err != nil {
		return response, err
	}
	response.CResponse = response.NewWithCriteria(ctr)
	if response.Items == nil {
		response.Items = make([]entity.ClientUser, 0)
	}
	return response, err
}

func (cur ClientUserRepository) UpdateStatus(clientID, userID, statusID int) error {
	command := `UPDATE client_user SET
						status = $3
					WHERE id = $2 
					  AND client_id = $1`

	stmt, err := cur.Conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(clientID, userID, statusID)

	return err
}
