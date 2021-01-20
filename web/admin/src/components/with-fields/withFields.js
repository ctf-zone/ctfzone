import React, { Component } from 'react'

export default (WrappedComponent) => {
  return class WithFiedls extends Component {

    static defaultState = {
      fields: {},
      updated: {},
    }

    state = {
      fields: {},
      updated: {},
    }

    handleFieldsChange = (changedFields) => {
      this.setState(({ fields, updated }) => ({
        fields: { ...fields, ...changedFields },
        updated: Object.assign(
          updated,
          ...Object.keys(changedFields).map((field) => ({ [field]: new Date() })),
        ),
      }))
    }

    render() {
      const props = {
        ...this.props,
        fields: this.state.fields,
        handleFieldsChange: this.handleFieldsChange,
      }

      return <WrappedComponent {...props} />
    }
  }
}
