import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App.tsx'
import CssBaseline from '@mui/material/CssBaseline';
import { SWRConfig } from 'swr'

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <CssBaseline>
      <SWRConfig value={{
        fetcher: (resource, init) => fetch(resource, init).then(res => res.json())
      }}>
        <App />
      </SWRConfig>
    </CssBaseline>
  </React.StrictMode>
)
