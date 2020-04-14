package main

import (
	"fmt"
	"math/rand"
	"os"
	"reflect"
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

func main() {
	var field = [3][3]string{}

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			field[i][j] = ""
		}
	}

	rand.Seed(time.Now().UnixNano())

	var currentCell []int
	var selectedUserCell []int
	var flag = false

	_, selectedCombination := generateNewCombination()
	leftCombination := howMuchToLeftCombination(currentCell, selectedCombination)
	currentCell = generateNewCell(leftCombination)
	leftCombination = howMuchToLeftCombination(currentCell, selectedCombination)

	for {
		if len(combinations) == 0 {
			break
		}

		result := checkAvailable(selectedUserCell, selectedCombination)

		// Если пользователь ввел ту комбинацию которую уничтожает бота,
		// то перегенирируем комбинацию
		if true == result {
			_, selectedCombination = generateNewCombination()

			leftCombination = howMuchToLeftCombination(currentCell, selectedCombination)
			currentCell = generateNewCell(leftCombination)
		}

		field, _ = changeField(field, currentCell, true)
		renderField(field)

		if len(leftCombination) == 0 {
			fmt.Println("bot win!")
			break
		}

		for {
			// Принимаем значение от пользователя
			var input int
			fmt.Fscan(os.Stdin, &input)

			selectedUserCell = numberCombinations[input-1]
			field, flag = changeField(field, selectedUserCell, false)

			if true == flag {
				break
			}

			fmt.Println("Wrong cell")
		}

		leftCombination = howMuchToLeftCombination(currentCell, leftCombination)

		if len(leftCombination) == 0 {
			fmt.Println("bot win!")
			break
		}

		currentCell = generateNewCell(leftCombination)
	}

	return
}

func generateNewCombination() (int, [][]int) {
	var keys []int
	var count = 0

	for i := range combinations {
		keys = append(keys, i)

		count++
	}

	key := rand.Intn(len(keys))

	return count, combinations[keys[key]]
}

// Проверяем и удаляем комбинации которые выбрал пользователь
func checkAvailable(cell []int, botCombination [][]int) bool {
	flag := false

	for i, value := range combinations {
		if reflect.DeepEqual(botCombination, value) == true {
			for _, value := range botCombination {
				if reflect.DeepEqual(value, cell) == true {
					flag = true
				}
			}
		}

		for _, value2 := range value {
			if reflect.DeepEqual(cell, value2) == true {
				delete(combinations, i)
			}
		}
	}

	return flag
}

func howMuchToLeftCombination(cell []int, selComb [][]int) [][]int {
	var array [][]int

	for _, value := range selComb {
		if reflect.DeepEqual(cell, value) != true {
			array = append(array, value)
		}
	}

	return array
}

func generateNewCell(availableComb [][]int) []int {
	cell := rand.Intn(len(availableComb))

	return availableComb[cell]
}

func renderField(field [3][3]string) {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			fmt.Print("[ ", field[i][j], " ]")
		}

		fmt.Println()
	}
}

func changeField(field [3][3]string, cell []int, isBot bool) ([3][3]string, bool) {
	flag := false

	if len(field[cell[0]][cell[1]]) == 0 {
		if isBot {
			field[cell[0]][cell[1]] = "X"
		} else {
			field[cell[0]][cell[1]] = "O"
		}

		flag = true
	}

	return field, flag
}
