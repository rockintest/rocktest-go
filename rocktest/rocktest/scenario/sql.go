package scenario

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"

	_ "github.com/alexbrainman/odbc"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/stdlib"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/proullon/ramsql/driver"
)

type sqlData struct {
	id *sql.DB
}

func (module *Module) Sql_connect(params map[string]interface{}, scenario *Scenario) error {

	paramsEx, err := scenario.ExpandMap(params)
	if err != nil {
		return err
	}
	driver, err := scenario.GetString(paramsEx, "driver", nil)

	if err != nil {
		return err
	}

	url, err := scenario.GetString(params, "url", nil)

	if err != nil {
		return err
	}

	name, _ := scenario.GetString(paramsEx, "name", "default")

	db, err := sql.Open(driver, url)

	if err != nil {
		return err
	}

	scenario.PutStore("sql."+name, sqlData{id: db})
	scenario.PutCleanup("sql", closeSQL)

	return nil

}

func closeSQL(scenario *Scenario) error {
	log.Info("Cleanup SQL")

	for k, v := range scenario.Store {
		if strings.HasPrefix(k, "sql.") {
			log.Debugf("Closing %s", k)
			data := v.(sqlData)
			data.id.Close()
		}
	}

	return nil
}

func (module *Module) Sql_request(params map[string]interface{}, scenario *Scenario) error {

	paramsEx, err := scenario.ExpandMap(params)
	if err != nil {
		return err
	}

	reqs, err := scenario.GetString(paramsEx, "request", nil)
	if err != nil {
		return err
	}

	as, err := scenario.GetString(paramsEx, "as", "request")
	if err != nil {
		return err
	}

	scenario.DeleteContextRegex(as + "\\..*")

	name, _ := scenario.GetString(paramsEx, "name", "default")
	data := scenario.GetStore("sql." + name).(sqlData)

	reqsArray := strings.Split(reqs, ";")

	for _, req := range reqsArray {

		if req == "" {
			continue
		}

		req2 := strings.ToUpper(req)
		req2 = strings.Trim(req2, " ")

		if !strings.HasPrefix(req2, "SELECT") {

			log.Debugf("Exec request: %s", req)

			_, err = data.id.Exec(req)
			if err != nil {
				return err
			}

			scenario.PutContextAs(paramsEx, "request", "nb", 0)

		} else {
			log.Debugf("Exec query: %s", req)
			rows, err := data.id.Query(req)

			if err != nil {
				return err
			}

			nbRet := 0
			lineRet := ""

			var jsonRet []map[string]string

			cols, _ := rows.Columns()
			for rows.Next() {
				nbRet++
				lineRet = ""
				row := make([]interface{}, len(cols))

				colMap := make(map[string]string)

				jsonRet = append(jsonRet, colMap)

				for i := range cols {
					var s string
					row[i] = &s
				}
				rows.Scan(row...)

				for i, v := range row {

					curr := *(v.(*string))

					lineRet += curr
					if i != len(row)-1 {
						lineRet += ","
					}

					scenario.PutContextAs(paramsEx, "request", fmt.Sprint(i+1), curr)
					scenario.PutContextAs(paramsEx, "request", cols[i], curr)

					colMap[cols[i]] = curr
					colMap[fmt.Sprint(i+1)] = curr
				}

				colMap["fullline"] = lineRet

			}

			rows.Close()
			j, err := json.Marshal(jsonRet)
			if err != nil {
				return nil
			}

			expect, err := scenario.GetList(paramsEx, "expect", nil)

			if err == nil {
				for _, cond := range expect {
					errsearch := module.search(jsonRet, fmt.Sprint(cond))
					if errsearch != nil {
						log.Debugf("Look for regex %s: NO", cond)
						return errsearch
					} else {
						log.Debugf("Look for regex %s: YES", cond)
					}
				}
			}

			scenario.PutContextAs(paramsEx, "request", "json", string(j))

			log.Debugf("LINE=%s", lineRet)

			scenario.PutContextAs(paramsEx, "request", "nb", nbRet)
			scenario.PutContextAs(paramsEx, "request", "0", lineRet)

		}

	}

	return nil

}

func (module *Module) search(info []map[string]string, pregex string) error {

	for _, v := range info {
		l, ok := v["fullline"]
		if ok {
			match, err := regexp.MatchString(pregex, l)
			if err != nil {
				return err
			}
			if match {
				return nil
			}
		}
	}

	return fmt.Errorf("no match found for %s", pregex)
}

func (module *Module) Sql_drivers(params map[string]interface{}, scenario *Scenario) error {
	ret := sql.Drivers()

	paramsEx, err := scenario.ExpandMap(params)
	if err != nil {
		return err
	}

	retJson, _ := json.Marshal(ret)

	scenario.PutContextAs(paramsEx, "drivers", "result", string(retJson))
	scenario.PutContext("??", string(retJson))

	return nil
}
