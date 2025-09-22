import {DeleteIcon} from '@chakra-ui/icons';
import {IconButton, Table, TableCaption, TableContainer, Tbody, Td, Th, Thead, Tr,} from '@chakra-ui/react';
import React, {FC} from 'react';
import {WifiNetwork} from '../../api/checkInSystemApi';

type Props = {
  networks: WifiNetwork[];
  onRemove: (ssid: string) => Promise<void>;
};

const WifiNetworkList: FC<Props> = ({networks, onRemove}) => {
  return (
    <TableContainer>
      <Table>
        <TableCaption>Configured Networks</TableCaption>
        <Thead>
          <Tr>
            <Th>SSID</Th>
            <Th>Password</Th>
            <Th textAlign="right">Actions</Th>
          </Tr>
        </Thead>
        <Tbody>
          {networks.map(network => (
            <Tr key={network.ssid}>
              <Td>{network.ssid}</Td>
              <Td>***</Td>
              <Td textAlign="right">
                <IconButton
                  colorScheme="red"
                  aria-label="Remove Network"
                  title="Remove Network"
                  icon={<DeleteIcon />}
                  onClick={() => onRemove(network.ssid)}
                />
              </Td>
            </Tr>
          ))}
          {networks.length === 0 && (
            <Tr>
              <Td colSpan={3}>No configured networks</Td>
            </Tr>
          )}
        </Tbody>
      </Table>
    </TableContainer>
  );
};

export default WifiNetworkList;
