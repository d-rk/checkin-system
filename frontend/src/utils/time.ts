import {Clock} from "../api/checkInSystemApi";
import {format, formatISO, parse, parseISO} from "date-fns";

const DATE_FORMAT = "dd.MM.yyyy HH:mm:ss";

export const toLocaleClock = (clock: Clock): Clock => ({
    refTimestamp: toLocaleString(clock.refTimestamp),
    timestamp: toLocaleString(clock.timestamp),
});

export const fromLocaleClock = (clock: Clock): Clock => ({
    refTimestamp: fromLocaleString(clock.refTimestamp),
    timestamp: fromLocaleString(clock.timestamp),
});

export const toLocaleString = (ts: string): string => format(parseISO(ts), DATE_FORMAT);

const fromLocaleString = (ts: string): string => formatISO(parse(ts, DATE_FORMAT, new Date()));