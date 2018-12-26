import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { createForm } from 'rc-form';

import { FormItem, Button } from '../../../components';

class SendTokenForm extends Component {
  static propTypes = {
    onSubmit: PropTypes.func.isRequired,
    form: PropTypes.object.isRequired,
    fields: PropTypes.object,
    className: PropTypes.string
  };

  componentDidMount() {
    this.emailInput.focus();
  }

  handleSubmit = e => {
    e.preventDefault();

    const { getFieldsValue } = this.props.form;
    this.props.onSubmit(getFieldsValue());
  };

  renderEmailField({ getFieldDecorator, getFieldError }) {
    return (
      <FormItem errors={getFieldError('email')} label="Email">
        {getFieldDecorator('email', {
          initialValue: '',
          rules: [
            {
              required: true,
              message: 'Email is required'
            },
            {
              type: 'email',
              message: 'Invalid email'
            }
          ]
        })(
          <input
            type="text"
            placeholder="Email"
            ref={el => (this.emailInput = el)}
          />
        )}
      </FormItem>
    );
  }

  renderSubmitButton({ getFieldError, isFieldTouched }) {
    const canSubmit = ['email'].reduce((result, field) => {
      return result && isFieldTouched(field) && !getFieldError(field);
    }, true);

    return (
      <FormItem>
        <Button disabled={!canSubmit} value="Send Token" />
      </FormItem>
    );
  }

  render() {
    const { form, className } = this.props;

    return (
      <form className={className} onSubmit={this.handleSubmit}>
        {this.renderEmailField(form)}
        {this.renderSubmitButton(form)}
      </form>
    );
  }
}

const onFieldsChange = (props, changed, all) => {
  props.onChange(all);
};

const mapPropsToFields = ({ fields }) => {
  return fields;
};

export default createForm({ onFieldsChange, mapPropsToFields })(SendTokenForm);
