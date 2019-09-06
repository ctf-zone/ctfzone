import React, { Component } from 'react'
import PropTypes from 'prop-types'
import { Icon } from 'antd'

import styles from './FilterIcon.module.css'

class FilterIcon extends Component {
  static propTypes = {
    isFiltered: PropTypes.bool,
    onClick: PropTypes.func,
  }

  render() {
    const { isFiltered, ...rest } = this.props

    return (
      <Icon
        type='filter'
        {...rest}
        className={isFiltered ? styles.active : styles.inactive}
      />
    )
  }
}

export default FilterIcon
