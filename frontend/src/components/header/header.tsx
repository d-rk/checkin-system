import {HamburgerIcon} from '@chakra-ui/icons';
import {
  Avatar,
  Box,
  Flex,
  HStack,
  IconButton,
  Link,
  Menu,
  MenuButton,
  MenuItem,
  MenuList,
  useColorModeValue,
} from '@chakra-ui/react';
import React from 'react';
import {NavLink} from 'react-router-dom';
import Logo from '../../assets/logo.svg?react';
import {useAuth} from '../auth/AuthProvider';

const links = [
  {name: 'Calendar', route: '/calendar'},
  {name: 'Users', route: '/users'},
  {name: 'CheckIns', route: '/checkins'},
];

const Header = () => {
  const hoverBg = useColorModeValue('gray.200', 'gray.700');
  const {user, isAuthenticated} = useAuth();

  const gray = useColorModeValue('gray.100', 'gray.900');

  if (!isAuthenticated) {
    return null;
  }

  return (
    <Box bg={gray} px={4}>
      <Flex h={16} alignItems={'center'} justifyContent={'space-between'}>
        <HStack spacing={8} alignItems={'center'}>
          <Box>
            <NavLink to="/">
              <Logo height={24} width={24} />
            </NavLink>
          </Box>
          <HStack as={'nav'} spacing={4} display={{base: 'none', md: 'flex'}}>
            {links.map(link => (
              <Link
                as={NavLink}
                px={2}
                py={1}
                rounded={'md'}
                _hover={{
                  textDecoration: 'none',
                  bg: hoverBg,
                }}
                _activeLink={{color: 'white', background: 'blue.500'}}
                key={link.route}
                to={link.route}
              >
                {link.name}
              </Link>
            ))}
          </HStack>
          <Menu>
            <MenuButton
              as={IconButton}
              aria-label="Options"
              icon={<HamburgerIcon />}
              variant="outline"
              display={{base: 'flex', md: 'none'}}
            />
            <MenuList>
              {links.map(link => (
                <MenuItem
                  key={link.route}
                  as={NavLink}
                  _activeLink={{color: 'white', background: 'blue.500'}}
                  to={link.route}
                >
                  {link.name}
                </MenuItem>
              ))}
            </MenuList>
          </Menu>
        </HStack>
        <Menu placement="bottom-start">
          <MenuButton>
            <Avatar />
          </MenuButton>
          <MenuList>
            <Box bg="blue.100" w="100%" p={2}>
              Logged in as <b>{user?.name}</b>
            </Box>
            <MenuItem as={NavLink} to="/logout">
              Logout
            </MenuItem>
          </MenuList>
        </Menu>
      </Flex>
    </Box>
  );
};

export default Header;
