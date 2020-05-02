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
	WhoStart string             // Кто начал игру player or bot
	Bot      *entity.Bot        // Объект сущности бота
	Player   *entity.Player     // Объект сущности игрока
	Field    [size][size]string // Внутриигровое виртуальное поле
	UID      string             // Уникальный идентификатор игры
}

// Обертка для handler'ов, для доступа к сессии и и логгеру
type Server struct {
	Store   *sessions.CookieStore
	Session *sessions.Session
	Logger  *Logger
}

// Структура всех json данных приходящих в запросе, которые необходимо распарсить
type JsonParsedDataFromRequest struct {
	WhoStart string
	Cell     int
}

// Домен фронта, для CORS
var frontServer = "http://" + os.Getenv("FRONT_DOMAIN") + ":" + os.Getenv("FRONT_PORT")

// Маппинг полей фронт:сервер
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

// Конвертируем комбинацию для фронта
func (server *Server) convertToFrontCombination(cell []int) string {
	for i, value := range mapping {
		if reflect.DeepEqual(value, cell) {
			return i
		}
	}

	return "0"
}

// Конвертируем комбинацию с фронта для бэка
func (server *Server) convertToBackedCombination(cell int) []int {
	var r []int

	for i, value := range mapping {
		if i == strconv.Itoa(cell) {
			return value
		}
	}

	return r
}

// Получаем из реквеста, кто стартанул игру
func (server *Server) getWhoStartFromRequest(request *http.Request) string {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		panic(err)
	}

	var parsedData JsonParsedDataFromRequest
	err = json.Unmarshal(body, &parsedData)
	if err != nil {
		log.Fatal(err)
	}

	if parsedData.WhoStart == "random" {
		var priority = rand.Intn(2)

		switch priority {
		case 0:
			parsedData.WhoStart = "bot"
		case 1:
			parsedData.WhoStart = "player"
		default:
			parsedData.WhoStart = "bot"
		}
	}

	return parsedData.WhoStart
}

// Получаем ячейку из реквеста для проставления
func (server *Server) getCellFromRequest(request *http.Request) []int {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		panic(err)
	}

	var parsedData JsonParsedDataFromRequest

	err = json.Unmarshal(body, &parsedData)
	if err != nil {
		log.Fatal(err)
	}

	return server.convertToBackedCombination(parsedData.Cell)
}

// Ставим заголовки
func (server *Server) setHeaders(w http.ResponseWriter) {
	w.Header().Add("Access-Control-Allow-Origin", frontServer)
	w.Header().Add("Access-Control-Allow-Credentials", "true")
}

// Ставим сессию
func (server *Server) setSessionGame(w http.ResponseWriter, r *http.Request) {
	session, err := server.Store.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	server.Session = session
}

// Сохраняем сессию
func (server *Server) saveSessionGame(w http.ResponseWriter, r *http.Request) {
	err := server.Session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Пишем в лог, когда стартанули игру
func (server *Server) writeInLogWhenStartedGame(currentGame *Game) {
	var message string

	server.Logger.SetPrefix(currentGame.UID)

	message = fmt.Sprintf(
		"Game started! Who start - %s; Current bot cell - %d; Current bot combination - %d; Current player cell - %d",
		currentGame.WhoStart,
		currentGame.Bot.GetCurrentCell(),
		currentGame.Bot.GetCurrentCombination(),
		currentGame.Player.GetCurrentCell(),
	)

	server.Logger.Write(message)
}

func (server *Server) sendResponse(w http.ResponseWriter, d []byte, c int) {
	_, err := w.Write(d)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(c)
}

// Метод выполняет несколько важных функций
// 1. Если игра уже была стартанута, то кидаем 403
// 2. Если игра ещё не была стартанута, то инициализуем изначальное состояние игры и сохраняем её в сессию
// 3. Если игру стартанул бот, то инициализируем бота, ставим рандомную ячейку и возвращаем её на фронт
func (server *Server) startGame(writer http.ResponseWriter, request *http.Request) {
	server.setHeaders(writer)
	server.setSessionGame(writer, request)

	// Если уже стартанули игру, то кидаем 403
	if server.Session.Values["game"] != nil {
		http.Error(writer, "", http.StatusForbidden)
	}

	currentGame := *new(Game)

	currentGame.WhoStart = server.getWhoStartFromRequest(request)
	currentGame.Player = new(entity.Player)
	currentGame.Bot = new(entity.Bot)

	currentGame.Bot.Initialize()
	currentGame.Player.Initialize()

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	currentGame.UID = uuid

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			currentGame.Field[i][j] = " "
		}
	}

	server.Session.Values["game"] = &currentGame
	server.saveSessionGame(writer, request)

	if server.Session.Values["game"] == nil {
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

		currentGame.changeField(true)

		response["cell"] = server.convertToFrontCombination(currentGame.Bot.GetCurrentCell())
	} else if currentGame.WhoStart == "player" {
		currentGame.Bot.SetCurrentCombination(0)
		currentGame.Bot.SetLeftCells()
	}

	server.writeInLogWhenStartedGame(&currentGame)

	server.Session.Values["game"] = &currentGame
	server.saveSessionGame(writer, request)

	response["status"] = "success"

	encodedData, _ := json.Marshal(response)

	server.sendResponse(writer, encodedData, http.StatusOK)
}

