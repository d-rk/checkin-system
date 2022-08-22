import {useEffect} from 'react';
import {useNavigate} from 'react-router-dom';

export const StartPage = () => {
  const navigate = useNavigate();

  useEffect(() => navigate('/users'), [navigate]);

  return null;
};
