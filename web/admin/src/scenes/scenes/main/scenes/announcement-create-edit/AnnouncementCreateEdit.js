import React, { Component } from 'react'
import PropTypes from 'prop-types'
import * as _ from 'lodash'
import { Row, Col } from 'antd'
import { connect } from 'react-redux'
import { withRouter } from 'react-router-dom'

import { addValues, addErrors } from '../../../../../utils/form'

import { AnnouncementForm } from './components'

class ChallengeCreateEdit extends Component {

  static propTypes = {
    // @connect
    announcement: PropTypes.object.isRequired,
    announcementsCreate: PropTypes.func,
    announcementsUpdate: PropTypes.func,
    announcementsGet: PropTypes.func,
    challenges: PropTypes.array,

    // @withRouter
    match: PropTypes.object.isRequired,
  }

  static defaultState = {
    fields: {},
  }

  state = {
    fields: {},
  }

  async componentDidMount() {
    const { announcementId } = this.props.match.params
    const { announcementsGet } = this.props

    if (announcementId) {
      const announcement = await announcementsGet(announcementId)

      if (!announcement) {
        return
      }

      this.setState(({ fields }) => ({
        fields: addValues(fields, {
          ...announcement,
          challengeId: announcement.challengeId
            ? announcement.challengeId.toString()
            : null,
        }),
      }))
    }

    const { challengesList } = this.props

    await challengesList({})
  }

  handleChange = (changedFields) => {
    this.setState(({ fields }) => ({
      fields: { ...fields, ...changedFields },
    }))
  }

  handleSubmit = async(announcement) => {
    const { announcementsCreate, announcementsUpdate, history } = this.props
    const { announcementId } = this.props.match.params

    if (!announcement.challengeId) {
      delete announcement.challengeId
    }

    try {
      if (announcementId) {
        await announcementsUpdate(
          {
            ...announcement,
            id: parseInt(announcementId),
          },
          { throw: true },
        )
      } else {
        await announcementsCreate(announcement, { throw: true })
      }

      history.push('/announcements')
    } catch (e) {
      this.setState(({ fields }) => ({
        fields: addErrors(fields, e.errors),
      }))
    }
  }

  render() {
    const { announcement, challenges } = this.props
    const isEdit = !!this.props.match.params.announcementId
    const { fields } = this.state

    return (
      <div>
        <Row>
          <Col offset={6}>
            <h1>
              {
                isEdit && !_.isEmpty(announcement)
                  ? 'Edit announcement'
                  : 'Create new announcement'
              }
            </h1>
          </Col>
        </Row>
        <AnnouncementForm
          isEdit={isEdit}
          onSubmit={this.handleSubmit}
          onChange={this.handleChange}
          fields={fields}
          challenges={challenges.map((v) => ({
            title: v.challenge.title,
            value: v.challenge.id,
          }))}
        />
      </div>
    )
  }
}

const mapStateToProps = state => ({
  announcement: state.announcements.item,
  challenges: state.challenges.challenges,
});

const mapDispatchToProps = dispatch => ({
  announcementsCreate: dispatch.announcements.create,
  announcementsUpdate: dispatch.announcements.update,
  announcementsGet: dispatch.announcements.get,
  challengesList: dispatch.challenges.list,
});

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(ChallengeCreateEdit))
