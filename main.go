package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/cfabrica46/proyecto-pokemon/open"
	"github.com/cfabrica46/proyecto-pokemon/pokedatabases"

	_ "github.com/mattn/go-sqlite3"
)

var (
	errUserExist       = errors.New("El username que usted escogió ya esta en uso")
	errPasswordInvalid = errors.New("Password Invalido")
	errNotPokemons     = errors.New("No tiene pokemons para acceder")
	errIncorrectID     = errors.New("ID seleccionada Incorrecta")
)

const (
	allUserPokemons = iota
	onlyPokeFromUser
	onlyPokeFromRival
	allPokemons
)

type usuario struct {
	id       int
	username string
	password string
	pokemons []pokemon
}

type pokemon struct {
	id      int
	name    string
	life    int
	tipo    string
	level   int
	ataques []ataque
}

type ataque struct {
	id       int
	name     string
	power    int
	accuracy int
}

func main() {

	var usernameScan, passwordScan, ingreso string

	databases, err := open.Open()

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
			log.Fatal(err)
		}

		err = pokedatabases.FuncIngreso(databases, *user)

		if err != nil {
			log.Fatal(err)
		}

	case "r":
		user, err := pokedatabases.FuncRegistro(databases)

		if err != nil {
			log.Fatal(err)
		}

		err = pokedatabases.FuncIngreso(databases, *user)

		if err != nil {
			log.Fatal(err)
		}

	default:
		log.Fatal("Error: ELECCIÓN INVALIDA")
	}

}
