import React, { Component } from "react";
import PropTypes from "prop-types";
import { Link } from "react-router-dom";
import { Menu, Icon } from "antd";
import { withRouter } from "react-router";

class Nav extends Component {
  static propTypes = {
    location: PropTypes.object,
  };


  render() {
    const { location } = this.props;

    const menuItems = [
      {
        title: 'Users',
        path: '/users',
        icon: 'team',
      },
      {
        title: 'Challenges',
        path: '/challenges',
        icon: 'bulb',
      },
      {
        title: 'Announcements',
        path: '/announcements',
        icon: 'notification',
      },
      {
        title: 'Files',
        path: '/files',
        icon: 'file',
      },
    ]

    let selectedKeys = [];
    menuItems.forEach((m, i) => {
      if (m.path === "/") {
        if (location.pathname === "/") selectedKeys.push(i.toString());
        return;
      }

      if (location.pathname.startsWith(m.path)) selectedKeys.push(i.toString());
    });

    return (
      <Menu selectedKeys={selectedKeys} theme="dark" style={{ height: "100%" }}>
        {menuItems.map((item, i) => (
          <Menu.Item key={i}>
            <Link to={item.path}>
              <Icon type={item.icon} />
              <span>{item.title}</span>
            </Link>
          </Menu.Item>
        ))}
      </Menu>
    );
  }
}

export default withRouter(Nav);
