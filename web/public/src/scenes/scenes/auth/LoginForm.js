import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { createForm } from 'rc-form';

import { FormItem, Button } from '../../../components';

class LoginForm extends Component {
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
              message: 'Необходимо ввести название команды'
            }
          ]
        })(
          <input
            type="text"
            placeholder="Название команды"
            ref={el => (this.emailInput = el)}
          />
        )}
      </FormItem>
    );
  }

  renderPaswordField({ getFieldDecorator, getFieldError }) {
    return (
      <FormItem errors={getFieldError('password')} label="Password">
        {getFieldDecorator('password', {
          initialValue: '',
          rules: [
            {
              required: true,
              message: 'Поле пароль обязательное'
            }
          ]
        })(<input type="password" placeholder="Пароль" />)}
      </FormItem>
    );
  }

  renderSubmitButton({ getFieldError, isFieldTouched }) {
    const canSubmit = ['email', 'password'].reduce((result, field) => {
      return result && isFieldTouched(field) && !getFieldError(field);
    }, true);

    return (
      <FormItem>
        <Button disabled={!canSubmit} value="Вход" />
      </FormItem>
    );
  }

  render() {
    const { form, className } = this.props;

    return (
      <form className={className} onSubmit={this.handleSubmit}>
        {this.renderEmailField(form)}
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

export default createForm({ onFieldsChange, mapPropsToFields })(LoginForm);
