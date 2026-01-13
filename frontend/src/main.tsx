import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './app/App';
import config from './app/config';

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App initialTitle={config.APP_TITLE} />
  </React.StrictMode>,
);
