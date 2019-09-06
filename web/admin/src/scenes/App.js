import React, { Component } from 'react'
import { Switch, Route, withRouter } from 'react-router-dom'
import { connect } from 'react-redux'

import * as scenes from './scenes'

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

export default withRouter(connect()(App))
