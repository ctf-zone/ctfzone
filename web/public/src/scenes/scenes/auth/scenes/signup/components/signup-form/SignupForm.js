import React, { Component } from 'react'
import PropTypes from 'prop-types'
import { createForm } from 'rc-form';

import { FormItem, Button } from '~/components'

@createForm({
  onFieldsChange: (props, changed, all) => {
    props.onChange(all)
  },
  mapPropsToFields: ({ fields }) => {
    return fields
  },
})
class SignupForm extends Component {

  static propTypes = {
    onSubmit: PropTypes.func.isRequired,
    form: PropTypes.object.isRequired,
    fields: PropTypes.object,
    className: PropTypes.string,
  }

  componentDidMount() {
    this.nameInput.focus()
  }

  handleSubmit = (e) => {
    e.preventDefault()

    const { getFieldsValue } = this.props.form
    this.props.onSubmit(getFieldsValue())
  }

  renderNameField({ getFieldDecorator, getFieldError }) {
    return (
      <FormItem
        errors={getFieldError('name')}
      >
        {getFieldDecorator('name', {
          initialValue: '',
          rules: [
            {
              required: true,
              message: 'Name is required',
            },
            {
              min: 3,
              message: 'Name is too short',
            },
          ],
        })(
          <input
            type='text'
            placeholder='Name'
            ref={(el) => this.nameInput = el}
          />
        )}
      </FormItem>
    )
  }

  renderEmailField({ getFieldDecorator, getFieldError }) {
    return (
      <FormItem
        errors={getFieldError('email')}
      >
        {getFieldDecorator('email', {
          initialValue: '',
          rules: [
            {
              required: true,
              message: 'Email is required',
            },
            {
              type: 'email',
              message: 'Invalid email',
            },
          ],
        })(
          <input
            type='text'
            placeholder='Email'
          />
        )}
      </FormItem>
    )
  }

  renderPaswordField({ getFieldDecorator, getFieldError }) {
    return (
      <FormItem
        errors={getFieldError('password')}
      >
        {getFieldDecorator('password', {
          initialValue: '',
          rules: [
            {
              required: true,
              message: 'Password is required',
            },
            {
              min: 8,
              message: 'Password is too short',
            },
          ],
        })(
          <input
            type='password'
            placeholder='Password'
          />
        )}
      </FormItem>
    )
  }

  renderSubmitButton({ getFieldError, isFieldTouched }) {
    const canSubmit = ['name', 'email', 'password'].reduce((result, field) => {
      return result && isFieldTouched(field) && !getFieldError(field)
    }, true);

    return (
      <FormItem>
        <Button
          disabled={!canSubmit}
          value='Sign Up'
        />
      </FormItem>
    )
  }

  render() {
    const { form, className } = this.props

    return (
      <form
        className={className}
        onSubmit={this.handleSubmit}
      >
        {this.renderNameField(form)}
        {this.renderEmailField(form)}
        {this.renderPaswordField(form)}
        {this.renderSubmitButton(form)}
      </form>
    )
  }
}

export default SignupForm
