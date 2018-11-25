import React, { Component } from 'react'
import PropTypes from 'prop-types'
import { Link } from 'react-router-dom'
import { Menu, Icon } from 'antd'
import { withRouter } from 'react-router'

@withRouter
class Nav extends Component {

  static propTypes = {
    menuItems: PropTypes.array,
    location: PropTypes.object,
  }

  static defaultProps = {
    menuItems: [],
  }

  render() {

    const { menuItems, location } = this.props

    let selectedKeys = []
    menuItems.forEach((m, i) => {
      if (m.path === '/') {
        if (location.pathname === '/')
          selectedKeys.push(i.toString())
        return
      }

      if (location.pathname.startsWith(m.path))
        selectedKeys.push(i.toString())
    })

    return (
      <Menu
        selectedKeys={selectedKeys}
        theme='dark'
      >
        { menuItems.map((item, i) => (
          <Menu.Item key={i}>
            <Link to={item.path}>
              <Icon type={item.icon} />
              <span>{item.title}</span>
            </Link>
          </Menu.Item>
        )) }
      </Menu>
    )
  }
}

export default Nav
