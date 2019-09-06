import React, { Component } from 'react'
import PropTypes from 'prop-types'
import { Button, Icon } from 'antd'

import styles from './Pagination.module.css'

class Pagination extends Component {

  static propTypes = {
    onNext: PropTypes.func,
    onPrev: PropTypes.func,
    hasNext: PropTypes.bool,
    hasPrev: PropTypes.bool,
  }

  render() {
    return (
      <div className={styles.pagination}>
        <Button.Group>
          <Button
            type='primary'
            ghost
            onClick={this.props.onPrev}
            disabled={!this.props.hasPrev}
          >
            <Icon type='left' />
          </Button>
          <Button
            type='primary'
            ghost
            onClick={this.props.onNext}
            disabled={!this.props.hasNext}
          >
            <Icon type='right' />
          </Button>
        </Button.Group>
      </div>
    )
  }
}

export default Pagination
