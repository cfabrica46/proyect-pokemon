package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

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

func migracion() (databases *sql.DB, err error) {
	archivoDB, err := os.Create("databases.db")

	if err != nil {
		return
	}
	archivoDB.Close()

	databases, err = sql.Open("sqlite3", "./databases.db?_foreign_keys=on")

	if err != nil {
		return
	}

	archivoSQL, err := os.Open("databases.sql")

	if err != nil {
		return
	}

	defer archivoSQL.Close()

	contendio, err := ioutil.ReadAll(archivoSQL)

	if err != nil {
		return
	}

	_, err = databases.Exec(string(contendio))
	if err != nil {
		return
	}

	return
}

func open() (databases *sql.DB, err error) {

	archivo, err := os.Open("databases.db")

	if err != nil {
		databases, err := migracion()

		if err != nil {
			return databases, err
		}

		return databases, err

	}

	defer archivo.Close()

	databases, err = sql.Open("sqlite3", "./databases.db?_foreign_keys=on")

	if err != nil {
		return
	}

	return
}

func main() {

	log.SetFlags(log.Llongfile)

	var ingreso string

	databases, err := open()

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
		funcIngreso(databases)
	case "r":
		funcRegistro(databases)
	default:
		log.Fatal("Error: ELECCIÓN INVALIDA")
	}

}

func funcIngreso(databases *sql.DB) {

	var usernameScan, passwordScan string
	var eleccionMenu int
	var salir bool
	var user usuario

	fmt.Println("Ingrese su username")
	fmt.Scan(&usernameScan)
	fmt.Println("Ingrese su password")
	fmt.Scan(&passwordScan)

	user, err := checkIngreso(databases, usernameScan, passwordScan)

	if err != nil {
		log.Fatal(err)
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
			jugar(databases, user)
		case 2:
			mostrarPokes(databases, user)
		case 3:
			añadirPoke(databases, user)
		case 4:
			liberarPokemon(databases, user)
		case 5:
			eliminarCuenta(databases, user, &salir)
		case 0:
			return
		default:
			fmt.Println("Opcion invalida")
		}
	}
}

func checkIngreso(databases *sql.DB, usernameScan, passwordScan string) (user usuario, err error) {

	rows, err := databases.Query("SELECT id,username,password FROM users WHERE username = ? AND password = ?", usernameScan, passwordScan)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {

		err = rows.Scan(&user.id, &user.username, &user.password)

		if err != nil {
			return
		}

	}
	if user.id == 0 {
		log.Fatal("Error: username y/o password invalidos")
	}

	return

}

func funcRegistro(databases *sql.DB) {

	var usernameScan, passwordScan string

	fmt.Println("Ingrese su username")
	fmt.Scan(&usernameScan)
	fmt.Println("Ingrese su password")
	fmt.Scan(&passwordScan)

	checkRegistro(databases, usernameScan, passwordScan)

}

func checkRegistro(databases *sql.DB, usernameScan, passwordScan string) {

	stmt, err := databases.Prepare("INSERT INTO users(username, password) VALUES(?,?)")

	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(usernameScan, passwordScan)

	if err != nil {
		log.Fatal(err, "El username que usted escogió ya esta en uso")
	}

	fmt.Println("nuevo usuario", usernameScan)

}

func mostrarPokes(databases *sql.DB, user usuario) {
	var newPokemon pokemon

	fmt.Println("Esta es tu lista de pokemons")

	rows, err := databases.Query("SELECT DISTINCT pokemons.id,pokemons.name,pokemons.life,pokemons.type,pokemons.level FROM users_pokemons INNER JOIN users ON users_pokemons.user_id = users.id INNER JOIN pokemons ON users_pokemons.poke_id = pokemons.id WHERE users.id = ?", user.id)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {

		err = rows.Scan(&newPokemon.id, &newPokemon.name, &newPokemon.life, &newPokemon.tipo, &newPokemon.level)

		if err != nil {
			log.Fatal(err)
		}

		user.pokemons = append(user.pokemons, newPokemon)

		fmt.Printf("%v -> %v\t", "id", newPokemon.id)
		fmt.Printf("%v -> %v\t", "name", newPokemon.name)
		fmt.Printf("%v -> %v\t", "life", newPokemon.life)
		fmt.Printf("%v -> %v\t", "type", newPokemon.tipo)
		fmt.Printf("%v -> %v\t", "level", newPokemon.level)

		mostrarPokesAux(databases, newPokemon)

	}
	if newPokemon.id == 0 {
		fmt.Println("No tienes ningun pokemon")
	}

}

