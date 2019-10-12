import api from '../utils/api';

export default {
  state: null,
  reducers: {
    set: (state, payload) => {
      return payload;
    }
  },
  // TODO: replace data with fields
  effects: {
    async register(data) {
      const { reCaptchaResponse, ...rest } = data;
      console.log(reCaptchaResponse, rest);
      await api.post('/auth/register', rest, {
        headers: {
          'X-G-Recaptcha-Response': reCaptchaResponse,
        }
      });
    },
    async login(data) {
      await api.post('/auth/login', data);
      this.set(true);
    },
    async logout() {
      await api.post('/auth/logout', {});
      this.set(false);
    },
    async check() {
      const response = await api.get('/auth/check');
      const { isLoggedIn } = response.data;
      this.set(isLoggedIn);
    },
    async resetPassword(data) {
      await api.post('/auth/reset-password', data);
    },
    async sendToken({ type, token, email }) {
      await api.post('/auth/send-token', { type, token, email });
    },
    async activate({ token }) {
      await api.post('/auth/activate', { token });
    }
  }
};
