import React, { Component } from 'react'
import PropTypes from 'prop-types'
import { connect } from 'react-redux'
import { withRouter, Link } from 'react-router-dom'

import { Loading, Window } from '~/components'

@withRouter
@connect(
  () => {},
  (dispatch) => ({
    authActivate: dispatch.auth.activate,
  })
)
class Login extends Component {

  static propTypes = {
    authActivate: PropTypes.func,
    match: PropTypes.object,
    history: PropTypes.object,
  }

  state = {
    success: false,
  }

  async componentDidMount() {
    const { authActivate, history } = this.props
    const { token } = this.props.match.params
    try {
      await authActivate({ token })
      this.setState({ success: true })
    } catch(e) {
      history.push('/auth/login')
    }
  }

  render() {
    const { success } = this.state

    return (
      <Loading
        loading={!success}
      >
        <Window
          title='Activate'
          type='activate'
        >
          <div className='ctf-form-message'>Your account was activated. Now you can log in.</div>
          <div className='ctf-form-links'>
            <Link to='/auth/login'>Log In</Link>
          </div>
        </Window>
      </Loading>
    )
  }
}

export default Login
