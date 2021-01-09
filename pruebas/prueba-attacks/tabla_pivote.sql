.mode column
.header ON

--Activamos los fereign_keys
PRAGMA foreign_keys = ON;

--Creamos la tabla usuarios
CREATE TABLE users(
    id INTEGER PRIMARY KEY,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    age INTEGER NOT NULL
);

--Creamos la tabla pokemons
CREATE TABLE pokemons (
    id INTEGER PRIMARY KEY, 
    name TEXT NOT NULL,
    type TEXT NOT NULL, 
    level INTEGER NOT NULL 
);

--Creamos la tabla attacks
CREATE TABLE attacks (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    power  INTEGER NOT NULL,
    accuracy INTEGER NOT NULL
);

--Creamos la tabla pivote
CREATE TABLE users_pokemons (
    user_id INTEGER NOT NULL, 
    poke_id INTEGER NOT NULL UNIQUE,
    
    FOREIGN KEY (user_id)
        REFERENCES users (id),

    FOREIGN KEY(poke_id)
        REFERENCES pokemons(id)
);

--Creamos la tabla pivote
CREATE TABLE pokemons_attacks (
    poke_id INTEGER NOT NULL ,
    attack_id INTEGER NOT NULL,
    
    FOREIGN KEY(poke_id)
        REFERENCES pokemons(id),
        
    FOREIGN KEY (attack_id)
        REFERENCES attacks (id)

);

--Le insertamos datos a ambas tablas
INSERT INTO users (username,password,age) 
    VALUES
        ('cesar','cfabrica46',16),
        ('arturo','01234',20),
        ('sebas','sinmanos',16),
        ('raiza','rai',9);

INSERT INTO pokemons (name,type,level) 
    VALUES 
        ('pikachu','electric',32),
        ('pikachu','electric',50),
        ('charizard','fire',64),
        ('seal','water',64),
        ('electabuzz','electric',12),
        ('charmander','fire',22),
        ('flareon','fire',8),
        ('mewtwo','psychic',100),
        ('mew','psychic',100);

INSERT INTO attacks (name,power,accuracy)
    VALUES
        ('thunderbolt',20,100),
        ('thunder',80,70),
        ('flamethrower',60,100),
        ('fire blast',100,70),
        ('water gun',30,100),
        ('ice beam',50,100),
        ('electric punch',20,100),
        ('embers',20,100),
        ('psycho shock',70,100),
        ('psycho destruction',150,100);    

INSERT INTO users_pokemons (user_id,poke_id)
    VALUES
        (1,1),
        (2,2),
        (2,3),
        (3,4),
        (4,5),
        (4,6),
        (1,7),
        (1,8),
        (2,9);

 INSERT INTO pokemons_attacks (poke_id,attack_id)
    VALUES
        (1,1),
        (1,2),
        (2,1),
        (2,7),
        (3,3),
        (3,4),
        (4,5),
        (4,6),
        (5,2),
        (5,7),
        (6,8),
        (6,3),
        (7,8),
        (7,4),
        (8,9),
        (8,10),
        (9,9),
        (9,10);
       

--Para sacarle todo el provecho a la tabla pivote utilizaremos JOIN
SELECT users_pokemons.user_id, users.name, users_pokemons.poke_id, pokemons.name
    FROM users_pokemons
        INNER JOIN users ON
            users_pokemons.user_id = users.id
        INNER JOIN pokemons ON
            users_pokemons.poke_id = pokemons.id
    ORDER BY
        user_id ASC;

--4to Pedido
 SELECT users.id,users.name,users.age, COUNT(*)
    FROM users_pokemons
        INNER JOIN users ON
            users_pokemons.user_id = users.id
        INNER JOIN pokemons ON
            users_pokemons.poke_id = pokemons.id
    GROUP BY users_pokemons.user_id
    HAVING COUNT(*)>2;


SELECT  pokemons.name,pokemons.type,pokemons.level,attacks.name,attacks.power,attacks.accuracy
    FROM users_pokemons
        INNER JOIN users ON
            users_pokemons.user_id = users.id
        INNER JOIN pokemons ON
            users_pokemons.poke_id = pokemons.id
        INNER JOIN pokemons_attacks ON
            users_pokemons.poke_id = pokemons_attacks.poke_id
        INNER JOIN attacks ON
            pokemons_attacks.attack_id = attacks.id 
    WHERE users.id=1;
 