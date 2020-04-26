import Request from "./repository/Request";
import React from "react";

import Selector from "./common/Selector";
import Field from "./common/Field";

import './style/App.scss'

export default class App extends React.Component {
  state = {
    info: 0,
    isGameStarted: false,
    isSavedGame: false,
    whoStartGame: 'bot',
    botCell: '',
    savedCells: [],
  }

  constructor(props) {
    super(props);

    this.onStartGame = this.onStartGame.bind(this)
    this.onChangeWhoStartGame = this.onChangeWhoStartGame.bind(this)
    this.onEndGame = this.onEndGame.bind(this)
    this.onContinueGame = this.onContinueGame.bind(this)
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
    Request.finishCurrentGame().then()

    this.setState({
      isGameStarted: false,
      isSavedGame: false,
      botCell: '',
      savedCells: [],
    })
  }

  onContinueGame() {
    this.setState({
      isSavedGame: false,
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
          isSavedGame: true,
          savedCells: res.data
        });
      }
    })
  }

  render() {
    return (
      <div className='wrapper'>
        {
          this.state.isGameStarted
            ?
            <Field
              isGameStarted={this.state.isGameStarted}
              isSavedGame={this.state.isSavedGame}
              botCell={this.state.botCell}
              savedCells={this.state.savedCells}
              onEndGame={this.onEndGame}
              onStartGame={this.onStartGame}
              onContinueGame={this.onContinueGame}
            />
            :
            <Selector
              whoStartGame={this.state.whoStartGame}
              onStartGame={this.onStartGame}
              onChangeWhoStartGame={this.onChangeWhoStartGame}
            />
        }
      </div>
    )
  }
}
