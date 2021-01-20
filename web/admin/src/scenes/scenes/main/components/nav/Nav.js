import React, { Component } from "react";
import PropTypes from "prop-types";
import { Link } from "react-router-dom";
import { Menu, Icon } from "antd";
import { withRouter } from "react-router";

class Nav extends Component {
  static propTypes = {
    menuItems: PropTypes.array,
    location: PropTypes.object,
  };

  static defaultProps = {
    menuItems: [],
  };

  render() {
    const { menuItems, location } = this.props;

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
        <Menu.Item key={0}>
          <Link to="/users">
            <Icon type="team" />
            <span>Users</span>
          </Link>
        </Menu.Item>
        <Menu.Item key={1}>
          <Link to="/challenges">
            <Icon type="bulb" />
            <span>Challenges</span>
          </Link>
        </Menu.Item>
        <Menu.Item key={2}>
          <Link to="/announcements">
            <Icon type="notification" />
            <span>Announcements</span>
          </Link>
        </Menu.Item>
        <Menu.Item key={3}>
          <Link to="/users">
            <Icon type="file" />
            <span>Files</span>
          </Link>
        </Menu.Item>
      </Menu>
    );
  }
}

export default withRouter(Nav);
