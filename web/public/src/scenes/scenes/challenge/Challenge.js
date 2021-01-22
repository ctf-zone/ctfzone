import React, { Component } from 'react';
import PropTypes from 'prop-types';
import classNames from 'classnames';
import { connect } from 'react-redux';
import { withRouter } from 'react-router-dom';
import ReactMarkdown from 'react-markdown';

import { addErrors } from '../../../utils/form';
import { Page, Loading, Divider } from '../../../components';
import FlagForm from './FlagForm';

class Challenge extends Component {
  static propTypes = {
    challenges: PropTypes.object,
    challengesGet: PropTypes.func,
    hints: PropTypes.object,
    hintsGet: PropTypes.func,
    userGetStats: PropTypes.func,
    match: PropTypes.object
  };

  state = {
    fields: {}
  };

  componentDidMount() {
    const { challengesGet, hintsGet } = this.props;
    const { id } = this.props.match.params;
    challengesGet(id);
    hintsGet(id);
  }

  handleChange = fields => {
    this.setState({ fields });
  };

  handleSubmit = async ({ flag }) => {
    const { createSolution, challengesGet, userGetStats } = this.props;
    const { id } = this.props.match.params;

    try {
      await createSolution({ id, flag }, { throw: true });
      await challengesGet(id);
      await userGetStats();
    } catch (e) {
      let errors = {};

      if (e.status === 400 || e.status === 409) {
        errors = e.errors || {};
      } else if (e.status === 418) {
        errors = { flag: ['Неправильный флаг'] };
      } else if (e.status === 403) {
        errors = { flag: [e.message] };
      }

      this.setState(({ fields }) => ({
        fields: addErrors(fields, errors)
      }));
    }
  };

  renderHints() {
    const { hints: allHints } = this.props;
    const { id } = this.props.match.params;

    if (!(id in allHints)) {
      return null;
    }

    const hints = allHints[id];

    return (
      <div>
        <Divider />
        <h2>Hints</h2>
        <ul>
          {hints.map((hint, i) => (
            <li key={i}>
              <ReactMarkdown source={hint.body} />
            </li>
          ))}
        </ul>
      </div>
    );
  }

  renderChallenge({ challenge, user, meta }) {
    const prefixClass = 'ctf-challenge';
    const className = classNames({
      [prefixClass]: true,
      [`${prefixClass}--solved`]: user.isSolved
    });
    const solved = user.isSolved ? 'ctf-page--challenge--solved' : '';
    const difficultyMap = {
      easy: "простое",
      medium: "среднее",
      hard: "сложное",
    }

    return (
      <Page type={`challenge ${solved}`} title={challenge.title}>
        <div className={className}>
          <div className={`${prefixClass}-info`}>
            <div className={`${prefixClass}-categories`}>
              {challenge.categories.map((category, i) => {
                return (
                  <div key={i} className={`${prefixClass}-category`}>
                    {category}
                  </div>
                );
              })}
            </div>
            <div className={`${prefixClass}-difficulty`}>
              {difficultyMap[challenge.difficulty]}
            </div>
            <div className={`${prefixClass}-solutions`}>
              {meta.solutionsCount}
            </div>
            <div className={`${prefixClass}-points`}>{challenge.points}</div>
          </div>
          <div className={`${prefixClass}-description`}>
            <ReactMarkdown source={challenge.description} />
          </div>
          {user.isSolved ? (
            <div className={`${prefixClass}-solved-message`}>
              Challenge solved.
            </div>
          ) : (
            <div className={`${prefixClass}-flag`}>
              <FlagForm
                onChange={this.handleChange}
                onSubmit={this.handleSubmit}
                fields={this.state.fields}
              />
            </div>
          )}
          {meta.hintsCount > 0 ? this.renderHints() : null}
        </div>
      </Page>
    );
  }

  render() {
    const { items: challenges } = this.props.challenges;
    const { id } = this.props.match.params;
    const isLoaded = id in challenges;
    const challenge = challenges[id];

    return (
      <Loading loading={!isLoaded}>
        {challenge ? this.renderChallenge(challenge) : null}
      </Loading>
    );
  }
}

const mapStateToProps = state => ({
  challenges: state.challenges,
  hints: state.news.hints
});

const mapDispatchToProps = dispatch => ({
  challengesGet: dispatch.challenges.get,
  createSolution: dispatch.user.createSolution,
  hintsGet: dispatch.news.getHints,
  userGetStats: dispatch.user.getStats
});

export default withRouter(
  connect(
    mapStateToProps,
    mapDispatchToProps
  )(Challenge)
);