func mostrarPokesAux(databases *sql.DB, newPokemon pokemon) {

	rows, err := databases.Query("SELECT DISTINCT attacks.id,attacks.name,attacks.power,attacks.accuracy FROM pokemons_attacks INNER JOIN attacks ON pokemons_attacks.attack_id=attacks.id WHERE pokemons_attacks.poke_id=?", newPokemon.id)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {

		var newAttack ataque

		err = rows.Scan(&newAttack.id, &newAttack.name, &newAttack.power, &newAttack.accuracy)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%v -> %v\t", "Attack Name", newAttack.name)
		fmt.Printf("%v -> %v\t", "Attack Power", newAttack.power)
		fmt.Printf("%v -> %v\t", "Attack Accuracy", newAttack.accuracy)

	}

	fmt.Println()
}

func añadirPoke(databases *sql.DB, user usuario) {

	var eleccionPoke, idAux int
	var sliceAux []int
	var newPokemon pokemon

	fmt.Println("¿Que pokemon Deseas?")

	mostrarTodosPokes(databases, &newPokemon)

	idAux = newPokemon.id + 1

	fmt.Scan(&eleccionPoke)

	añadirSeleccionarPoke(databases, &newPokemon, eleccionPoke)

	insertarNuevoPoke(databases, newPokemon)

	insertarRelacionNuevoPoke(databases, user, newPokemon, idAux)

	añadirSeleccionAtaque(databases, &sliceAux, eleccionPoke)

	insertarRelacionAtaques(databases, sliceAux, idAux)

	fmt.Println("Ahora el sera tu nuevo amigo :D")
}

func mostrarTodosPokes(databases *sql.DB, newPokemon *pokemon) {

	rows, err := databases.Query("SELECT DISTINCT pokemons.id,pokemons.name,pokemons.life,pokemons.type,pokemons.level FROM pokemons")

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {

		err = rows.Scan(&newPokemon.id, &newPokemon.name, &newPokemon.life, &newPokemon.tipo, &newPokemon.level)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%v -> %v\t", "id", newPokemon.id)
		fmt.Printf("%v -> %v\t", "name", newPokemon.name)
		fmt.Printf("%v -> %v\t", "life", newPokemon.life)
		fmt.Printf("%v -> %v\t", "type", newPokemon.tipo)
		fmt.Printf("%v -> %v\t", "level", newPokemon.level)

		fmt.Println()
	}
}

func añadirSeleccionarPoke(databases *sql.DB, newPokemon *pokemon, eleccionPoke int) {

	rows, err := databases.Query("SELECT DISTINCT pokemons.name,pokemons.life,pokemons.type,pokemons.level FROM pokemons WHERE pokemons.id=?", eleccionPoke)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {

		rows.Scan(&newPokemon.name, &newPokemon.life, &newPokemon.tipo, &newPokemon.level)

	}

}

func insertarNuevoPoke(databases *sql.DB, newPokemon pokemon) {

	stmt, err := databases.Prepare("INSERT INTO pokemons(name,life,type,level) VALUES (?,?,?,?)")

	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(newPokemon.name, newPokemon.life, newPokemon.tipo, newPokemon.level)

	if err != nil {
		log.Fatal(err)
	}

}

func insertarRelacionNuevoPoke(databases *sql.DB, user usuario, newPokemon pokemon, idAux int) {
	stmt, err := databases.Prepare("INSERT INTO users_pokemons (user_id,poke_id) VALUES (?,?)")

	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(user.id, idAux)

	if err != nil {
		log.Fatal(err)
	}
}

func añadirSeleccionAtaque(databases *sql.DB, sliceAux *[]int, eleccionPoke int) {
	var idAttack int

	rows, err := databases.Query("SELECT pokemons_attacks.attack_id FROM pokemons_attacks WHERE pokemons_attacks.poke_id=?", eleccionPoke)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {

		rows.Scan(&idAttack)

		*sliceAux = append(*sliceAux, idAttack)

	}
}

