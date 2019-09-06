import React, { Component } from 'react'
import PropTypes from 'prop-types'
import { Row, Col } from 'antd'
import { connect } from 'react-redux'
import { withRouter } from 'react-router-dom'

import { addValues, addErrors } from '../../../../../utils/form'
import { UserForm } from './components'

class UserCreateEdit extends Component {

  static propTypes = {
    // @connect
    user: PropTypes.object.isRequired,

    usersGet: PropTypes.func.isRequired,
    usersCreate: PropTypes.func.isRequired,
    usersUpdate: PropTypes.func.isRequired,

    // @withRouter
    match: PropTypes.object.isRequired,
    history: PropTypes.object.isRequired,
  }

  static defaultState = {
    fields: {},
  }

  state = {
    fields: {},
  }

  async componentDidMount() {
    const { userId } = this.props.match.params
    const { usersGet } = this.props

    if (userId) {
      const user = await usersGet(userId)

      this.setState(({ fields }) => ({
        fields: addValues(fields, {
          ...user,
          extra: JSON.stringify(user.extra),
        }),
      }))
    }
  }

  handleChange = changedFields => {
    this.setState(({ fields }) => ({
      fields: { ...fields, ...changedFields },
    }))
  }

  handleSubmit = async user => {
    const { usersUpdate, usersCreate, history } = this.props
    const { userId } = this.props.match.params

    try {
      if (userId) {
        await usersUpdate(
          {
            ...user,
            id: parseInt(userId),
          },
          {
            throw: true,
          },
        )
      } else {
        await usersCreate(user, { throw: true })
      }

      history.push('/users')
    } catch (err) {
      this.setState(({ fields }) => ({
        fields: addErrors(fields, err.errors),
      }))
    }
  }

  render() {
    const { user } = this.props
    const isEdit = !!this.props.match.params.userId
    const { fields } = this.state

    return (
      <div>
        <Row>
          <Col offset={6}>
            <h1>
              {
                isEdit
                  ? `Edit user "${user.name}"`
                  : 'Create new user'
              }
            </h1>
          </Col>
        </Row>
        <UserForm
          isEdit={isEdit}
          onSubmit={this.handleSubmit}
          onChange={this.handleChange}
          fields={fields}
        />
      </div>
    )
  }
}

const mapStateToProps = state => ({
  user: state.users.user,
});

const mapDispatchToProps = dispatch => ({
  usersCreate: dispatch.users.create,
  usersUpdate: dispatch.users.update,
  usersGet: dispatch.users.get,
});


export default withRouter(connect(mapStateToProps, mapDispatchToProps)(UserCreateEdit))
