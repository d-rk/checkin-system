import {DownloadIcon} from '@chakra-ui/icons';
import {Box, Button, Center, useToast} from '@chakra-ui/react';
import {format, parse} from 'date-fns';
import React, {FC, useCallback, useEffect, useState} from 'react';
import {useParams} from 'react-router-dom';
import {
  createWebsocket,
  deleteCheckIn,
  downloadAllCheckIns,
  downloadCheckInList,
  isCheckInMessage,
  useCheckInList,
} from '../api/checkInSystemApi';
import {CheckInFilter} from '../components/checkin/CheckinFilter';
import CheckInList from '../components/checkin/CheckInList';
import {errorToast} from '../utils/toast';
import { LoadingPage } from "./LoadingPage";

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
  const {data: checkIns,isLoading, error, mutate} = useCheckInList(date);

  React.useMemo(
    () =>
      createWebsocket((payload: any) => {
        if (isCheckInMessage(payload) && payload.check_in) {
          if (
            format(new Date(payload.check_in.date), 'yyyy-MM-dd') ===
            format(date, 'yyyy-MM-dd')
          ) {
            mutate();
          } else {
            console.log(
              `received checking for different date ${payload.check_in.date} != ${date}`
            );
          }
        }
      }),
    [date, mutate]
  );

  const handleDateChanged = useCallback(
    (date: Date) => {
      setDate(date);
    },
    [setDate]
  );

  if (!date) {
    return null;
  }

  if (isLoading) {
    return <LoadingPage />
  }

  const handleDownload = () => {
    downloadCheckInList(date);
  };

  const handleDownloadAll = () => {
    downloadAllCheckIns();
  };

  const onDeleteCheckIn = async (checkinId: number) => {
    if (await confirm('Are you sure?')) {
      try {
        console.log(`delete checkin with id: ${checkinId}`);
        await deleteCheckIn(checkinId);
        mutate(checkIns?.filter(c => c.id !== checkinId));
      } catch (error) {
        toast(errorToast('unable to delete checkin', error));
      }
    }
  };

  if (error) {
    toast(errorToast('unable to list checkIns', error));
  }

  return (
    <Center>
      <Box>
        <CheckInFilter
          date={date}
          onDateChange={handleDateChanged}
          onDownload={handleDownload}
        />
        <CheckInList checkIns={checkIns ?? []} onDelete={onDeleteCheckIn} />
        <Center>
          <Button
            leftIcon={<DownloadIcon />}
            colorScheme="teal"
            aria-label="Download all CheckIns as .csv"
            title="Download all CheckIns as .csv"
            variant="solid"
            onClick={handleDownloadAll}
          >
            Download all
          </Button>
        </Center>
      </Box>
    </Center>
  );
};
