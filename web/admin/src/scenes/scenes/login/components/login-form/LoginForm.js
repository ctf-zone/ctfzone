import React, { Component } from 'react'
import PropTypes from 'prop-types'
import { connect } from 'react-redux'
import { withRouter } from 'react-router-dom'

import { Form, Input, Icon, Button } from 'antd'

// import logo from '~/images/logo-big.svg'
import logo from '~/images/mctf.svg'

import { hasErrors } from '~/utils/form'

import './LoginForm.css'

@withRouter
@connect(
  state => ({
    ...state.api.effects.auth.login,
  }),
  dispatch => ({
    authLogin: dispatch.auth.login,
  }),
)
@Form.create({
  mapPropsToFields: props => {
    const { error } = props

    if (error && error.response && error.response.status == 401) {
      return {
        password: Form.createFormField({
          value: '',
          errors: [new Error('Invalid password')],
        }),
      }
    }
  },
})
class LoginForm extends Component {

  static propTypes = {
    // @connect
    success: PropTypes.bool.isRequired,
    error: PropTypes.object,
    loading: PropTypes.bool.isRequired,

    authLogin: PropTypes.func.isRequired,

    // @Form.create
    form: PropTypes.object,

    // @withRouter
    match: PropTypes.object.isRequired,
    location: PropTypes.object.isRequired,
    history: PropTypes.object.isRequired,
  }

  componentDidMount() {
    setTimeout(() => {
      this.passwordInput.focus()
    })
  }

  componentDidUpdate() {
    const { history, success } = this.props

    if (success) {
      history.push('/')
    }
  }

  handleSubmit = e => {
    e.preventDefault()

    const { authLogin } = this.props
    const { getFieldValue } = this.props.form

    authLogin({ password: getFieldValue('password') })
  }

  renderPasswordField({ getFieldDecorator }) {
    return (
      <Form.Item
        hasFeedback
      >
        {getFieldDecorator('password', {
          rules: [
            {
              required: true,
              message: 'Password is required',
            },
          ],
        })(
          <Input
            ref={el => this.passwordInput = el}
            prefix={<Icon type='lock'/>}
            type='password'
            placeholder='Password'
          />
        )}
      </Form.Item>
    )
  }

  renderSubmitButton({ isFieldTouched, getFieldsError }) {
    const { loading } = this.props

    return (
      <Form.Item>
        <Button
          styleName='submit'
          type='primary'
          htmlType='submit'
          loading={loading}
          disabled={!isFieldTouched('password') || hasErrors(getFieldsError())}
        >Log in</Button>
      </Form.Item>
    )
  }

  render() {
    const { form } = this.props

    return (
      <div>

        <div styleName='logo'>
          <img alt='logo' src={logo} />
        </div>


        <Form onSubmit={this.handleSubmit}>
          {this.renderPasswordField(form)}
          {this.renderSubmitButton(form)}
        </Form>
      </div>
    )
  }
}

export default LoginForm
