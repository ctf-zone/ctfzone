import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { withRouter } from 'react-router-dom';

const RequireAuth = Wrapped => {
  class Wrapper extends Component {
    static propTypes = {
      isLoggedIn: PropTypes.bool,
      history: PropTypes.object
    };

    componentDidMount() {
      const { isLoggedIn, history } = this.props;

      if (!isLoggedIn) {
        history.push('/auth/login');
      }
    }

    render() {
      return <Wrapped />;
    }
  }

  const mapStateToProps = state => ({
    isLoggedIn: state.auth
  });

  return withRouter(connect(mapStateToProps)(Wrapper));
};

export default RequireAuth;
