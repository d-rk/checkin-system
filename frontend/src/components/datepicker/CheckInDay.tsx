import {Button} from '@chakra-ui/button';
import {useToast} from '@chakra-ui/react';
import {useCalendarDay} from '@uselessdev/datepicker';
import {format} from 'date-fns';
import React from 'react';
import {useCheckInDates} from '../../api/checkInSystemApi';
import {errorToast} from '../../utils/toast';

export const CheckInDay = () => {
  const toast = useToast();
  const {day, onSelectDates, isDisabled} = useCalendarDay();

  const {data, error} = useCheckInDates();

  if (error) {
    toast(errorToast('unable to fetch checkin dates', error));
  }

  const dates = data ? data.map(d => new Date(d.date).toDateString()) : [];

  const hasCheckIns = dates.includes(new Date(day).toDateString());
  const isToday = new Date(day).toDateString() === new Date().toDateString();

  const today = isToday
    ? {
        bgColor: 'gray.300',
        color: 'black',
        rounded: 4,
        _hover: {
          bgColor: 'gray.100',
        },
      }
    : {};

  const withCheckins = hasCheckIns
    ? {
        bgColor: 'teal.300',
        color: 'white',
        rounded: 4,
        _hover: {
          bgColor: 'teal.200',
        },
      }
    : {};

  return (
    <Button
      variant="ghost"
      disabled={isDisabled}
      onClick={() => onSelectDates(day)}
      sx={{...today, ...withCheckins}}
    >
      {format(day, 'd')}
    </Button>
  );
};
