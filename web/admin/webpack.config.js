const webpack = require('webpack')
const path = require('path')
const fs = require('fs')

const HtmlWebPackPlugin = require('html-webpack-plugin')
const MiniCssExtractPlugin = require('mini-css-extract-plugin')

module.exports = (env, argv) => ({
  entry: './src/index.js',
  output: {
    publicPath: '/',
    filename: '[name].[hash].js',
  },
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
                'react-css-modules',
                {
                  exclude: 'node_modules',
                  webpackHotModuleReloading: true,
                  generateScopedName: '[name]--[local]--[hash:base64:5]',
                },
              ],
              [
                'import',
                {
                  libraryName: 'antd',
                  libraryDirectory: 'es',
                  style: true,
                },
              ],
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
              plugins: () => [
                require('postcss-modules')({
                  globalModulePaths: ['codemirror', 'flag-icon-css'],
                  generateScopedName: '[name]--[local]--[hash:base64:5]',
                  // Get rid of 'css.json' files near all '.css' files.
                  getJSON: (cssFileName, json) => {
                    const cssName = path.basename(cssFileName, '.css')
                    const jsonFileName = path.resolve('./.cache/' + cssName + '.json')
                    fs.writeFileSync(jsonFileName, JSON.stringify(json))
                  },
                }),
                require('postcss-color-mod-function')(),
                require('postcss-custom-properties')({
                  // The only way to avoid duplicates in output css.
                  // preserve: 'computed',
                  preserve: true,
                }),
                require('autoprefixer')(),
              ],
            },
          },
        ],
      },
      {
        test: /\.less$/,
        use: [
          argv.mode === 'development' ? 'style-loader' : MiniCssExtractPlugin.loader,
          {
            loader: 'css-loader',
            options: {
              importLoaders: 1,
            },
          },
          {
            loader: 'less-loader',
            options: {
              javascriptEnabled: true,
            },
          },
        ],
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
    host: 'admin.ctfzone.test',
    port: 3001,
    contentBase: path.join(__dirname, 'dist'),
    historyApiFallback: true,
    hot: true,
    proxy: {
      '/api': {
        changeOrigin: true,
        target: 'http://localhost:8080',
        onProxyReq: (req) => {
          // Change host header to allow host routing.
          req.setHeader('Host', 'admin.ctfzone.test:8080')
        },
        onProxyRes: (proxyRes, req, res) => {
          const end = res.end
          let body = ''

          proxyRes.on('data', (data) => {
            body += data
          })

          res.write = () => {}

          res.end = () => {
            if ('link' in proxyRes.headers) {
              let link = proxyRes.headers.link
              res.setHeader('link', link.replace(':8080', ':3001'))
            }
            end.apply(res, [body])
          }
        },
      },
    },
  },
  plugins: [
    new webpack.HotModuleReplacementPlugin(),
    new webpack.EnvironmentPlugin({
      NODE_ENV: 'development',
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
})
