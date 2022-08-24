import {Box, Center, Flex, Heading} from '@chakra-ui/react';
import {
  Calendar,
  CalendarControls,
  CalendarDate,
  CalendarDays,
  CalendarMonth,
  CalendarMonthName,
  CalendarMonths,
  CalendarNextButton,
  CalendarPrevButton,
  CalendarWeek,
} from '@uselessdev/datepicker';
import {format} from 'date-fns';
import React, {FC} from 'react';
import {useNavigate} from 'react-router-dom';
import {CheckInDay} from '../components/datepicker/CheckInDay';

export const CalendarPage: FC = () => {
  const navigate = useNavigate();

  const handleSelectDate = (date: CalendarDate) => {
    navigate(`/checkins/${format(date, 'yyyy-MM-dd')}`);
  };

  return (
    <Center>
      <Box>
        <Flex>
          <Heading as="h5">CheckIn Calendar</Heading>
        </Flex>
        <Calendar
          value={{}}
          onSelectDate={date => handleSelectDate(date as CalendarDate)}
          singleDateSelection
          disableFutureDates
        >
          <CalendarControls>
            <CalendarPrevButton />
            <CalendarNextButton />
          </CalendarControls>

          <CalendarMonths>
            <CalendarMonth>
              <CalendarMonthName />
              <CalendarWeek />
              <CalendarDays>
                <CheckInDay />
              </CalendarDays>
            </CalendarMonth>
          </CalendarMonths>
        </Calendar>
      </Box>
    </Center>
  );
};
