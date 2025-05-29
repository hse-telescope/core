package models

type Project struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type Graph struct {
	ID             int    `db:"id"`
	ProjectID      int    `db:"project_id"`
	Name           string `db:"name"`
	CurrentVersion int    `db:"current_version"`
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

type GraphEvent struct {
	ID        int                    `db:"id"`
	GraphID   int                    `db:"graph_id"`
	EventType string                 `db:"event_type"`
	EventData map[string]interface{} `db:"event_data"`
	Version   int                    `db:"version"`
}
