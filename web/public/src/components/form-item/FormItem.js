import React, { Component } from 'react'
import PropTypes from 'prop-types'
import classNames from 'classnames'

class FormItem extends Component {

  static propTypes = {
    label: PropTypes.string,
    children: PropTypes.node,
    errors: PropTypes.array,
  }

  render() {
    const { children, errors, label } = this.props
    const valid = !errors || errors.length === 0
    const baseClass = 'ctf-form-item'
    const className = classNames({
      [baseClass]: true,
      [baseClass + '--valid']: valid,
      [baseClass + '--invalid']: !valid,
    })

    return (
      <div className={className}>
        <div className={baseClass + '-field'}>
          {children}
          <label>{label}</label>
        </div>
        <div className={baseClass + '-error'}>
          {(errors || []).join(', ')}
        </div>
      </div>
    )
  }
}

export default FormItem
