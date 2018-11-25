import React, { Component } from 'react'
import PropTypes from 'prop-types'
import { connect } from 'react-redux'
import { withRouter } from 'react-router-dom'

@withRouter
@connect(
  (state) => ({
    ...state.api.effects.auth.logout,
  }),
  (dispatch) => ({
    authLogout: dispatch.auth.logout,
  })
)
class Logout extends Component {

  static propTypes = {
    authLogout: PropTypes.func,
    history: PropTypes.object,
  }

  componentDidMount() {
    const { authLogout, history } = this.props

    authLogout()
    history.push('/auth/login')
  }

  render() {
    return (
      <div>
      </div>
    )
  }
}

export default Logout