// Метод выполняет несколько функций:
// 1. Принимает ячейку от пользователя и ставит её
// 2. Проверяет, что у пользователя собралась комбинация, если да - то завершает игру
// 3. Проверяет, если у бота нарушилась комбинация из-за пользователя, то -
// 	  если есть доступные комбинации генерирует новую комбинацию для бота,
//	  в противном случае ничья и завершает игру
// 4. Ставит новую ячейку для бота и смотрит, собралась комбнация, если да - победил бот и завершает игру,
//    в противном случае возвращает ячейку проставленную ботом
func (server *Server) setCell(writer http.ResponseWriter, request *http.Request) {
	server.setHeaders(writer)
	server.setSessionGame(writer, request)

	if server.Session.Values["game"] == nil {
		writer.WriteHeader(http.StatusForbidden)
		return
	}

	currentGame := server.Session.Values["game"].(*Game)

	playerCell := server.getCellFromRequest(request)
	currentGame.Player.SetCurrentCell(playerCell)

	server.Logger.SetPrefix(currentGame.UID)
	server.Logger.Write(fmt.Sprintf("Set player cell - %d", playerCell))

	var response = map[string]string{}

	if currentGame.checkCell() == true {
		currentGame.changeField(false)

		if currentGame.Player.CheckCombination(&currentGame.Field) == true {
			response["win"] = "player"
			response["cell"] = server.convertToFrontCombination(currentGame.Player.GetCurrentCell())

			encodedData, _ := json.Marshal(response)

			server.Session.Values["game"] = nil
			err := server.Session.Save(request, writer)
			if err != nil {
				log.Fatal(err)
			}

			server.Logger.Write("Player win")

			server.sendResponse(writer, encodedData, http.StatusOK)

			return
		}
	}

	if currentGame.Bot.CheckAndMaybeDeleteAvailableCombination(currentGame.Player.GetCurrentCell()) == true {
		if len(currentGame.Bot.GetAllCombinations()) == 0 {
			currentGame.changeField(true)

			if currentGame.WhoStart == "player" {
				response["cell"] = server.convertToFrontCombination(currentGame.Bot.GetCurrentCell())
			}

			response["win"] = "draw"

			encodedData, _ := json.Marshal(response)

			server.Session.Values["game"] = nil
			err := server.Session.Save(request, writer)
			if err != nil {
				log.Fatal(err)
			}

			server.Logger.Write("Draw")

			server.sendResponse(writer, encodedData, http.StatusOK)

			return
		}

		currentGame.Bot.GenerateNewCurrentCombination()
		currentGame.Bot.SetLeftCells()
	}

	currentGame.Bot.CheckPreSetCombination(&currentGame.Field)
	currentGame.Bot.SetCurrentCell()

	currentGame.changeField(true)

	server.Logger.Write(fmt.Sprintf("Bot set cell - %d", currentGame.Bot.GetCurrentCell()))

	if len(currentGame.Bot.GetLeftCells()) == 0 {
		currentGame.changeField(true)

		response["win"] = "bot"
		response["cell"] = server.convertToFrontCombination(currentGame.Bot.GetCurrentCell())

		encodedData, _ := json.Marshal(response)

		server.Session.Values["game"] = nil
		err := server.Session.Save(request, writer)
		if err != nil {
			log.Fatal(err)
		}

		server.Logger.Write("Bot win")

		server.sendResponse(writer, encodedData, http.StatusOK)

		return
	}

	server.Session.Values["game"] = &currentGame
	server.saveSessionGame(writer, request)

	response["cell"] = server.convertToFrontCombination(currentGame.Bot.GetCurrentCell())

	encodedData, _ := json.Marshal(response)

	server.sendResponse(writer, encodedData, http.StatusOK)
}

