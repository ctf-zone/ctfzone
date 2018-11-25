import React, { Component } from 'react'
import PropTypes from 'prop-types'
import { connect } from 'react-redux'
import { Link } from 'react-router-dom'
import classNames from 'classnames'

import { Page, Divider } from '~/components'

@connect(
  (state) => ({
    challenges: state.challenges,
    effects: state.api.effects.challenges.list,
  }),
  (dispatch) => ({
    challengesList: dispatch.challenges.list,
  })
)
class Challenges extends Component {

  static propTypes = {
    challenges: PropTypes.object,
    challengesList: PropTypes.func,
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
    const { challengesList } = this.props
    challengesList()
  }

  renderChallenge({ challenge, user, meta }, index) {
    const prefixClass = 'ctf-challenge-list-item'

    const className = classNames({
      [prefixClass]: true,
      [prefixClass + '--solved']: user.isSolved,
    })

    return (
      <Link
        key={index}
        to={`/challenges/${challenge.id}`}
      >
        <div className={className}>
          <div className={`${prefixClass}-icon`}></div>
          <div className={`${prefixClass}-title`}>{challenge.title}</div>
          <div className={`${prefixClass}-difficulty ${prefixClass}-difficulty--${challenge.difficulty}`}>{challenge.difficulty}</div>
          <div className={`${prefixClass}-solutions`}>{meta.solutionsCount}</div>
          <div className={`${prefixClass}-points`}>{challenge.points}</div>
          <div className={`${prefixClass}-categories`}>
            {challenge.categories.map((category, i) => {
              return <div key={i} className={`${prefixClass}-category`}>{category}</div>
            })}
          </div>
          {meta.hintsCount > 0
              ? <div className={`${prefixClass}-hints`}></div>
              : ''
          }
        </div>
        <Divider />
      </Link>
    )
  }

  render() {
    const { items: challenges } = this.props.challenges

    return (
      <Page
        title='Challenges'
        type='challenges'
      >
        {Object
            .values(challenges)
            .map(this.renderChallenge)}
      </Page>
    )
  }
}

export default Challenges
