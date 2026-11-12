CREATE TABLE alfajor(
	id VARCHAR PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	relleno VARCHAR(255) NOT NULL,
	cobertura VARCHAR(255),
	precio INT NOT NULL,
	tapas INT NOT NULL CHECK(tapas IN (2,3)),
	marca VARCHAR(255) NOT NULL,
	descripcion TEXT,
	imagen VARCHAR(255)
);

CREATE TABLE usuario(
	id SERIAL PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	email VARCHAR(255) UNIQUE NOT NULL
);


CREATE TABLE review(
	id_alfajor VARCHAR(255) NOT NULL REFERENCES alfajor(id) ON DELETE CASCADE,
	id_user INT NOT NULL REFERENCES usuario(id) ON DELETE CASCADE,
	calificacion INT NOT NULL CHECK (calificacion BETWEEN 1 AND 10),
	comentario TEXT,
	fecha DATE NOT NULL DEFAULT CURRENT_DATE,
	PRIMARY KEY (id_alfajor, id_user)
);


