CREATE TABLE IF NOT EXISTS projects (
    id SERIAL PRIMARY KEY,
    name TEXT
);

CREATE TABLE IF NOT EXISTS graphs (
    id SERIAL PRIMARY KEY,
    project_id INTEGER REFERENCES projects(id) ON DELETE CASCADE,
    name TEXT
);

CREATE TABLE IF NOT EXISTS services (
    id SERIAL PRIMARY KEY,
    graph_id INTEGER REFERENCES graphs(id) ON DELETE CASCADE,
    name TEXT,
    description TEXT,
    x REAL,
    y REAL
);

CREATE TABLE IF NOT EXISTS relations (
    id SERIAL PRIMARY KEY,
    graph_id INTEGER REFERENCES graphs(id) ON DELETE CASCADE,
    name TEXT,
    description TEXT,
    from_service INTEGER REFERENCES services(id) ON DELETE CASCADE,
    to_service INTEGER REFERENCES services(id) ON DELETE CASCADE
);
