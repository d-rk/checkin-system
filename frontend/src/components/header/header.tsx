import {Box, Flex, HStack, Link, useColorModeValue} from '@chakra-ui/react';
import React, {ReactNode} from 'react';

const links = [
  {name: 'Users', route: '/users'},
  {name: 'CheckIns', route: '/checkins'},
];

const NavLink = ({children}: {children: ReactNode}) => (
  <Link
    px={2}
    py={1}
    rounded={'md'}
    _hover={{
      textDecoration: 'none',
      bg: useColorModeValue('gray.200', 'gray.700'),
    }}
    href={'#'}
  >
    {children}
  </Link>
);

const Header = () => {
  return (
    <Box bg={useColorModeValue('gray.100', 'gray.900')} px={4}>
      <Flex h={16} alignItems={'center'} justifyContent={'space-between'}>
        <HStack spacing={8} alignItems={'center'}>
          <Box>Logo</Box>
          <HStack as={'nav'} spacing={4} display={{base: 'none', md: 'flex'}}>
            {links.map(link => (
              <NavLink key={link.route}>{link.name}</NavLink>
            ))}
          </HStack>
        </HStack>
      </Flex>
    </Box>
  );
};

export default Header;
