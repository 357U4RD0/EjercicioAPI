CREATE TABLE IF NOT EXISTS incidents (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'pendiente',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO incidents (title, description, status) VALUES
('Fuga de agua', 'Hay una fuga en el baño del segundo piso', 'pendiente'),
('Fallo eléctrico', 'Corte de energía en el área de servidores', 'en progreso'),
('Accidente menor', 'Empleado se resbaló', 'resuelto'),
('Incendio', 'La cafetera se ha prendido en fuego', 'resuelto'),
('API', 'Problemas con el ejercicio', 'en progreso');