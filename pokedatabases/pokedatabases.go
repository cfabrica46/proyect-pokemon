package pokedatabases

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/cfabrica46/proyecto-pokemon/selectpoke"
)

//Estos errores seran utilizados a lo largo del paquete
var (
	ErrUserExist       = errors.New("El username que usted escogió ya esta en uso")
	ErrPasswordInvalid = errors.New("Password Invalido")
)

//FuncIngreso posee el menu
func FuncIngreso(databases *sql.DB, user selectpoke.User) (err error) {

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
			Jugar(databases, user)
		case 2:
			pokes, err := selectpoke.SeleccionarPokemons(databases, selectpoke.AllUserPokemons, user.ID, 0)

			if err != nil {
				fmt.Println(err.Error())
			}

			selectpoke.ImprimirPokemons(pokes)
		case 3:
			err := AñadirPoke(databases, user)

			if err != nil {
				fmt.Println(err.Error())
			}
		case 4:
			err := LiberarPokemon(databases, user)

			if err != nil {
				fmt.Println(err.Error())
			}
		case 5:
			err := EliminarCuenta(databases, user, &salir)

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

//GetUser vericara si existe un user registrado con los parametros predefinidos
func GetUser(databases *sql.DB, usernameScan, passwordScan string) (user *selectpoke.User, err error) {

	var userAux selectpoke.User

	row := databases.QueryRow("SELECT id,username,password FROM users WHERE username = ? AND password = ?", usernameScan, passwordScan)

	err = row.Scan(&userAux.ID, &userAux.Username, &userAux.Password)

	if err != nil {
		return
	}

	user = &userAux

	return

}

//FuncRegistro registra usuarios
func FuncRegistro(databases *sql.DB) (user *selectpoke.User, err error) {

	var usernameScan, passwordScan string

	fmt.Println("Ingrese su username")
	fmt.Scan(&usernameScan)
	fmt.Println("Ingrese su password")
	fmt.Scan(&passwordScan)

	err = CheckRegistro(databases, usernameScan, passwordScan)

	if err != nil {
		return
	}

	user, err = GetUser(databases, usernameScan, passwordScan)

	if err != nil {
		return
	}

	return
}

//CheckRegistro verifica si ya existe un usuario registrado con el mismo username
func CheckRegistro(databases *sql.DB, usernameScan, passwordScan string) (err error) {

	stmt, err := databases.Prepare("INSERT INTO users(username, password) VALUES(?,?)")

	if err != nil {
		return
	}
	_, err = stmt.Exec(usernameScan, passwordScan)

	if err != nil {
		err = ErrUserExist
		return
	}

	fmt.Println("nuevo usuario", usernameScan)

	return
}

//AñadirPoke añade pokemones a la databases
func AñadirPoke(databases *sql.DB, user selectpoke.User) (err error) {

	var eleccionPoke, idAux int

	fmt.Println("¿Que pokemon Deseas?")

	allPokes, err := selectpoke.SeleccionarPokemons(databases, selectpoke.AllPokemons, 0, 0)

	if err != nil {
		return
	}

	selectpoke.ImprimirPokemons(allPokes)

	idAux = len(allPokes) + 1

	fmt.Scan(&eleccionPoke)

	pokemonSelected, err := selectpoke.SeleccionarPokemons(databases, selectpoke.OnlyPokeFromRival, 0, eleccionPoke)

	if err != nil {
		return
	}

	newPokemon := pokemonSelected[0]

	err = InsertarNuevoPoke(databases, newPokemon)

	if err != nil {
		return
	}

	err = InsertarRelacionNuevoPoke(databases, user, newPokemon, idAux)
	if err != nil {
		return
	}

	err = InsertarRelacionAtaques(databases, newPokemon.Ataques, idAux)

	if err != nil {
		return
	}

	fmt.Println("Ahora el sera tu nuevo amigo :D")
	return
}

//InsertarNuevoPoke inserta un nuevo pokemon
func InsertarNuevoPoke(databases *sql.DB, newPokemon selectpoke.Pokemon) (err error) {

	stmt, err := databases.Prepare("INSERT INTO pokemons(name,life,type,level) VALUES (?,?,?,?)")

	if err != nil {
		return
	}

	_, err = stmt.Exec(newPokemon.Name, newPokemon.Life, newPokemon.Tipo, newPokemon.Level)

	if err != nil {
		return
	}
	return
}

//InsertarRelacionNuevoPoke inserta datos en la tabla pivote users_pokemons
func InsertarRelacionNuevoPoke(databases *sql.DB, user selectpoke.User, newPokemon selectpoke.Pokemon, idAux int) (err error) {
	stmt, err := databases.Prepare("INSERT INTO users_pokemons (user_id,poke_id) VALUES (?,?)")

	if err != nil {
		return
	}

	_, err = stmt.Exec(user.ID, idAux)

	if err != nil {
		return
	}
	return
}

//InsertarRelacionAtaques inserta datos en la tabla pivote pokemons_attacks
func InsertarRelacionAtaques(databases *sql.DB, attacks []selectpoke.Attack, idAux int) (err error) {

	for i := range attacks {

		stmtAttack, err := databases.Prepare("INSERT INTO pokemons_attacks (poke_id,attack_id) VALUES (?,?)")

		if err != nil {
			return err
		}

		_, err = stmtAttack.Exec(idAux, attacks[i].ID)

		if err != nil {
			return err
		}

	}
	return
}

//LiberarPokemon libera un pokemon
func LiberarPokemon(databases *sql.DB, user selectpoke.User) (err error) {

	var pokeEliminar int

	allPokes, err := selectpoke.SeleccionarPokemons(databases, selectpoke.AllUserPokemons, user.ID, 0)

	if err != nil {
		return
	}

	selectpoke.ImprimirPokemons(allPokes)

	fmt.Println("Escribe el ID del pokemon que deseas eliminar")

	fmt.Scan(&pokeEliminar)

	sliceAux, err := selectpoke.SeleccionarPokemons(databases, selectpoke.OnlyPokeFromUser, user.ID, pokeEliminar)

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

		_, err = stmtDelete.Exec(sliceAux[i].ID)

		if err != nil {
			log.Fatal(err)
		}

	}

	fmt.Println("Hasta pronto amiguito :'D")
	return
}

