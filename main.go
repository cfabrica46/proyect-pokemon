package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	log.SetFlags(log.Llongfile)

	var id, edad int
	var ingreso, usernameScan, passwordScan, username, password string
	var acceso bool

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

		rows, err := databases.Query("SELECT id,username , password FROM users")

		if err != nil {
			log.Fatal(err)
		}

		for rows.Next() {

			err = rows.Scan(&id, &username, &password)

			if err != nil {
				log.Fatal(err)
			}

			if username == usernameScan && password == passwordScan {

				acceso = true

				break

			}

		}

		if acceso == true {
			fmt.Println("Bienvenido", username)

			queryAux := "SELECT  pokemons.id,pokemons.name,pokemons.type,pokemons.level,attacks.name,attacks.power,attacks.accuracy FROM users_pokemons INNER JOIN users ON users_pokemons.user_id = users.id INNER JOIN pokemons ON users_pokemons.poke_id = pokemons.id INNER JOIN pokemons_attacks ON users_pokemons.poke_id = pokemons_attacks.poke_id INNER JOIN attacks ON pokemons_attacks.attack_id = attacks.id WHERE users.id = " + strconv.Itoa(id)

			rows1, err := databases.Query(queryAux)

			if err != nil {
				log.Fatal(err)
			}

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
				fmt.Println("")
			}
		} else {
			log.Fatal("Error: username y/o password invalidos")
		}

		rows.Close()
	case "r":
		var intAux int

		fmt.Println("Ingrese su username")
		fmt.Scan(&usernameScan)
		fmt.Println("Ingrese su password")
		fmt.Scan(&passwordScan)
		fmt.Println("Ingrese su edad")
		fmt.Scan(&edad)

		stringSelect := "SELECT username FROM users WHERE username GLOB '" + usernameScan + "*'"

		rows, err := databases.Query(stringSelect)

		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("check")
		for rows.Next() {

			err = rows.Scan(&username)

			if err != nil {
				log.Fatal(err)
			}

			if usernameScan == username {

				intAux++

				usernameScan = usernameScan + strconv.Itoa(intAux)

				fmt.Println("El username que escogio ya fue no esta disponible")

				fmt.Println("Se le asignara el siguiente username ", usernameScan)
			}

		}

		stmt, err := databases.Prepare("INSERT INTO users(username, password,age) VALUES(?,?,?)")
		fmt.Println("check")

		if err != nil {
			log.Fatal(err)
		}
		_, err = stmt.Exec(usernameScan, passwordScan, edad)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("nuevo usuario", usernameScan)
		break
	default:
		log.Fatal("Error: ELECCIÓN INVALIDA")
	}

}
