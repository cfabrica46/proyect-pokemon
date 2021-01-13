package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

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
		if strings.Contains(err.Error(), "open databases.db: El sistema no puede encontrar el archivo especificado.") {

			databases, err := migracion()

			if err != nil {
				return databases, err
			}

			return databases, err
		}
		return
	}
	defer archivo.Close()

	databases, err = sql.Open("sqlite3", "./databases.db?_foreign_keys=on")

	if err != nil {
		return
	}

	return
}

func main() {

	var usernameScan, passwordScan, ingreso string

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

		fmt.Println("Ingrese su username")
		fmt.Scan(&usernameScan)
		fmt.Println("Ingrese su password")
		fmt.Scan(&passwordScan)

		user, err := getUser(databases, usernameScan, passwordScan)

		if err != nil {
			log.Fatal(err)
		}

		err = funcIngreso(databases, *user)

		if err != nil {
			log.Fatal(err)
		}

	case "r":
		user, err := funcRegistro(databases)

		if err != nil {
			log.Fatal(err)
		}

		err = funcIngreso(databases, *user)

		if err != nil {
			log.Fatal(err)
		}

	default:
		log.Fatal("Error: ELECCIÓN INVALIDA")
	}

}

func funcIngreso(databases *sql.DB, user usuario) (err error) {

	var eleccionMenu int
	var salir bool

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
			pokes, err := seleccionarPokemons(databases, allUserPokemons, user.id, 0)

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

func getUser(databases *sql.DB, usernameScan, passwordScan string) (user *usuario, err error) {

	var userAux usuario

	row := databases.QueryRow("SELECT id,username,password FROM users WHERE username = ? AND password = ?", usernameScan, passwordScan)

	err = row.Scan(&userAux.id, &userAux.username, &userAux.password)

	if err != nil {
		return
	}

	user = &usuario{userAux.id, userAux.username, userAux.password, userAux.pokemons}

	return

}

func funcRegistro(databases *sql.DB) (user *usuario, err error) {

	var usernameScan, passwordScan string

	fmt.Println("Ingrese su username")
	fmt.Scan(&usernameScan)
	fmt.Println("Ingrese su password")
	fmt.Scan(&passwordScan)

	err = checkRegistro(databases, usernameScan, passwordScan)

	if err != nil {
		return
	}

	user, err = getUser(databases, usernameScan, passwordScan)

	if err != nil {
		return
	}

	return
}

func checkRegistro(databases *sql.DB, usernameScan, passwordScan string) (err error) {

	stmt, err := databases.Prepare("INSERT INTO users(username, password) VALUES(?,?)")

	if err != nil {
		return
	}
	_, err = stmt.Exec(usernameScan, passwordScan)

	if err != nil {
		err = errUserExist
		return
	}

	fmt.Println("nuevo usuario", usernameScan)

	return
}

func seleccionarPokemons(databases *sql.DB, flag int, idUser, idPoke int) (pokes []pokemon, err error) {

	var check bool

	switch flag {
	case allUserPokemons:

		rows, err := databases.Query("SELECT DISTINCT pokemons.id,pokemons.name,pokemons.life,pokemons.type,pokemons.level FROM users_pokemons INNER JOIN users ON users_pokemons.user_id = users.id INNER JOIN pokemons ON users_pokemons.poke_id = pokemons.id WHERE users.id = ?", idUser)

		if err != nil {
			return pokes, err
		}

		defer rows.Close()

		for rows.Next() {
			var newPokemon pokemon

			err = rows.Scan(&newPokemon.id, &newPokemon.name, &newPokemon.life, &newPokemon.tipo, &newPokemon.level)

			if err != nil {
				return pokes, err
			}

			err = seleccionAtaques(databases, &newPokemon)

			if err != nil {
				return pokes, err
			}

			pokes = append(pokes, newPokemon)

			check = true

		}

		if check == false {
			err = errNotPokemons
			return pokes, err
		}

	case onlyPokeFromUser:

		var newPokemon pokemon

		row := databases.QueryRow("SELECT DISTINCT pokemons.id,pokemons.name,pokemons.life,pokemons.type,pokemons.level FROM users_pokemons INNER JOIN pokemons ON users_pokemons.poke_id = pokemons.id WHERE users_pokemons.poke_id = ? AND users_pokemons.user_id = ?", idPoke, idUser)

		if err != nil {
			return
		}

		err = row.Scan(&newPokemon.id, &newPokemon.name, &newPokemon.life, &newPokemon.tipo, &newPokemon.level)

		if err != nil {
			err = errIncorrectID
			return
		}

		err = seleccionAtaques(databases, &newPokemon)

		if err != nil {
			return
		}

		pokes = append(pokes, newPokemon)

	case onlyPokeFromRival:

		var newPokemon pokemon

		row := databases.QueryRow("SELECT DISTINCT pokemons.id,pokemons.name,pokemons.life,pokemons.type,pokemons.level FROM users_pokemons INNER JOIN pokemons ON users_pokemons.poke_id = pokemons.id WHERE users_pokemons.poke_id = ?", idPoke)

		if err != nil {
			return
		}

		err = row.Scan(&newPokemon.id, &newPokemon.name, &newPokemon.life, &newPokemon.tipo, &newPokemon.level)

		if err != nil {
			return
		}

		err = seleccionAtaques(databases, &newPokemon)

		if err != nil {
			return
		}

		pokes = append(pokes, newPokemon)

	case allPokemons:

		rows, err := databases.Query("SELECT DISTINCT pokemons.id,pokemons.name,pokemons.life,pokemons.type,pokemons.level FROM pokemons")

		if err != nil {
			return pokes, err
		}

		defer rows.Close()

		for rows.Next() {
			var newPokemon pokemon

			err = rows.Scan(&newPokemon.id, &newPokemon.name, &newPokemon.life, &newPokemon.tipo, &newPokemon.level)

			if err != nil {
				return pokes, err
			}

			err = seleccionAtaques(databases, &newPokemon)

			if err != nil {
				return pokes, err
			}

			pokes = append(pokes, newPokemon)

			check = true

		}

	}

	return
}

func seleccionAtaques(databases *sql.DB, newPokemon *pokemon) (err error) {

	rows, err := databases.Query("SELECT DISTINCT attacks.id,attacks.name,attacks.power,attacks.accuracy FROM pokemons_attacks INNER JOIN attacks ON pokemons_attacks.attack_id=attacks.id WHERE pokemons_attacks.poke_id=?", newPokemon.id)

	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {

		var newAttack ataque

		err = rows.Scan(&newAttack.id, &newAttack.name, &newAttack.power, &newAttack.accuracy)

		if err != nil {
			return
		}

		newPokemon.ataques = append(newPokemon.ataques, newAttack)
	}

	return
}

func imprimirPokemons(pokes []pokemon) {

	for i := range pokes {
		fmt.Println(pokes[i].name)
		fmt.Printf(" ID -> %v\n", pokes[i].id)
		fmt.Printf(" type -> %v\n", pokes[i].tipo)
		fmt.Printf(" level -> %v\n", pokes[i].level)
		fmt.Printf(" ataques:\n")
		for i2 := range pokes[i].ataques {
			fmt.Printf("\t%v. Attack name -> %v\n", i2+1, pokes[i].ataques[i2].name)
			fmt.Printf("\t%v. Attack power-> %v\n", i2+1, pokes[i].ataques[i2].power)
			fmt.Printf("\t%v. Attack accuracy -> %v\n", i2+1, pokes[i].ataques[i2].accuracy)
			fmt.Println()
		}
	}
}

func añadirPoke(databases *sql.DB, user usuario) (err error) {

	var eleccionPoke, idAux int

	fmt.Println("¿Que pokemon Deseas?")

	allPokes, err := seleccionarPokemons(databases, allPokemons, 0, 0)

	if err != nil {
		return
	}

	imprimirPokemons(allPokes)

	idAux = len(allPokes) + 1

	fmt.Scan(&eleccionPoke)

	pokemonSelected, err := seleccionarPokemons(databases, onlyPokeFromRival, 0, eleccionPoke)

	if err != nil {
		return
	}

	newPokemon := pokemonSelected[0]

	err = insertarNuevoPoke(databases, newPokemon)

	if err != nil {
		return
	}

	err = insertarRelacionNuevoPoke(databases, user, newPokemon, idAux)
	if err != nil {
		return
	}

	err = insertarRelacionAtaques(databases, newPokemon.ataques, idAux)

	if err != nil {
		return
	}

	fmt.Println("Ahora el sera tu nuevo amigo :D")
	return
}

func insertarNuevoPoke(databases *sql.DB, newPokemon pokemon) (err error) {

	stmt, err := databases.Prepare("INSERT INTO pokemons(name,life,type,level) VALUES (?,?,?,?)")

	if err != nil {
		return
	}

	_, err = stmt.Exec(newPokemon.name, newPokemon.life, newPokemon.tipo, newPokemon.level)

	if err != nil {
		return
	}
	return
}

func insertarRelacionNuevoPoke(databases *sql.DB, user usuario, newPokemon pokemon, idAux int) (err error) {
	stmt, err := databases.Prepare("INSERT INTO users_pokemons (user_id,poke_id) VALUES (?,?)")

	if err != nil {
		return
	}

	_, err = stmt.Exec(user.id, idAux)

	if err != nil {
		return
	}
	return
}

func insertarRelacionAtaques(databases *sql.DB, attacks []ataque, idAux int) (err error) {

	for i := range attacks {

		stmtAttack, err := databases.Prepare("INSERT INTO pokemons_attacks (poke_id,attack_id) VALUES (?,?)")

		if err != nil {
			return err
		}

		_, err = stmtAttack.Exec(idAux, attacks[i].id)

		if err != nil {
			return err
		}

	}
	return
}

func liberarPokemon(databases *sql.DB, user usuario) (err error) {

	var pokeEliminar int

	allPokes, err := seleccionarPokemons(databases, allUserPokemons, user.id, 0)

	if err != nil {
		return
	}

	imprimirPokemons(allPokes)

	fmt.Println("Escribe el ID del pokemon que deseas eliminar")

	fmt.Scan(&pokeEliminar)

	sliceAux, err := seleccionarPokemons(databases, onlyPokeFromUser, user.id, pokeEliminar)

	if err != nil {
		return
	}

	if len(sliceAux) == 0 {
		fmt.Println("id no valido")
		return
	}

	for i := range sliceAux {

		stmtDelete, err := databases.Prepare("DELETE FROM users_pokemons WHERE users_pokemons.poke_id = ?;")

		if err != nil {
			log.Fatal(err)
		}

		_, err = stmtDelete.Exec(sliceAux[i].id)

		if err != nil {
			log.Fatal(err)
		}

	}

	fmt.Println("Hasta pronto amiguito :'D")
	return
}

func eliminarCuenta(databases *sql.DB, user usuario, salir *bool) (err error) {

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
					return err
				}

				_, err = stmtDelete.Exec(user.id)

				if err != nil {
					return err
				}

				fmt.Println("Hasta Pronto :'D")
				*salir = true
				return err

			}
			err = errPasswordInvalid
			return

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

	pokes, err := seleccionarPokemons(databases, allUserPokemons, user.id, 0)

	if err != nil {
		return
	}

	imprimirPokemons(pokes)

	fmt.Println("Escribe el ID del Pokemon con el que quieras jugar")

	fmt.Scan(&idPokeJugador)

	sliceAux, err := seleccionarPokemons(databases, onlyPokeFromUser, user.id, idPokeJugador)

	if err != nil {
		return
	}

	j1 := sliceAux[0]

	time.Sleep(2 * time.Second)

	fmt.Println("Tu oponente sera...")

	time.Sleep(time.Second * 2)

	posiblesRivales, err := seleccionarPokemons(databases, allPokemons, 0, 0)

	if err != nil {
		return
	}

	rand.Seed(time.Now().UnixNano())

	idRival := rand.Intn(len(posiblesRivales) + 1)

	pokeRival, err := seleccionarPokemons(databases, onlyPokeFromRival, 0, idRival+1)

	if err != nil {
		return
	}

	j2 := pokeRival[0]

	fmt.Println(j2.name)

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

			if ataquej1 >= len(j1.ataques) {

				fmt.Println("Opcion invalida")
				ataquej1 = -1

			}

			mostrarBatalla(&j1, &j2, ataquej1)

		} else {

			fmt.Printf("Turno %v\n", t)
			fmt.Println("Es turno del rival")

			ataquej2 := rand.Intn(len(j2.ataques))

			mostrarBatalla(&j2, &j1, ataquej2)
		}

		time.Sleep(time.Second * 3)

		t = t + 1

		if j1.life <= 0 || j2.life <= 0 {
			break
		}

		turno = !turno
		fmt.Println()
	}

	if j1.life <= 0 {

		fmt.Println("Oh no, perdiste")

	} else {

		fmt.Println("Genial ganaste esta batalla")

	}
}

func mostrarVidas(j1 pokemon, j2 pokemon) {

	fmt.Printf("Vida Actual: %v => %v	||	%v => %v\n", j1.name, j1.life, j2.name, j2.life)

}

func mostrarAtaques(poke pokemon) {

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
