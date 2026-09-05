package services

import (
	"bitbucket.org/lyndus/backend/domain/client/contracts"
	"bitbucket.org/lyndus/backend/domain/client/entity"
	"bitbucket.org/lyndus/backend/domain/client/repository"
	"bitbucket.org/lyndus/backend/global/aggregate"
	"bitbucket.org/lyndus/backend/infra/criteria"
	"strconv"
)

//A enquete ficará disponível somente para os USERS vinculados ao CLIENT;
//O CLIENT poderá criar uma enquete para todos os USERS da empresa ou para categorias específicas;
//Quando o CLIENT for definir por categoria ele poderá selecionar 1 ou + categorias, por exemplo, lançar uma enquete para as categorias RH e TI por exemplo.
//No momento da criação da enquete a empresa determina se o usuário deverá se identificar
//se não for necessário o usuário escolhe se vai se identificar ou irá responder como anônimo;
//Caso o CLIENT queira reaproveitar as perguntas de uma enquete já criada ele poderá duplicar a enquete em questão e editar/utilizar as perguntas que deseja;
//O CLIENT poderá encerrar manualmente a enquete antes da data de término da mesma;
//O CLIENT poder cancelar a enquete;
//Os status que enquete terá serão os seguintes:
//Criada - A partir do momento que ela é criada mas não enviada para aceitar respostas;
//Em andamento - Quando ela foi disparada para os users responderem;
//Encerrada - Quando termina o período de respostas ou o CLIENT clica para encerrar;
//Cancelada - Quando o CLIENT clica em cancelar.
//A diferença entre os status Encerrada e cancelada é: Quando uma enquete é encerrada automaticamente ou manualmente começa a contabilizar as respostas,
//quando ela é cancelada não são feitas as contabilizações de respostas (independente se tenha respostas ou não).

type ClientPollService struct {
	repo contracts.ClientPollRepositoryInterface
}


var Poll contracts.ClientPollServiceInterface

func NewClientPollService() * ClientPollService{
	repo := repository.NewClientPollRepository()
	return &ClientPollService{
		repo : repo,
	}
}

func (cps ClientPollService) Create(poll entity.ClientPoll)(id int, err error) {
	//poll entity.ClientPoll  o que preciso para criar uma Poll ?
	return cps.repo.Create(poll)
}

func (cps ClientPollService) Update(poll entity.ClientPoll)(err error) {
	return cps.repo.Update(poll)
}

func (cps ClientPollService)GetById(pollID int)(poll aggregate.PollAgregate,err error){
	return cps.repo.GetById(pollID)
}

func (cps ClientPollService) GetAll(clientID, ActivePage, ItemsPerPage int) (response criteria.CPollResponse, err error){
	criteria := criteria.Criteria{
		ActivePage:   ActivePage,
		ItemsPerPage: ItemsPerPage,
	}

	response.TotalPages, err = criteria.TotalPagesWithWhere(cps.repo,"WHERE client_id="+strconv.Itoa(clientID))

	if err != nil {
		return response, err
	}

	response.ActivePage = ActivePage
	response.Items, err = cps.repo.GetAllByIdClient(clientID,criteria)

	if err != nil {
		return response, err
	}

	return response, err
}
func (cps ClientPollService) ValidateClient(pollID, clientID int) (bool, error) {
	return cps.repo.ValidateClient(pollID, clientID)
}



func (cps ClientPollService) Close(poll entity.ClientPoll) {

}

func (cps ClientPollService) Cancel(poll entity.ClientPoll) {

}
