package entity

var playerCombination = map[int][][]int{
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

type Player struct {
	currentCell []int
}

func (player *Player) SetCurrentCell(cell []int) {
	player.currentCell = cell
}

func (player *Player) GetCurrentCell() []int {
	return player.currentCell
}

// Для игрока мы просто проверяем чтобы комбинации ноликов была в ряд
// @todo fix this
func (player *Player) CheckCombination(field [size][size]string) bool {
	var position []int

	for _, value := range playerCombination {
		for _, value2 := range value {
			if len(value2) != 0 {
				position = append(position, value2[0], value2[1])
			}
		}

		// По факту первое условия не нужно, но приходится его добавлять, так как ругается шторм
		if len(position) == 0 || len(position) < 6 {
			panic("Error!")
		}

		one := position[0]
		two := position[1]
		three := position[2]
		four := position[3]
		five := position[4]
		six := position[5]

		if field[one][two] == "O" && field[three][four] == "O" && field[five][six] == "O" {
			return true
		}

		position = nil
	}

	return false
}
