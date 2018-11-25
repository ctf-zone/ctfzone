import React, { Component } from 'react'
import PropTypes from 'prop-types'

class Window extends Component {
  static propTypes = {
    title: PropTypes.string,
    type: PropTypes.string,
    children: PropTypes.node,
  }

  render() {
    return (
      <div className={'ctf-window ctf-window' + this.props.type}>
        <h1 className='ctf-window-title'>{this.props.title}</h1>
        {this.props.children}
      </div>
    )
  }
}

export default Window
