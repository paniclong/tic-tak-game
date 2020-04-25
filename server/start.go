package server

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/paniclong/tic-tak-game/entity"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"strconv"
)

// Размер поля
const size = 3

type Game struct {
	// Статус игры
	Status   bool
	WhoStart string
	Bot      *entity.Bot
	Player   *entity.Player
	Field    [size][size]string
}

type JsonParsedDataFromRequest struct {
	WhoStart string
	Cell     int
}

var game = *new(Game)
var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECURE_KEY")))
var frontServer = "http://" + os.Getenv("FRONT_DOMAIN") + ":" + os.Getenv("FRONT_PORT")

var mapping = map[string][]int{
	"0": {0, 0},
	"1": {0, 1},
	"2": {0, 2},
	"3": {1, 0},
	"4": {1, 1},
	"5": {1, 2},
	"6": {2, 0},
	"7": {2, 1},
	"8": {2, 2},
}

func convertToFrontCombination(cell []int) string {
	for i, value := range mapping {
		if reflect.DeepEqual(value, cell) {
			return i
		}
	}

	return "0"
}

func convertToBackedCombination(cell int) []int {
	var r []int

	for i, value := range mapping {
		if i == strconv.Itoa(cell) {
			return value
		}
	}

	return r
}

func getWhoStartFromRequest(request *http.Request) string {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		panic(err)
	}

	var parsedData JsonParsedDataFromRequest
	err = json.Unmarshal(body, &parsedData)
	if err != nil {

	}

	return parsedData.WhoStart
}

func getCellFromRequest(request *http.Request) []int {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		panic(err)
	}

	var parsedData JsonParsedDataFromRequest

	err = json.Unmarshal(body, &parsedData)
	if err != nil {

	}

	return convertToBackedCombination(parsedData.Cell)
}

func startGame(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Access-Control-Allow-Origin", frontServer)
	writer.Header().Add("Access-Control-Allow-Credentials", "true")

	session, err := store.Get(request, "session")
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	var currentGame = &game

	// Стартуем игру и записываем в сессию всю инфу
	if session.Values["game"] == nil {
		var whoStart = getWhoStartFromRequest(request)

		if whoStart == "random" {
			var priority = rand.Intn(2)

			switch priority {
			case 0:
				whoStart = "bot"
			case 1:
				whoStart = "player"
			}
		}

		currentGame.Status = true
		currentGame.WhoStart = whoStart
		currentGame.Player = new(entity.Player)
		currentGame.Bot = new(entity.Bot)

		currentGame.Bot.Initialize()
		game.Player.Initialize()

		for i := 0; i < size; i++ {
			for j := 0; j < size; j++ {
				currentGame.Field[i][j] = " "
			}
		}

		session.Values["game"] = &currentGame
	} else {
		currentGame = session.Values["game"].(*Game)
	}

	err = session.Save(request, writer)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	if currentGame.Status == false {
		fmt.Print("Failed to start game, check code")
		return
	}

	var response = map[string]string{}

	// Если стартует бот
	if currentGame.WhoStart == "bot" {
		currentGame.Bot.SetCurrentCombination(0)
		currentGame.Bot.SetLeftCells()
		currentGame.Bot.CheckPreSetCombination(&currentGame.Field)
		currentGame.Bot.SetCurrentCell()

		changeField(&currentGame.Field, currentGame.Bot.GetCurrentCell(), true)

		response["cell"] = convertToFrontCombination(currentGame.Bot.GetCurrentCell())
	} else if currentGame.WhoStart == "player" {
		currentGame.Bot.SetCurrentCombination(0)
		currentGame.Bot.SetLeftCells()
	}

	session.Values["game"] = &currentGame

	err = session.Save(request, writer)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Set bot current combination", currentGame.Bot.GetCurrentCombination())

	response["status"] = "success"

	encodedData, _ := json.Marshal(response)

	_, err = writer.Write(encodedData)
	if err != nil {
		log.Fatal(err)
	}
}

