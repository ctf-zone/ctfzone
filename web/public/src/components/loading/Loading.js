import React, { Component } from 'react'
import PropTypes from 'prop-types'

import Spinner from '../spinner/Spinner'

class Loading extends Component {

  static propTypes = {
    children: PropTypes.node,
    loading: PropTypes.bool,
  }

  render() {
    const { loading, children } = this.props

    if (loading) {
      return (
        <div className='ctf-loading'>
          <Spinner />
        </div>
      )
    }

    return children
  }
}

export default Loading
