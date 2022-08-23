import {Box, Center, useToast} from '@chakra-ui/react';
import React, {FC, useState} from 'react';
import {useCheckInList} from '../api/checkInSystemApi';
import {CheckInFilter} from '../components/checkin/CheckinFilter';
import CheckInList from '../components/checkin/CheckInList';
import {errorToast} from '../utils/toast';

export const CheckInListPage: FC = () => {
  const [date, setDate] = useState<Date>(new Date());

  const toast = useToast();
  const {data: checkIns, error, mutate} = useCheckInList(date);

  if (error) {
    toast(errorToast('unable to list checkIns', error));
  }

  const onDateChange = (date: Date) => {
    setDate(date);
  };

  return (
    <Center>
      <Box>
        <CheckInFilter onDateChange={onDateChange} />
        <CheckInList checkIns={checkIns ?? []} />
      </Box>
    </Center>
  );
};
