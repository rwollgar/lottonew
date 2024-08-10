import React from 'react'
import ReactDOM from 'react-dom/client'
import App from '../ui/App'
import '../ui/index.css'

import { RouterProvider, createRouter } from '@tanstack/react-router'

// Import the generated route tree
import { routeTree } from '../ui/routeTree.gen';

// Create a new router instance
const router = createRouter({ routeTree })

ReactDOM.createRoot(document.getElementById('root')).render(
    <React.StrictMode>
        {/* <App /> */}
        <RouterProvider router={router} />
    </React.StrictMode>,
)