func insertarRelacionAtaques(databases *sql.DB, sliceAux []int, idAux int) {

	for i := range sliceAux {

		stmtAttack, err := databases.Prepare("INSERT INTO pokemons_attacks (poke_id,attack_id) VALUES (?,?)")

		if err != nil {
			log.Fatal(err)
		}

		_, err = stmtAttack.Exec(idAux, sliceAux[i])

		if err != nil {
			log.Fatal(err)
		}

	}
}

func liberarPokemon(databases *sql.DB, user usuario) {

	var pokeEliminar, idAux int

	check := checkPokemons(databases, user)

	if check == false {
		fmt.Println("No tienes ningun Pokemon para liberar")
		return
	}

	mostrarPokes(databases, user)

	fmt.Println("Escribe el ID del pokemon que deseas eliminar")

	fmt.Scan(&pokeEliminar)

	sliceAux := seleccionarPokeEliminar(databases, user, pokeEliminar, &idAux)

	if len(sliceAux) == 0 {
		fmt.Println("id no valido")
	}

	for i := range sliceAux {

		stmtDelete, err := databases.Prepare("DELETE FROM users_pokemons WHERE users_pokemons.poke_id = ?;")

		if err != nil {
			log.Fatal(err)
		}

		_, err = stmtDelete.Exec(sliceAux[i])

		if err != nil {
			log.Fatal(err)
		}

	}

	fmt.Println("Hasta pronto amiguito :'D")
}

func seleccionarPokeEliminar(databases *sql.DB, user usuario, pokeEliminar int, idAux *int) (sliceAux []int) {

	rows, err := databases.Query("SELECT DISTINCT pokemons.id FROM users_pokemons INNER JOIN users ON users_pokemons.user_id = users.id INNER JOIN pokemons ON users_pokemons.poke_id = pokemons.id WHERE pokemons.id = ? AND users.id = ?", pokeEliminar, user.id)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {

		err = rows.Scan(&idAux)

		if err != nil {
			log.Fatal(err)
		}

		sliceAux = append(sliceAux, *idAux)

	}
	return
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

			}

			log.Fatal("Error: password invalido")

		case "n":
			return

		default:
			fmt.Println("Opncion invalida")
		}
	}
}

func jugar(databases *sql.DB, user usuario) {

	var idPokeJugador int
	var turno bool

	check := checkPokemons(databases, user)

	if check == false {
		fmt.Println("No tienes ningun Pokemon para jugar")
		return
	}

	mostrarPokes(databases, user)

	fmt.Println("Escribe el ID del Pokemon con el que quieras jugar")

	fmt.Scan(&idPokeJugador)

	j1 := seleccionarPokemon(databases, user, idPokeJugador)

	if j1.id == 0 {

		fmt.Println("Seleccion Invalida")
		return

	}

	time.Sleep(2 * time.Second)

	fmt.Println("Tu oponente sera...")

	time.Sleep(time.Second * 2)

	lenPoke := seleccionTodosPokes(databases)

	rand.Seed(time.Now().UnixNano())

	idRival := rand.Intn(lenPoke)

	rival := seleccionarRival(databases, idRival+1)

	fmt.Println(rival.name)

	time.Sleep(time.Second * 3)

	t := 1

	for {
		var attack1 int

		if turno == false {

			fmt.Printf("Turno %v\n", t)
			fmt.Println("Es tu turno")
			mostrarVidas(j1, rival)
			fmt.Println("Que ataque desar usar?")

			mostrarAtaques(j1)

			fmt.Scan(&attack1)
			ataquej1 := attack1 - 1

			if ataquej1 >= len(j1.ataques) {

				fmt.Println("Opcion invalida")
				ataquej1 = -1

			}

			mostrarBatalla(j1, rival, ataquej1)

		} else {

			fmt.Printf("Turno %v\n", t)
			fmt.Println("Es turno del rival")

			ataquej2 := rand.Intn(len(rival.ataques))

			mostrarBatalla(rival, j1, ataquej2)
		}

		time.Sleep(time.Second * 3)

		t = t + 1

		if j1.life <= 0 || rival.life <= 0 {
			break
		}

		turno = !turno
	}

	if j1.life <= 0 {

		fmt.Println("Oh no, perdiste")

	} else {

		fmt.Println("Genial ganaste esta batalla")

	}
}

