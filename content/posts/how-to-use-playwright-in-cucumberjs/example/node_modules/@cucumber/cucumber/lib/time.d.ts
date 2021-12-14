import * as messages from '@cucumber/messages';
declare const methods: any;
export declare function durationBetweenTimestamps(startedTimestamp: messages.Timestamp, finishedTimestamp: messages.Timestamp): messages.Duration;
export declare function wrapPromiseWithTimeout<T>(promise: Promise<T>, timeoutInMilliseconds: number, timeoutMessage?: string): Promise<T>;
export default methods;
