import React, { Component } from 'react'
import { Switch, Route, withRouter } from 'react-router-dom'
import { connect } from 'react-redux'
import { hot } from 'react-hot-loader'

import * as scenes from './scenes'

@hot(module)
@withRouter
@connect()
class App extends Component {
  render() {
    return (
      <Switch>
        <Route exact path='/login' component={scenes.Login} />
        <Route path='/' component={scenes.Main} />
      </Switch>
    )
  }
}

export default App
