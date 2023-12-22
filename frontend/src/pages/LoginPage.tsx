import React, {FC} from 'react';
import {useLocation, useNavigate} from 'react-router-dom';
import {useAuth} from '../components/auth/AuthProvider';
import {LoginForm, LoginFields} from '../components/login/LoginForm';
import {errorToast} from '../utils/toast';
import {useToast} from '@chakra-ui/react';

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

  const from = (location.state as LocationState)?.from?.pathname || '/';

  async function handleSubmit({username, password, rememberMe}: LoginFields) {
    try {
      await auth.login(username, password, rememberMe);
      navigate(from, {replace: true});
    } catch (error) {
      toast(errorToast('login failed', error));
    }
  }

  return <LoginForm onSubmit={handleSubmit} />;
};

export default LoginPage;
