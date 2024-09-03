package postgres

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Хранилище данных.
type Storage struct {
	db *pgxpool.Pool
}

// Connect устанавливает соединение с базой данных и возвращает объект DB.
func Connect(constr string) (*Storage, error) {

	// Открываем соединение с базой данных
	db, err := pgxpool.Connect(context.Background(), constr)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
		return nil, err
	}
	s := Storage{
		db: db,
	}
	log.Println("Подключение к базе данных успешно установлено")
	return &s, nil
}

// Конструктор, принимает строку подключения к БД.
func New(constr string) (*Storage, error) {
	db, err := pgxpool.Connect(context.Background(), constr)
	if err != nil {
		return nil, err
	}
	s := Storage{
		db: db,
	}
	return &s, nil
}

// Задача.
type Task struct {
	ID         int
	Opened     int64
	Closed     int64
	AuthorID   int
	AssignedID int
	Title      string
	Content    string
}

// Tasks возвращает список задач из БД.
func (s *Storage) Tasks(taskID, authorID, labelID int) ([]Task, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT 
			t.id,
			t.opened,
			t.closed,
			t.author_id,
			t.assigned_id,
			t.title,
			t.content
		FROM tasks t
		LEFT JOIN tasks_labels tl ON t.id = tl.task_id
		LEFT JOIN labels l ON l.id = tl.label_id
		WHERE
			($1 = 0 OR t.id = $1) AND
			($2 = 0 OR t.author_id = $2) AND
			($3 = 0 OR l.id = $3)
		ORDER BY t.id;
	`,
		taskID,
		authorID,
		labelID,
	)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		tasks = append(tasks, t)

	}
	// ВАЖНО не забыть проверить rows.Err()
	return tasks, rows.Err()
}

// NewTask создаёт новую задачу и возвращает её id.
func (s *Storage) NewTask(t Task) (int, error) {
	var id int
	err := s.db.QueryRow(context.Background(), `
		INSERT INTO tasks (title, content)
		VALUES ($1, $2) RETURNING id;
		`,
		t.Title,
		t.Content,
	).Scan(&id)
	return id, err
}

// UpdateTask обновляет задачу в базе данных по ID.
func (s *Storage) UpdateTask(t Task) error {
	_, err := s.db.Exec(context.Background(), `
		UPDATE tasks
		SET
			title = $1,
			content = $2,
			closed = $3
		WHERE id = $4;
	`,
		t.Title,
		t.Content,
		t.Closed,
		t.ID,
	)
	return err
}

// DeleteTask удаляет задачу из базы данных по ID.
func (s *Storage) DeleteTask(taskID int) error {
	// Сначала удаляем все записи из tasks_labels, связанные с задачей
	_, err := s.db.Exec(context.Background(), `
		DELETE FROM tasks_labels WHERE task_id = $1;
	`, taskID)
	if err != nil {
		return err
	}

	// Затем удаляем саму задачу
	_, err = s.db.Exec(context.Background(), `
		DELETE FROM tasks WHERE id = $1;
	`, taskID)
	return err
}
