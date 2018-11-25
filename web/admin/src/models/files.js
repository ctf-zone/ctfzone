import api from '~/utils/api'

export default {
  state: {
    files: [],
  },
  reducers: {
    set: (state, payload) => {
      return { ...state, ...payload }
    },
  },
  effects: {
    async list() {
      const response = await api.get('files')
      const files = response.data
      this.set({ files })
      return files
    },
    async upload(file) {
      let formData = new FormData()

      formData.append('file', file)

      const response = await api.post('files', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      })

      const fileInfo = response.data

      return fileInfo
    },
  },
}
