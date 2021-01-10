CREATE TABLE IF NOT EXISTS users(
    id INTEGER PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);


CREATE TABLE IF NOT EXISTS pokemons (
    id INTEGER PRIMARY KEY, 
    name TEXT NOT NULL,
    life INTEGER NOT NULL,
    type TEXT NOT NULL, 
    level INTEGER NOT NULL 
);

CREATE TABLE IF NOT EXISTS attacks (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    power  INTEGER NOT NULL,
    accuracy INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS users_pokemons (
    user_id INTEGER NOT NULL, 
    poke_id INTEGER NOT NULL UNIQUE,
    
    FOREIGN KEY (user_id)
        REFERENCES users (id),

    FOREIGN KEY(poke_id)
        REFERENCES pokemons(id)
);


CREATE TABLE IF NOT EXISTS pokemons_attacks (
    poke_id INTEGER NOT NULL ,
    attack_id INTEGER NOT NULL,
    
    FOREIGN KEY(poke_id)
        REFERENCES pokemons(id),
        
    FOREIGN KEY (attack_id)
        REFERENCES attacks (id)
);

INSERT INTO users (username,password) 
    VALUES
        ('cesar','cfabrica46'),
        ('arturo','01234'),
        ('sebas','sinmanos'),
        ('raiza','rai');

INSERT INTO pokemons (name,life,type,level) 
    VALUES 
        ('pikachu',200,'electric',32),
        ('pikachu',200,'electric',50),
        ('charizard',300,'fire',64),
        ('seal','water',250,64),
        ('electabuzz',260,'electric',12),
        ('charmander',150,'fire',22),
        ('flareon',250,'fire',8),
        ('mewtwo',500,'psychic',100),
        ('mew',500,'psychic',100);

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
       