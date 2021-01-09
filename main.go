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

	rows1, err := databases.Query("SELECT  pokemons.id,pokemons.name,pokemons.type,pokemons.level,attacks.name,attacks.power,attacks.accuracy FROM users_pokemons INNER JOIN users ON users_pokemons.user_id = users.id INNER JOIN pokemons ON users_pokemons.poke_id = pokemons.id INNER JOIN pokemons_attacks ON users_pokemons.poke_id = pokemons_attacks.poke_id INNER JOIN attacks ON pokemons_attacks.attack_id = attacks.id WHERE users.id = ?", user.id)

	if err != nil {
		log.Fatal(err)
	}

	defer rows1.Close()

	var pokeName, tipo, attackName string
	var pokeid, level, power, accuracy int

	for rows1.Next() {

		err = rows1.Scan(&pokeid, &pokeName, &tipo, &level, &attackName, &power, &accuracy)

		if err != nil {
			log.Fatal(err)
		}

		column, err := rows1.ColumnTypes()

		if err != nil {
			log.Fatal(err)
		}

		campos := []interface{}{pokeid, pokeName, tipo, attackName, level, power, accuracy}

		for i := range column {
			fmt.Printf("%v -> %v\t", column[i].Name(), campos[i])
		}
		fmt.Println()
	}
}