//EliminarCuenta te permite eliminar tu cuenta
func EliminarCuenta(databases *sql.DB, user selectpoke.User, salir *bool) (err error) {

	for {
		var preguntaSeguridad1, preguntaSeguridad2 string

		fmt.Println("¿Esta seguro de eliminar su cuenta? [S/N]")
		fmt.Scan(&preguntaSeguridad1)

		preguntaSeguridad1 = strings.ToLower(preguntaSeguridad1)

		switch preguntaSeguridad1 {
		case "s":
			fmt.Println("Introduzca su password")
			fmt.Scan(&preguntaSeguridad2)

			if preguntaSeguridad2 == user.Password {

				stmtDelete, err := databases.Prepare("DELETE FROM users where id = ?")

				if err != nil {
					return err
				}

				_, err = stmtDelete.Exec(user.ID)

				if err != nil {
					return err
				}

				fmt.Println("Hasta Pronto :'D")
				*salir = true
				return err

			}
			err = ErrPasswordInvalid
			return

		case "n":
			return

		default:
			fmt.Println("Opncion invalida")
		}
	}

}

//Jugar para JUGAR :D
func Jugar(databases *sql.DB, user selectpoke.User) {

	var idPokeJugador int
	var turno bool

	pokes, err := selectpoke.SeleccionarPokemons(databases, selectpoke.AllUserPokemons, user.ID, 0)

	if err != nil {
		return
	}

	selectpoke.ImprimirPokemons(pokes)

	fmt.Println("Escribe el ID del Pokemon con el que quieras jugar")

	fmt.Scan(&idPokeJugador)

	sliceAux, err := selectpoke.SeleccionarPokemons(databases, selectpoke.OnlyPokeFromUser, user.ID, idPokeJugador)

	if err != nil {
		return
	}

	j1 := sliceAux[0]

	time.Sleep(2 * time.Second)

	fmt.Println("Tu oponente sera...")

	time.Sleep(time.Second * 2)

	posiblesRivales, err := selectpoke.SeleccionarPokemons(databases, selectpoke.AllPokemons, 0, 0)

	if err != nil {
		return
	}

	rand.Seed(time.Now().UnixNano())

	idRival := rand.Intn(len(posiblesRivales) + 1)

	pokeRival, err := selectpoke.SeleccionarPokemons(databases, selectpoke.OnlyPokeFromRival, 0, idRival+1)

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
			MostrarVidas(j1, j2)
			fmt.Println("Que ataque desar usar?")

			MostrarAtaques(j1)

			fmt.Scan(&attack1)
			ataquej1 := attack1 - 1

			if ataquej1 >= len(j1.Ataques) {

				fmt.Println("Opcion invalida")
				ataquej1 = -1

			}

			MostrarBatalla(&j1, &j2, ataquej1)

		} else {

			fmt.Printf("Turno %v\n", t)
			fmt.Println("Es turno del rival")

			ataquej2 := rand.Intn(len(j2.Ataques))

			MostrarBatalla(&j2, &j1, ataquej2)
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

//MostrarVidas imprime la salud de ambos combatientes
func MostrarVidas(j1 selectpoke.Pokemon, j2 selectpoke.Pokemon) {

	fmt.Printf("Vida Actual: %v => %v	||	%v => %v\n", j1.Name, j1.Life, j2.Name, j2.Life)

}

//MostrarAtaques te muestra las opciones de ataques
func MostrarAtaques(poke selectpoke.Pokemon) {

	for i := range poke.Ataques {

		fmt.Printf("%v.%v => %v\n", i+1, poke.Ataques[i].Name, poke.Ataques[i].Power)

	}

}

//MostrarBatalla imprime lo sucedido en combate
func MostrarBatalla(atacante *selectpoke.Pokemon, receptor *selectpoke.Pokemon, eleccion int) {

	if eleccion == -1 {
		fmt.Println("fallaste")
	} else {

		fmt.Printf("%v uso %v contra %v\n", atacante.Name, atacante.Ataques[eleccion].Name, receptor.Name)

		receptor.Life = receptor.Life - atacante.Ataques[eleccion].Power
	}
}
