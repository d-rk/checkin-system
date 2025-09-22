import {Heading, SimpleGrid, Table, TableContainer, Tbody, Td, Th, Thead, Tr,} from '@chakra-ui/react';
import React, {FC} from 'react';
import {VersionInfo,} from '../../api/checkInSystemApi';
import {FRONTEND_VERSION} from "../../api/config";
import {toLocaleString} from "../../utils/time";

type Props = {
  backend?: VersionInfo;
};

export const VersionsTable: FC<Props> = ({backend}) => {

  return (
      <SimpleGrid spacing={5} columns={1}>
        <Heading>Versions</Heading>
        <TableContainer>
          <Table variant='simple'>
            <Thead>
              <Tr>
                <Th>Module</Th>
                <Th>Version</Th>
                <Th>Date</Th>
                <Th>GitCommit</Th>
              </Tr>
            </Thead>
            <Tbody>
              <Tr>
                <Td>Backend</Td>
                <Td>{ backend?.version ?? '-' }</Td>
                <Td>{ backend?.buildTime !== undefined ? toLocaleString(backend!.buildTime) : '-' }</Td>
                <Td>{ backend?.gitCommit ?? '-' }</Td>
              </Tr>
              <Tr>
                <Td>Frontend</Td>
                <Td>{ FRONTEND_VERSION.version ?? '-' }</Td>
                <Td>{ FRONTEND_VERSION.buildTime ? toLocaleString(FRONTEND_VERSION.buildTime) : '-' }</Td>
                <Td>{ FRONTEND_VERSION.gitCommit ?? '-' }</Td>
              </Tr>
            </Tbody>
          </Table>
        </TableContainer>
      </SimpleGrid>
  );
};
