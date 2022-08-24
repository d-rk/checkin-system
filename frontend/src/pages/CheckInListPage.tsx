import {Box, Center, useToast} from '@chakra-ui/react';
import {parse} from 'date-fns';
import React, {FC, useEffect, useState} from 'react';
import {useParams} from 'react-router-dom';
import {Websocket} from 'websocket-ts/lib';
import {
  createWebsocket,
  isCheckInMessage,
  useCheckInList,
} from '../api/checkInSystemApi';
import {CheckInFilter} from '../components/checkin/CheckinFilter';
import CheckInList from '../components/checkin/CheckInList';
import {errorToast} from '../utils/toast';

const parseDay = (day?: string) => {
  if (day) {
    return parse(day, 'yyyy-MM-dd', new Date(0));
  } else {
    return new Date();
  }
};

export const CheckInListPage: FC = () => {
  const {day} = useParams();

  const [date, setDate] = useState<Date>(parseDay(day));

  useEffect(() => {
    setDate(parseDay(day));
  }, [day]);

  const toast = useToast();
  const {data: checkIns, error, mutate} = useCheckInList(date);

  React.useState<Websocket>(
    createWebsocket((payload: any) => {
      if (isCheckInMessage(payload) && payload.check_in) {
        if (new Date(payload.check_in.date).getTime() === date.getTime()) {
          mutate();
        }
      }
    })
  );

  if (!date) {
    return null;
  }

  if (error) {
    toast(errorToast('unable to list checkIns', error));
  }

  const onDateChange = (date: Date) => {
    setDate(date);
  };

  return (
    <Center>
      <Box>
        <CheckInFilter date={date} onDateChange={onDateChange} />
        <CheckInList checkIns={checkIns ?? []} />
      </Box>
    </Center>
  );
};