// Проверяем, что игра уже была запущена и возвращаем ячейки, которые заполены ботом и игроком
func (server *Server) checkCurrentSessionGame(writer http.ResponseWriter, request *http.Request) {
	server.setHeaders(writer)
	server.setSessionGame(writer, request)

	if server.Session.Values["game"] == nil {
		writer.WriteHeader(http.StatusNoContent)

		return
	}

	currentGame := server.Session.Values["game"].(*Game)

	server.Logger.Write(fmt.Sprintf("Check current session for game - %s", currentGame.UID))

	var response = map[string][]string{}
	var cell = []int{0, 0}

	for i, value := range currentGame.Field {
		for j := range value {
			if currentGame.Field[i][j] == "X" {
				cell[0] = i
				cell[1] = j

				response["bot"] = append(response["bot"], server.convertToFrontCombination(cell))
			}

			if currentGame.Field[i][j] == "O" {
				cell[0] = i
				cell[1] = j

				response["player"] = append(response["player"], server.convertToFrontCombination(cell))
			}
		}
	}

	server.Logger.Write(fmt.Sprintf("Found cells - %s", response))

	encodedData, _ := json.Marshal(response)

	server.sendResponse(writer, encodedData, http.StatusOK)
}

// Принудительно завершаем игру, если нужно
func (server *Server) forceFinishCurrentGame(writer http.ResponseWriter, request *http.Request) {
	server.setHeaders(writer)
	server.setSessionGame(writer, request)

	if server.Session.Values["game"] == nil {
		writer.WriteHeader(http.StatusNoContent)

		return
	}

	currentGame := server.Session.Values["game"].(*Game)

	server.Logger.Write(fmt.Sprintf("Force finish game - %s", currentGame.UID))

	server.Session.Values["game"] = nil
	server.saveSessionGame(writer, request)

	writer.WriteHeader(http.StatusOK)
}

// Входная точка
func Run() {
	// Нужно для того, чтобы в структуру можно было записать в сессию
	gob.Register(&Game{})

	// Генерируем куки и ставим настройки
	var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECURE_KEY")))
	store.Options = &sessions.Options{
		Domain:   "localhost",
		Path:     "/",
		MaxAge:   3600 * 8, // 8 hours
		HttpOnly: true,
	}

	// Создаем логгер
	err, logger := CreateLogger()
	if err != nil {
		log.Fatal("Cannot create log")

		return
	}

	var server = *new(Server)

	server.Logger = logger
	server.Store = store

	// Инициализурем роуты
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/start", server.startGame)
	http.HandleFunc("/set", server.setCell)
	http.HandleFunc("/check", server.checkCurrentSessionGame)
	http.HandleFunc("/finish", server.forceFinishCurrentGame)

	fmt.Println("Server successfully started!")

	// Сам старт сервера
	log.Fatal(http.ListenAndServe(":"+os.Getenv("SERVER_PORT"), nil))
}

// Метод заглушка для корневого запроса
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

// Проверяем свободная ли ячейка или нет для игрока
func (game *Game) checkCell() bool {
	var cell = game.Player.GetCurrentCell()

	if game.Field[cell[0]][cell[1]] == " " {
		return true
	}

	return false
}

// Меняем внутриигровое поле
func (game *Game) changeField(isBot bool) {
	var cell []int

	if isBot == true {
		cell = game.Bot.GetCurrentCell()
	} else {
		cell = game.Player.GetCurrentCell()
	}

	if game.Field[cell[0]][cell[1]] == " " {
		if isBot {
			game.Field[cell[0]][cell[1]] = "X"
		} else {
			game.Field[cell[0]][cell[1]] = "O"
		}
	}
}
