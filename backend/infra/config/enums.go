package config

const (
	USERSTATUSACTIVE   = 1
	USERSTATUSINACTIVE = 2
	USERSTATUSDELETED  = 3
)

var UserStatusPt = map[int]string{
	USERSTATUSACTIVE:   "Ativo",
	USERSTATUSINACTIVE: "Inativo",
	USERSTATUSDELETED:  "Deletado",
}
