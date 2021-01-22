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
  Upload,
  message,
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
    challengesCreate: PropTypes.func.isRequired,
    challengesUpdate: PropTypes.func.isRequired,
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

  state = {
    selectedRowKeys: [],
  };

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

  handleSelectChange = (selectedRowKeys) => {
    this.setState({ selectedRowKeys });
  }

  handleUnlockClick = () => {
    const { challengesUpdate, challengesList, filters } = this.props;
    const { challenges } = this.props.challenges;
    const { selectedRowKeys } = this.state;

    challenges.forEach(async({ challenge }) => {
      if (selectedRowKeys.includes(challenge.id)) {
        try {
          const data = { ...challenge, isLocked: false };
          delete data.createdAt;
          delete data.updatedAt;
          await challengesUpdate(data, { throw: true });
          await challengesList({ filters });
        } catch (e) {
          message.error(e.message);
        }
      }
    })

    this.setState({ selectedRowKeys: [] })
  }

  handleLockClick = () => {
    const { challengesUpdate, challengesList, filters } = this.props;
    const { challenges } = this.props.challenges;
    const { selectedRowKeys } = this.state;

    challenges.forEach(async({ challenge }) => {
      if (selectedRowKeys.includes(challenge.id)) {
        try {
          const data = { ...challenge, isLocked: true };
          delete data.createdAt;
          delete data.updatedAt;
          await challengesUpdate(data, { throw: true });
          await challengesList({ filters });
        } catch (e) {
          message.error(e.message);
        }
      }
    })

    this.setState({ selectedRowKeys: [] })
  }


  renderControls() {
    const { extra, ...restFilters } = this.props.filters
    const { selectedRowKeys } = this.state;
    const isFiltered = !(_.isEmpty(restFilters) && _.isEmpty(extra))

    return (
      <div className={styles.controls}>
        <Row>
          <Col span={12} className={styles.controlsLeft}>
            <Link to='/challenges/create'>
              <Button type='primary'>Create</Button>
            </Link>
            {selectedRowKeys.length > 0 ?
              (
                <Button style={{ marginLeft: '10px' }} onClick={this.handleUnlockClick}>
                  <Icon type="primary" /> Unlock
                </Button>
              ) :
              ''
            }
            {selectedRowKeys.length > 0 ?
              (
                <Button style={{ marginLeft: '10px' }} onClick={this.handleLockClick}>
                  <Icon type="primary" /> Lock
                </Button>
              ) :
              ''
            }
            <Upload
              className={styles.upload}
              showUploadList={false}
              customRequest={(e) => {
                let reader = new FileReader();
                const { challengesCreate, challengesList, filters } = this.props;

                reader.onload = async(e) => {
                  const json = e.target.result;
                  try {
                    const challenges = JSON.parse(json);
                    challenges.forEach(async(challenge) => {
                      try {
                        await challengesCreate(challenge, { throw: true });
                      } catch (e) {
                        const message = e.message;
                        message.error(message);
                      }
                    });
                    await challengesList({ filters });
                  } catch (e) {
                    message.error("Invalid JSON");
                  }
                };
                reader.readAsText(e.file, "UTF-8");
              }}
            >
              <Button>
                <Icon type="upload" /> Import
              </Button>
            </Upload>
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
        width: '21%',
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
        width: '7%',
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
        title: 'Depends On',
        dataIndex: 'challenge.dependsOn',
        width: '7%',
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
              <a href='#'>Delete</a>
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
          rowSelection={{ selectedRowKeys: this.state.selectedRowKeys, onChange: this.handleSelectChange }}
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
  challengesCreate: dispatch.challenges.create,
  challengesUpdate: dispatch.challenges.update,
  gameGet: dispatch.game.get,
});


export default withFilters(connect(mapStateToProps, mapDispatchToProps)(ChallengesTable))
