package entity

import (
	"math/rand"
	"reflect"
)

const size = 3

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
	CurrentCell        []int
	LeftCells          [][]int
	CurrentCombination [][]int
	Combinations       map[int][][]int
}

func (bot *Bot) Initialize() *Bot {
	bot.Combinations = combinations

	return bot
}

// Ставим новую комбинацию боту
func (bot *Bot) SetCurrentCombination(key int) *Bot {
	if key != 0 {
		bot.CurrentCombination = bot.Combinations[key]

		return bot
	}

	var keys []int
	var count = 0

	for i := range bot.Combinations {
		keys = append(keys, i)

		count++
	}

	if len(keys) == 0 {
		panic("Unknown error!")
	}

	randKey := rand.Intn(len(keys))

	bot.CurrentCombination = bot.Combinations[keys[randKey]]

	return bot
}

// Ставим остаток
func (bot *Bot) SetLeftCells() *Bot {
	bot.LeftCells = bot.CurrentCombination
	return bot
}

// Ставим новую текущую ячейку для бота
func (bot *Bot) SetCurrentCell() {
	randomCell := rand.Intn(len(bot.LeftCells))
	bot.CurrentCell = bot.LeftCells[randomCell]

	// Удаляем из остатка текущую ячейку
	bot.LeftCells[len(bot.LeftCells)-1],
		bot.LeftCells[randomCell] = bot.LeftCells[randomCell],
		bot.LeftCells[len(bot.LeftCells)-1]

	bot.LeftCells = bot.LeftCells[:len(bot.LeftCells)-1]
}

func (bot *Bot) GetCurrentCombination() [][]int {
	return bot.CurrentCombination
}

func (bot *Bot) GetLeftCells() [][]int {
	return bot.LeftCells
}

func (bot *Bot) GetCurrentCell() []int {
	return bot.CurrentCell
}

// Перед проставлением новой комбинации проверяем что,
// там уже не стояло предыдущее значение с предыдущей комбинации
// и если стоит, то в оставшиеся ячейки записываем остаток
func (bot *Bot) CheckPreSetCombination(field *[size][size]string) *Bot {
	var array [][]int

	for _, value := range bot.LeftCells {
		if field[value[0]][value[1]] == " " {
			array = append(array, value)
		}
	}

	if len(array) == 0 {
		panic("Unknown error!")
	}

	bot.LeftCells = array

	return bot
}

// Если предыдущая комбинация бота была нарушена пользователем, то удаляем её из всех возможных комбинаций
func (bot *Bot) CheckAndMaybeDeleteAvailableCombination(userCell []int) bool {
	isDeleted := false

	for i, value := range bot.Combinations {
		if reflect.DeepEqual(bot.CurrentCombination, value) == true {
			for _, value := range bot.CurrentCombination {
				if reflect.DeepEqual(value, userCell) == true {
					isDeleted = true
				}
			}
		}

		for _, value2 := range value {
			if reflect.DeepEqual(userCell, value2) == true {
				delete(bot.Combinations, i)
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

	for i, value := range bot.Combinations {
		for _, comb := range value {
			if reflect.DeepEqual(bot.CurrentCell, comb) == true {
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

func (bot *Bot) GetAllCombinations() map[int][][]int {
	return bot.Combinations
}
