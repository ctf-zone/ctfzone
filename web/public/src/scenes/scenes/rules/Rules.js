import React, { Component } from 'react';
import ReactMarkdown from 'react-markdown';

import { Page } from '../../../components';

// TODO: get from API
const rules = `
Online-stage, starts on 20.10.2018 at 09:00 UTC and will last for 24 hours.
Stage format – Jeopardy. Top 10 teams will be invited to the final.

Rules:

- register;
- tasks will be available at the link;
- those who earn more points wins the CTF;
- in case of equal points, the winner will be the one who earned them first.

It's forbidden to:

- attack the organizer's infrastructure;
- generate large amounts of traffic (DDoS);
- attack computers of the jury or other participants;
- share flags with other participants.

The amount of points that every team gets for each task depends on how many times this task was solved by all teams.

Each task is marked as easy, medium or hard.
This difficulty level doesn’t affect scoring formula of the task,
which only depends on the amount of teams that has submitted the flag.

Flag format: \`mctf{[a-f0-9]{32}}\`

The organizers reserve the right to disqualify participants for violating the rules.
`;

class Rules extends Component {
  static propTypes = {};

  componentDidMount() {}

  render() {
    return (
      <Page title="Rules" type="rules">
        <ReactMarkdown source={rules} />
      </Page>
    );
  }
}

export default Rules;
