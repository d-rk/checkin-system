import {ChakraProvider, Spinner} from '@chakra-ui/react';
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
import {AuthProvider, RequireAuth} from './components/auth/AuthProvider';
import LoginPage from './pages/LoginPage';
import LogoutPage from './pages/LogoutPage';
import {SettingsPage} from './pages/SettingsPage';
import {VersionsPage} from './pages/VersionsPage';

const rootElement = document.getElementById('root');

const root = ReactDOM.createRoot(rootElement!);

root.render(
  <React.StrictMode>
    <Suspense fallback={<Spinner />}>
      <ChakraProvider theme={CalendarDefaultTheme}>
        <AuthProvider>
          <ModalProvider>
            <Router>
              <Routes>
                <Route path="/" element={<App />}>
                  <Route index element={<StartPage />} />
                  <Route path="/login" element={<LoginPage />} />
                  <Route path="/logout" element={<LogoutPage />} />
                  <Route
                    path="/calendar"
                    element={
                      <RequireAuth>
                        <CalendarPage />
                      </RequireAuth>
                    }
                  />
                  <Route
                    path="/users"
                    element={
                      <RequireAuth>
                        <UserListPage />
                      </RequireAuth>
                    }
                  />
                  <Route
                    path="/users/:userId/checkins"
                    element={
                      <RequireAuth>
                        <UserCheckInListPage />
                      </RequireAuth>
                    }
                  />
                  <Route
                    path="/checkins"
                    element={
                      <RequireAuth>
                        <CheckInListPage />
                      </RequireAuth>
                    }
                  />
                  <Route
                    path="/checkins/:day"
                    element={
                      <RequireAuth>
                        <CheckInListPage />
                      </RequireAuth>
                    }
                  />
                  <Route
                    path="/settings"
                    element={
                      <RequireAuth>
                        <SettingsPage />
                      </RequireAuth>
                    }
                  />
                  <Route
                    path="/versions"
                    element={
                      <RequireAuth>
                        <VersionsPage />
                      </RequireAuth>
                    }
                  />
                  <Route
                    path="*"
                    element={
                      <RequireAuth>
                        <NotFoundPage />
                      </RequireAuth>
                    }
                  />
                </Route>
              </Routes>
            </Router>
          </ModalProvider>
        </AuthProvider>
      </ChakraProvider>
    </Suspense>
  </React.StrictMode>
);
