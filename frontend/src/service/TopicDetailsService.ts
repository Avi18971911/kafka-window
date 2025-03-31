import {TopicDetails} from "../model/TopicDetails.ts";
import {ModelTopicDetails} from "../backend_api";

export const mapModelTopicDetailsToTopicDetails = (model: ModelTopicDetails[]): TopicDetails[] => {
    return  model.map((modelDetails) => {
        return {
            topic: modelDetails.name,
            numPartitions: modelDetails.numPartitions,
            replicationFactor: modelDetails.replicationFactor,
            isInternal: modelDetails.isInternal,
            cleanupPolicy: modelDetails.cleanupPolicy,
            retentionMs: modelDetails.retentionMs,
            retentionBytes: modelDetails.retentionBytes,
            additionalConfigs: modelDetails.additionalConfigs
        }
    });
}