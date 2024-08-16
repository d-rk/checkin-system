import {FC, useState} from 'react';
import {Navigate, useLocation, useNavigate} from 'react-router-dom';
import {useAuth} from '../components/auth/AuthProvider';
import {LoginForm, LoginFields} from '../components/login/LoginForm';
import {errorToast} from '../utils/toast';
import {useToast} from '@chakra-ui/react';
import {LoadingPage} from './LoadingPage';
import * as React from 'react';

interface LocationState {
  from: {
    pathname: string;
  };
}

const LoginPage: FC = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const auth = useAuth();
  const toast = useToast();
  const [isLoading, setLoading] = useState<boolean>(false);

  const from = (location.state as LocationState)?.from?.pathname || '/';

  async function handleSubmit({username, password, rememberMe}: LoginFields) {
    try {
      setLoading(true);
      await auth.login(username, password, rememberMe);
      navigate(from, {replace: true});
    } catch (error) {
      toast(errorToast('login failed', error));
    } finally {
      setLoading(false);
    }
  }

  if (isLoading) {
    return <LoadingPage />;
  }

  if (auth.isAuthenticated) {
    return <Navigate to="/" replace />;
  }

  return <LoginForm onSubmit={handleSubmit} />;
};

export default LoginPage;
