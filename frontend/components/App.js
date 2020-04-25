import Request from "./repository/Request";
import React from "react";

import Selector from "./common/Selector";
import Field from "./common/Field";

export default class App extends React.Component {
  state = {
    info: 0,
    isGameStarted: false,
    whoStartGame: 'bot',
    botCell: '',
    savedCells: [],
  }

  constructor(props) {
    super(props);

    this.onStartGame = this.onStartGame.bind(this)
    this.onChangeWhoStartGame = this.onChangeWhoStartGame.bind(this)
    this.onEndGame = this.onEndGame.bind(this)
  }

  onStartGame(e) {
    e.preventDefault();

    Request.startGame(this.state.whoStartGame).then(res => {
      const cell = res.data.hasOwnProperty('cell') ? res.data.cell : ""

      this.setState({
        isGameStarted: true,
        botCell: cell,
        savedCells: [],
      })
    })
  }

  onEndGame() {
    this.setState({
      isGameStarted: false,
      botCell: '',
      savedCells: [],
    })
  }

  onChangeWhoStartGame(e) {
    this.setState({
      whoStartGame: e.target.value
    });
  }

  componentDidMount() {
    Request.checkCurrentGame().then(res => {
      if (res.status === 200) {
        this.setState({
          isGameStarted: true,
          savedCells: res.data
        });
      }
    })
  }

  render() {
    if (this.state.isGameStarted) {
      return (
        <div>
          <Field
            isGameStarted={this.state.isGameStarted}
            botCell={this.state.botCell}
            savedCells={this.state.savedCells}
            onEndGame={this.onEndGame}
            onStartGame={this.onStartGame}
          />
        </div>
      )
    }

    return (
      <div>
        <p>Выбери кто начнет игру:</p>
        <Selector
          isGameStarted={this.state.isGameStarted}
          whoStartGame={this.state.whoStartGame}
          onStartGame={this.onStartGame}
          onChangeWhoStartGame={this.onChangeWhoStartGame}
        />
      </div>
    )
  }
}
