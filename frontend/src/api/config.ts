export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '';

export const WEBSOCKET_BASE_URL =
  API_BASE_URL !== ''
    ? API_BASE_URL.replace(/^http/, 'ws')
    : location.origin.replace(/^http/, 'ws');

export const ADMIN_USER = import.meta.env.VITE_ADMIN_USER || '';
export const ADMIN_PASSWORD = import.meta.env.VITE_ADMIN_PASSWORD || '';

export const AUTO_LOGIN = ADMIN_USER && ADMIN_PASSWORD;
