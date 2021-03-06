import React, { Component } from 'react'
import PropTypes from 'prop-types'
import * as _ from 'lodash'
import { Row, Col } from 'antd'
import { connect } from 'react-redux'
import { withRouter } from 'react-router-dom'

import { addValues, addErrors } from '../../../../../utils/form'

import { ChallengeForm } from './components'

class ChallengeCreateEdit extends Component {

  static propTypes = {
    // @connect
    challenge: PropTypes.object.isRequired,
    game: PropTypes.object.isRequired,
    gameGet: PropTypes.func.isRequired,
    challengesGet: PropTypes.func.isRequired,
    challengesCreate: PropTypes.func.isRequired,
    challengesUpdate: PropTypes.func.isRequired,

    // @withRouter
    match: PropTypes.object.isRequired,
    history: PropTypes.object.isRequired,
  }

  state = {
    fields: {},
  }

  async componentDidMount() {
    const { game, gameGet } = this.props

    if (_.isEmpty(game)) {
      await gameGet()
    }

    const { challengeId } = this.props.match.params
    const { challengesGet } = this.props

    if (challengeId) {
      const challenge = await challengesGet(challengeId)
      this.setState(({ fields }) => ({
        fields: addValues(fields, challenge.challenge),
      }))
    }
  }

  static defaultState = {
    fields: {},
  }

  handleChange = (fields) => {
    this.setState({ fields })
  }

  handleSubmit = async(challenge) => {
    const { challengesCreate, challengesUpdate, history } = this.props
    const { challengeId } = this.props.match.params

    try {
      if (challengeId) {
        await challengesUpdate(
          {
            ...challenge,
            id: parseInt(challengeId),
          },
          { throw: true },
        )
      } else {
        await challengesCreate(challenge, { throw: true })
      }

      history.push('/challenges')
    } catch (err) {
      this.setState(({ fields }) => ({
        fields: addErrors(fields, err.errors),
      }))
    }
  }

  render() {
    const { challenge, game } = this.props
    const isEdit = !!this.props.match.params.challengeId
    const { fields } = this.state

    return (
      <div>
        <Row>
          <Col offset={6}>
            <h1>
              {
                isEdit && !_.isEmpty(challenge) && !_.isEmpty(challenge.challenge)
                  ? `Edit challenge "${challenge.challenge.title}"`
                  : 'Create new challenge'
              }
            </h1>
          </Col>
        </Row>
        <ChallengeForm
          isEdit={isEdit}
          onSubmit={this.handleSubmit}
          onChange={this.handleChange}
          fields={fields}
          categories={game.categories || []}
        />
      </div>
    )
  }
}

const mapStateToProps = (state) => ({
  challenge: state.challenges.challenge,
  game: state.game,
});

const mapDispatchToProps = (dispatch) => ({
  challengesCreate: dispatch.challenges.create,
  challengesUpdate: dispatch.challenges.update,
  challengesGet: dispatch.challenges.get,
  gameGet: dispatch.game.get,
});


export default withRouter(connect(mapStateToProps, mapDispatchToProps)(ChallengeCreateEdit))
