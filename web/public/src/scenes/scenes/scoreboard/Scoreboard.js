import React, { Component } from 'react'
import PropTypes from 'prop-types'
import { connect } from 'react-redux'

import dayjs from '~/utils/date'
import { Page } from '~/components'

@connect(
  (state) => ({
    scores: state.scores,
  }),
  (dispatch) => ({
    scoresList: dispatch.scores.list,
  })
)
class Scoreboard extends Component {

  static propTypes = {
    scores: PropTypes.object,
    scoresList: PropTypes.func,
  }

  state = {
    intervalId: 0,
  }

  componentDidMount() {
    this.fetchData()
    const intervalId = setInterval(this.fetchData, 30000)
    this.setState({ intervalId })
  }

  componentWillUnmount() {
    clearInterval(this.state.intervalId)
  }

  fetchData = () => {
    const { scoresList } = this.props
    scoresList()
  }

  render() {
    const { items: scores } = this.props.scores

    return (
      <Page
        title='Scoreboard'
        type='scoreboard'
      >
        <table className='ctf-scores'>
          <thead>
            <tr>
              <th>Rank</th>
              <th>Team</th>
              <th>Score</th>
              <th>Last flag submit</th>
            </tr>
          </thead>
          <tbody>
            {scores.map((score, i) => {
              const updatedAt = dayjs(score.updatedAt)
              return (
                <tr key={i}>
                  <td>{score.rank}</td>
                  <td>{score.user.name}</td>
                  <td>{score.score}</td>
                  <td>{score.updatedAt ? `${updatedAt.format('HH:mm DD.MM.YYYY')} (${updatedAt.fromNow()})` : 'N/A'}</td>
                </tr>
              )
            })}
          </tbody>
        </table>
      </Page>
    )
  }
}

export default Scoreboard
