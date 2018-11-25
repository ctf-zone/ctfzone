import { createFormField } from 'rc-form'

export const addErrors = (fields, errors) => {
  return Object.assign(
    fields,
    ...Object.keys(errors).map((field) => ({
      [field]: createFormField({
        ...fields[field],
        errors: errors[field].map((e) => new Error(e)),
      }),
    })),
  )
}

export const mapPropsToFields = ({ fields }) => {
  return Object.assign(
    {},
    ...Object.keys(fields || {}).map((field) => {
      return {
        [field]: createFormField({
          ...fields[field],
        }),
      }
    })
  )
}
