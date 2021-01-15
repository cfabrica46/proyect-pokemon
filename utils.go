package main

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/cfabrica46/proyecto-pokemon/pokedatabases"
)

func imprimirPokemons(pokes []pokedatabases.Pokemon) {

	for i := range pokes {
		fmt.Println(pokes[i].Name)
		fmt.Printf(" ID -> %v\n", pokes[i].ID)
		fmt.Printf(" type -> %v\n", pokes[i].Tipo)
		fmt.Printf(" level -> %v\n", pokes[i].Level)
		fmt.Printf(" ataques:\n")
		for i2 := range pokes[i].Ataques {
			fmt.Printf("\t%v. Attack name -> %v\n", i2+1, pokes[i].Ataques[i2].Name)
			fmt.Printf("\t%v. Attack power-> %v\n", i2+1, pokes[i].Ataques[i2].Power)
			fmt.Printf("\t%v. Attack accuracy -> %v\n", i2+1, pokes[i].Ataques[i2].Accuracy)
			fmt.Println()
		}
	}
}

func añadirPoke(databases *sql.DB, user pokedatabases.User) (err error) {

	var eleccionPoke, idAux int

	fmt.Println("¿Que pokemon Deseas?")

	allPokes, err := pokedatabases.SeleccionarPokemons(databases, pokedatabases.AllPokemons, 0, 0)

	if err != nil {
		return
	}

	imprimirPokemons(allPokes)

	idAux = len(allPokes) + 1

	fmt.Scan(&eleccionPoke)

	pokemonSelected, err := pokedatabases.SeleccionarPokemons(databases, pokedatabases.OnlyPokeFromRival, 0, eleccionPoke)

	if err != nil {
		return
	}

	newPokemon := pokemonSelected[0]

	tx, err := databases.Begin()

	if err != nil {
		return
	}

	err = pokedatabases.InsertarNuevoPoke(tx, newPokemon)

	if err != nil {
		tx.Rollback()
		return
	}

	err = pokedatabases.InsertarRelacionNuevoPoke(tx, user, newPokemon, idAux)
	if err != nil {
		tx.Rollback()
		return
	}

	err = pokedatabases.InsertarRelacionAtaques(tx, newPokemon.Ataques, idAux)

	if err != nil {
		tx.Rollback()
		return
	}

	if tx.Commit() != nil {
		return
	}

	fmt.Println("Ahora el sera tu nuevo amigo :D")
	return
}

func liberarPokemon(databases *sql.DB, user pokedatabases.User) (err error) {

	var pokeEliminar int

	allPokes, err := pokedatabases.SeleccionarPokemons(databases, pokedatabases.AllUserPokemons, user.ID, 0)

	if err != nil {
		return
	}

	imprimirPokemons(allPokes)

	fmt.Println("Escribe el ID del pokemon que deseas eliminar")

	fmt.Scan(&pokeEliminar)

	sliceAux, err := pokedatabases.SeleccionarPokemons(databases, pokedatabases.OnlyPokeFromUser, user.ID, pokeEliminar)

	if err != nil {
		return
	}

	if len(sliceAux) == 0 {
		fmt.Println("id no valido")
		return
	}

	err = pokedatabases.DeletePoke(databases, user, sliceAux)

	if err != nil {
		return
	}

	fmt.Println("Hasta pronto amiguito :'D")
	return
}

func eliminarCuenta(databases *sql.DB, user pokedatabases.User) (err error) {

	for {
		var preguntaSeguridad1, preguntaSeguridad2 string

		fmt.Println("¿Esta seguro de eliminar su cuenta? [S/N]")
		fmt.Scan(&preguntaSeguridad1)

		preguntaSeguridad1 = strings.ToLower(preguntaSeguridad1)

		switch preguntaSeguridad1 {
		case "s":
			fmt.Println("Introduzca su password")
			fmt.Scan(&preguntaSeguridad2)

			check, err := pokedatabases.BorrarCuenta(databases, user.ID, user.Password, preguntaSeguridad2)

			if err != nil {
				return err
			}

			if check == true {
				fmt.Println("Hasta Pronto :'D")
				return err
			}

			return err

		case "n":
			return

		default:
			fmt.Println("Opncion invalida")
		}
	}

}
