import React, { Component } from 'react'
import PropTypes from 'prop-types'
import { connect } from 'react-redux'
import { withRouter } from 'react-router-dom'

const requireAuth = (Wrapped) => {

  @withRouter
  @connect(
    (state) => ({
      isLoggedIn: state.auth,
    }),
  )
  class Wrapper extends Component {

    static propTypes = {
      isLoggedIn: PropTypes.bool,
      history: PropTypes.object,
    }

    componentDidMount() {
      const { isLoggedIn, history } = this.props

      if (!isLoggedIn) {
        history.push('/auth/login')
      }
    }

    render() {
      return <Wrapped />
    }
  }

  return Wrapper
}

export default requireAuth
