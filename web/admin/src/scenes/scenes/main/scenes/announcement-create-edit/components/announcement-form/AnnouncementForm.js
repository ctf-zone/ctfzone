import React, { Component } from 'react'
import PropTypes from 'prop-types'
import { connect } from 'react-redux'
import {
  Form,
  Input,
  Button,
  Select,
  Icon,
} from 'antd'

import { MarkdownEditor } from '../../../../../../../components'
import { mapPropsToFields } from '../../../../../../../utils/form'


class AnnouncementForm extends Component {

  static propTypes = {
    isEdit: PropTypes.bool.isRequired,
    onSubmit: PropTypes.func.isRequired,
    onChange: PropTypes.func.isRequired,
    fields: PropTypes.object.isRequired,
    challenges: PropTypes.array,
    filesUpload: PropTypes.func,

    // @form
    form: PropTypes.object,
  }

  formItemLayout = {
    labelCol: {
      sm: { span: 6 },
    },
    wrapperCol: {
      sm: { span: 14 },
    },
  }

  handleSubmit = (e) => {
    e.preventDefault()

    const { getFieldsValue } = this.props.form
    const { fields } = this.props
    const announcement = getFieldsValue()
    this.props.onSubmit({
      ...announcement,
      body: fields.body.value,
      challengeId: parseInt(announcement.challengeId) || null,
    })
  }

  renderTitleField({ getFieldDecorator }) {
    return (
      <Form.Item
        hasFeedback
        label='Title'
        {...this.formItemLayout}
      >
        {getFieldDecorator('title', {
          rules: [
            {
              required: true,
              message: 'Title is required!',
            },
          ],
        })(
          <Input
            prefix={<Icon type='profile' />}
          />,
        )}
      </Form.Item>
    )
  }

  renderBodyField({ getFieldValue }) {
    const { filesUpload } = this.props

    return (
      <Form.Item
        wrapperCol={{ sm: { span: 14, offset: 6 } }}
      >
        <MarkdownEditor
          value={getFieldValue('body')}
          onChange={(editor, data, value) => {
            this.props.onChange({ body: { value } })
          }}
          onBeforeChange={(editor, data, value) => {
            this.props.onChange({ body: { value } })
          }}
          onDrop={async(editor, e) => {
            // TODO: refactor
            e.preventDefault()

            const files = Array.from(e.dataTransfer.files)

            files.forEach(async(file) => {
              const { url } = await filesUpload(file)
              const doc = editor.getDoc()
              const cursor = doc.getCursor()

              doc.replaceRange(`[](${url})`, cursor)
            })
          }}
        />
      </Form.Item>
    )
  }

  renderChallengeField({ getFieldDecorator }) {
    const { challenges } = this.props

    const options = challenges.map(({ title, value }) =>
      <Select.Option key={value}>{title}</Select.Option>)

    return (
      <Form.Item
        hasFeedback
        label='Challenge'
        {...this.formItemLayout}
      >
        {getFieldDecorator('challengeId')(
          <Select
            optionFilterProp='children'
            allowClear={true}
          >
            {options}
          </Select>,
        )}
      </Form.Item>
    )
  }

  renderSubmitButton({ isFieldTouched, getFieldValue, getFieldError }) {
    const canSubmit = ['title', 'body'].reduce((result, field) => {
      return result &&
        (isFieldTouched(field) || getFieldValue(field)) &&
        !getFieldError(field)
    }, true);

    return (
      <Form.Item
        wrapperCol={{ span: 12, offset: 6 }}
      >
        <Button
          type='primary'
          htmlType='submit'
          disabled={!canSubmit}
        >Save</Button>
      </Form.Item>
    )
  }

  render() {
    const { form } = this.props

    return (
      <Form onSubmit={this.handleSubmit}>
        {this.renderTitleField(form)}
        {this.renderBodyField(form)}
        {this.renderChallengeField(form)}
        {this.renderSubmitButton(form)}
      </Form>
    )
  }
}

const onFieldsChange = (props, changedFields) => {
  props.onChange(changedFields)
};

const mapStateToProps = () => ({});

const mapDispatchToProps = (dispatch) => ({
  filesUpload: dispatch.files.upload,
});


export default Form.create({ mapPropsToFields, onFieldsChange })(connect(mapStateToProps, mapDispatchToProps)(AnnouncementForm))
