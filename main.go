package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
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

	log.SetFlags(log.Llongfile)

	var ingreso, usernameScan, passwordScan string
	var eleccionMenu int
	var salir bool
	var user usuario

	databases, err := sql.Open("sqlite3", "./databases.db?_foreign_keys=on")

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

		rows, err := databases.Query("SELECT id,username,password FROM users WHERE username = ? AND password = ?", usernameScan, passwordScan)

		if err != nil {
			log.Fatal(err)
		}

		defer rows.Close()

		for rows.Next() {

			err = rows.Scan(&user.id, &user.username, &user.password)

			if err != nil {
				log.Fatal(err)
			}

		}
		if user.id == 0 {
			log.Fatal("Error: username y/o password invalidos")
		}

		fmt.Printf("Bienvenido %v tu ID es: %v\n", user.username, user.id)

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
				fmt.Println("juga")
			case 2:
				mostrarPokes(databases, user)
			case 3:
				fmt.Println("otro uwu")
			case 4:
				fmt.Println("liberar :v")
			case 5:
				eliminarCuenta(databases, user, &salir)

			case 0:
				return
			default:
				fmt.Println("Opcion invalida")
			}
		}
	case "r":

		fmt.Println("Ingrese su username")
		fmt.Scan(&usernameScan)
		fmt.Println("Ingrese su password")
		fmt.Scan(&passwordScan)

		stmt, err := databases.Prepare("INSERT INTO users(username, password) VALUES(?,?)")

		if err != nil {
			log.Fatal(err)
		}
		_, err = stmt.Exec(usernameScan, passwordScan)

		if err != nil {
			log.Fatal(err, "El username que usted escogió ya esta en uso")
		}

		fmt.Println("nuevo usuario", usernameScan)

	default:
		log.Fatal("Error: ELECCIÓN INVALIDA")
	}

}

func mostrarPokes(databases *sql.DB, user usuario) {

	fmt.Println("Esta es tu lista de pokemons")

	rows1, err := databases.Query("SELECT DISTINCT pokemons.id,pokemons.name,pokemons.type,pokemons.level FROM users_pokemons INNER JOIN users ON users_pokemons.user_id = users.id INNER JOIN pokemons ON users_pokemons.poke_id = pokemons.id WHERE users.id = ?", user.id)

	if err != nil {
		log.Fatal(err)
	}

	defer rows1.Close()

	for rows1.Next() {

		var newPokemon pokemon

		err = rows1.Scan(&newPokemon.id, &newPokemon.name, &newPokemon.tipo, &newPokemon.level)

		if err != nil {
			log.Fatal(err)
		}

		user.pokemons = append(user.pokemons, newPokemon)

		fmt.Printf("%v -> %v\t", "id", newPokemon.id)
		fmt.Printf("%v -> %v\t", "name", newPokemon.name)
		fmt.Printf("%v -> %v\t", "type", newPokemon.tipo)
		fmt.Printf("%v -> %v\t", "level", newPokemon.level)

		rows2, err := databases.Query("SELECT DISTINCT attacks.id,attacks.name,attacks.power,attacks.accuracy FROM pokemons_attacks INNER JOIN attacks ON pokemons_attacks.attack_id=attacks.id WHERE pokemons_attacks.poke_id=?", newPokemon.id)

		if err != nil {
			log.Fatal(err)
		}

		defer rows2.Close()

		for rows2.Next() {

			var newAttack ataque

			err = rows2.Scan(&newAttack.id, &newAttack.name, &newAttack.power, &newAttack.accuracy)

			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("%v -> %v\t", "Attack Name", newAttack.name)
			fmt.Printf("%v -> %v\t", "Attack Type", newAttack.power)
			fmt.Printf("%v -> %v\t", "Attack Level", newAttack.accuracy)

		}

		fmt.Println()
	}

}

func eliminarCuenta(databases *sql.DB, user usuario, salir *bool) {
	for {
		var preguntaSeguridad1, preguntaSeguridad2 string

		fmt.Println("¿Esta seguro de eliminar su cuenta? [S/N]")
		fmt.Scan(&preguntaSeguridad1)

		preguntaSeguridad1 = strings.ToLower(preguntaSeguridad1)

		switch preguntaSeguridad1 {
		case "s":
			fmt.Println("Introduzca su password")
			fmt.Scan(&preguntaSeguridad2)

			if preguntaSeguridad2 == user.password {

				stmtDelete, err := databases.Prepare("DELETE FROM users where id = ?")

				if err != nil {
					log.Fatal(err)
				}

				_, err = stmtDelete.Exec(user.id)

				if err != nil {
					log.Fatal(err)
				}

				fmt.Println("Hasta Pronto :'D")
				*salir = true
				return

			} else {

				log.Fatal("Error: password invalido")

			}
		case "n":
			return

		default:
			fmt.Println("Opncion invalida")
		}
	}
}

func liberarPokemon() {

}