func setCell(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Access-Control-Allow-Origin", frontServer)
	writer.Header().Add("Access-Control-Allow-Credentials", "true")

	session, err := store.Get(request, "session")
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	if session.Values["game"] == nil {
		writer.WriteHeader(http.StatusForbidden)
		return
	}

	currentGame := session.Values["game"].(*Game)

	playerCell := getCellFromRequest(request)
	currentGame.Player.SetCurrentCell(playerCell)

	var response = map[string]string{}

	if checkCell(&currentGame.Field, currentGame.Player.GetCurrentCell()) == true {
		changeField(&currentGame.Field, currentGame.Player.GetCurrentCell(), false)

		if currentGame.Player.CheckCombination(&currentGame.Field) == true {
			response["win"] = "player"
			response["cell"] = convertToFrontCombination(currentGame.Player.GetCurrentCell())

			encodedData, _ := json.Marshal(response)

			session.Values["game"] = nil
			err = session.Save(request, writer)

			_, err = writer.Write(encodedData)
			if err != nil {
				log.Fatal(err)
			}

			writer.WriteHeader(http.StatusOK)

			return
		}
	}

	if currentGame.Bot.CheckAndMaybeDeleteAvailableCombination(currentGame.Player.GetCurrentCell()) == true {
		if len(currentGame.Bot.GetAllCombinations()) == 0 {
			changeField(&currentGame.Field, currentGame.Bot.GetCurrentCell(), true)

			response["win"] = "draw"

			encodedData, _ := json.Marshal(response)

			session.Values["game"] = nil
			err = session.Save(request, writer)

			_, err = writer.Write(encodedData)
			if err != nil {
				log.Fatal(err)
			}

			writer.WriteHeader(http.StatusOK)

			return
		}

		currentGame.Bot.GenerateNewCurrentCombination()
		currentGame.Bot.SetLeftCells()

		fmt.Println("generate bot current combination", currentGame.Bot.GetCurrentCombination())
	}

	currentGame.Bot.CheckPreSetCombination(&currentGame.Field)
	currentGame.Bot.SetCurrentCell()

	changeField(&currentGame.Field, currentGame.Bot.GetCurrentCell(), true)

	if len(currentGame.Bot.GetLeftCells()) == 0 {
		changeField(&currentGame.Field, currentGame.Bot.GetCurrentCell(), true)

		response["win"] = "bot"
		response["cell"] = convertToFrontCombination(currentGame.Bot.GetCurrentCell())

		encodedData, _ := json.Marshal(response)

		session.Values["game"] = nil
		err = session.Save(request, writer)

		_, err = writer.Write(encodedData)
		if err != nil {
			log.Fatal(err)
		}

		writer.WriteHeader(http.StatusOK)

		return
	}

	session.Values["game"] = &currentGame

	err = session.Save(request, writer)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	response["cell"] = convertToFrontCombination(currentGame.Bot.GetCurrentCell())

	encodedData, _ := json.Marshal(response)

	_, err = writer.Write(encodedData)
	if err != nil {
		log.Fatal(err)
	}

	writer.WriteHeader(http.StatusOK)
}

func checkCurrentSessionGame(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Access-Control-Allow-Origin", frontServer)
	writer.Header().Add("Access-Control-Allow-Credentials", "true")

	session, err := store.Get(request, "session")
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	if session.Values["game"] == nil {
		writer.WriteHeader(http.StatusNoContent)

		return
	}

	currentGame := session.Values["game"].(*Game)

	var response = map[string][]string{}
	var cell = []int{0, 0}

	for i, value := range currentGame.Field {
		for j := range value {
			if currentGame.Field[i][j] == "X" {
				cell[0] = i
				cell[1] = j

				response["bot"] = append(response["bot"], convertToFrontCombination(cell))
			}

			if currentGame.Field[i][j] == "O" {
				cell[0] = i
				cell[1] = j

				response["player"] = append(response["player"], convertToFrontCombination(cell))
			}
		}
	}

	encodedData, _ := json.Marshal(response)

	_, err = writer.Write(encodedData)
	if err != nil {
		log.Fatal(err)
	}

	writer.WriteHeader(http.StatusOK)
}

func Run() {
	gob.Register(&Game{})

	store.Options = &sessions.Options{
		Domain:   "localhost",
		Path:     "/",
		MaxAge:   3600 * 8, // 8 hours
		HttpOnly: true,
	}

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/start", startGame)
	http.HandleFunc("/set", setCell)
	http.HandleFunc("/check", checkCurrentSessionGame)

	fmt.Println("Server successfully started!")

	log.Fatal(http.ListenAndServe(":"+os.Getenv("SERVER_PORT"), nil))
}

// Метод заглушка
func rootHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writer.Header().Set("Connection", "close")
		writer.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	writer.Header().Add("Access-Control-Allow-Origin", "*")

	var info = map[string]string{
		"connect": "successfully",
	}

	encodedData, _ := json.Marshal(info)

	_, err := writer.Write(encodedData)
	if err != nil {
		log.Fatal(err)
	}
}

func checkCell(field *[size][size]string, cell []int) bool {
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
