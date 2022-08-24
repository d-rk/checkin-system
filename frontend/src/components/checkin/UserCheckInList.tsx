import {
  Table,
  TableCaption,
  TableContainer,
  Tbody,
  Td,
  Th,
  Thead,
  Tr,
} from '@chakra-ui/react';
import {formatInTimeZone} from 'date-fns-tz';
import React, {FC} from 'react';
import {CheckIn} from '../../api/checkInSystemApi';

type Props = {
  checkIns: CheckIn[];
};

const UserCheckInList: FC<Props> = ({checkIns}) => {
  return (
    <TableContainer>
      <Table>
        <TableCaption>CheckIns</TableCaption>
        <Thead>
          <Tr>
            <Th isNumeric>ID</Th>
            <Th>Date</Th>
            <Th>Time</Th>
          </Tr>
        </Thead>
        <Tbody>
          {checkIns.map(checkIn => (
            <Tr key={checkIn.id}>
              <Td isNumeric>{checkIn.id}</Td>
              <Td>
                {formatInTimeZone(
                  new Date(checkIn.timestamp),
                  'Europe/Berlin',
                  'yy-MM-dd'
                )}
              </Td>
              <Td>
                {formatInTimeZone(
                  new Date(checkIn.timestamp),
                  'Europe/Berlin',
                  'HH:mm:ss'
                )}
              </Td>
            </Tr>
          ))}
          {checkIns.length === 0 && (
            <Tr>
              <Td></Td>
              <Td>-</Td>
              <Td></Td>
            </Tr>
          )}
        </Tbody>
      </Table>
    </TableContainer>
  );
};

export default UserCheckInList;
