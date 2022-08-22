import {ChakraProvider} from '@chakra-ui/react';
import React, {Suspense} from 'react';
import ReactDOM from 'react-dom/client';
import {BrowserRouter as Router, Route, Routes} from 'react-router-dom';
import './api/axiosConfig';
import {App} from './App';
import {ModalProvider} from './components/useModals';
import {CheckInListPage} from './pages/CheckInListPage';
import NotFoundPage from './pages/NotFoundPage';
import {StartPage} from './pages/StartPage';
import {UserListPage} from './pages/UserListPage';
import './style.scss';

const rootElement = document.getElementById('root');

const root = ReactDOM.createRoot(rootElement!);

root.render(
  <React.StrictMode>
    <Suspense fallback={<div>Loading...</div>}>
      <ChakraProvider>
        <ModalProvider>
          <Router>
            <Routes>
              <Route path="/" element={<App />}>
                <Route index element={<StartPage />} />
                <Route path="/users" element={<UserListPage />} />
                <Route path="/checkins" element={<CheckInListPage />} />
                <Route path="*" element={<NotFoundPage />} />
              </Route>
            </Routes>
          </Router>
        </ModalProvider>
      </ChakraProvider>
    </Suspense>
  </React.StrictMode>
);
