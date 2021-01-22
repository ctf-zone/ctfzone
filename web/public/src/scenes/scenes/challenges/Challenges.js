import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { Link } from 'react-router-dom';
import classNames from 'classnames';

import { Page, Divider } from '../../../components';

class Challenges extends Component {
  static propTypes = {
    challenges: PropTypes.object,
    challengesList: PropTypes.func
  };

  componentDidMount() {
    this.fetchData();
    this.intervalId = setInterval(this.fetchData, 30000);
  }

  componentWillUnmount() {
    clearInterval(this.intervalId);
  }

  fetchData = () => {
    const { challengesList } = this.props;
    challengesList();
  };

  renderChallenge({ challenge, user, meta }, index) {
    const prefixClass = 'ctf-challenge-list-item';

    const className = classNames({
      [prefixClass]: true,
      [prefixClass + '--solved']: user.isSolved
    });

    const difficultyMap = {
      easy: "простое",
      medium: "среднее",
      hard: "сложное",
    }

    return (
      <Link key={index} to={`/challenges/${challenge.id}`}>
        <div className={className}>
          <div className={`${prefixClass}-icon`} />
          <div className={`${prefixClass}-title`}>{challenge.title}</div>
          <div
            className={`${prefixClass}-difficulty ${prefixClass}-difficulty--${
              challenge.difficulty
            }`}
          >
            {difficultyMap[challenge.difficulty]}
          </div>
          <div className={`${prefixClass}-solutions`}>
            {meta.solutionsCount}
          </div>
          <div className={`${prefixClass}-points`}>{challenge.points}</div>
          <div className={`${prefixClass}-categories`}>
            {challenge.categories.map((category, i) => {
              return (
                <div key={i} className={`${prefixClass}-category`}>
                  {category}
                </div>
              );
            })}
          </div>
          {meta.hintsCount > 0 ? (
            <div className={`${prefixClass}-hints`} />
          ) : (
            ''
          )}
        </div>
        <Divider />
      </Link>
    );
  }

  render() {
    const { items: challenges } = this.props.challenges;

    return (
      <Page title="Задания" type="challenges">
        {Object.values(challenges).map(this.renderChallenge)}
      </Page>
    );
  }
}

const mapStateToProps = state => ({
  challenges: state.challenges,
  effects: state.api.effects.challenges.list
});

const mapDispatchToProps = dispatch => ({
  challengesList: dispatch.challenges.list
});

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(Challenges);
