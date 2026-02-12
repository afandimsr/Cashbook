import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './app/App';
import config from './app/config';
import { registerSW } from 'virtual:pwa-register'

// Register PWA Service Worker
registerSW({ immediate: true })

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App initialTitle={config.APP_TITLE} />
  </React.StrictMode>,
);