func seleccionarPokemon(databases *sql.DB, user usuario, idPoke int) *pokemon {

	var newPokemon pokemon
	rows, err := databases.Query("SELECT DISTINCT pokemons.id,pokemons.name,pokemons.life,pokemons.type,pokemons.level FROM users_pokemons INNER JOIN pokemons ON users_pokemons.poke_id = pokemons.id WHERE pokemons.id = ? AND users_pokemons.user_id = ?", idPoke, user.id)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {

		err = rows.Scan(&newPokemon.id, &newPokemon.name, &newPokemon.life, &newPokemon.tipo, &newPokemon.level)

		if err != nil {
			log.Fatal(err)
		}

		seleccionarAtaques(databases, idPoke, &newPokemon)

		fmt.Println()
	}
	return &newPokemon
}

func seleccionarAtaques(databases *sql.DB, idPoke int, newPokemon *pokemon) {
	var newAttack ataque

	rows, err := databases.Query("SELECT DISTINCT attacks.id,attacks.name,attacks.power,attacks.accuracy FROM pokemons_attacks INNER JOIN attacks ON pokemons_attacks.attack_id=attacks.id WHERE pokemons_attacks.poke_id=?", idPoke)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {

		err = rows.Scan(&newAttack.id, &newAttack.name, &newAttack.power, &newAttack.accuracy)

		if err != nil {
			log.Fatal(err)
		}

		newPokemon.ataques = append(newPokemon.ataques, newAttack)

	}
}

func seleccionTodosPokes(databases *sql.DB) (lenPoke int) {

	rows, err := databases.Query("SELECT id FROM pokemons")

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {

		err = rows.Scan(&lenPoke)

		if err != nil {
			log.Fatal(err)
		}

	}

	return
}

func seleccionarRival(databases *sql.DB, idPoke int) *pokemon {

	var pokeRival pokemon

	rowsPoke, err := databases.Query("SELECT DISTINCT pokemons.id,pokemons.name,pokemons.life,pokemons.level FROM pokemons WHERE id=? ", idPoke)

	if err != nil {
		log.Fatal(err)
	}

	for rowsPoke.Next() {

		rowsPoke.Scan(&pokeRival.id, &pokeRival.name, &pokeRival.life, &pokeRival.level)

	}

	seleccionarRivalAtaques(databases, &pokeRival, idPoke)

	return &pokeRival
}

func seleccionarRivalAtaques(databases *sql.DB, pokeRival *pokemon, idPoke int) {

	var attack ataque

	rows, err := databases.Query("SELECT DISTINCT attacks.id,attacks.name,attacks.power FROM pokemons_attacks INNER JOIN attacks ON pokemons_attacks.attack_id=attacks.id WHERE pokemons_attacks.poke_id=? ", idPoke)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {

		err = rows.Scan(&attack.id, &attack.name, &attack.power)

		if err != nil {
			log.Fatal(err)
		}

		pokeRival.ataques = append(pokeRival.ataques, attack)
	}
}

func mostrarVidas(j1 *pokemon, j2 *pokemon) {

	fmt.Printf("Vida Actual: %v => %v	||	%v => %v\n", j1.name, j1.life, j2.name, j2.life)

}

func mostrarAtaques(poke *pokemon) {

	for i := range poke.ataques {

		fmt.Printf("%v.%v => %v\n", i+1, poke.ataques[i].name, poke.ataques[i].power)

	}

}

func mostrarBatalla(atacante *pokemon, receptor *pokemon, eleccion int) {

	if eleccion == -1 {
		fmt.Println("fallaste")
	} else {

		fmt.Printf("%v uso %v contra %v\n", atacante.name, atacante.ataques[eleccion].name, receptor.name)

		receptor.life = receptor.life - atacante.ataques[eleccion].power
	}
}

func checkPokemons(databases *sql.DB, user usuario) (check bool) {

	var idAux int

	rows, err := databases.Query("SELECT users_pokemons.poke_id FROM users_pokemons WHERE users_pokemons.user_id = ?", user.id)

	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		err = rows.Scan(&idAux)

		if err != nil {
			log.Fatal(err)
		}
	}

	if idAux == 0 {
		check = false
	} else {
		check = true
	}

	return
}
