import axios from 'axios';

export default axios.create({
  baseURL: process.env.API_URL || '/api',
  xsrfCookieName: 'csrf-token',
  xsrfHeaderName: 'X-CSRF-Token',
  withCredentials: true
});

export class ApiError extends Error {
  constructor(status, message, errors) {
    super('ApiError');

    this.status = status;
    this.message = message;
    this.errors = errors;
  }
}
