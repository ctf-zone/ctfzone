import React, { Component } from 'react'
import PropTypes from 'prop-types'
import Moment from 'react-moment'
import * as _ from 'lodash'
import { Link } from 'react-router-dom'
import { connect } from 'react-redux'
import {
  Table,
  Icon,
  Spin,
  Button,
  Row,
  Col,
  Divider,
  Popconfirm,
  Upload,
  message,
} from 'antd'

// import { countries } from '~/components/flag-icon/FlagIcon'
import {
  // FlagIcon,
  Pagination,
  withFilters,
  textFilter,
  radioFilter,
  // selectFilter,
  dateRangeFilter,
} from '../../../../../components'

import styles from './UsersTable.module.css'

class UsersTable extends Component {

  static propTypes = {
    dispatch: PropTypes.func,

    // @connect
    users: PropTypes.object.isRequired,
    usersListResult: PropTypes.object.isRequired,
    usersDeleteResult: PropTypes.object.isRequired,
    usersList: PropTypes.func.isRequired,
    usersDelete: PropTypes.func.isRequired,
    usersCreate: PropTypes.func.isRequired,

    // @withFilters
    filters: PropTypes.object,
    handleFilterSet: PropTypes.func,
    handleFilterReset: PropTypes.func,
    handleResetFilters: PropTypes.func,
    isFiltered: PropTypes.func,
  }

  componentDidMount() {
    const { usersList } = this.props

    usersList({})
  }

  componentDidUpdate(prevProps) {
    const { usersList, filters } = this.props

    if (!_.isEqual(filters, prevProps.filters)) {
      usersList({ filters })
    }
  }

  handleUserDelete = (userId) => async() => {
    const { usersDelete, usersList, filters } = this.props

    await usersDelete(userId)
    await usersList({ filters })
  }

  handleNextButtonClick = () => {
    const { usersList } = this.props
    const { next } = this.props.users.links

    usersList({ link: next.url })
  }

  handlePrevButtonClick = () => {
    const { usersList } = this.props
    const { prev } = this.props.users.links

    usersList({ link: prev.url })
  }

  renderControls() {
    const { extra, ...restFilters } = this.props.filters
    const isFiltered = !(_.isEmpty(restFilters) && _.isEmpty(extra))

    return (
      <div className={styles.controls}>
        <Row>
          <Col span={12} className={styles.controlsLeft}>
            <Link to='/users/create'>
              <Button type='primary'>Create</Button>
            </Link>
            <Upload
              showUploadList={false}
              customRequest={(e) => {
                let reader = new FileReader()
                const { usersCreate } = this.props

                reader.onload = (e) => {
                  const json = e.target.result
                  try {
                    const users = JSON.parse(json)
                    users.forEach(async(user) => {
                      try {
                        await usersCreate(user)
                      } catch(e) {
                        message.error(e.message)
                      }
                    })
                  } catch(e) {
                    message.error('Invalid JSON')
                  }
                }
                reader.readAsText(e.file, 'UTF-8')
              }}
            >
              <Button>
                <Icon type='upload' /> Import
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
    const { users, links } = this.props.users
    const { prev, next } = links

    const columns = [
      {
        title: 'ID',
        dataIndex: 'id',
      },
      {
        title: 'Name',
        dataIndex: 'name',
        ...textFilter('name', this.props),
      },
      {
        title: 'Email',
        dataIndex: 'email',
        ...textFilter('email', this.props),
      },
      // TODO: read settings from server
      // to know that country field is present in extra
      // {
      //   title: 'Country',
      //   dataIndex: 'extra',
      //   render: (data) => (
      //     <FlagIcon
      //       size='lg'
      //       code={data.country}
      //       className={styles.flag}
      //     />
      //   ),
      //   ...selectFilter(
      //     'extra.country',
      //     countries.map(({ name, code }) => ({ data: name, value: code })),
      //     this.props,
      //   ),
      // },
      {
        title: 'Activated',
        dataIndex: 'isActivated',
        render: (data) => {
          return (
            data
              ? <Icon type='check' style={{ 'color': 'green' }} />
              : <Icon type='close' style={{ 'color': 'red' }}/>
          )
        },
        ...radioFilter(
          'isActivated',
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
        title: 'Registered',
        dataIndex: 'createdAt',
        render: (data) => <Moment fromNow>{data}</Moment>,
        ...dateRangeFilter('createdAt', this.props),
      },
      {
        title: 'Action',
        render: (data, record) => (
          <span>
            <Link to={`users/${record.id}/edit`}>
              <span>Edit</span>
            </Link>
            <Divider type='vertical' />
            <Popconfirm
              title={`Are you sure that you want to delete user "${record.name}"?`}
              okText='Yes'
              cancelText='No'
              placement='leftTop'
              onConfirm={this.handleUserDelete(record.id)}
            >
              <a href='javascript:;'>Delete</a>
            </Popconfirm>
          </span>
        ),
      },
    ]

    return (
      <div>
        <Table
          rowKey='id'
          bordered={true}
          columns={columns}
          dataSource={users}
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
    const { loading } = this.props.usersListResult

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


const mapStateToProps = state => ({
    users: state.users,
    usersListResult: state.api.effects.users.list,
    usersDeleteResult: state.api.effects.users.delete,
});

const mapDispatchToProps = dispatch => ({
    usersList: dispatch.users.list,
    usersDelete: dispatch.users.delete,
    usersCreate: dispatch.users.create,
});


export default withFilters(connect(mapStateToProps, mapDispatchToProps)(UsersTable))
