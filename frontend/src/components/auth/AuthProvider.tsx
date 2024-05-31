import * as React from 'react';
import {useEffect} from 'react';
import {Navigate, useLocation} from 'react-router-dom';
import {apiLogin, getAuthenticatedUser, User} from '../../api/checkInSystemApi';
import axios from 'axios';
import {LoadingPage} from '../../pages/LoadingPage';
import {ADMIN_USER, ADMIN_PASSWORD, AUTO_LOGIN} from '../../api/config';

interface AuthContextType {
  token: string | null;
  user: User | null;
  login: (user: string, password: string, rememberMe: boolean) => Promise<void>;
  logout: () => Promise<void>;
  isUserLoading: boolean;
  isTokenExpired: boolean;
  isAuthenticated: boolean;
}

const AuthContext = React.createContext<AuthContextType>(null!);

export const AuthProvider = ({children}: {children: React.ReactNode}) => {
  const [token, setToken] = React.useState<string | null>(
    localStorage.getItem('token')
  );
  const [user, setUser] = React.useState<User | null>(null);
  const [isUserLoading, setUserLoading] = React.useState<boolean>(true);
  const [isTokenExpired, setTokenExpired] = React.useState<boolean>(false);

  useEffect(() => {
    if (token === null && !AUTO_LOGIN) {
      delete axios.defaults.headers.common['Authorization'];
      localStorage.removeItem('token');
      setUserLoading(false);
      return;
    }

    const fetchUser = async () => {
      try {
        setUserLoading(true);

        let loginToken = token;

        if (AUTO_LOGIN && token === null) {
          const {token} = await apiLogin({
            username: ADMIN_USER,
            password: ADMIN_PASSWORD,
          });
          loginToken = token;
        }

        axios.defaults.headers.common['Authorization'] = `Bearer ${loginToken}`;
        const user = await getAuthenticatedUser();
        setToken(loginToken);
        setUser(user);
      } catch (ex) {
        setTokenExpired(true);
        await logout();
      } finally {
        setUserLoading(false);
      }
    };
    fetchUser();
  }, [token]);

  const login = async (
    username: string,
    password: string,
    rememberMe: boolean
  ): Promise<void> => {
    await logout();

    const {token} = await apiLogin({
      username: username,
      password: password,
    });

    if (rememberMe) {
      localStorage.setItem('token', token);
    }

    setUserLoading(true);
    setToken(token);
  };

  const logout = (): Promise<void> => {
    setToken(null);
    setUser(null);
    return Promise.resolve();
  };

  const isAuthenticated = token !== null;

  const value = {
    token,
    user,
    login,
    logout,
    isUserLoading,
    isTokenExpired,
    isAuthenticated,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuth = () => {
  return React.useContext(AuthContext);
};

export const RequireAuth = ({children}: {children: JSX.Element}) => {
  const {isUserLoading, isAuthenticated} = useAuth();
  const location = useLocation();

  if (isUserLoading) {
    return <LoadingPage />;
  }

  if (!isAuthenticated) {
    // Redirect them to the /login page, but save the current location they were
    // trying to go to when they were redirected. This allows us to send them
    // along to that page after they log in, which is a nicer user experience
    // than dropping them off on the home page.
    return <Navigate to="/login" state={{from: location}} replace />;
  }

  return children;
};
