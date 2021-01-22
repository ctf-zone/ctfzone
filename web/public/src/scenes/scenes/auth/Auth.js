import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { withRouter, Switch, Route } from 'react-router-dom';

import { NotFound } from '../../../components';
import Login from './Login';
import Signup from './Signup';
import Logout from './Logout';
import SendToken from './SendToken';
import ResetPassword from './ResetPassword';
import Activate from './Activate';

class Auth extends Component {
  static propTypes = {
    match: PropTypes.object,
    history: PropTypes.object
  };

  componentDidMount() {
    const { match, history } = this.props;

    if (match.isExact) {
      history.push('/auth/login');
    }
  }

  render() {
    const { match } = this.props;

    return (
      <Switch>
        <Route exact path={`${match.url}/login`} component={Login} />
        <Route exact path={`${match.url}/logout`} component={Logout} />
        <Route component={NotFound} />
      </Switch>
    );
  }
}

export default withRouter(Auth);
