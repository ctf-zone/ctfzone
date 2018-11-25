import React, { Component } from 'react'
import PropTypes from 'prop-types'

import Divider from '../divider/Divider'

class Page extends Component {

  static propTypes = {
    children: PropTypes.node,
    type: PropTypes.string,
    title: PropTypes.string,
  }

  render() {
    const { type, title, children } = this.props
    return (
      <div className={`ctf-page ctf-page--${type}`}>
        <h1 className='ctf-page-title'>{title}</h1>
        <Divider />
        {children}
      </div>
    )
  }
}

export default Page
