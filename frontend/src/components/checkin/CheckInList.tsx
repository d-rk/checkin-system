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
import {CheckInWithUser} from '../../api/checkInSystemApi';

type Props = {
  checkIns: CheckInWithUser[];
};

const CheckInList: FC<Props> = ({checkIns}) => {
  return (
    <TableContainer>
      <Table>
        <TableCaption>CheckIns</TableCaption>
        <Thead>
          <Tr>
            <Th isNumeric>ID</Th>
            <Th>Time</Th>
            <Th>User</Th>
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
                  'HH:mm:ss'
                )}
              </Td>
              <Td>{checkIn.user.name}</Td>
            </Tr>
          ))}
          {checkIns.length === 0 && (
            <Tr>
              <Td isNumeric></Td>
              <Td>NO CHECK INS FOUND</Td>
              <Td></Td>
            </Tr>
          )}
        </Tbody>
      </Table>
    </TableContainer>
  );
};

export default CheckInList;
