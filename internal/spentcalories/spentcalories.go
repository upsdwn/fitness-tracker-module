// Пакет spentcalories обрабатывает, рассчитывает потраченные калории.
//
// Расчет в зависимости от вида активности — бега или ходьбы.
// Возвращает информацию обо всех тренировках.
package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

// parseTraining принимает строку с данными формата "3456,Ходьба,3h00m",
// которая содержит количество шагов, вид активности и продолжительность активности.
//
// Возвращает:
// int — количество шагов.
// string — вид активности.
// time.Duration — продолжительность активности.
// error — ошибку, если что-то пошло не так.
func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, errors.New("bad data format")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("failed to extract steps: %w", err)
	}

	if steps <= 0 {
		return 0, "", 0, errors.New("steps is not positive")
	}

	activity := parts[1]

	d, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("failed to extract duration: %w", err)
	}

	if d <= 0 {
		return 0, "", 0, errors.New("duration is not positive")
	}

	return steps, activity, d, nil
}

// distance принимает количество шагов и рост пользователя в метрах,
// а возвращает дистанцию в километрах.
func distance(steps int, height float64) float64 {
	length := height * stepLengthCoefficient
	return float64(steps) * length / mInKm
}

// meanSpeed принимает количество шагов steps,
// рост пользователя height и продолжительность активности duration
// и возвращает среднюю скорость.
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0.0
	}

	dist := distance(steps, height)
	return dist / duration.Hours()
}

// TrainingInfo принимает:
// data string — строку с данными формата "3456,Ходьба,3h00m".
// weight, height float64 — вес (кг.) и рост (м.) пользователя.
//
// Возвращает:
// string — строка с информацией о тренировке в формате, приведенном ниже.
// error — ошибку, при ее возникновении внутри функции.
func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, d, err := parseTraining(data)
	if err != nil {
		return "", fmt.Errorf("parseTraining: %w", err)
	}

	var calories float64 = 0.0

	switch activity {
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, d)
		if err != nil {
			return "", fmt.Errorf("RunningSpentCalories: %w", err)
		}
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, d)
		if err != nil {
			return "", fmt.Errorf("WalkingSpentCalories: %w", err)
		}
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	text := `Тип тренировки: %s
Длительность: %.2f ч.
Дистанция: %.2f км.
Скорость: %.2f км/ч
Сожгли калорий: %.2f
`

	speed := meanSpeed(steps, height, d)
	dist := distance(steps, height)

	return fmt.Sprintf(text, activity, d.Hours(), dist, speed, calories), nil
}

// RunningSpentCalories принимает:
// steps int — количество шагов.
// weight, height float64 — вес(кг.) и рост(м.) пользователя.
// duration time.Duration — продолжительность бега.
//
// Возвращает:
// float64 — количество калорий, потраченных при беге.
// error — ошибку, если входные параметры некорректны.
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0.0, errors.New("steps is not positive")
	}

	if weight <= 0 {
		return 0.0, errors.New("weight is not positive")
	}

	if duration <= 0 {
		return 0.0, errors.New("duration is not positive")
	}

	ms := meanSpeed(steps, height, duration)

	return (weight * ms * duration.Minutes()) / minInH, nil
}

// WalkingSpentCalories принимает:
// steps int — количество шагов.
// weight, height float64 — вес(кг.) и рост(м.) пользователя.
// duration time.Duration — продолжительность ходьбы.
//
// Возвращает:
// float64 — количество калорий, потраченных при ходьбе.
// error — ошибку, если входные параметры некорректны.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0.0, errors.New("steps is not positive")
	}

	if weight <= 0 {
		return 0.0, errors.New("weight is not positive")
	}

	if height <= 0 {
		return 0.0, errors.New("height is not positive")
	}

	ms := meanSpeed(steps, height, duration)
	calories := weight * ms * duration.Minutes()
	caloriesSpent := calories / minInH

	return caloriesSpent * walkingCaloriesCoefficient, nil
}
