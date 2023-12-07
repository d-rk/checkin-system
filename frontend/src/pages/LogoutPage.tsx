import React, {FC, useEffect} from 'react';
import {useNavigate} from 'react-router-dom';
import {useAuth} from '../components/auth/AuthProvider';
import {LoadingPage} from './LoadingPage';

const LogoutPage: FC = () => {
  const navigate = useNavigate();
  const auth = useAuth();

  useEffect(() => {
    const logout = async () => {
      await auth.logout();
      navigate('/');
    };
    logout();
  }, [auth, navigate]);

  return <LoadingPage />;
};

export default LogoutPage;
