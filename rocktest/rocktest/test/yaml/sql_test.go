package yamlTest

import (
	"testing"
)

func TestSQLBase(t *testing.T) {
	shouldPass(t, "sql/sql.yaml")
}

func TestSQLBasic(t *testing.T) {
	shouldPass(t, "sql/basic.yaml")
}

func TestSQLBadDriver(t *testing.T) {
	shouldFail(t, "sql/badDriver.yaml")
}

func TestSQLBadRequest(t *testing.T) {
	shouldFail(t, "sql/badRequest.yaml")
}

func TestSQLCheckFail(t *testing.T) {
	shouldFail(t, "sql/checkFail.yaml")
}

func TestSQLMulti(t *testing.T) {
	shouldPass(t, "sql/sqlMulti.yaml")
}

func TestSQLlite(t *testing.T) {
	shouldPass(t, "sql/sqlite.yaml")
}

// This test needs MySQL to be installed on localhost
// Create a user 'rocktest' with all privileges
// Create a database 'rocktest' (use sudo mysql to connect as root on mysqsl/mariadb)
// Here is the SQL script to create the database :
//
// CREATE USER 'admin'@'localhost' IDENTIFIED BY 'some_pass';
// GRANT ALL PRIVILEGES ON *.* TO 'admin'@'localhost' WITH GRANT OPTION;
// FLUSH PRIVILEGES;
// CREATE DATABASE rocktest

func TestSQLMySQL(t *testing.T) {
	shouldPass(t, "sql/mysql.yaml")
}

// Download appropriate ODBC driver at:
// https://mariadb.com/downloads/#connectors
// Install it (copy lib64/mariadb dir in /usr/lib64)
// Install UnixODBC and configure the datasource
// https://mariadb.com/kb/en/about-mariadb-connector-odbc/
// Needs unixodbc-dev package to be installed

func TestSQLodbc(t *testing.T) {
	shouldPass(t, "sql/odbc.yaml")
}

/*
Setup postgres : sudo apt install postgresql
$ sudo  -i -u postgres
# psql
postgres=# CREATE ROLE rocktest LOGIN;
CREATE ROLE
postgres=# ALTER ROLE rocktest CREATEDB;
ALTER ROLE
postgres=# CREATE DATABASE rocktest OWNER rocktest;
CREATE DATABASE
postgres=# ALTER ROLE rocktest WITH ENCRYPTED PASSWORD 'rocktest';
ALTER ROLE
postgres=# \q
# psql rocktest
*/

func TestSQLpostgres(t *testing.T) {
	shouldPass(t, "sql/postgres.yaml")
}
