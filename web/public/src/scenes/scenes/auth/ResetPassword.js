import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { withRouter, Link } from 'react-router-dom';

import { Window } from '../../../components';
import { addErrors } from '../../../utils/form';
import ResetPasswordForm from './ResetPasswordForm';

class ResetPassword extends Component {
  static propTypes = {
    authResetPassword: PropTypes.func,
    match: PropTypes.object,
    history: PropTypes.object
  };

  state = {
    fields: {}
  };

  handleSubmit = async values => {
    const { token } = this.props.match.params;
    const { authResetPassword, history } = this.props;

    try {
      await authResetPassword({ token, ...values }, { throw: true });
      history.push('/auth/login');
    } catch (e) {
      let errors = {};

      if (e.status === 400) {
        errors = { password: e.errors.token || e.errors.password };
      } else if (e.status === 404) {
        errors = { password: ['Token not found'] };
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
    return (
      <Window title="Reset Password" type="reset-password">
        <ResetPasswordForm
          className="ctf-form ctf-form--reset-password"
          fields={this.state.fields}
          onSubmit={this.handleSubmit}
          onChange={this.handleChange}
        />
        <div className="ctf-form-links">
          <Link to="/auth/login">Log In</Link>
        </div>
      </Window>
    );
  }
}

const mapDispatchToProps = dispatch => ({
  authResetPassword: dispatch.auth.resetPassword
});

export default withRouter(
  connect(
    () => {},
    mapDispatchToProps
  )(ResetPassword)
);
