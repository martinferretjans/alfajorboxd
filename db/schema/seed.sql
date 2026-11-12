INSERT INTO alfajor (id, name, relleno, cobertura, precio, tapas, marca, descripcion, imagen) VALUES
('havanna-chocolate', 'Havanna de chocolate', 'Dulce de leche', 'Chocolate semiamargo', 2500, 2, 'Havanna', 'El clasico de la costa atlantica.', 'havanna-chocolate.webp'),
('cachafaz-maicena', 'Cachafaz de maicena', 'Dulce de leche', 'Coco rallado', 2100, 2, 'Cachafaz', 'Tapas de maicena que se deshacen en la boca.', 'cachafaz-maicena.webp'),
('guaymallen-negro', 'Guaymallen negro', 'Dulce de leche', 'Chocolate', 800, 2, 'Guaymallen', 'El alfajor de kiosco por excelencia.', 'guaymallen-negro.webp'),
('jorgito-chocolate', 'Jorgito de chocolate', 'Dulce de leche', 'Chocolate', 1100, 2, 'Jorgito', 'Compañero de recreos.', 'jorgito-chocolate.webp'),
('fantoche-triple', 'Fantoche triple', 'Dulce de leche', 'Chocolate', 1300, 3, 'Fantoche', 'Tres tapas y mucho dulce de leche.', 'fantoche-triple.webp'),
('capitan-del-espacio', 'Capitan del Espacio', 'Dulce de leche', 'Chocolate', 1000, 2, 'Capitan del Espacio', 'Un clasico de culto.', 'capitan-del-espacio.webp'),
('terrabusi-blanco', 'Terrabusi blanco triple', 'Dulce de leche', 'Chocolate blanco', 1700, 3, 'Terrabusi', NULL, 'terrabusi-blanco.webp'),
('suchard-triple', 'Suchard mousse', 'Mousse de chocolate', 'Chocolate', 1500, 3, 'Suchard', NULL, 'suchard-triple.webp');

INSERT INTO usuario (name, email) VALUES
('martin', 'martin@alfajorboxd.com'),
('franco', 'franco@alfajorboxd.com'),
('lucia', 'lucia@alfajorboxd.com');

INSERT INTO review (id_alfajor, id_user, calificacion, comentario, fecha) VALUES
('havanna-chocolate', 1, 10, 'Perfecto, no le falta nada.', '2026-11-01'),
('havanna-chocolate', 2, 9, 'Muy bueno, aunque un poco caro.', '2026-11-03'),
('havanna-chocolate', 3, 9, 'El dulce de leche es lo mejor.', '2026-11-05'),
('cachafaz-maicena', 2, 8, 'El coco le da un toque unico.', '2026-11-02'),
('guaymallen-negro', 1, 6, 'Barato y cumple.', '2026-11-04'),
('guaymallen-negro', 3, 5, 'Nostalgico pero nada especial.', '2026-11-06'),
('fantoche-triple', 2, 7, 'Mucho chocolate, me encanto.', '2026-11-07'),
('capitan-del-espacio', 1, 8, 'Mejor de lo que esperaba.', '2026-11-07');
