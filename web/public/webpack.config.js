const webpack = require('webpack')
const path = require('path')

const HtmlWebPackPlugin = require('html-webpack-plugin')
const MiniCssExtractPlugin = require('mini-css-extract-plugin')
const UglifyJsPlugin = require('uglifyjs-webpack-plugin')

module.exports = (env, argv) => ({
  entry: './src/index.js',
  output: {
    publicPath: '/',
    filename: '[name].[hash].js',
  },
  devtool: 'source-map',
  module: {
    rules: [
      {
        test: /\.js$/,
        exclude: /node_modules/,
        use: {
          loader: 'babel-loader',
          options: {
            presets: [
              '@babel/preset-env',
              '@babel/preset-react',
            ],
            plugins: [
              '@babel/plugin-transform-runtime',
              'react-hot-loader/babel',
              [
                'root-import',
                {
                  rootPathSuffix: 'src',
                },
              ],
              [
                '@babel/plugin-proposal-decorators',
                {
                  legacy: true,
                },
              ],
              '@babel/plugin-proposal-class-properties',
            ],
          },
        },
      },
      {
        test: /\.css$/,
        use: [
          argv.mode === 'development' ? 'style-loader' : MiniCssExtractPlugin.loader,
          {
            loader: 'css-loader',
            options: {
              importLoaders: 1,
            },
          },
          {
            loader: 'postcss-loader',
            options: {
              plugins: () => {
                let plugins = [
                  require('postcss-color-mod-function')(),
                  require('postcss-custom-properties')({
                    preserve: true,
                  }),
                  require('postcss-nested')(),
                  require('postcss-mixins')(),
                  require('autoprefixer')(),
                ]

                if (argv.mode === 'production') {
                  plugins.push(require('cssnano')())
                }

                return plugins
              },
            },
          },
        ],
      },
      {
        test: /\.jpg/,
        loader: 'file-loader',
      },
      {
        test: /\.svg$/,
        loader: 'file-loader',
      },
      {
        test: /\.html$/,
        loader: 'html-loader',
      },
    ],
  },
  devServer: {
    host: 'ctfzone.test',
    port: 3000,
    contentBase: path.join(__dirname, 'dist'),
    historyApiFallback: true,
    hot: true,
    proxy: {
      '/api': {
        changeOrigin: true,
        target: 'http://localhost:8080',
        onProxyReq: (req) => {
          // Change host header to allow host routing.
          req.setHeader('Host', 'ctfzone.test:8080')
        },
      },
      '/files': {
        target: 'http://localhost:8080',
      },
    },
  },
  plugins: [
    ...(argv.mode !== 'production' ? [new webpack.HotModuleReplacementPlugin()] : []),
    new webpack.EnvironmentPlugin({
      NODE_ENV: argv.mode === 'development' ? 'development' : 'production',
      API_URL: '/api',
    }),
    new MiniCssExtractPlugin({
      filename: '[name].[hash].css',
      chunkFilename: '[id].[hash].css',
    }),
    new HtmlWebPackPlugin({
      template: './src/index.html',
    }),
  ],
  optimization: {
    minimizer: [new UglifyJsPlugin({
      sourceMap: true,
    })],
  },
})
