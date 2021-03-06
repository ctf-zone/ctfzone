import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { createForm } from 'rc-form';
import Recaptcha from 'react-recaptcha';

import { FormItem, Button } from '../../../components';

class SignupForm extends Component {
  static propTypes = {
    onSubmit: PropTypes.func.isRequired,
    form: PropTypes.object.isRequired,
    fields: PropTypes.object,
    className: PropTypes.string
  };

  componentDidMount() {
    this.nameInput.focus();
  }

  handleSubmit = e => {
    e.preventDefault();

    const { getFieldsValue } = this.props.form;
    this.props.onSubmit(getFieldsValue());
  };

  renderNameField({ getFieldDecorator, getFieldError }) {
    return (
      <FormItem errors={getFieldError('name')}>
        {getFieldDecorator('name', {
          initialValue: '',
          rules: [
            {
              required: true,
              message: 'Name is required'
            },
            {
              min: 3,
              message: 'Name is too short'
            }
          ]
        })(
          <input
            type="text"
            placeholder="Name"
            ref={el => (this.nameInput = el)}
          />
        )}
      </FormItem>
    );
  }

  renderEmailField({ getFieldDecorator, getFieldError }) {
    return (
      <FormItem errors={getFieldError('email')}>
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
        })(<input type="text" placeholder="Email" />)}
      </FormItem>
    );
  }

  renderPaswordField({ getFieldDecorator, getFieldError }) {
    return (
      <FormItem errors={getFieldError('password')}>
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
        })(<input type="password" placeholder="Password" />)}
      </FormItem>
    );
  }

  renderReCaptcha(verifyCallback) {
    return (
      <FormItem>
        <div style={{'paddingBottom': '10px'}}>
          <Recaptcha
            sitekey='6LeIxAcTAAAAAJcZVRqyHh71UMIEGNQ_MXjiZKhI'
            render='explicit'
            verifyCallback={verifyCallback}
            onloadCallback={() => {
              console.log('loaded');
            }}
          />
        </div>
      </FormItem>
    );
  }

  renderSubmitButton({ getFieldError, isFieldTouched }) {
    const canSubmit = ['name', 'email', 'password'].reduce((result, field) => {
      return result && isFieldTouched(field) && !getFieldError(field);
    }, true);

    return (
      <FormItem>
        <Button disabled={!canSubmit} value="Sign Up" />
      </FormItem>
    );
  }

  render() {
    const { form, className, verifyCallback } = this.props;

    return (
      <form className={className} onSubmit={this.handleSubmit}>
        {this.renderNameField(form)}
        {this.renderEmailField(form)}
        {this.renderPaswordField(form)}
        {this.renderReCaptcha(verifyCallback)}
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

export default createForm({ onFieldsChange, mapPropsToFields })(SignupForm);
