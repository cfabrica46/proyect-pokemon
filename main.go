package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/cfabrica46/proyecto-pokemon/pokedatabases"

	_ "github.com/mattn/go-sqlite3"
)

var (
	errUsernamePasswordIncorrect = errors.New("Username y/o Password incorrectos")
)

func main() {

	log.SetFlags(log.Llongfile)
	var usernameScan, passwordScan, ingreso string

	databases, err := pokedatabases.Open()

	if err != nil {
		log.Fatal(err)
	}

	defer databases.Close()

	fmt.Println("Bienvenido")

	fmt.Println("Ingresar || Registrarse [I/R]")

	fmt.Scan(&ingreso)

	ingreso = strings.ToLower(ingreso)

	switch ingreso {
	case "i":

		fmt.Println("Ingrese su username")
		fmt.Scan(&usernameScan)
		fmt.Println("Ingrese su password")
		fmt.Scan(&passwordScan)

		user, err := pokedatabases.GetUser(databases, usernameScan, passwordScan)

		if err != nil {
			if err == sql.ErrNoRows && user == nil {
				log.Fatal(errUsernamePasswordIncorrect)
			}
			log.Fatal(err)
		}

		err = ingresar(databases, *user)

		if err != nil {
			log.Fatal(err)
		}

	case "r":
		user, err := registrar(databases)

		if err != nil {
			log.Fatal(err)
		}

		err = ingresar(databases, *user)

		if err != nil {
			log.Fatal(err)
		}

	default:
		log.Fatal("Error: ELECCIÓN INVALIDA")
	}

}
