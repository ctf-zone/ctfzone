import React, { Component } from 'react'
import PropTypes from 'prop-types'
import { Form, Input, Icon, Switch, Button, Checkbox } from 'antd'

import { hasErrors, mapPropsToFields } from '~/utils/form'

@Form.create({
  onFieldsChange: (props, changedFields) => {
    props.onChange(changedFields)
  },
  mapPropsToFields,
})
class UserForm extends Component {

  static propTypes = {
    isEdit: PropTypes.bool.isRequired,
    onSubmit: PropTypes.func.isRequired,
    onChange: PropTypes.func.isRequired,
    fields: PropTypes.object.isRequired,

    // @form
    form: PropTypes.object,
  }

  formItemLayout = {
    labelCol: {
      sm: { span: 6 },
    },
    wrapperCol: {
      sm: { span: 14 },
    },
  }

  handleSubmit = e => {
    e.preventDefault()

    const { isEdit } = this.props
    const { getFieldsValue } = this.props.form
    const user = getFieldsValue()

    if (user.extra && user.extra.length !== 0) {
      user.extra = JSON.parse(user.extra)
    } else {
      user.extra = {}
    }

    if (isEdit && !user.updatePassword) {
      delete user.password
    }

    delete user.updatePassword

    this.props.onSubmit(user)
  }

  renderNameField({ getFieldDecorator }) {
    return (
      <Form.Item
        hasFeedback
        label='Name'
        {...this.formItemLayout}
      >
        {getFieldDecorator('name', {
          rules: [
            {
              required: true,
              message: 'Name is required!',
            },
          ],
        })(
          <Input
            prefix={<Icon type='user' />}
          />
        )}
      </Form.Item>
    )
  }

  renderEmailField({ getFieldDecorator }) {
    return (
      <Form.Item
        hasFeedback
        label='Email'
        {...this.formItemLayout}
      >
        {getFieldDecorator('email', {
          rules: [
            {
              required: true,
              message: 'Email is required!',
            },
            {
              type: 'email',
              message: 'The input is not valid email',
            },
          ],
        })(
          <Input
            prefix={<Icon type='mail' />}
          />
        )}
      </Form.Item>
    )
  }

  renderUpdatePasswordField({ getFieldDecorator }) {
    return (
      <Form.Item
        label='Update password'
        {...this.formItemLayout}
      >
        {getFieldDecorator('updatePassword')(
          <Checkbox onChange={() => {
            const { resetFields } = this.props.form
            resetFields('password')
          }} />
        )}
      </Form.Item>
    )
  }

  renderPasswordField({ getFieldDecorator, getFieldValue }) {
    const { isEdit } = this.props
    const updatePassword = getFieldValue('updatePassword')

    let rules = [
      {
        min: 8,
        message: 'Minimum password length is 8!',
      },
    ]

    if (!isEdit || updatePassword) {
      rules.push(
        {
          required: true,
          message: 'Password is required!',
        },
      )
    }

    return (
      <Form.Item
        hasFeedback
        label='Password'
        {...this.formItemLayout}
      >
        {getFieldDecorator('password', {
          rules,
        })(
          <Input
            prefix={<Icon type='lock' />}
            type='password'
            disabled={isEdit && !updatePassword}
          />
        )}
      </Form.Item>
    )
  }

  renderIsActivatedField({ getFieldDecorator, getFieldValue }) {
    return (
      <Form.Item
        label='Activated'
        {...this.formItemLayout}
      >
        {getFieldDecorator('isActivated', {
          initialValue: false,
          rules: [],
        })(
          <Switch checked={getFieldValue('isActivated')} />
        )}
      </Form.Item>
    )
  }

  renderExtraField({ getFieldDecorator }) {
    const validateExtra = (rule, value, callback) => {
      if (value && value.length !== 0) {
        try {
          JSON.parse(value)
        } catch (err) {
          callback(err)
        }
      }
      callback()
    }

    return (
      <Form.Item
        hasFeedback
        label='Extra'
        {...this.formItemLayout}
      >
        {getFieldDecorator('extra', {
          rules: [
            {
              validator: validateExtra,
            },
          ],
        })(
          <Input.TextArea rows={4} />
        )}
      </Form.Item>
    )
  }

  renderSubmitButton({ isFieldTouched, getFieldValue, getFieldsError }) {
    const { isEdit } = this.props
    let required = [ 'name', 'email' ]

    if (!isEdit) {
      required.push('password')
    }

    let requiredIsTouched = true

    for (let field of required) {
      requiredIsTouched = requiredIsTouched &&
        (isFieldTouched(field) || getFieldValue(field))
    }

    return (
      <Form.Item
        wrapperCol={{ span: 12, offset: 6 }}
      >
        <Button
          type='primary'
          htmlType='submit'
          disabled={!requiredIsTouched || hasErrors(getFieldsError())}
        >Save</Button>
      </Form.Item>
    )
  }

  render() {
    const { form, isEdit } = this.props

    return (
      <Form
        onSubmit={this.handleSubmit}
      >
        {this.renderNameField(form)}
        {this.renderEmailField(form)}
        {
          isEdit
            ? this.renderUpdatePasswordField(form)
            : null
        }
        {this.renderPasswordField(form)}
        {this.renderIsActivatedField(form)}
        {this.renderExtraField(form)}
        {this.renderSubmitButton(form)}
      </Form>
    )
  }
}

export default UserForm
