import React, { Component } from 'react'
import PropTypes from 'prop-types'
import classNames from 'classnames'

class Button extends Component {
  static propTypes = {
    children: PropTypes.node,
    disabled: PropTypes.bool,
    onClick: PropTypes.func,
  }

  render() {
    const { disabled, ...rest } = this.props
    const className = classNames({
      'ctf-button': true,
      'ctf-button--disabled': disabled,
    })

    return (
      <input
        type="submit"
        className={className}
        disabled={disabled}
        {...rest}
      />
    )
  }
}

export default Button
