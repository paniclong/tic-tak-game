import React from 'react';
import PropTypes from 'prop-types';

import Input from "./Input";
import Button from "./Button";

import '../style/Selector.scss'

const inputs = [
  {"value": 'bot', 'text': 'Бот'},
  {"value": 'player', 'text': 'Игрок'},
  {"value": 'random', 'text': 'Случайно'},
]

export default class Selector extends React.Component {
  static propTypes = {
    whoStartGame: PropTypes.string,
    onStartGame: PropTypes.func,
    onChangeWhoStartGame: PropTypes.func,
  }

  static defaultProps = {
    whoStartGame: 'bot',
    onStartGame: null,
    onChangeWhoStartGame: null,
  }

  render() {
    return (
      <div className='selector'>
        <p>Выберите кто начнет игру:</p>
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
