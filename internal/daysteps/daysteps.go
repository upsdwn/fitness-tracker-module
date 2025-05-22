// Пакет daysteps отвечает за учёт активности в течение дня.
//
// Он собирает переданную информацию в виде строк, парсит их и выводит
// информацию о количестве шагов, пройденной дистанции и потраченных калориях.
package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

// parsePackage парсит строку формата "678,0h50m",
// в которой 678 - шаги, 0h50m - продолжительность.
//
// Возвращает:
// int — количество шагов
// time.Duration — продолжительность прогулки.
// error — ошибку, если что-то пошло не так.
func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, errors.New("bad data format")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to extract steps: %w", err)
	}

	if steps <= 0 {
		return 0, 0, errors.New("steps must be positive")
	}

	d, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to extract duration: %w", err)
	}

	if d <= 0 {
		return 0, 0, errors.New("duration is not positive")
	}

	return steps, d, nil
}

// DayActionInfo вычисляет дистанцию в километрах и количество потраченных калорий,
// возвращает отформатированную строку с данными.
func DayActionInfo(data string, weight, height float64) string {
	steps, d, err := parsePackage(data)
	if err != nil {
		log.Printf("parsePackage: %v", err)
		return ""
	}

	if steps <= 0 {
		return ""
	}

	meters := float64(steps) * stepLength
	kilometers := meters / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, d)
	if err != nil {
		log.Printf("WalkingSpentCalories: %v", err)
		return ""
	}

	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps, kilometers, calories,
	)
}
