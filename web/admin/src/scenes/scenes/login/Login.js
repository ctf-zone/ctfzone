import React, { Component } from 'react'

import { LoginForm } from './components'

import styles from './Login.module.css'

class Login extends Component {
  render() {
    return (
      <div className={styles.page}>
        <div className={styles.form}>
          <LoginForm />
        </div>
      </div>
    )
  }
}

export default Login
