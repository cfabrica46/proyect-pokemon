package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/cfabrica46/proyecto-pokemon/pokedatabases"
)

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
