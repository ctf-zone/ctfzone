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
class FlagForm extends Component {

  static propTypes = {
    onSubmit: PropTypes.func.isRequired,
    form: PropTypes.object.isRequired,
    fields: PropTypes.object,
    className: PropTypes.string,
  }

  handleSubmit = (e) => {
    e.preventDefault()

    const { getFieldsValue } = this.props.form
    this.props.onSubmit(getFieldsValue())
  }

  renderFlagField({ getFieldDecorator, getFieldError }) {
    return (
      <FormItem
        errors={getFieldError('flag')}
      >
        {getFieldDecorator('flag', {
          initialValue: '',
          rules: [
            {
              required: true,
              message: 'Flag is required',
            },
          ],
        })(
          <input
            type="text"
            placeholder="Flag"
          />
        )}
      </FormItem>
    )
  }

  renderSubmitButton({ getFieldError, isFieldTouched }) {
    const canSubmit = ['flag'].reduce((result, field) => {
      return result && isFieldTouched(field) && !getFieldError(field)
    }, true);

    return (
      <FormItem>
        <Button
          disabled={!canSubmit}
          value="Submit"
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
        {this.renderFlagField(form)}
        {this.renderSubmitButton(form)}
      </form>
    )
  }
}

export default FlagForm
