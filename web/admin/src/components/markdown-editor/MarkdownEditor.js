import React, { Component } from 'react'
import PropTypes from 'prop-types'

import { Tabs } from 'antd'
import { Controlled as CodeMirror } from 'react-codemirror2'
import ReactMarkdown from 'react-markdown'

import 'codemirror/mode/gfm/gfm'

import './MarkdownEditor.css'

class MarkdownEditor extends Component {

  static propTypes = {
    value: PropTypes.string,
  }

  render() {
    const { value, ...rest } = this.props

    return (
      <Tabs defaultActiveKey='1' animated={false}>
        <Tabs.TabPane tab='Edit' key='1'>
          <div styleName='edit'>
            <CodeMirror
              value={value}
              options={{
                mode: 'gfm',
                theme: 'idea',
              }}
              {...rest}
            />
          </div>
        </Tabs.TabPane>
        <Tabs.TabPane tab='Preview' key='2'>
          <div styleName='preview'>
            <ReactMarkdown source={this.props.value} />
          </div>
        </Tabs.TabPane>
      </Tabs>
    )
  }
}

export default MarkdownEditor
