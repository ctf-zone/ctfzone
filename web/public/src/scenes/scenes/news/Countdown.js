import React, { Component } from 'react';
import PropTypes from 'prop-types';

class Countdown extends Component {
  static propTypes = {
    startTime: PropTypes.object,
    onStart: PropTypes.func
  };

  state = {
    intervalId: 0,
    currentTime: new Date().getTime()
  };

  componentDidMount() {
    const intervalId = setInterval(this.updateTime, 1000);
    this.setState({ intervalId });
  }

  componentWillUnmount() {
    clearInterval(this.state.intervalId);
  }

  updateTime = () => {
    const now = new Date().getTime();
    const { startTime } = this.props;

    if (now <= startTime) {
      this.setState({ currentTime: now });
    } else {
      this.setState({ currentTime: startTime });
      clearInterval(this.state.intervalId);
      this.props.onStart();
    }
  };

  render() {
    const { startTime } = this.props;
    const { currentTime } = this.state;

    const distance = startTime - currentTime;

    const days = Math.floor(distance / (1000 * 60 * 60 * 24));
    const hours = Math.floor(
      (distance % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60)
    );
    const minutes = Math.floor((distance % (1000 * 60 * 60)) / (1000 * 60));
    const seconds = Math.floor((distance % (1000 * 60)) / 1000);

    const prefixClass = 'ctf-countdown';

    return (
      <div className={prefixClass}>
        <div className={`${prefixClass}-message`}>Contest starts in</div>
        <span className={`${prefixClass}-days`}>{('0' + days).slice(-2)}</span>
        <span className={`${prefixClass}-hours`}>
          {('0' + hours).slice(-2)}
        </span>
        <span className={`${prefixClass}-minutes`}>
          {('0' + minutes).slice(-2)}
        </span>
        <span className={`${prefixClass}-seconds`}>
          {('0' + seconds).slice(-2)}
        </span>
      </div>
    );
  }
}

export default Countdown;
