import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { withRouter } from 'react-router-dom';

class Logout extends Component {
  static propTypes = {
    authLogout: PropTypes.func,
    history: PropTypes.object
  };

  componentDidMount() {
    const { authLogout, history } = this.props;

    authLogout();
    history.push('/auth/login');
  }

  render() {
    return <div />;
  }
}

const mapStateToProps = state => ({
  ...state.api.effects.auth.logout
});

const mapDispatchToProps = dispatch => ({
  authLogout: dispatch.auth.logout
});

export default withRouter(
  connect(
    mapStateToProps,
    mapDispatchToProps
  )(Logout)
);
