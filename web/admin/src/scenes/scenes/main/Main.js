import React, { Component } from 'react'
import PropTypes from 'prop-types'
import { Switch, Route, withRouter } from 'react-router-dom'
import { Layout, Spin } from 'antd'
import { connect } from 'react-redux'

import { Nav, Sidebar } from './components'

import * as scenes from './scenes'

import styles from './Main.module.css'

class Main extends Component {

  static propTypes = {
    // @connect
    isChecking: PropTypes.bool.isRequired,
    isLoggedIn: PropTypes.bool.isRequired,
    authCheck: PropTypes.func.isRequired,

    // @withRouter
    history: PropTypes.object.isRequired,
  }

  componentDidMount() {
    const { authCheck } = this.props
    authCheck()
  }

  componentDidUpdate(prevProps) {
    const { history, isChecking, isLoggedIn } = this.props
    const wasChecking = prevProps.isChecking

    if (!isChecking) {
      if (wasChecking && !isLoggedIn) {
        history.push('/login')
      }
    }
  }

  renderSidebar() {
    return (
      <Sidebar>
        <Nav />
      </Sidebar>
    )
  }

  renderHeader() {
    return (
      <div>
      </div>
    )
  }

  renderContent() {
    return (
      <Switch>
        <Route exact path='/' component={scenes.Dashboard} />
        <Route exact path='/users' component={scenes.UsersTable} />
        <Route exact path='/users/create' component={scenes.UserCreateEdit} />
        <Route exact path='/users/:userId/edit' component={scenes.UserCreateEdit} />
        <Route exact path='/challenges' component={scenes.ChallengesTable} />
        <Route exact path='/challenges/create' component={scenes.ChallengeCreateEdit} />
        <Route exact path='/challenges/:challengeId/edit' component={scenes.ChallengeCreateEdit} />
        <Route exact path='/announcements' component={scenes.AnnouncementsTable} />
        <Route exact path='/announcements/create' component={scenes.AnnouncementCreateEdit} />
        <Route exact path='/announcements/:announcementId/edit' component={scenes.AnnouncementCreateEdit} />
        <Route exact path='/files' component={scenes.FilesTree} />
      </Switch>
    )
  }

  renderFooter() {
    return (
      <span>CTFZone 2018</span>
    )
  }

  renderAll() {
    return (
      <Layout style={{ minHeight:"100vh" }}>

        {this.renderSidebar()}

        <Layout>

          <Layout.Header>
            {this.renderHeader()}
          </Layout.Header>

          <Layout.Content className={styles.content}>
            <div className={styles.body}>
              {this.renderContent()}
            </div>
          </Layout.Content>

          <Layout.Footer>
            {this.renderFooter()}
          </Layout.Footer>

        </Layout>
      </Layout>
    )
  }

  renderSpin() {
    const { isChecking } = this.props

    return (
      <div className={styles.spin}>
        <Spin spinning={isChecking} />
      </div>
    )
  }

  render() {
    const { isLoggedIn } = this.props

    return (
      isLoggedIn
        ? this.renderAll()
        : this.renderSpin()
    )
  }
}

const mapStateToProps = (state) => ({
  isLoggedIn: state.api.models.auth.success,
  isChecking: state.api.models.auth.loading,
});

const mapDispatchToProps = (dispatch) => ({
  authCheck: dispatch.auth.check,
});

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(Main))
