import api from '../utils/api'

export default {
  state: {},
  reducers: {
    set: (state, payload) => {
      return { ...state, ...payload }
    },
  },
  effects: {
    async get() {
      const response = await api.get('game')
      const game = response.data
      this.set(game)
      return game
    },
  },
}
