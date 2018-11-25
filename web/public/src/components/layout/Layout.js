import React, { Component } from 'react'
import PropTypes from 'prop-types'
import classNames from 'classnames'

const generator = (props) => {
  return (BasicComponent) => {
    return class Adapter extends Component {
      render() {
        const { prefixClass } = props
        return <BasicComponent prefixClass={prefixClass} {...this.props} />
      }
    }
  }
}

class Basic extends Component {
  static propTypes = {
    className: PropTypes.string,
    prefixClass: PropTypes.string,
    children: PropTypes.node,
  }

  render() {
    const { prefixClass, className, children, ...others } = this.props
    const divClass = classNames(className, prefixClass)
    return (
      <div className={divClass} {...others}>{children}</div>
    );
  }
}

const Layout = generator({ prefixClass: 'ctf-layout' })(Basic)
const Header = generator({ prefixClass: 'ctf-layout-header' })(Basic)
const Content = generator({ prefixClass: 'ctf-layout-content' })(Basic)
const Container = generator({ prefixClass: 'ctf-layout-container' })(Basic)
const Footer = generator({ prefixClass: 'ctf-layout-footer' })(Basic)

Layout.Header = Header
Layout.Content = Content
Layout.Footer = Footer
Layout.Container = Container

export default Layout
