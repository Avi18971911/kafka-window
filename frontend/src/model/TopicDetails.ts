export type TopicDetails = {
    topic: string;
    numPartitions: number;
    replicationFactor: number;
    isInternal: boolean;
    cleanupPolicy: string;
    retentionMs: RetentionMs
    retentionBytes: number;
    additionalConfigs: Record<string, string>;
}

export type RetentionMs = {
    value: number;
    indefinite: boolean;
}