import {Button} from '@chakra-ui/button';
import {Box, Circle, Text} from '@chakra-ui/react';
import {useCalendarDay} from '@uselessdev/datepicker';
import {format} from 'date-fns';
import React from 'react';

export const CheckInDay = () => {
  const {day, onSelectDates, isSelected, isInRange} = useCalendarDay();

  const selected = isSelected
    ? {
        bgColor: 'teal.400',
        color: 'white',
        rounded: 0,
        _hover: {
          bgColor: 'teal.300',
        },
      }
    : {};

  const range = isInRange
    ? {
        bgColor: 'teal.300',
        color: 'white',
        rounded: 'none',
        _hover: {
          bgColor: 'teal.200',
        },
      }
    : {};

  return (
    <Button
      variant="ghost"
      onClick={() => onSelectDates(day)}
      sx={{...selected, ...range}}
    >
      {new Date(day).getDate() < 8 ? (
        <Box flexDirection="column" alignItems="center">
          <Text>{format(day, 'd')}</Text>
          <Circle size="4px" bgColor="pink.300" />
        </Box>
      ) : (
        format(day, 'd')
      )}
    </Button>
  );
};
