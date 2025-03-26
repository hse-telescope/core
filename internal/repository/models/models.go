package models

type Project struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type Graph struct {
	ID        int    `db:"id"`
	ProjectID int    `db:"project_id"`
	Name      string `db:"name"`
}

type Service struct {
	ID          int     `db:"id"`
	GraphID     int     `db:"graph_id"`
	Name        string  `db:"name"`
	Description string  `db:"description"`
	X           float32 `db:"x"`
	Y           float32 `db:"y"`
}

type Relation struct {
	ID          int    `db:"id"`
	GraphID     int    `db:"graph_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	FromService int    `db:"from_service"`
	ToService   int    `db:"to_service"`
}
