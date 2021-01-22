import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { Link } from 'react-router-dom';
import ReactMarkdown from 'react-markdown';

import { Page, Divider } from '../../../components';
import dayjs from '../../../utils/date';
import Countdown from './Countdown';

class News extends Component {
  static propTypes = {
    newsList: PropTypes.func,
    challengesList: PropTypes.func,
    news: PropTypes.object,
    challenges: PropTypes.object,
    contestGetStatus: PropTypes.func,
    contestStatus: PropTypes.object
  };

  state = {
    intervalId: 0
  };

  componentDidMount() {
    this.fetchData();

    // Update data every 30 seconds
    const intervalId = setInterval(this.fetchData, 30000);
    this.setState({ intervalId });
  }

  componentWillUnmount() {
    clearInterval(this.state.intervalId);
  }

  fetchData = () => {
    const { newsList, challengesList, contestGetStatus } = this.props;

    contestGetStatus();
    newsList({});
    challengesList();
  };

  handleStart = () => {
    this.fetchData();
  };

  renderAnnouncement(announcement, index) {
    const prefixClass = 'ctf-announcement';
    const { items: challenges } = this.props.challenges;
    const challengeId = announcement.challengeId;

    return (
      <div key={index}>
        <div className={prefixClass}>
          <div className={`${prefixClass}-title`}>{announcement.title}</div>
          {challengeId && challengeId in challenges ? (
            <div className={`${prefixClass}-hint`}>
              <Link to={`/challenges/${challengeId}`}>
                {challenges[challengeId].challenge.title}
              </Link>
            </div>
          ) : null}
          <div className={`${prefixClass}-date`}>
            {dayjs(announcement.createdAt).fromNow()}
          </div>
          <div className={`${prefixClass}-body`}>
            <ReactMarkdown source={announcement.body} />
          </div>
        </div>
        <Divider />
      </div>
    );
  }

  render() {
    const { items: announcements } = this.props.news;
    const { contestStatus } = this.props;
    const countdown = contestStatus.status === 'countdown';

    return (
      <Page title="Новости" type="news">
        {countdown ? (
          <div>
            <Countdown
              startTime={new Date(contestStatus.start)}
              onStart={this.handleStart}
            />
            <Divider />
          </div>
        ) : null}
        {announcements.map(this.renderAnnouncement.bind(this))}
      </Page>
    );
  }
}

const mapStateToProps = state => ({
  news: state.news,
  challenges: state.challenges,
  contestStatus: state.contest.status
});

const mapDispatchToProps = dispatch => ({
  newsList: dispatch.news.list,
  challengesList: dispatch.challenges.list,
  contestGetStatus: dispatch.contest.getStatus
});

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(News);
