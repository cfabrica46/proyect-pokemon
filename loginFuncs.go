package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/cfabrica46/proyecto-pokemon/pokedatabases"
)

func ingresar(databases *sql.DB, user pokedatabases.User) (err error) {

	var eleccionMenu int
	var salir bool

	fmt.Printf("Bienvenido %v tu ID es: %v\n", user.Username, user.ID)

	for salir == false {
		fmt.Println("¿Qué Desea Hacer?")
		fmt.Println("1.	Jugar")
		fmt.Println("2.	Ver tus Pokemones")
		fmt.Println("3. 	Añadir un Pokemon")
		fmt.Println("4.	Liberar un Pokemon")
		fmt.Println("5.	Eliminar tu cuenta")
		fmt.Println("0.	Salir")

		fmt.Scan(&eleccionMenu)

		switch eleccionMenu {
		case 1:
			jugar(databases, user)
		case 2:
			pokes, err := pokedatabases.GetPokemonsFromUser(databases, user.ID)

			if err != nil {
				fmt.Println(err.Error())
			}

			imprimirPokemons(pokes)
		case 3:
			err := añadirPoke(databases, user)

			if err != nil {
				fmt.Println(err.Error())
			}
		case 4:
			err := liberarPokemon(databases, user)

			if err != nil {
				fmt.Println(err.Error())
			}
		case 5:
			err := eliminarCuenta(databases, user)

			if err != nil {
				fmt.Println(err.Error())
			} else {
				salir = true
			}

		case 0:
			return
		default:
			fmt.Println("Opcion invalida")
		}
	}
	return
}

func registrar(databases *sql.DB) (user *pokedatabases.User, err error) {

	var usernameScan, passwordScan string

	fmt.Println("Ingrese su username")
	fmt.Scan(&usernameScan)
	fmt.Println("Ingrese su password")
	fmt.Scan(&passwordScan)

	check, err := pokedatabases.CheckIfUserAlreadyExist(databases, usernameScan)

	if err != nil {
		return
	}

	if check == true {
		err = pokedatabases.ErrUserExist
		return
	}

	err = pokedatabases.InsertUser(databases, usernameScan, passwordScan)
	if err != nil {
		return
	}

	user, err = pokedatabases.GetUser(databases, usernameScan, passwordScan)

	if err != nil {
		if err == sql.ErrNoRows {
			if user == nil {
				log.Fatal(errUsernamePasswordIncorrect)
			}
			log.Fatal(err)
		}
		log.Fatal(err)
	}

	fmt.Println("nuevo usuario", usernameScan)

	return
}
