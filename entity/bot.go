package entity

import (
	"math/rand"
	"reflect"
)

var combinations = map[int][][]int{
	1: {
		{0, 0},
		{0, 1},
		{0, 2},
	},
	2: {
		{1, 0},
		{1, 1},
		{1, 2},
	},
	3: {
		{2, 0},
		{2, 1},
		{2, 2},
	},
	4: {
		{0, 0},
		{1, 0},
		{2, 0},
	},
	5: {
		{0, 1},
		{1, 1},
		{2, 1},
	},
	6: {
		{0, 2},
		{1, 2},
		{2, 2},
	},
	7: {
		{0, 0},
		{1, 1},
		{2, 2},
	},
	8: {
		{0, 2},
		{1, 1},
		{2, 0},
	},
}

type Bot struct {
	currentCell        []int
	leftCells          [][]int
	currentCombination [][]int
}

// Ставим новую комбинацию боту
func (bot *Bot) SetCurrentCombination(key int) *Bot {
	if key != 0 {
		bot.currentCombination = combinations[key]

		return bot
	}

	var keys []int
	var count = 0

	for i := range combinations {
		keys = append(keys, i)

		count++
	}

	if len(keys) == 0 {
		panic("Unknown error!")
	}

	randKey := rand.Intn(len(keys))

	bot.currentCombination = combinations[keys[randKey]]

	return bot
}

// Ставим остаток
func (bot *Bot) SetLeftCells() *Bot {
	bot.leftCells = bot.currentCombination

	return bot
}

// Ставим новую текущую ячейку для бота
func (bot *Bot) SetCurrentCell() {
	randomCell := rand.Intn(len(bot.leftCells))
	bot.currentCell = bot.leftCells[randomCell]

	// Удаляем из остатка текущую ячейку
	bot.leftCells[len(bot.leftCells)-1],
		bot.leftCells[randomCell] = bot.leftCells[randomCell],
		bot.leftCells[len(bot.leftCells)-1]

	bot.leftCells = bot.leftCells[:len(bot.leftCells)-1]
}

func (bot *Bot) GetCurrentCombination() [][]int {
	return bot.currentCombination
}

func (bot *Bot) GetLeftCells() [][]int {
	return bot.leftCells
}

func (bot *Bot) GetCurrentCell() []int {
	return bot.currentCell
}

const size = 3

// Перед проставлением новой комбинации проверяем что,
// там уже не стояло предыдущее значение с предыдущей комбинации
// и если стоит, то в оставшиеся ячейки записываем остаток
func (bot *Bot) CheckPreSetCombination(field [size][size]string) *Bot {
	var array [][]int

	for _, value := range bot.leftCells {
		if field[value[0]][value[1]] == " " {
			array = append(array, value)
		}
	}

	if len(array) == 0 {
		panic("Unknown error!")
	}

	bot.leftCells = array

	return bot
}

// Если предыдущая комбинация бота была нарушена пользователем, то удаляем её из всех возможных комбинаций
func (bot *Bot) CheckAndMaybeDeleteAvailableCombination(userCell []int) bool {
	isDeleted := false

	for i, value := range combinations {
		if reflect.DeepEqual(bot.currentCombination, value) == true {
			for _, value := range bot.currentCombination {
				if reflect.DeepEqual(value, userCell) == true {
					isDeleted = true
				}
			}
		}

		for _, value2 := range value {
			if reflect.DeepEqual(userCell, value2) == true {
				delete(combinations, i)
			}
		}
	}

	return isDeleted
}

// Метод генерируем новую комбинацию
// Выбираем рядом доступные комбинации
// Если их нет, выберем рандомную
func (bot *Bot) GenerateNewCurrentCombination() {
	var keys []int

	for i, value := range combinations {
		for _, comb := range value {
			if reflect.DeepEqual(bot.currentCell, comb) == true {
				keys = append(keys, i)
			}
		}
	}

	if len(keys) > 0 {
		bot.SetCurrentCombination(keys[rand.Intn(len(keys))])
	} else {
		bot.SetCurrentCombination(0)
	}
}
