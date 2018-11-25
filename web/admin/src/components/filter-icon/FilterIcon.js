import React, { Component } from 'react'
import PropTypes from 'prop-types'
import { Icon } from 'antd'

import './FilterIcon.css'

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
        styleName={isFiltered ? 'active' : 'inactive'}
      />
    )
  }
}

export default FilterIcon
