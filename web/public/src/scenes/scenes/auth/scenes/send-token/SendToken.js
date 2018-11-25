import React, { Component } from 'react'
import PropTypes from 'prop-types'
import { connect } from 'react-redux'
import { withRouter, Link } from 'react-router-dom'

import { Window } from '~/components'
import { SendTokenForm } from './components'
import { addErrors } from '~/utils/form'

@withRouter
@connect(
  () => {},
  (dispatch) => ({
    authSendToken: dispatch.auth.sendToken,
  })
)
class SendToken extends Component {

  static propTypes = {
    authSendToken: PropTypes.func,
    match: PropTypes.object,
  }

  state = {
    fields: {},
    tokenIsSent: false,
  }

  handleSubmit = async(values) => {
    const { authSendToken } = this.props
    const { type } = this.props.match.params

    try {
      await authSendToken({ type, ...values }, { throw: true })
      this.setState({ tokenIsSent: true })
    } catch (e) {
      let errors = {}

      if (e.status === 400) {
        errors = e.errors
      }

      this.setState(({ fields }) => ({
        fields: addErrors(fields, errors),
      }))
    }
  }

  handleChange = (fields) => {
    this.setState({ fields })
  }

  render() {
    const { tokenIsSent } = this.state

    return (
      <Window
        title='Send Token'
        type='send-token'
      >
        {tokenIsSent
            ? <div className='ctf-form-message'>Link was sent to your email. Please, check you inbox.</div>
            : (
              <SendTokenForm
                className='ctf-form ctf-form--send-token'
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

export default SendToken
