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
  Tag,
  Popconfirm,
  Divider,
  Icon,
} from 'antd'

import {
  Pagination,
  textFilter,
  radioFilter,
  selectFilter,
  withFilters,
} from '../../../../../components'

import styles from './ChallengesTable.module.css'

class ChallengesTable extends Component {

  static propTypes = {
    // @connect
    challenges: PropTypes.object.isRequired,
    game: PropTypes.object.isRequired,
    challengesList: PropTypes.func.isRequired,
    challengesDelete: PropTypes.func.isRequired,
    challengesListResult: PropTypes.object.isRequired,
    gameGet: PropTypes.func.isRequired,

    // @withFilters
    filters: PropTypes.object,
    handleFilterSet: PropTypes.func,
    handleFilterReset: PropTypes.func,
    handleResetFilters: PropTypes.func,
    isFiltered: PropTypes.func,
  }

  componentDidMount() {
    const { challengesList, gameGet } = this.props

    challengesList({})
    gameGet()
  }

  componentDidUpdate(prevProps) {
    const { challengesList, filters } = this.props

    if (!_.isEqual(filters, prevProps.filters)) {
      challengesList({ filters })
    }
  }

  handleChallengeDelete = (challengeId) => async() => {
    const { challengesDelete, challengesList, filters } = this.props

    await challengesDelete(challengeId)
    await challengesList({ filters })
  }

  handleNextButtonClick = () => {
    const { challengesList } = this.props
    const { next } = this.props.challenges.links

    challengesList({ link: next.url })
  }

  handlePrevButtonClick = () => {
    const { challengesList } = this.props
    const { prev } = this.props.challenges.links

    challengesList({ link: prev.url })
  }

  renderControls() {
    const { extra, ...restFilters } = this.props.filters
    const isFiltered = !(_.isEmpty(restFilters) && _.isEmpty(extra))

    return (
      <div className={styles.controls}>
        <Row>
          <Col span={12} className={styles.controlsLeft}>
            <Link to='/challenges/create'>
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
    const { challenges, links } = this.props.challenges
    const { prev, next } = links
    const { game } = this.props
    const categories = game.categories || []

    const columns = [
      {
        title: 'ID',
        dataIndex: 'challenge.id',
        width: '10%',
      },
      {
        title: 'Title',
        dataIndex: 'challenge.title',
        width: '20%',
        ...textFilter('title', this.props),
      },
      {
        title: 'Categories',
        dataIndex: 'challenge.categories',
        width: '25%',
        render: (data) => data.map((category, i) => <Tag color='blue' key={i}>{category}</Tag>),
        ...selectFilter(
          'categories',
          categories.map((category) => ({ data: category, value: category })),
          this.props,
        ),
      },
      {
        title: 'Points',
        dataIndex: 'challenge.points',
        width: '20%',
        ...textFilter('points', this.props),
      },
      {
        title: 'Locked',
        dataIndex: 'challenge.isLocked',
        width: '10%',
        render: (data) => {
          return (data
            ? <Icon type='lock' />
            : ''
          )
        },
        ...radioFilter(
          'isLocked',
          [
            {
              data: 'Yes',
              value: true,
            },
            {
              data: 'No',
              value: false,
            },
          ],
          this.props,
        ),
      },
      {
        title: 'Action',
        render: (data, record) => (
          <span>
            <Link to={`challenges/${record.challenge.id}/edit`}>
              <span>Edit</span>
            </Link>
            <Divider type='vertical' />
            <Popconfirm
              title={`Are you sure that you want to delete challenge "${record.challenge.title}"?`}
              okText='Yes'
              cancelText='No'
              placement='leftTop'
              onConfirm={this.handleChallengeDelete(record.challenge.id)}
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
          rowKey={(record) => record.challenge.id}
          bordered={true}
          columns={columns}
          dataSource={challenges}
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
    const { loading } = this.props.challengesListResult

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
  challenges: state.challenges,
  game: state.game,
  challengesListResult: state.api.effects.challenges.list,
  challengesDeleteResult: state.api.effects.challenges.delete,
});

const mapDispatchToProps = (dispatch) => ({
  challengesList: dispatch.challenges.list,
  challengesDelete: dispatch.challenges.delete,
  gameGet: dispatch.game.get,
});


export default withFilters(connect(mapStateToProps, mapDispatchToProps)(ChallengesTable))
