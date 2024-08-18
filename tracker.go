package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Определяем структуру задачи
type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Структура для хранения всех задач и следующего ID
type TaskStorage struct {
	Tasks  []Task `json:"tasks"`
	NextID int    `json:"next_id"`
}

// Инициализация хранилища задач
var taskStorage = TaskStorage{NextID: 1}

// Функция для загрузки задач из файла
func loadTasks(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Если файл не существует, возвращаем nil
			return nil
		}
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(&taskStorage)
}

// Функция для сохранения задач в файл
func saveTasks(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(taskStorage)
}

// Функция для добавления новой задачи
func addTask(description string) {
	task := Task{
		ID:          taskStorage.NextID,
		Description: description,
		Status:      "todo",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	taskStorage.NextID++
	taskStorage.Tasks = append(taskStorage.Tasks, task)
	fmt.Printf("Task added successfully (ID: %d)\n", task.ID)
}

// Функция для обновления задачи
func updateTask(id int, newDescription string) {
	for i, task := range taskStorage.Tasks {
		if task.ID == id {
			taskStorage.Tasks[i].Description = newDescription
			taskStorage.Tasks[i].UpdatedAt = time.Now()
			fmt.Println("Task updated successfully")
			return
		}
	}
	fmt.Println("Task not found")
}

// Функция для удаления задачи
func deleteTask(id int) {
	for i, task := range taskStorage.Tasks {
		if task.ID == id {
			taskStorage.Tasks = append(taskStorage.Tasks[:i], taskStorage.Tasks[i+1:]...)
			fmt.Println("Task deleted successfully")
			return
		}
	}
	fmt.Println("Task not found")
}

// Функция для изменения статуса задачи
func markTaskStatus(id int, status string) {
	for i, task := range taskStorage.Tasks {
		if task.ID == id {
			taskStorage.Tasks[i].Status = status
			taskStorage.Tasks[i].UpdatedAt = time.Now()
			fmt.Printf("Task marked as %s\n", status)
			return
		}
	}
	fmt.Println("Task not found")
}

// Функция для вывода списка задач по статусу
func listTasks(status string) {
	for _, task := range taskStorage.Tasks {
		if status == "" || task.Status == status {
			fmt.Printf("ID: %d, Description: %s, Status: %s, CreatedAt: %s, UpdatedAt: %s\n",
				task.ID, task.Description, task.Status,
				task.CreatedAt.Format(time.RFC3339),
				task.UpdatedAt.Format(time.RFC3339))
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: task-cli <command> [arguments]")
		return
	}

	command := os.Args[1]
	filePath := "tasks.json"

	if err := loadTasks(filePath); err != nil {
		fmt.Println("Error loading tasks:", err)
		return
	}

	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: task-cli add <description>")
			return
		}
		description := strings.Join(os.Args[2:], " ")
		addTask(description)
	case "update":
		if len(os.Args) < 4 {
			fmt.Println("Usage: task-cli update <id> <new_description>")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid task ID")
			return
		}
		newDescription := strings.Join(os.Args[3:], " ")
		updateTask(id, newDescription)
	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: task-cli delete <id>")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid task ID")
			return
		}
		deleteTask(id)
	case "mark-in-progress":
		if len(os.Args) < 3 {
			fmt.Println("Usage: task-cli mark-in-progress <id>")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid task ID")
			return
		}
		markTaskStatus(id, "in-progress")
	case "mark-done":
		if len(os.Args) < 3 {
			fmt.Println("Usage: task-cli mark-done <id>")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid task ID")
			return
		}
		markTaskStatus(id, "done")
	case "list":
		if len(os.Args) == 2 {
			listTasks("")
		} else if len(os.Args) == 3 {
			status := os.Args[2]
			if status == "todo" || status == "in-progress" || status == "done" {
				listTasks(status)
			} else {
				fmt.Println("Invalid status. Use 'todo', 'in-progress', or 'done'.")
			}
		} else {
			fmt.Println("Usage: task-cli list [status]")
		}
	default:
		fmt.Println("Unknown command:", command)
	}

	if err := saveTasks(filePath); err != nil {
		fmt.Println("Error saving tasks:", err)
	}
}
