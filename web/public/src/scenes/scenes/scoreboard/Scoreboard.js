import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';

import dayjs from '../../../utils/date';
import { Page } from '../../../components';

class Scoreboard extends Component {
  static propTypes = {
    scores: PropTypes.object,
    scoresList: PropTypes.func
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
    const { scoresList } = this.props;
    scoresList();
  };

  render() {
    const { items: scores } = this.props.scores;

    return (
      <Page title="Таблица результатов" type="scoreboard">
        <table className="ctf-scores">
          <thead>
            <tr>
              <th>Позиция</th>
              <th>Имя</th>
              <th>Очки</th>
              <th>Последняя сдача флага</th>
            </tr>
          </thead>
          <tbody>
            {scores.map((score, i) => {
              const updatedAt = dayjs(score.updatedAt);
              return (
                <tr key={i}>
                  <td>{score.rank}</td>
                  <td>{score.user.name}</td>
                  <td>{score.score}</td>
                  <td>
                    {score.updatedAt
                      ? `${updatedAt.format(
                          'HH:mm DD.MM.YYYY'
                        )} (${updatedAt.fromNow()})`
                      : 'N/A'}
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </Page>
    );
  }
}

const mapStateToProps = state => ({
  scores: state.scores
});

const mapDispatchToProps = dispatch => ({
  scoresList: dispatch.scores.list
});

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(Scoreboard);
