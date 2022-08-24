import {Box} from '@chakra-ui/react';
import React from 'react';
import {Outlet} from 'react-router-dom';
import Header from './components/header/header';
export const App = () => {
  return (
    <>
      <Header />
      <Box p={4}>
        <Outlet />
      </Box>
    </>
  );
};
