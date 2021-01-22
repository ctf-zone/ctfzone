import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { withRouter, Link } from 'react-router-dom';

import { Window } from '../../../components';
import { addErrors } from '../../../utils/form';
import LoginForm from './LoginForm';

class Login extends Component {
  static propTypes = {
    authLogin: PropTypes.func,
    authSendToken: PropTypes.func
  };

  state = {
    fields: {},
    activationError: false
  };

  handleSubmit = async values => {
    const { authLogin, history } = this.props;

    try {
      await authLogin(values, { throw: true });
      history.push('/news');
    } catch (e) {
      let errors = {};

      if (e.status === 400) {
        errors = e.errors;
      } else if (e.status === 401) {
        errors = { password: [e.message] };
      } else if (e.status === 422) {
        this.setState({ activationError: true });
        return;
      }

      this.setState(({ fields }) => ({
        fields: addErrors(fields, errors)
      }));
    }
  };

  handleChange = fields => {
    this.setState({ fields });
  };

  render() {
    const { activationError } = this.state;

    return (
      <Window title="Вход" type="login">
        {activationError ? (
          <div className="ctf-form-message">
            Account is not activated.
            <br />
            <Link to="/auth/send-token/activate">Resend token</Link>
          </div>
        ) : (
          <LoginForm
            className="ctf-form ctf-form--login"
            fields={this.state.fields}
            onSubmit={this.handleSubmit}
            onChange={this.handleChange}
          />
        )}
        <div className="ctf-form-links">
          <Link to="/auth/send-token/reset">Reset Password</Link>
          <Link to="/auth/signup">Sign Up</Link>
        </div>
      </Window>
    );
  }
}

const mapDispatchToProps = dispatch => ({
  authLogin: dispatch.auth.login,
  authSendToken: dispatch.auth.sendToken
});

export default withRouter(
  connect(
    () => ({}),
    mapDispatchToProps
  )(Login)
);
