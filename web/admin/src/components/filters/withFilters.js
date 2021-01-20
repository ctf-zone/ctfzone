import React, { Component } from 'react'
import * as _ from 'lodash'

import { set, unset } from '../../utils/object'

export default (WrappedComponent) => {
  return class WithFilters extends Component {
    state = {
      filters: {},
      resetFilters: {},
    }

    static defaultState = {
      filters: {},
      resetFilters: {},
    }

    isFiltered = (path) => {
      return _.has(this.state.filters, path)
    }

    handleFilterSet = (path, value, confirm, clearFilters) => () => {
      confirm()

      const resetFilters = { ...this.state.resetFilters }

      if (!([path] in resetFilters)) {
        resetFilters[path] = clearFilters
      }

      this.setState({
        filters: set(path, value, this.state.filters),
        resetFilters,
      })
    }

    handleFilterReset = (path, clearFilters) => () => {
      clearFilters()

      const resetFilters = { ...this.state.resetFilters }

      if ([path] in resetFilters) {
        delete resetFilters[path]
      }

      this.setState({
        filters: unset(path, this.state.filters),
        resetFilters,
      })
    }

    handleResetFilters = () => {
      Object.values(this.state.resetFilters).
        forEach((clear) => clear())

      this.setState({
        filters: {},
      })
    }

    render() {
      const props = {
        ...this.props,
        filters: this.state.filters,
        handleFilterSet: this.handleFilterSet,
        handleFilterReset: this.handleFilterReset,
        handleResetFilters: this.handleResetFilters,
        isFiltered: this.isFiltered,
      }

      return <WrappedComponent {...props} />
    }
  }
}
