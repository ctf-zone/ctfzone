import React, { Component } from 'react'
import PropTypes from 'prop-types'

import styles from './FilterDropdown.module.css'

class FilterDropdown extends Component {

  static propTypes = {
    onConfirm: PropTypes.func,
    onReset: PropTypes.func,
    children: PropTypes.node,
  }

  render() {
    return (
      <div className={styles.dropdown}>
        <div className={styles.container}>
          {this.props.children}
        </div>
        <div className={styles.buttons}>
          <a href='javascript:;' className={styles.ok} onClick={this.props.onConfirm}>OK</a>
          <a href='javascript:;' className={styles.reset} onClick={this.props.onReset}>Reset</a>
        </div>
      </div>
    )
  }
}

export default FilterDropdown
