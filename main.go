package main

import (
	"fmt"
	"github.com/paniclong/tic-tak-game/entity"
	"github.com/paniclong/tic-tak-game/server"
	"log"
	"math/rand"
	"os"
	"time"
)

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
	rand.Seed(time.Now().UnixNano())

	if os.Getenv("TYPE_GAME") == "web" {
		fmt.Println("Starting server...")

		server.Run()
	}

	var field [size][size]string

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			field[i][j] = " "
		}
	}

	rand.Seed(time.Now().UnixNano())

	bot := *new(entity.Bot)
	player := *new(entity.Player)

	bot.Initialize()
	player.Initialize()

	bot.SetCurrentCombination(0)
	bot.SetLeftCells()
	bot.CheckPreSetCombination(&field)
	bot.SetCurrentCell()

	for {
		changeField(&field, bot.GetCurrentCell(), true)
		renderField(field)

		for {
			// Принимаем значение от пользователя
			var input int
			fmt.Println("Введите номер ячейки(от 1 до 9): ")

			_, err := fmt.Fscan(os.Stdin, &input)
			if err != nil {
				log.Fatalf("Error: %v", err)
			}

			player.SetCurrentCell(numberCombinations[input-1])

			if checkCell(field, player.GetCurrentCell()) == true {
				changeField(&field, player.GetCurrentCell(), false)

				if player.CheckCombination(&field) == true {
					renderField(field)

					fmt.Println("Победил игрок!")

					return
				}

				break
			}

			fmt.Println("Некорректный номер ячейки")
		}

		if bot.CheckAndMaybeDeleteAvailableCombination(player.GetCurrentCell()) == true {
			if len(bot.GetAllCombinations()) == 0 {
				changeField(&field, bot.GetCurrentCell(), true)
				renderField(field)

				fmt.Println("Ничья!")

				return
			}

			bot.GenerateNewCurrentCombination()
			bot.SetLeftCells()
		}

		bot.CheckPreSetCombination(&field)
		bot.SetCurrentCell()

		if len(bot.GetLeftCells()) == 0 {
			changeField(&field, bot.GetCurrentCell(), true)
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
