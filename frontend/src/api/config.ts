export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '';

export const WEBSOCKET_BASE_URL =
  API_BASE_URL !== ''
    ? API_BASE_URL.replace(/^http/, 'ws')
    : location.origin.replace(/^http/, 'ws');

export const API_USER = import.meta.env.VITE_API_USER || '';
export const API_PASSWORD = import.meta.env.VITE_API_PASSWORD || '';

export const AUTO_LOGIN = API_USER && API_PASSWORD;
