import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { Switch, Route, withRouter, Link } from 'react-router-dom';
import { connect } from 'react-redux';

import { Layout, Loading, NotFound } from '../components';

import { Auth, Challenge, Challenges, News, Rules, Scoreboard } from './scenes';
import { Nav, UserStats, RequireAuth } from './components';

class Main extends Component {
  static propTypes = {
    isLoggedIn: PropTypes.bool,
    authCheck: PropTypes.func,
    contestStatus: PropTypes.object,
    contestGetStatus: PropTypes.func
  };

  componentDidMount() {
    const { authCheck, contestGetStatus } = this.props;
    authCheck();
    contestGetStatus();
  }

  renderHeader() {
    const { isLoggedIn, contestStatus } = this.props;
    const countdown = contestStatus.status === 'countdown';

    return (
      <Layout.Header className="ctf-header">
        <Layout.Container>
          <Link to="/">
            <div className="ctf-header-logo" />
          </Link>
          <Nav
            leftItems={[
              { title: 'Новости', path: '/news'},
              { title: 'Правила', path: '/rules' },
              { title: 'Таблица результатов', path: '/scoreboard'},
              {
                title: 'Задания',
                path: '/challenges',
                hidden: !isLoggedIn || countdown
              }
            ]}
            rightItems={[
              { title: 'Вход', path: '/auth/login', hidden: isLoggedIn},
              { title: 'Выход', path: '/auth/logout', hidden: !isLoggedIn }
            ]}
            rightPrefix={isLoggedIn ? <UserStats /> : null}
          />
        </Layout.Container>
      </Layout.Header>
    );
  }

  renderContent() {
    return (
      <Layout.Content className="ctf-content">
        <Layout.Container>
          <Switch>
            <Route path="/auth" component={Auth} />
            <Route exact path="/" component={News} />
            <Route exact path="/news" component={News} />
            <Route exact path="/scoreboard" component={Scoreboard} />
            <Route exact path="/rules" component={Rules} />
            <Route
              exact
              path="/challenges"
              component={RequireAuth(Challenges)}
            />
            <Route
              exact
              path="/challenges/:id"
              component={RequireAuth(Challenge)}
            />
            <Route component={NotFound} />
          </Switch>
        </Layout.Container>
      </Layout.Content>
    );
  }

  renderFooter() {
    return (
      <Layout.Footer className="ctf-footer">
        <Layout.Container>Cyber Polygon 2021</Layout.Container>
      </Layout.Footer>
    );
  }

  render() {
    const { isLoggedIn } = this.props;

    return (
      <Loading loading={isLoggedIn === null}>
        <Layout className="ctf">
          {this.renderHeader()}
          {this.renderContent()}
          {this.renderFooter()}
        </Layout>
      </Loading>
    );
  }
}

const mapStateToProps = state => ({
  isLoggedIn: state.auth,
  contestStatus: state.contest.status
});

const mapDispatchToProps = dispatch => ({
  authCheck: dispatch.auth.check,
  contestGetStatus: dispatch.contest.getStatus
});

export default withRouter(
  connect(
    mapStateToProps,
    mapDispatchToProps
  )(Main)
);
