import React from "react";
import PropTypes from "prop-types";

import Button from "./Button";

import Request from "../repository/Request";

const styles = {
  width: '33%',
  minWidth: '33%',
  flexGrow: '1',
  height: '2rem'
}

const divStyle = {
  display: "flex",
  flexWrap: "wrap",
  height: "5rem"
}

export default class Field extends React.Component {
  static propTypes = {
    isGameStarted: PropTypes.bool.isRequired,
    botCell: PropTypes.string.isRequired,
    savedCells: PropTypes.array.isRequired,
    onEndGame: PropTypes.func.isRequired,
    onStartGame: PropTypes.func.isRequired,
  }

  static defaultProps = {
    isGameStarted: false,
    botCell: '',
    savedCells: [],
    onEndGame: null,
    onStartGame: null,
  }

  state = {
    field: null,
    isEndGame: false,
    whoWin: '',
  }

  constructor(props) {
    super(props);

    this.state.field = this.getInitialStateField()

    this.onClick = this.onClick.bind(this);
  }

  /**
   * Маппинг игрового поля
   * Выносим в отдельный метод, так как в случае рестарта нам необходимо привести стейт к дефолтному значению
   *
   * @returns {({isDisabled: boolean, value: string, order: number})[]}
   */
  getInitialStateField() {
    return [
      {order: 0, value: " ", isDisabled: false},
      {order: 1, value: " ", isDisabled: false},
      {order: 2, value: " ", isDisabled: false},
      {order: 3, value: " ", isDisabled: false},
      {order: 4, value: " ", isDisabled: false},
      {order: 5, value: " ", isDisabled: false},
      {order: 6, value: " ", isDisabled: false},
      {order: 7, value: " ", isDisabled: false},
      {order: 8, value: " ", isDisabled: false},
    ]
  }

  /**
   * Обрабатываем клик на ячейку
   *
   * 1. Стучимся на сервер
   * 2. Проставляем ячейку для игрока и для бота (приходит в ответе)
   * 3. Если игра завершилась, пришло поле "win" в ответе,
   * парсим кто выиграл и в зависимости от этого меняем стейт isEndGame и whoWin
   *
   * @param {string} playerCellId
   */
  onClick(playerCellId) {
    const copyField = this.state.field

    Request.setCells(playerCellId).then(res => {
      const botCellId = res.data.hasOwnProperty('cell') ? res.data.cell : ""
      const whoWin = res.data.hasOwnProperty('win') ? res.data.win : ""

      let textOfWinner = ''
      let isEndGame = false

      this.blockField(copyField, playerCellId, false)

      if (typeof botCellId !== "undefined") {
        this.blockField(copyField, botCellId, true)
      }

      if (whoWin !== "") {
        copyField.map(obj => {
          obj.isDisabled = true
        })

        switch (whoWin) {
          case "bot":
            textOfWinner = 'К сожалению выиграл бот, попробуйте ещё раз!';
            break
          case "player":
            textOfWinner = 'Вы выиграли, поздравляю!';
            break
          case "draw":
            textOfWinner = 'Ничья!';
            break
        }

        isEndGame = true
      }

      this.setState({
        field: copyField,
        isEndGame: isEndGame,
        whoWin: textOfWinner
      })
    }).catch(err => {
      console.log(err)
    })
  }

  /**
   * Метод обертка, для избежания дубляжа в коде
   *
   * @param {array}   copyField
   * @param {string}  cellId
   * @param {boolean} isBot
   */
  blockField(copyField, cellId, isBot) {
    if (copyField[cellId].isDisabled === false) {
      copyField[cellId].isDisabled = true;
      copyField[cellId].value = isBot ? 'X' : 'O'
    }
  }

  /**
   * Обрабатываем событие, когда пользователь
   * после окончании игры, решил вернуться в меню
   *
   * @param e
   */
  onBackToMenu(e) {
    e.preventDefault()

    this.setState({ isEndGame: false })

    this.props.onEndGame()
  }

  /**
   * Обрабатываем событие, когда пользователь
   * после окончании игры, решил начать заново
   *
   * @param e
   */
  onReplayGame(e) {
    e.preventDefault()

    this.props.onEndGame()
    this.props.onStartGame(e)

    this.setState({
      isEndGame: false,
      field: this.getInitialStateField(),
    })
  }

  /**
   * Метод-монстр, выполняет две выжные функции
   *
   * 1. В начали инициализации, когда стартанули игру, проставляем ячейку для бота
   * 2. Если в сессии сохранена текущая игра, то парсим ячейки бота и игрока и проставляем в массив
   */
  componentDidMount() {
    const botCell = this.props.botCell
    const copyField = this.state.field

    if (this.props.isGameStarted && botCell !== "") {
      this.blockField(copyField, botCell, true)
    }

    if (this.props.savedCells.hasOwnProperty('bot')) {
      this.props.savedCells.bot.map(value => {
        this.blockField(copyField, value, true)
      })
    }

    if (this.props.savedCells.hasOwnProperty('player')) {
      this.props.savedCells.player.map(value => {
        this.blockField(copyField, value, false)
      })
    }

    this.setState({
      field: copyField,
    })
  }

  render() {
    return (
      <div>
        {
          this.state.isEndGame
            ?
            <div>
              <p>{this.state.whoWin} </p>
              <Button
                value={'Вернуться в меню'}
                onClick={(e) => this.onBackToMenu(e)}
              />
              <Button
                value={'Начать заново'}
                onClick={(e) => this.onReplayGame(e)}
              />
            </div>
            : ''
        }
        <div style={divStyle}>
          {
            this.state.field.map((el) =>
              <Button
                key={el.order}
                disabled={el.isDisabled}
                onClick={() => this.onClick(el.order)}
                value={el.value}
                styles={styles}
              />
            )
          }
        </div>
      </div>
    )
  }
}
