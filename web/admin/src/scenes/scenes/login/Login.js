import React, { Component } from 'react'

import { LoginForm } from './components'

import './Login.css'

class Login extends Component {
  render() {
    return (
      <div styleName='page'>
        <div styleName='form'>
          <LoginForm />
        </div>
      </div>
    )
  }
}

export default Login
