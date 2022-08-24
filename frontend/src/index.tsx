import {ChakraProvider} from '@chakra-ui/react';
import {CalendarDefaultTheme} from '@uselessdev/datepicker';
import React, {Suspense} from 'react';
import ReactDOM from 'react-dom/client';
import {BrowserRouter as Router, Route, Routes} from 'react-router-dom';
import './api/axiosConfig';
import {App} from './App';
import {ModalProvider} from './components/useModals';
import {CalendarPage} from './pages/CalendarPage';
import {CheckInListPage} from './pages/CheckInListPage';
import NotFoundPage from './pages/NotFoundPage';
import {StartPage} from './pages/StartPage';
import {UserCheckInListPage} from './pages/UserCheckInListPage';
import {UserListPage} from './pages/UserListPage';
import './style.scss';

const rootElement = document.getElementById('root');

const root = ReactDOM.createRoot(rootElement!);

root.render(
  <React.StrictMode>
    <Suspense fallback={<div>Loading...</div>}>
      <ChakraProvider theme={CalendarDefaultTheme}>
        <ModalProvider>
          <Router>
            <Routes>
              <Route path="/" element={<App />}>
                <Route index element={<StartPage />} />
                <Route path="/calendar" element={<CalendarPage />} />
                <Route path="/users" element={<UserListPage />} />
                <Route
                  path="/users/:userId/checkins"
                  element={<UserCheckInListPage />}
                />
                <Route path="/checkins" element={<CheckInListPage />} />
                <Route path="/checkins/:day" element={<CheckInListPage />} />
                <Route path="*" element={<NotFoundPage />} />
              </Route>
            </Routes>
          </Router>
        </ModalProvider>
      </ChakraProvider>
    </Suspense>
  </React.StrictMode>
);
