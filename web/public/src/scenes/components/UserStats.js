import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';

class UserStats extends Component {
  static propTypes = {
    userGetStats: PropTypes.func,
    userStats: PropTypes.object
  };

  state = {
    intervalId: 0
  };

  componentDidMount() {
    this.fetchData();

    const intervalId = setInterval(this.fetchData, 30000);
    this.setState({ intervalId });
  }

  componentWillUnmount() {
    clearInterval(this.state.intervalId);
  }

  fetchData = () => {
    const { userGetStats } = this.props;
    userGetStats();
  };

  render() {
    const { userStats } = this.props;
    const prefixClass = 'ctf-user-stats';
    return (
      <div className={prefixClass}>
        <div className={`${prefixClass}-rank`}>{userStats.rank}</div>
        <div className={`${prefixClass}-score`}>{userStats.score}</div>
      </div>
    );
  }
}

const mapStateToProps = state => ({
  userStats: state.user.stats
});

const mapDispatchToProps = dispatch => ({
  userGetStats: dispatch.user.getStats
});

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(UserStats);
