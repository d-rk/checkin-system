import axios from 'axios';
import {API_BASE_URL} from './config';
import {refreshAccessToken} from './checkInSystemApi';
import {initializeAuth} from './tokenService';

axios.defaults.baseURL = API_BASE_URL;

axios.defaults.headers.common = {
  'Content-Type': 'application/json',
};

// Initialize auth header from stored token on app startup
initializeAuth();

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
      originalRequest.headers.Authorization = `Bearer ${token}`;
      return axios(originalRequest);
    } else if (error.response.status === 401) {
      originalRequest._retryFailed = true;
    }
    return Promise.reject(error);
  }
);
