import React, { Component } from 'react'
import PropTypes from 'prop-types'
import { connect } from 'react-redux'
import { withRouter, Link } from 'react-router-dom'

import { Window } from '~/components'
import { SignupForm } from './components'
import { addErrors } from '~/utils/form'

@withRouter
@connect(
  (state) => ({
    ...state.api.effects.auth.register,
  }),
  (dispatch) => ({
    authRegister: dispatch.auth.register,
  })
)
class Signup extends Component {

  static propTypes = {
    authRegister: PropTypes.func,
    history: PropTypes.object,
  }

  state = {
    fields: {},
    tokenIsSent: false,
  }

  handleSubmit = async(values) => {
    const { authRegister } = this.props

    try {
      await authRegister(values, { throw: true })
      this.setState({ tokenIsSent: true })
    } catch (e) {
      if (e.status === 400 || e.status === 409) {
        this.setState(({ fields }) => ({
          fields: addErrors(fields, e.errors),
        }))
      }
    }
  }

  handleChange = (fields) => {
    this.setState({ fields })
  }

  render() {
    const { tokenIsSent } = this.state

    return (
      <Window
        title='Sign Up'
        type='signup'
      >
        {tokenIsSent
            ? <div className='ctf-form-message'>Link was sent to your email. Please, check you inbox.</div>
            : (
              <SignupForm
                className='ctf-form ctf-form--signup'
                fields={this.state.fields}
                onSubmit={this.handleSubmit}
                onChange={this.handleChange}
              />
            )
        }
        <div className='ctf-form-links'>
          <Link to='/auth/login'>Log In</Link>
        </div>
      </Window>
    )
  }
}

export default Signup
