import api from '../utils/api'

export default {
  state: false,
  reducers: {
    set: (state, payload) => {
      return payload
    },
  },
  effects: {
    async login(data) {
      await api.post('auth/login', data)
      this.set(true)
    },
    async check() {
      await api.get('auth/check')
      this.set(true)
    },
    async logout() {
      await api.post('auth/logout')
      this.set(false)
    },
  },
}
