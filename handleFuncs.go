package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/cfabrica46/proyecto-pokemon/pokedatabases"
)

func funcIngreso(databases *sql.DB, user pokedatabases.User) (err error) {

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
			pokes, err := pokedatabases.SeleccionarPokemons(databases, pokedatabases.AllUserPokemons, user.ID, 0)

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
			err := eliminarCuenta(databases, user, &salir)

			if err != nil {
				fmt.Println(err.Error())
			}
		case 0:
			return
		default:
			fmt.Println("Opcion invalida")
		}
	}
	return
}

func funcRegistro(databases *sql.DB) (user pokedatabases.User, err error) {

	var usernameScan, passwordScan string

	fmt.Println("Ingrese su username")
	fmt.Scan(&usernameScan)
	fmt.Println("Ingrese su password")
	fmt.Scan(&passwordScan)

	err = pokedatabases.CheckRegistro(databases, usernameScan, passwordScan)

	if err != nil {
		return
	}

	err = pokedatabases.InsertRegistro(databases, usernameScan, passwordScan)
	if err != nil {
		return
	}

	user, err = pokedatabases.GetUser(databases, usernameScan, passwordScan)

	if err != nil {
		return
	}
	fmt.Println("nuevo usuario", usernameScan)

	return
}

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

	err = pokedatabases.InsertarNuevoPoke(databases, newPokemon)

	if err != nil {
		return
	}

	err = pokedatabases.InsertarRelacionNuevoPoke(databases, user, newPokemon, idAux)
	if err != nil {
		return
	}

	err = pokedatabases.InsertarRelacionAtaques(databases, newPokemon.Ataques, idAux)

	if err != nil {
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

func eliminarCuenta(databases *sql.DB, user pokedatabases.User, salir *bool) (err error) {

	for {
		var preguntaSeguridad1, preguntaSeguridad2 string

		fmt.Println("¿Esta seguro de eliminar su cuenta? [S/N]")
		fmt.Scan(&preguntaSeguridad1)

		preguntaSeguridad1 = strings.ToLower(preguntaSeguridad1)

		switch preguntaSeguridad1 {
		case "s":
			fmt.Println("Introduzca su password")
			fmt.Scan(&preguntaSeguridad2)

			check, err := pokedatabases.BorrarCuenta(databases, user, salir, preguntaSeguridad2)

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

func jugar(databases *sql.DB, user pokedatabases.User) {

	var idPokeJugador int
	var turno bool

	pokes, err := pokedatabases.SeleccionarPokemons(databases, pokedatabases.AllUserPokemons, user.ID, 0)

	if err != nil {
		return
	}

	imprimirPokemons(pokes)

	fmt.Println("Escribe el ID del Pokemon con el que quieras jugar")

	fmt.Scan(&idPokeJugador)

	sliceAux, err := pokedatabases.SeleccionarPokemons(databases, pokedatabases.OnlyPokeFromUser, user.ID, idPokeJugador)

	if err != nil {
		return
	}

	j1 := sliceAux[0]

	time.Sleep(2 * time.Second)

	fmt.Println("Tu oponente sera...")

	time.Sleep(time.Second * 2)

	posiblesRivales, err := pokedatabases.SeleccionarPokemons(databases, pokedatabases.AllPokemons, 0, 0)

	if err != nil {
		return
	}

	rand.Seed(time.Now().UnixNano())

	idRival := rand.Intn(len(posiblesRivales) + 1)

	pokeRival, err := pokedatabases.SeleccionarPokemons(databases, pokedatabases.OnlyPokeFromRival, 0, idRival+1)

	if err != nil {
		return
	}

	j2 := pokeRival[0]

	fmt.Println(j2.Name)

	time.Sleep(time.Second * 3)

	t := 1

	for {
		var attack1 int

		if turno == false {

			fmt.Printf("Turno %v\n", t)
			fmt.Println("Es tu turno")
			mostrarVidas(j1, j2)
			fmt.Println("Que ataque desar usar?")

			mostrarAtaques(j1)

			fmt.Scan(&attack1)
			ataquej1 := attack1 - 1

			if ataquej1 >= len(j1.Ataques) {

				fmt.Println("Opcion invalida")
				ataquej1 = -1

			}

			mostrarBatalla(&j1, &j2, ataquej1)

		} else {

			fmt.Printf("Turno %v\n", t)
			fmt.Println("Es turno del rival")

			ataquej2 := rand.Intn(len(j2.Ataques))

			mostrarBatalla(&j2, &j1, ataquej2)
		}

		time.Sleep(time.Second * 3)

		t = t + 1

		if j1.Life <= 0 || j2.Life <= 0 {
			break
		}

		turno = !turno
		fmt.Println()
	}

	if j1.Life <= 0 {

		fmt.Println("Oh no, perdiste")

	} else {

		fmt.Println("Genial ganaste esta batalla")

	}
}

func mostrarVidas(j1 pokedatabases.Pokemon, j2 pokedatabases.Pokemon) {

	fmt.Printf("Vida Actual: %v => %v	||	%v => %v\n", j1.Name, j1.Life, j2.Name, j2.Life)

}

func mostrarAtaques(poke pokedatabases.Pokemon) {

	for i := range poke.Ataques {

		fmt.Printf("%v.%v => %v\n", i+1, poke.Ataques[i].Name, poke.Ataques[i].Power)

	}

}

func mostrarBatalla(atacante *pokedatabases.Pokemon, receptor *pokedatabases.Pokemon, eleccion int) {

	if eleccion == -1 {
		fmt.Println("fallaste")
	} else {

		fmt.Printf("%v uso %v contra %v\n", atacante.Name, atacante.Ataques[eleccion].Name, receptor.Name)

		receptor.Life = receptor.Life - atacante.Ataques[eleccion].Power
	}
}
