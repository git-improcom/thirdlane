/*
 * Copyright 2019. Improcom Inc
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 *  IMPROCOM, Inc
 *  1803 Gravesend Neck road
 *  Brooklyn, NY, 11229 USA
 *  Phone: +1-718-769-3400
 *  @author Alex Romanchuk - alexr@improcom.com - Improcom INC
 */

//Package thdbutil implements mysql database functions. Should be abandoned in favour of API
package thdbutil

import (
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strconv"
)

type DirectoryEntrySql struct {
	TenantID     int32
	Firstname    sql.NullString
	Lastname     sql.NullString
	Email        sql.NullString
	Extension    sql.NullString
	Phone_office sql.NullString
	Phone_home   sql.NullString
	Phone_mobile sql.NullString
	Phone_other  sql.NullString
}

type DirectoryEntry struct {
	TenantID     int32
	Firstname    string
	Lastname     string
	Email        string
	Extension    string
	Phone_office string
	Phone_home   string
	Phone_mobile string
	Phone_other  string
}

func GetAllEntries(err error) {

}

var (
	vipe, vipe_dev                                *viper.Viper
	db_host, db_name, db_user, db_pass, DB_STRING string
	db_port                                       int
	Db                                            *sql.DB
)

func init() {
	vipe = viper.New()
	vipe_dev = viper.New()
	vipe.SetConfigName("config")
	vipe.SetConfigType("yaml")
	vipe.AddConfigPath(".")
	dev_env, _ := strconv.ParseBool(os.Getenv("DEV_ENVIRONMENT"))
	if !dev_env {
		vipe.AddConfigPath("/usr/local/utils/thirdlane_modules/thirdlane_directory_generator")
	}
	err := vipe.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if dev_env {
		vipe_dev.SetConfigName(".env")
		vipe_dev.SetConfigType("yaml")
		vipe_dev.AddConfigPath(".")
		err := vipe_dev.ReadInConfig()
		if err != nil {
			panic(err)
		}
	}

	if dev_env {
		db_host = vipe_dev.GetString("db_host")
		db_port = vipe_dev.GetInt("db_port")
		db_name = vipe_dev.GetString("db_name")
		db_user = vipe_dev.GetString("db_user")
		db_pass = vipe_dev.GetString("db_pass")
	} else {
		db_host = vipe.GetString("db_host")
		db_port = vipe.GetInt("db_port")
		db_name = vipe.GetString("db_name")
		db_user = vipe.GetString("db_user")
		db_pass = vipe.GetString("db_pass")
	}

	DB_STRING = fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&allowOldPasswords=1&parseTime=true&loc=%v", db_user, db_pass, db_host, db_port, db_name, "Local")

}

func GetDirectoryEntries(tenantname string) (dirEntries []DirectoryEntry, err error) {
	Db, err = sql.Open("mysql", DB_STRING)
	if err != nil {
		panic(err)
	}
	defer Db.Close()
	queryString := `
select 
d.firstname,
d.lastname,
d.email,
d.ext,
d.office,
d.home,
d.mobile,
d.other
FROM directory d JOIN tenants t ON d.tenantid = t.id 
WHERE t.tenant=?
AND owner=''
`
	stmt, err := Db.Prepare(queryString)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(tenantname)

	for rows.Next() {
		var dirEntrySql DirectoryEntrySql
		var dirEntry DirectoryEntry
		err = rows.Scan(&dirEntrySql.Firstname, &dirEntrySql.Lastname, &dirEntrySql.Email, &dirEntrySql.Extension, &dirEntrySql.Phone_office, &dirEntrySql.Phone_home, &dirEntrySql.Phone_mobile, &dirEntrySql.Phone_other)
		if err != nil {
			panic(err)
		}

		dirEntry.TenantID = dirEntrySql.TenantID
		dirEntry.Firstname = dirEntrySql.Firstname.String
		dirEntry.Lastname = dirEntrySql.Lastname.String
		dirEntry.Email = dirEntrySql.Email.String
		dirEntry.Extension = dirEntrySql.Extension.String
		dirEntry.Phone_office = dirEntrySql.Phone_office.String
		dirEntry.Phone_home = dirEntrySql.Phone_home.String
		dirEntry.Phone_mobile = dirEntrySql.Phone_mobile.String
		dirEntry.Phone_other = dirEntrySql.Phone_other.String

		dirEntries = append(dirEntries, dirEntry)
	}

	return dirEntries, nil
}
