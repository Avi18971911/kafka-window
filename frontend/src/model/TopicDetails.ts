export type TopicDetails = {
    topic: string;
    numPartitions: number;
    replicationFactor: number;
    isInternal: boolean;
    cleanupPolicy: string;
    retentionMs: RetentionMs | undefined;
    retentionBytes: number | undefined;
    additionalConfigs: Record<string, string>;
}

export type RetentionMs = {
    value: number;
    indefinite: boolean;
}