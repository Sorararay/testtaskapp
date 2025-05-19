package db

import "fmt"

type Datasource interface {
	Driver() string
	DSN() string
}

type mysqlDatasource struct {
	username string
	password string
	host     string
	dbname   string
}

func NewMySQLDatasource(username, password, host, dbname string) Datasource {
	return &mysqlDatasource{
		username: username,
		password: password,
		host:     host,
		dbname:   dbname,
	}
}

func (d *mysqlDatasource) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", d.username, d.password, d.host, d.dbname)
}

func (d mysqlDatasource) Driver() string {
	return "mysql"
}

// pgDatasource は PostgreSQL のデータソース情報を保持します。
type pgDatasource struct {
	driverName string
	dsn        string
}

// NewPgDatasource は新しい PostgreSQL データソースを作成します。
// command.go で組み立て済みのDSNを受け取ります。
func NewPgDatasource(driverName string, dsn string) Datasource {
	return &pgDatasource{
		driverName: driverName,
		dsn:        dsn,
	}
}

// DSN は PostgreSQL のデータソース名 (DSN) を返します。
func (d *pgDatasource) DSN() string {
	return d.dsn
}

// Driver は PostgreSQL のドライバ名を返します。
func (d *pgDatasource) Driver() string {
	return d.driverName
}
