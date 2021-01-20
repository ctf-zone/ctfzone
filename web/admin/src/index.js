import 'regenerator-runtime/runtime'

import React from 'react'
import ReactDOM from 'react-dom'

import { init } from '@rematch/core'
import apiPlugin from './utils/rematch-api'
import { Provider } from 'react-redux'
import { createLogger } from 'redux-logger'

import { BrowserRouter as Router } from 'react-router-dom';

import models from './models'
import App from './scenes/App'

import './index.css'

let middlewares = []

if (process.env.NODE_ENV === 'development') {
  middlewares.push(createLogger())
}

const store = init({
  models,
  redux: {
    middlewares: middlewares,
  },
  plugins: [
    apiPlugin({}),
  ],
})

ReactDOM.render(
  <Provider store={store}>
    <Router>
      <App />
    </Router>
  </Provider>,
  document.getElementById('root'),
)
