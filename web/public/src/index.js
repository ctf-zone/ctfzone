import React from 'react';
import ReactDOM from 'react-dom';

import { init } from '@rematch/core';
import apiPlugin from './utils/rematch-api';
import { Provider } from 'react-redux';

import { BrowserRouter as Router } from 'react-router-dom';

import models from './models';
import Main from './scenes/Main';

import './index.scss';

const store = init({
  models,
  plugins: [apiPlugin({})]
});

ReactDOM.render(
  <Provider store={store}>
    <Router>
      <Main />
    </Router>
  </Provider>,
  document.getElementById('root')
);
