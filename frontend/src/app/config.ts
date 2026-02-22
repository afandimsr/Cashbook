/**
 * Application configuration centralized from environment variables.
 * Uses Vite's import.meta.env for loading .env files.
 */

export const config = {
    API_URL: import.meta.env.VITE_API_URL || '/api/v1',
    APP_TITLE: import.meta.env.VITE_APP_TITLE || 'CashBook',
    APP_VERSION: '1.16.0', // Standard application version
    IS_DEV: import.meta.env.DEV,
    IS_PROD: import.meta.env.PROD,
};

export type Config = typeof config;
export default config;
