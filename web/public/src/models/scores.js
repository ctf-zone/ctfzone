import api from '~/utils/api'

export default {
  state: {
    items: [],
  },
  reducers: {
    set: (state, payload) => {
      return { ...state, ...payload }
    },
  },
  effects: {
    async list() {
      const response = await api.get('scores')
      const items = response.data
      this.set({ items })
      return items
    },
  },
}
