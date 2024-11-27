package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"time"
)

type StudyPlan struct {
	Sections map[string][]string `json:"sections"`
}

func main() {
	// Шаг 1: Считать JSON-файл
	filePath := "study_plan.json"
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read JSON file: %v", err)
	}

	var plan StudyPlan
	err = json.Unmarshal(data, &plan)
	if err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	// Шаг 2: Параметры генерации плана
	days := 5 // Количество дней для плана
	sectionsToInclude := []string{"reading", "speaking", "listening", "writing"}
	tasksPerSection := 1 // Количество задач из каждой секции

	// Шаг 3: Генерация плана
	generatedPlan := generateUniquePlanForDays(plan, sectionsToInclude, tasksPerSection, days)

	// Шаг 4: Вывод плана
	fmt.Println("Your 5-Day English Study Plan:")
	for day, dailyPlan := range generatedPlan {
		fmt.Printf("\nDay %d:\n", day+1)
		for section, tasks := range dailyPlan {
			fmt.Printf("  %s:\n", capitalize(section))
			for i, task := range tasks {
				fmt.Printf("    %d. %s\n", i+1, task)
			}
		}
	}
}

// Функция для генерации уникального плана на несколько дней
func generateUniquePlanForDays(plan StudyPlan, sections []string, tasksPerSection int, days int) []map[string][]string {
	rand.Seed(time.Now().UnixNano())
	var generatedPlan []map[string][]string
	usedTasks := make(map[string]map[string]bool) // Отслеживание использованных задач по секциям

	// Инициализация usedTasks для всех секций
	for section := range plan.Sections {
		usedTasks[section] = make(map[string]bool)
	}

	for day := 0; day < days; day++ {
		dailyPlan := make(map[string][]string)
		for _, section := range sections {
			if tasks, exists := plan.Sections[section]; exists {
				availableTasks := filterUnusedTasks(tasks, usedTasks[section])
				if len(availableTasks) == 0 {
					// Если задачи закончились, сбрасываем использованные
					log.Printf("Resetting tasks for section: %s", section)
					usedTasks[section] = make(map[string]bool)
					availableTasks = tasks
				}

				// Перемешать задачи и выбрать нужное количество
				rand.Shuffle(len(availableTasks), func(i, j int) { availableTasks[i], availableTasks[j] = availableTasks[j], availableTasks[i] })
				limit := min(tasksPerSection, len(availableTasks))
				selectedTasks := availableTasks[:limit]
				dailyPlan[section] = selectedTasks

				// Обновить использованные задачи
				for _, task := range selectedTasks {
					usedTasks[section][task] = true
				}
			} else {
				log.Printf("Warning: Section %s not found in JSON file", section)
			}
		}
		generatedPlan = append(generatedPlan, dailyPlan)
	}
	return generatedPlan
}

// Функция для фильтрации неиспользованных задач
func filterUnusedTasks(tasks []string, used map[string]bool) []string {
	var unusedTasks []string
	for _, task := range tasks {
		if !used[task] {
			unusedTasks = append(unusedTasks, task)
		}
	}
	return unusedTasks
}

// Вспомогательная функция для нахождения минимального значения
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Вспомогательная функция для преобразования строки в заглавный регистр
func capitalize(str string) string {
	if len(str) == 0 {
		return ""
	}
	return string(str[0]-32) + str[1:]
}
