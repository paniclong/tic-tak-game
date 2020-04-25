import React from 'react';
import PropTypes from 'prop-types';

import Input from "./Input";
import Button from "./Button";

const inputs = [
  {"value": 'bot', 'text': 'Бот'},
  {"value": 'player', 'text': 'Игрок'},
  {"value": 'random', 'text': 'Рандом'},
]

export default class Selector extends React.Component {
  static propTypes = {
    isGameStarted: PropTypes.bool,
    whoStartGame: PropTypes.string,
    onStartGame: PropTypes.func,
    onChangeWhoStartGame: PropTypes.func,
  }

  static defaultProps = {
    isGameStarted: false,
    whoStartGame: 'bot',
    onStartGame: null,
    onChangeWhoStartGame: null,
  }

  render() {
    if (this.props.isGameStarted === true) {
      return (
        <div/>
      )
    }

    return (
      <div>
        {
          inputs.map((robot) =>
            <Input
              key={robot.value}
              type={'radio'}
              value={robot.value}
              isChecked={this.props.whoStartGame === robot.value}
              text={robot.text}
              onChange={(e) => { this.props.onChangeWhoStartGame(e) }}
            />
          )
        }
        <Button value={'Начать'} onClick={(e) => {this.props.onStartGame(e)}}/>
      </div>
    )
  }
}
