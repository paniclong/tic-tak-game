package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
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

var numberCombinations = [][]int{
	{0, 0},
	{0, 1},
	{0, 2},
	{1, 0},
	{1, 1},
	{1, 2},
	{2, 0},
	{2, 1},
	{2, 2},
}

// Размер поля
const size = 3

func main() {
	var field [size][size]string

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			field[i][j] = " "
		}
	}

	rand.Seed(time.Now().UnixNano())

	bot := *new(Bot)
	player := *new(Player)

	bot.setCurrentCombination(0)
	bot.setLeftCells()
	bot.checkPreSetCombination(field)
	bot.setCurrentCell()

	for {
		changeField(&field, bot.getCurrentCell(), true)
		renderField(field)

		for {
			// Принимаем значение от пользователя
			var input int
			fmt.Println("Введите номер ячейки(от 1 до 9): ")
			_, err := fmt.Fscan(os.Stdin, &input)

			if err != nil {
				panic("Unknown error!")
			}

			player.currentCell = numberCombinations[input-1]

			if checkCell(field, player.currentCell) == true {
				changeField(&field, player.currentCell, false)

				if player.checkCombination(field) == true {
					renderField(field)

					fmt.Println("Победил игрок!")

					return
				}

				break
			}

			fmt.Println("Некорректный номер ячейки")
		}

		if bot.checkAndMaybeDeleteAvailableCombination(player.currentCell) == true {
			if len(combinations) == 0 {
				changeField(&field, bot.getCurrentCell(), true)
				renderField(field)

				fmt.Println("Ничья!")

				return
			}

			bot.generateNewCurrentCombination()
		}

		bot.setLeftCells()
		bot.checkPreSetCombination(field)
		bot.setCurrentCell()

		if len(bot.getLeftCells()) == 0 {
			changeField(&field, bot.getCurrentCell(), true)
			renderField(field)

			fmt.Println("Бот выиграл!")

			return
		}
	}
}

func renderField(field [size][size]string) {
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			fmt.Print("[ ", field[i][j], " ]")
		}

		fmt.Println()
	}
}

func checkCell(field [size][size]string, cell []int) bool {
	if field[cell[0]][cell[1]] == " " {
		return true
	}

	return false
}

func changeField(field *[size][size]string, cell []int, isBot bool) {
	if field[cell[0]][cell[1]] == " " {
		if isBot {
			field[cell[0]][cell[1]] = "X"
		} else {
			field[cell[0]][cell[1]] = "O"
		}
	}
}
