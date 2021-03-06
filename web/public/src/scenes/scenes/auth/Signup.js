import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { withRouter, Link } from 'react-router-dom';

import { Window } from '../../../components';
import { addErrors } from '../../../utils/form';
import SignupForm from './SignupForm';

class Signup extends Component {
  static propTypes = {
    authRegister: PropTypes.func,
    history: PropTypes.object
  };

  state = {
    fields: {},
    reCaptchaResponse: '',
    tokenIsSent: false
  };

  handleSubmit = async values => {
    const { authRegister } = this.props;
    const { reCaptchaResponse } = this.state

    try {
      console.log({...values, reCaptchaResponse});
      await authRegister({ ...values, reCaptchaResponse }, { throw: true });
      this.setState({ tokenIsSent: true });
    } catch (e) {
      if (e.status === 400 || e.status === 409) {
        this.setState(({ fields }) => ({
          fields: addErrors(fields, e.errors)
        }));
      }
    }
  };

  handleChange = fields => {
    this.setState({ fields });
  };

  handleReCaptcha = response => {
    console.log(response);
    this.setState({ reCaptchaResponse: response });
  }

  render() {
    const { tokenIsSent } = this.state;

    return (
      <Window title="Sign Up" type="signup">
        {tokenIsSent ? (
          <div className="ctf-form-message">
            Link was sent to your email. Please, check you inbox.
          </div>
        ) : (
          <SignupForm
            className="ctf-form ctf-form--signup"
            fields={this.state.fields}
            onSubmit={this.handleSubmit}
            onChange={this.handleChange}
            verifyCallback={this.handleReCaptcha}
          />
        )}
        <div className="ctf-form-links">
          <Link to="/auth/login">Log In</Link>
        </div>
      </Window>
    );
  }
}

const mapStateToProps = state => ({
  ...state.api.effects.auth.register
});

const mapDispatchToProps = dispatch => ({
  authRegister: dispatch.auth.register
});

export default withRouter(
  connect(
    mapStateToProps,
    mapDispatchToProps
  )(Signup)
);
