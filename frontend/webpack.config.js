const path = require("path");
const webpack = require('webpack');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');

module.exports = {
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
        test: /\.scss$/,
        use: [
          'style-loader',
          MiniCssExtractPlugin.loader,
          {
            loader: 'css-loader',
            options: {sourceMap: true}
          }, {
            loader: 'sass-loader',
            options: {sourceMap: true}
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
    new MiniCssExtractPlugin({
      filename: "[name].css",
      chunkFilename: "[id].css",
    })
  ]
}
