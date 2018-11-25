import React, { Component } from 'react'
import PropTypes from 'prop-types'
import { Tree } from 'antd'
import { connect } from 'react-redux'

@connect(
  state => ({
    files: state.files,
    filesListResult: state.api.effects.files.list,
  }),
  dispatch => ({
    filesList: dispatch.files.list,
  }),
)
class FilesTree extends Component {

  static propTypes = {
    files: PropTypes.object.isRequired,
    filesList: PropTypes.func.isRequired,
  }

  componentDidMount() {
    const { filesList } = this.props

    filesList()
  }

  renderNodes(nodes, keyPrefix) {
    const listNodes = nodes.map((node, i) =>
      <Tree.TreeNode
        title={node.name}
        key={`${keyPrefix}-${i}`}
        isLeaf={!node.isDirectory || node.items.length === 0}
      >
        {this.renderNodes(node.items || [], `${keyPrefix}-${i}`)}
      </Tree.TreeNode>
    )

    return listNodes
  }

  render() {
    const { files } = this.props.files

    return (
      <div>
        <h1>Files</h1>
        <Tree.DirectoryTree
          multiple
        >
          {this.renderNodes(files, '0')}
        </Tree.DirectoryTree>
      </div>
    )
  }
}

export default FilesTree
