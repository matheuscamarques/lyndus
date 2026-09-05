package criteria

import (
	"strconv"

	"bitbucket.org/lyndus/backend/global/contracts"
)

type Criteria struct {
	ActivePage   int
	ItemsPerPage int
	TotalPages   int
	TotalItems   int
	Search       string
	Command      string
	Args         []interface{}
}

func (c *Criteria) Exec(repo contracts.RepositoryInterface) error {
	c.CheckPages()
	return c.total(repo)
}

func (c *Criteria) ExecWithQuery(repo contracts.RepositoryInterface, command string, args ...interface{}) error {
	c.CheckPages()
	c.Command = command
	return c.totalWithQuery(repo, command, args...)
}

func (c *Criteria) CheckPages() {
	if c.ItemsPerPage <= 0 || c.ItemsPerPage >= 100 {
		c.ItemsPerPage = 10
	}
	if c.ActivePage <= 0 {
		c.ActivePage = 1
	}
}

func (c Criteria) pages(total int) int {
	rest := total % c.ItemsPerPage
	total = total / c.ItemsPerPage
	if rest >= 1 {
		total++
	}
	return total
}

func (c *Criteria) totalWithQuery(repo contracts.RepositoryInterface, command string, args ...interface{}) (err error) {
	c.TotalItems, err = c.totalItemsWithQuery(repo, command, args...)
	c.TotalPages = c.pages(c.TotalItems)
	return err
}
func (c *Criteria) total(repo contracts.RepositoryInterface) (err error) {
	c.TotalItems, err = c.totalItems(repo)
	c.TotalPages = c.pages(c.TotalItems)
	return err
}

func (c *Criteria) totalItems(repo contracts.RepositoryInterface) (total int, err error) {
	table := repo.GetTable()
	command := `SELECT count(*) FROM ` + table
	err = repo.GetConnection().Get(&total, command)
	return total, err
}

func (c *Criteria) TotalPagesWithWhere(repo contracts.RepositoryInterface, queryWhere string) (total int, err error) {
	table := repo.GetTable()
	c.CheckPages()
	command := `SELECT count(*) as pages FROM ` + table + " " + queryWhere
	err = repo.GetConnection().Get(&total, command)
	total = c.pages(total)
	return total, err
}

func (c Criteria) Query(query string) string {
	if c.ActivePage == 0 && c.ItemsPerPage == 0 {
		return query
	}
	return query + " OFFSET " + strconv.Itoa((c.ActivePage-1)*c.ItemsPerPage) + " LIMIT " + strconv.Itoa(c.ItemsPerPage)
}

func (c *Criteria) totalItemsWithQuery(repo contracts.RepositoryInterface, command string, args ...interface{}) (total int, err error) {
	table := repo.GetTable()
	command = `SELECT count(*) FROM ` + table + " " + command
	err = repo.GetConnection().Get(&total, command, args...)
	return total, err
}
