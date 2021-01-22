import React, { Component } from 'react'
import PropTypes from 'prop-types'
import { connect } from 'react-redux'
import {
  Form,
  Input,
  Select,
  InputNumber,
  Checkbox,
  Switch,
  Button,
  Icon,
} from 'antd'

import { MarkdownEditor } from '../../../../../../../components'
import { mapPropsToFields, hasErrors } from '../../../../../../../utils/form'

//import styles from './ChallengeForm.module.css'

class ChallengeForm extends Component {

  static propTypes = {
    isEdit: PropTypes.bool.isRequired,
    onSubmit: PropTypes.func.isRequired,
    onChange: PropTypes.func.isRequired,
    fields: PropTypes.object.isRequired,
    categories: PropTypes.array.isRequired,
    challenges: PropTypes.array.isRequired,
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

    const { isEdit, fields } = this.props
    const { getFieldsValue } = this.props.form
    const challenge = getFieldsValue()

    if (isEdit && !challenge.updateFlag) {
      delete challenge.flag
    }

    delete challenge.updateFlag

    challenge.description = fields.description.value

    this.props.onSubmit(challenge)
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

  renderDescriptionField({ getFieldValue }) {
    const { fields, filesUpload } = this.props

    return (
      <Form.Item
        wrapperCol={{ sm: { span: 14, offset: 6 } }}
      >
        <MarkdownEditor
          value={getFieldValue('description')}
          onChange={(editor, data, value) => {
            this.props.onChange({
              ...fields,
              description: { value },
            })
          }}
          onBeforeChange={(editor, data, value) => {
            this.props.onChange({
              ...fields,
              description: { value },
            })
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

  renderCategoriesField({ getFieldDecorator }) {
    const { categories } = this.props

    const options = categories.map((category) =>
      <Select.Option key={category}>{category}</Select.Option>)

    return (
      <Form.Item
        hasFeedback
        label='Categories'
        {...this.formItemLayout}
      >
        {getFieldDecorator('categories')(
          <Select
            mode='multiple'
            optionFilterProp='children'
          >
            {options}
          </Select>,
        )}
      </Form.Item>
    )
  }

  renderPointsField({ getFieldDecorator }) {
    return (
      <Form.Item
        hasFeedback
        label='Points'
        {...this.formItemLayout}
      >
        {getFieldDecorator('points', {
          initialValue: 100,
        })(
          <InputNumber
            min={1}
            step={100}
          />,
        )}
      </Form.Item>
    )
  }

  renderDifficultyField({ getFieldDecorator }) {
    return (
      <Form.Item
        hasFeedback
        label='Difficulty'
        {...this.formItemLayout}
      >
        {getFieldDecorator('difficulty', {
          initialValue: 'easy',
          rules: [
            {
              required: true,
              message: 'Difficulty is required!',
            },
          ],
        })(
          <Select
            optionFilterProp='children'
          >
            <Select.Option key='easy'>easy</Select.Option>
            <Select.Option key='medium'>medium</Select.Option>
            <Select.Option key='hard'>hard</Select.Option>
          </Select>,
        )}
      </Form.Item>
    )
  }

  renderUpdateFlagField({ getFieldDecorator }) {
    return (
      <Form.Item
        label='Update flag'
        {...this.formItemLayout}
      >
        {getFieldDecorator('updateFlag')(
          <Checkbox onChange={() => {
            const { resetFields } = this.props.form
            resetFields('flag')
          }} />,
        )}
      </Form.Item>
    )
  }

  renderFlagField({ getFieldDecorator, getFieldValue }) {
    const { isEdit } = this.props
    const updateFlag = getFieldValue('updateFlag')

    let rules = []

    if (!isEdit || updateFlag) {
      rules.push(
        {
          required: true,
          message: 'Flag is required!',
        },
      )
    }

    return (
      <Form.Item
        hasFeedback
        label='Flag'
        {...this.formItemLayout}
      >
        {getFieldDecorator('flag', {
          rules,
        })(
          <Input
            prefix={<Icon type='flag' />}
            type='password'
            disabled={isEdit && !updateFlag}
          />,
        )}
      </Form.Item>
    )
  }

  renderIsLockedField({ getFieldDecorator, getFieldValue }) {
    return (
      <Form.Item
        label='Locked'
        {...this.formItemLayout}
      >
        {getFieldDecorator('isLocked', {
          initialValue: false,
          rules: [],
        })(
          <Switch checked={getFieldValue('isLocked')} />,
        )}
      </Form.Item>
    )
  }

  renderDependsOnField({ getFieldDecorator }) {
    const { challenges } = this.props;

    return (
      <Form.Item
        hasFeedback
        label='DependsOn'
        {...this.formItemLayout}
      >
        {getFieldDecorator('dependsOn')(
          <Select
            optionFilterProp='children'
            allowClear
          >
            {challenges.map(({ challenge }) => (
              <Select.Option key={challenge.id} value={challenge.id}>{challenge.title}</Select.Option>
            ))}
          </Select>,
        )}
      </Form.Item>
    )
  }

  renderSubmitButton({ isFieldTouched, getFieldValue, getFieldsError }) {
    const { isEdit } = this.props
    let required = [ 'title', 'points' ]

    if (!isEdit) {
      required.push('flag')
    }

    let requiredIsTouched = true

    for (let field of required) {
      requiredIsTouched = requiredIsTouched &&
        (isFieldTouched(field) || getFieldValue(field))
    }

    return (
      <Form.Item
        wrapperCol={{ span: 12, offset: 6 }}
      >
        <Button
          type='primary'
          htmlType='submit'
          disabled={!requiredIsTouched || hasErrors(getFieldsError())}
        >Save</Button>
      </Form.Item>
    )
  }

  render() {
    const { form, isEdit } = this.props

    return (
      <Form onSubmit={this.handleSubmit}>
        {this.renderTitleField(form)}
        {this.renderDescriptionField(form)}
        {this.renderCategoriesField(form)}
        {this.renderPointsField(form)}
        {this.renderDifficultyField(form)}
        {
          isEdit
            ? this.renderUpdateFlagField(form)
            : ''
        }
        {this.renderFlagField(form)}
        {this.renderIsLockedField(form)}
        {this.renderDependsOnField(form)}
        {this.renderSubmitButton(form)}
      </Form>
    )
  }
}

const mapStateToProps = () => ({});

const mapDispatchToProps = (dispatch) => ({
  filesUpload: dispatch.files.upload,
});

const onFieldsChange = (props, changedFields, allFields) => {
  props.onChange(allFields)
}


export default Form.create({ onFieldsChange, mapPropsToFields })(connect(mapStateToProps, mapDispatchToProps)(ChallengeForm))
