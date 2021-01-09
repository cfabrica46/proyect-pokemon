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
	var user usuario

	databases, err := sql.Open("sqlite3", "./databases.db")

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

		mostrarPokes(databases, user)

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
