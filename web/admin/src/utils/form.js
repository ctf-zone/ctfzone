import { Form } from 'antd'

export const hasErrors = (fieldsError) => {
  return Object
    .keys(fieldsError)
    .some((field) => fieldsError[field]);
}

export const addErrors = (fields, errors) => {
  return Object.assign(
    fields,
    ...Object.keys(errors).map((field) => ({
      [field]: {
        ...fields[field],
        errors: errors[field].map((e) => new Error(e)),
      },
    })),
  )
}

export const addValues = (fields, values) => {
  return Object.assign(
    fields,
    ...Object.keys(values).map((field) => ({
      [field]: {
        ...fields[field],
        value: values[field],
      },
    })),
  )
}

export const mapPropsToFields = ({ fields }) => {
  return Object.assign(
    {},
    ...Object.keys(fields || {}).map((field) => {
      return {
        [field]: Form.createFormField({
          ...fields[field],
        }),
      }
    }),
  )
}
