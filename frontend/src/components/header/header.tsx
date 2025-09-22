import {HamburgerIcon, WarningTwoIcon} from '@chakra-ui/icons';
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
  Tooltip,
  useColorModeValue,
  useToast,
} from '@chakra-ui/react';
import React, {useCallback, useEffect, useMemo} from 'react';
import {NavLink} from 'react-router-dom';
import Logo from '../../assets/logo.svg?react';
import {useAuth} from '../auth/AuthProvider';
import {useClock} from '../../api/checkInSystemApi';
import {parseISO, startOfMinute} from 'date-fns';
import {errorToast} from '../../utils/toast';
import {toLocaleString} from "../../utils/time";

const links = [
  {name: 'Calendar', route: '/calendar'},
  {name: 'Users', route: '/users'},
  {name: 'CheckIns', route: '/checkins'},
];

const Header = () => {
  const hoverBg = useColorModeValue('gray.200', 'gray.700');
  const {user, isAuthenticated} = useAuth();

  const now = startOfMinute(new Date());

  const {data: clock} = useClock(now);
  const toast = useToast();

  const clockDiff = useMemo(() => {
    if (clock) {
      const refTime = parseISO(clock?.refTimestamp);
      const clockTime = parseISO(clock?.timestamp);
      return Math.abs((refTime.getTime() - clockTime.getTime()) / 1000);
    }
  }, [clock]);

  const clockOutOfSync = useMemo(() => {
    return clockDiff !== undefined && clockDiff > 60;
  }, [clockDiff]);

  const showOutOfSyncToast = useCallback(() => {
    if (clock && clockDiff) {
      toast(
        errorToast(
          `hardware clock shows ${toLocaleString(clock.timestamp)}, which is out of sync by ${Math.round(clockDiff / 60)} minutes`
        )
      );
    }
  }, [clock, clockDiff, toast]);

  useEffect(() => {
    if (clockOutOfSync) {
      showOutOfSyncToast();
    }
  }, [showOutOfSyncToast, clockOutOfSync]);

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
        <Box flexGrow={1} textAlign={'end'} px={2}>
          {clockOutOfSync && (
            <Tooltip label="Hardware Clock out of sync">
              <IconButton aria-label="warn" onClick={showOutOfSyncToast}>
                <WarningTwoIcon color="red.400" />
              </IconButton>
            </Tooltip>
          )}
        </Box>

        <Menu placement="bottom-start">
          <MenuButton>
            <Avatar />
          </MenuButton>
          <MenuList>
            <Box bg="blue.100" w="100%" p={2}>
              Logged in as <b>{user?.name}</b>
            </Box>
            <MenuItem as={NavLink} to="/settings">
              Settings
            </MenuItem>
            <MenuItem as={NavLink} to="/versions">
              Versions
            </MenuItem>
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
