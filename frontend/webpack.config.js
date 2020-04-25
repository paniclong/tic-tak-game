const path = require("path");
const webpack = require('webpack');

module.exports = env = {
  entry: "./index.js",
  mode: "development",
  output: {
    filename: "./main.js"
  },
  devServer: {
    contentBase: path.join(__dirname, "static"),
    compress: true,
    host: process.env.FRONT_DOMAIN,
    port: process.env.FRONT_PORT,
    watchContentBase: true,
    progress: true
  },
  module: {
    rules: [
      {
        test: /\.m?js$/,
        exclude: /(node_modules|bower_components)/,
        use: {
          loader: "babel-loader"
        }
      },
      {
        test: /\.css$/,
        use: [
          "style-loader",
          {
            loader: "css-loader",
            options: {
              modules: true
            }
          }
        ]
      },
      {
        test: /\.(png|svg|jpg|gif)$/,
        use: ["file-loader"]
      }
    ]
  },
  plugins: [
    new webpack.DefinePlugin({
      'process.env': {
        FRONT_DOMAIN: JSON.stringify(process.env.FRONT_DOMAIN),
        FRONT_PORT: JSON.stringify(process.env.FRONT_PORT),
        SERVER_DOMAIN: JSON.stringify(process.env.SERVER_DOMAIN),
        SERVER_PORT: JSON.stringify(process.env.SERVER_PORT),
      },
    }),
  ]
}
