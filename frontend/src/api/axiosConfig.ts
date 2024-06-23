import axios from 'axios';
import {API_BASE_URL} from './config';
import {refreshAccessToken} from './checkInSystemApi';

axios.defaults.baseURL = API_BASE_URL;

axios.defaults.headers.common = {
  'Content-Type': 'application/json',
};

// Response interceptor to automatically refresh token
axios.interceptors.response.use(
  response => {
    return response;
  },
  // eslint-disable-next-line promise/prefer-await-to-callbacks
  async error => {
    const originalRequest = error.config;
    if (error.response.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;
      const {token} = await refreshAccessToken();
      axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
      originalRequest.headers.Authorization = `Bearer ${token}`;
      return axios(originalRequest);
    } else if (error.response.status === 401) {
      originalRequest._retryFailed = true;
    }
    return Promise.reject(error);
  }
);
