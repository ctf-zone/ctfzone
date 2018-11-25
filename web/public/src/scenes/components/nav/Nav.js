import React, { Component } from 'react'
import PropTypes from 'prop-types'
import { Link } from 'react-router-dom'

class Nav extends Component {

  static propTypes = {
    leftItems: PropTypes.array,
    rightItems: PropTypes.array,
    rightPrefix: PropTypes.node,
  }

  static defaultProps = {
    leftItems: [],
    rightItems: [],
  }

  renderLinks(items) {
    return items.filter(({ hidden }) => !hidden).map(({ title, path }, i) => {
      return (
        <div key={i} className='ctf-header-nav-item'>
          <Link to={path}>{title}</Link>
        </div>
      )
    })
  }

  render() {
    const { leftItems, rightItems, rightPrefix } = this.props

    return (
      <div className='ctf-header-nav'>
        <div className='ctf-header-nav-left'>
          {this.renderLinks(leftItems)}
        </div>
        <div className='ctf-header-nav-right'>
          {rightPrefix}
          {this.renderLinks(rightItems)}
        </div>
      </div>
    )
  }
}

export default Nav
