package selectpoke

import (
	"database/sql"
	"errors"
	"fmt"
)

//User representa los datos de la tabla users de databases.db
type User struct {
	ID       int
	Username string
	Password string
	Pokemons []Pokemon
}

//Pokemon representa los datos de la tabla pokemons de databases.db
type Pokemon struct {
	ID      int
	Name    string
	Life    int
	Tipo    string
	Level   int
	Ataques []Attack
}

//Attack representa los datos de la tabla attack de databases.db
type Attack struct {
	ID       int
	Name     string
	Power    int
	Accuracy int
}

//Estos errores seran utilizados a lo largo del paquete
var (
	ErrNotPokemons = errors.New("No tiene pokemons para acceder")
	ErrIncorrectID = errors.New("ID seleccionada Incorrecta")
)

//Parte fundamental de la funcion SeleccionarPokemons
const (
	AllUserPokemons = iota
	OnlyPokeFromUser
	OnlyPokeFromRival
	AllPokemons
)

//SeleccionarPokemons Selecciona los pokemons deseados
func SeleccionarPokemons(databases *sql.DB, flag int, idUser, idPoke int) (pokes []Pokemon, err error) {

	var check bool

	switch flag {
	case AllUserPokemons:

		rows, err := databases.Query("SELECT DISTINCT pokemons.id,pokemons.name,pokemons.life,pokemons.type,pokemons.level FROM users_pokemons INNER JOIN users ON users_pokemons.user_id = users.id INNER JOIN pokemons ON users_pokemons.poke_id = pokemons.id WHERE users.id = ?", idUser)

		if err != nil {
			return pokes, err
		}

		defer rows.Close()

		for rows.Next() {
			var newPokemon Pokemon

			err = rows.Scan(&newPokemon.ID, &newPokemon.Name, &newPokemon.Life, &newPokemon.Tipo, &newPokemon.Level)

			if err != nil {
				return pokes, err
			}

			err = SeleccionAtaques(databases, &newPokemon)

			if err != nil {
				return pokes, err
			}

			pokes = append(pokes, newPokemon)

			check = true

		}

		if check == false {
			err = ErrNotPokemons
			return pokes, err
		}

	case OnlyPokeFromUser:

		var newPokemon Pokemon

		row := databases.QueryRow("SELECT DISTINCT pokemons.id,pokemons.name,pokemons.life,pokemons.type,pokemons.level FROM users_pokemons INNER JOIN pokemons ON users_pokemons.poke_id = pokemons.id WHERE users_pokemons.poke_id = ? AND users_pokemons.user_id = ?", idPoke, idUser)

		if err != nil {
			return
		}

		err = row.Scan(&newPokemon.ID, &newPokemon.Name, &newPokemon.Life, &newPokemon.Tipo, &newPokemon.Level)

		if err != nil {
			err = ErrIncorrectID
			return
		}

		err = SeleccionAtaques(databases, &newPokemon)

		if err != nil {
			return
		}

		pokes = append(pokes, newPokemon)

	case OnlyPokeFromRival:

		var newPokemon Pokemon

		row := databases.QueryRow("SELECT DISTINCT pokemons.id,pokemons.name,pokemons.life,pokemons.type,pokemons.level FROM users_pokemons INNER JOIN pokemons ON users_pokemons.poke_id = pokemons.id WHERE users_pokemons.poke_id = ?", idPoke)

		if err != nil {
			return
		}

		err = row.Scan(&newPokemon.ID, &newPokemon.Name, &newPokemon.Life, &newPokemon.Tipo, &newPokemon.Level)

		if err != nil {
			return
		}

		err = SeleccionAtaques(databases, &newPokemon)

		if err != nil {
			return
		}

		pokes = append(pokes, newPokemon)

	case AllPokemons:

		rows, err := databases.Query("SELECT DISTINCT pokemons.id,pokemons.name,pokemons.life,pokemons.type,pokemons.level FROM pokemons")

		if err != nil {
			return pokes, err
		}

		defer rows.Close()

		for rows.Next() {
			var newPokemon Pokemon

			err = rows.Scan(&newPokemon.ID, &newPokemon.Name, &newPokemon.Life, &newPokemon.Tipo, &newPokemon.Level)

			if err != nil {
				return pokes, err
			}

			err = SeleccionAtaques(databases, &newPokemon)

			if err != nil {
				return pokes, err
			}

			pokes = append(pokes, newPokemon)

			check = true

		}

	}

	return
}

//SeleccionAtaques complementa a SelecionarPokemons dandole los ataques
func SeleccionAtaques(databases *sql.DB, newPokemon *Pokemon) (err error) {

	rows, err := databases.Query("SELECT DISTINCT attacks.id,attacks.name,attacks.power,attacks.accuracy FROM pokemons_attacks INNER JOIN attacks ON pokemons_attacks.attack_id=attacks.id WHERE pokemons_attacks.poke_id=?", newPokemon.ID)

	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {

		var newAttack Attack

		err = rows.Scan(&newAttack.ID, &newAttack.Name, &newAttack.Power, &newAttack.Accuracy)

		if err != nil {
			return
		}

		newPokemon.Ataques = append(newPokemon.Ataques, newAttack)
	}

	return
}

//ImprimirPokemons imprime pokes con un formato predefinido
func ImprimirPokemons(pokes []Pokemon) {

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
