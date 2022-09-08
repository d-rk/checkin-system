import axios from 'axios';

const API_BASE_URL = process.env.REACT_APP_API_BASE_URL || '';

axios.defaults.baseURL = API_BASE_URL;

axios.defaults.headers.common = {
  'Content-Type': 'application/json',
};
