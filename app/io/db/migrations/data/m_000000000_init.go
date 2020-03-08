// data Migration file
package data

import (
	"github.com/dimonrus/godb"
	"gost/app/base"
)

type m_000000000_init struct{}

func init() {
	base.App.GetMigration().Registry["data"] = append(base.App.GetMigration().Registry["data"], m_000000000_init{})
}

func (m m_000000000_init) GetVersion() string {
	return "m_000000000_init"
}

func (m m_000000000_init) Up(tx *godb.SqlTx) error {
	return nil
}

func (m m_000000000_init) Down(tx *godb.SqlTx) error {
	return nil
}
