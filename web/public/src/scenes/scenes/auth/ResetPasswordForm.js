import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { createForm } from 'rc-form';

import { FormItem, Button } from '../../../components';

class ResetPasswordForm extends Component {
  static propTypes = {
    onSubmit: PropTypes.func.isRequired,
    form: PropTypes.object.isRequired,
    fields: PropTypes.object,
    className: PropTypes.string
  };

  componentDidMount() {
    this.passwordInput.focus();
  }

  handleSubmit = e => {
    e.preventDefault();

    const { getFieldsValue } = this.props.form;
    this.props.onSubmit(getFieldsValue());
  };

  renderPaswordField({ getFieldDecorator, getFieldError }) {
    return (
      <FormItem errors={getFieldError('password')} label="Password">
        {getFieldDecorator('password', {
          initialValue: '',
          rules: [
            {
              required: true,
              message: 'Password is required'
            },
            {
              min: 8,
              message: 'Password is too short'
            }
          ]
        })(
          <input
            type="password"
            placeholder="New password"
            ref={el => (this.passwordInput = el)}
          />
        )}
      </FormItem>
    );
  }

  renderSubmitButton({ getFieldError, isFieldTouched }) {
    const canSubmit = ['password'].reduce((result, field) => {
      return result && isFieldTouched(field) && !getFieldError(field);
    }, true);

    return (
      <FormItem>
        <Button disabled={!canSubmit} value="Reset" />
      </FormItem>
    );
  }

  render() {
    const { form, className } = this.props;

    return (
      <form className={className} onSubmit={this.handleSubmit}>
        {this.renderPaswordField(form)}
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

export default createForm({ onFieldsChange, mapPropsToFields })(
  ResetPasswordForm
);
