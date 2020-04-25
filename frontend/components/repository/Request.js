import axios from "axios";
import React from "react";

class Request {
  constructor() {
    this.axios = axios.create({
      baseURL: 'http://' + process.env.SERVER_DOMAIN + ":" + process.env.SERVER_PORT,
      withCredentials: true,
    })
  }

  // noinspection JSValidateJSDoc
  /**
   * @param {string} whoStart
   *
   * @returns {Promise<AxiosResponse<any>>}
   */
  startGame(whoStart) {
    if (whoStart === "") {
      throw Error("Empty data!");
    }

    return this.axios.post('start', JSON.stringify({whoStart: whoStart}))
  }

  // noinspection JSValidateJSDoc
  /**
   * @param {string} cell
   *
   * @returns {Promise<AxiosResponse<any>>}
   */
  setCells(cell) {
    if (cell < 0 || cell > 8) {
      throw Error("Wrong cell");
    }

    return this.axios.post('set', JSON.stringify({Cell: cell}))
  }

  // noinspection JSValidateJSDoc
  /**
   * @returns {Promise<AxiosResponse<any>>}
   */
  checkCurrentGame() {
    return this.axios.get('check')
  }
}

export default new Request
