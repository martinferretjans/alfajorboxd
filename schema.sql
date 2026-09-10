CREATE TABLE alfajor(
	id VARCHAR PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	relleno VARCHAR(255) NOT NULL,
	cobertura VARCHAR(255) NOT NULL,
	precio INT NOT NULL,
	tapas INT DEFAULT 2 CHECK(tapas IN (2,3)),
	marca VARCHAR(255),
	descripcion TEXT
);

CREATE TABLE usuario(
	id SERIAL PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	email VARCHAR(255) UNIQUE NOT NULL
);


CREATE TABLE review(
	id_alfajor VARCHAR(255) REFERENCES alfajor(id) ON DELETE CASCADE,
	id_user INT REFERENCES usuario(id) ON DELETE CASCADE,
	calificacion DECIMAL(3,2),
	comentario TEXT,
	fecha DATE,
	PRIMARY KEY (id_alfajor, id_user)
);


