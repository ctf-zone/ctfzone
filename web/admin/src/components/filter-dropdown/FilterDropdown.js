import React, { Component } from 'react'
import PropTypes from 'prop-types'

import './FilterDropdown.css'

class FilterDropdown extends Component {

  static propTypes = {
    onConfirm: PropTypes.func,
    onReset: PropTypes.func,
    children: PropTypes.node,
  }

  render() {
    return (
      <div styleName='dropdown'>
        <div styleName='container'>
          {this.props.children}
        </div>
        <div styleName='buttons'>
          <a href='javascript:;' styleName='ok' onClick={this.props.onConfirm}>OK</a>
          <a href='javascript:;' styleName='reset' onClick={this.props.onReset}>Reset</a>
        </div>
      </div>
    )
  }
}

export default FilterDropdown
