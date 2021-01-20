import React, { Component } from 'react'
import PropTypes from 'prop-types'
import * as _ from 'lodash'
import { Link } from 'react-router-dom'
import { connect } from 'react-redux'
import {
  Table,
  Spin,
  Button,
  Row,
  Col,
  Popconfirm,
  Divider,
} from 'antd'

import {
  Pagination,
  textFilter,
  withFilters,
} from '../../../../../components'

import styles from './AnnouncementsTable.module.css'

class AnnouncementsTable extends Component {

  static propTypes = {
    // @connect
    announcements: PropTypes.object.isRequired,
    challenges: PropTypes.array,
    effects: PropTypes.object,
    challengesList: PropTypes.func,
    announcementsList: PropTypes.func.isRequired,
    announcementsDelete: PropTypes.func.isRequired,

    // @withFilters
    filters: PropTypes.object,
    handleFilterSet: PropTypes.func,
    handleFilterReset: PropTypes.func,
    handleResetFilters: PropTypes.func,
    isFiltered: PropTypes.func,
  }

  componentDidMount() {
    const { announcementsList, challengesList } = this.props
    announcementsList({})
    challengesList({})
  }

  componentDidUpdate(prevProps) {
    const { announcementsList, filters } = this.props

    if (!_.isEqual(filters, prevProps.filters)) {
      announcementsList({ filters })
    }
  }

  handleAnnouncementDelete = (announcementId) => async() => {
    const { announcementsDelete, announcementsList, filters } = this.props

    await announcementsDelete(announcementId)
    await announcementsList({ filters })
  }

  handleNextButtonClick = () => {
    const { announcementsList } = this.props
    const { next } = this.props.announcements.links

    announcementsList({ link: next.url })
  }

  handlePrevButtonClick = () => {
    const { announcementsList } = this.props
    const { prev } = this.props.announcements.links

    announcementsList({ link: prev.url })
  }

  renderControls() {
    const { extra, ...restFilters } = this.props.filters
    const isFiltered = !(_.isEmpty(restFilters) && _.isEmpty(extra))

    return (
      <div className={styles.controls}>
        <Row>
          <Col span={12} className={styles.controlsLeft}>
            <Link to='/announcements/create'>
              <Button type='primary'>Create</Button>
            </Link>
          </Col>
          <Col span={12} className={styles.controlsRight}>
            <Button
              onClick={this.props.handleResetFilters}
              disabled={!isFiltered}
            >Reset filters</Button>
          </Col>
        </Row>
      </div>
    )
  }

  renderTable() {
    const { items: announcements, links } = this.props.announcements
    const { challenges } = this.props
    const { prev, next } = links

    const columns = [
      {
        title: 'ID',
        dataIndex: 'id',
        width: '10%',
      },
      {
        title: 'Title',
        dataIndex: 'title',
        width: '20%',
        ...textFilter('title', this.props),
      },
      {
        title: 'Challenge',
        dataIndex: 'challengeId',
        width: '10%',
        render: (challengeId) => {
          const challenge = _.find(challenges, { challenge: { id: challengeId } })

          if (!challenge) {
            return ''
          }

          return challenge.challenge.title
        },
      },
      {
        title: 'Action',
        render: (data, record) => (
          <span>
            <Link to={`announcements/${record.id}/edit`}>
              <span>Edit</span>
            </Link>
            <Divider type='vertical' />
            <Popconfirm
              title={`Are you sure that you want to delete announcement "${record.title}"?`}
              okText='Yes'
              cancelText='No'
              placement='leftTop'
              onConfirm={this.handleAnnouncementDelete(record.id)}
            >
              <a href='javascript:;'>Delete</a>
            </Popconfirm>
          </span>
        ),
        width: '10%',
      },
    ]

    return (
      <div>
        <Table
          rowKey={(record) => record.id}
          bordered={true}
          columns={columns}
          dataSource={announcements}
          pagination={false}
        />
        <Pagination
          onNext={this.handleNextButtonClick}
          onPrev={this.handlePrevButtonClick}
          hasNext={!!next}
          hasPrev={!!prev}
        />
      </div>
    )
  }

  render() {
    const { loading } = this.props.effects

    return (
      <Spin spinning={loading}>
        <div>
          {this.renderControls()}
          {this.renderTable()}
        </div>
      </Spin>
    )
  }
}

const mapStateToProps = (state) => ({
    announcements: state.announcements,
    challenges: state.challenges.challenges,
    effects: state.api.effects.announcements.list,
});

const mapDispatchToProps = (dispatch) => ({
    challengesList: dispatch.challenges.list,
    announcementsList: dispatch.announcements.list,
    announcementsDelete: dispatch.announcements.delete,
});


export default withFilters(connect(mapStateToProps, mapDispatchToProps)(AnnouncementsTable))
