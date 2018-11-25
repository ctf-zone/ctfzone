import React, { Component } from 'react'
import PropTypes from 'prop-types'
import { Switch, Route, withRouter, Link } from 'react-router-dom'
import { hot } from 'react-hot-loader'
import { connect } from 'react-redux'

import { Layout, Loading, NotFound } from '~/components'
import { Nav, UserStats, requireAuth } from './components'
import * as scenes from './scenes'

import './App.css'

@hot(module)
@withRouter
@connect(
  (state) => ({
    isLoggedIn: state.auth,
    contestStatus: state.contest.status,
  }),
  (dispatch) => ({
    authCheck: dispatch.auth.check,
    contestGetStatus: dispatch.contest.getStatus,
  })
)
class App extends Component {

  static propTypes = {
    isLoggedIn: PropTypes.bool,
    authCheck: PropTypes.func,
    contestStatus: PropTypes.object,
    contestGetStatus: PropTypes.func,
  }

  componentDidMount() {
    const { authCheck, contestGetStatus } = this.props
    authCheck()
    contestGetStatus()
  }

  renderHeader() {
    const { isLoggedIn, contestStatus } = this.props
    const countdown = contestStatus.status === 'countdown'

    return (
      <Layout.Header className='ctf-header'>
        <Layout.Container>
          <Link to='/'>
            <div className='ctf-header-logo'></div>
          </Link>
          <Nav
            leftItems={[
              { title: 'News', path: '/news' },
              { title: 'Rules', path: '/rules' },
              { title: 'Scoreboard', path: '/scoreboard' },
              { title: 'Challenges', path: '/challenges', hidden: !isLoggedIn || countdown },
            ]}
            rightItems={[
              { title: 'Log In', path: '/auth/login', hidden: isLoggedIn },
              { title: 'Log Out', path: '/auth/logout', hidden: !isLoggedIn },
            ]}
            rightPrefix={isLoggedIn ? <UserStats /> : null}
          />
        </Layout.Container>
      </Layout.Header>
    )
  }

  renderContent() {
    return (
      <Layout.Content className='ctf-content'>
        <Layout.Container>
          <Switch>
            <Route path='/auth' component={scenes.Auth} />
            <Route exact path='/' component={scenes.News} />
            <Route exact path='/news' component={scenes.News} />
            <Route exact path='/scoreboard' component={scenes.Scoreboard} />
            <Route exact path='/rules' component={scenes.Rules} />
            <Route exact path='/challenges' component={requireAuth(scenes.Challenges)} />
            <Route exact path='/challenges/:id' component={requireAuth(scenes.Challenge)} />
            <Route component={NotFound} />
          </Switch>
        </Layout.Container>
      </Layout.Content>
    )
  }

  renderFooter() {
    return (
      <Layout.Footer className='ctf-footer'>
        <Layout.Container>
          M⭑CTF 2018
        </Layout.Container>
      </Layout.Footer>
    )
  }

  render() {
    const { isLoggedIn } = this.props

    return (
      <Loading
        loading={isLoggedIn === null}
      >
        <Layout className='ctf'>
          {this.renderHeader()}
          {this.renderContent()}
          {this.renderFooter()}
        </Layout>
      </Loading>
    )
  }
}

export default App
