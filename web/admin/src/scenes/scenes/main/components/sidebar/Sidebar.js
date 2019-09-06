import React, { Component } from 'react'
import PropTypes from 'prop-types'
import { Layout, Icon } from 'antd'

import logoBig from '../../../../../images/logo-big.svg'
import logoSmall from '../../../../../images/logo-small.svg'

import styles from './Sidebar.module.css'

class Sidebar extends Component {

  static propTypes = {
    children: PropTypes.node,
  }

  state = {
    collapsed: false,
  }

  toggleSiderbar = () => {
    const { collapsed } = this.state

    this.setState({
      collapsed: !collapsed,
    })
  }

  render() {
    const { collapsed } = this.state
    const { children } = this.props

    const logo = collapsed
      ? logoSmall
      : logoBig

    const toggleIcon = collapsed ? 'right' : 'left'

    return (
      <Layout.Sider
        className={styles.sider}
        collapsed={collapsed}
      >
        <Layout
          className={styles.layout}
        >
          <Layout.Header style={{ padding: 0 }}>
            <div className={styles.logo}>
              <img alt='logo' src={logo} />
            </div>
          </Layout.Header>
          <Layout.Content>
            {children}
          </Layout.Content>
          <Layout.Footer style={{ padding: 0 }}>
            <div
              className={styles.toggle}
              onClick={this.toggleSiderbar}
            >
              <Icon type={toggleIcon} />
            </div>
          </Layout.Footer>
        </Layout>
      </Layout.Sider>
    )
  }
}

export default Sidebar
