import React, { Component } from 'react'
import PropTypes from 'prop-types'
import { withRouter, Switch, Route } from 'react-router-dom'

import { NotFound } from '~/components'
import * as scenes from './scenes'

@withRouter
class Auth extends Component {

  static propTypes = {
    match: PropTypes.object,
    history: PropTypes.object,
  }

  componentDidMount() {
    const { match, history } = this.props

    if (match.isExact) {
      history.push('/auth/login')
    }
  }

  render() {
    const { match } = this.props

    return (
      <Switch>
        <Route exact path={`${match.url}/login`} component={scenes.Login} />
        <Route exact path={`${match.url}/logout`} component={scenes.Logout} />
        <Route exact path={`${match.url}/signup`} component={scenes.Signup} />
        <Route exact path={`${match.url}/send-token/:type(activate|reset)`} component={scenes.SendToken} />
        <Route exact path={`${match.url}/reset-password/:token([0-9a-f]{64})`} component={scenes.ResetPassword} />
        <Route exact path={`${match.url}/activate/:token([0-9a-f]{64})`} component={scenes.Activate} />
        <Route component={NotFound} />
      </Switch>
    )
  }
}

export default Auth
