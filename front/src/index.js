import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';
import axios from 'axios';

const currentURL = document.URL

if (currentURL.includes("176.109.99.197") || currentURL.includes("enchik-pet")) {
    axios.defaults.baseURL = "http://176.109.99.197:8008/"
} else {
    axios.defaults.baseURL = "http://localhost:8008/"
}

axios.defaults.withCredentials = true

const root = ReactDOM.createRoot(document.getElementById('root'))
root.render(
    <App />
)
