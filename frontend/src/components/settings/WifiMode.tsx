import {Heading, Switch} from '@chakra-ui/react';
import React, {FC} from 'react';

type Props = {
    isHotspot: boolean;
    onToggle: () => Promise<void>;
};

export const WifiMode: FC<Props> = ({isHotspot, onToggle}) => {
    return (
        <>
            <Heading as='h4' size='md'>Hotspot Mode</Heading>
            <Switch isChecked={isHotspot} onChange={() => onToggle()} />
        </>
    );
};
