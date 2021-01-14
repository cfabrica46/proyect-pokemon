package open

import (
	"database/sql"
	"io/ioutil"
	"os"
	"strings"

	//necesitamos importar para tener el driver
	_ "github.com/mattn/go-sqlite3"
)

//Migracion es una funcion complementaria de Open
//Al ejecutarla se migraran los datos del archivo .sql a un archivo .db
func Migracion() (databases *sql.DB, err error) {
	archivoDB, err := os.Create("databases.db")

	if err != nil {
		return
	}
	archivoDB.Close()

	databases, err = sql.Open("sqlite3", "./databases.db?_foreign_keys=on")

	if err != nil {
		return
	}

	archivoSQL, err := os.Open("databases.sql")

	if err != nil {
		return
	}

	defer archivoSQL.Close()

	contendio, err := ioutil.ReadAll(archivoSQL)

	if err != nil {
		return
	}

	_, err = databases.Exec(string(contendio))
	if err != nil {
		return
	}

	return
}

//Open Abrira el archivo .db o en su defecto ejecutara Migracion
func Open() (databases *sql.DB, err error) {

	archivo, err := os.Open("databases.db")

	if err != nil {
		if strings.Contains(err.Error(), "open databases.db: El sistema no puede encontrar el archivo especificado.") {

			databases, err := Migracion()

			if err != nil {
				return databases, err
			}

			return databases, err
		}
		return
	}
	defer archivo.Close()

	databases, err = sql.Open("sqlite3", "./databases.db?_foreign_keys=on")

	if err != nil {
		return
	}

	return
}
