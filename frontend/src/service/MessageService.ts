import {ModelJSONValue, ModelMessage, ModelPayloadType} from "../backend_api";
import {JSONValue, MessageDetails, PayloadType} from "../model/MessageDetails.ts";

export const mapModelMessageToMessage = (model: ModelMessage[]): MessageDetails[] => {
    return  model.map((modelDetails) => {
        return {
            key: modelDetails.key,
            keyJsonPayload: mapJSONValueToNativeType(modelDetails.keyJsonPayload),
            keyPayloadType: mapModelPayloadTypeToPayloadType(modelDetails.keyPayloadType),
            offset: modelDetails.offset,
            partition: modelDetails.partition,
            timestamp: modelDetails.timestamp,
            topic: modelDetails.topic,
            value: modelDetails.value,
            valueJsonPayload: mapJSONValueToNativeType(modelDetails.valueJsonPayload),
            valuePayloadType: mapModelPayloadTypeToPayloadType(modelDetails.valuePayloadType)
        }
    });
}

const mapJSONValueToNativeType = (value: ModelJSONValue | undefined): JSONValue => {
    if (value === undefined) {
        return null;
    } else if (value.arrayVal) {
        return value.arrayVal.map(mapJSONValueToNativeType);
    } else if (value.boolVal !== undefined) {
        return value.boolVal;
    } else if (value.nullVal !== undefined && value.nullVal) {
        return null;
    } else if (value.numberVal !== undefined) {
        return value.numberVal;
    } else if (value.objectVal) {
        const obj: { [key: string]: JSONValue } = {};
        for (const key in value.objectVal) {
            obj[key] = mapJSONValueToNativeType(value.objectVal[key]);
        }
        return obj;
    } else if (value.stringVal !== undefined) {
        return value.stringVal;
    }
    console.log("Invalid ModelJSONValue:", value);
    return null;
}

const mapModelPayloadTypeToPayloadType = (type: ModelPayloadType): PayloadType => {
    switch (type) {
        case ModelPayloadType.StringPayload:
            return 'string';
        case ModelPayloadType.JSONPayload:
            return 'json';
        case ModelPayloadType.ConsumerOffsetPayload:
            return 'consumerOffset';
        default:
            throw new Error(`Unknown ModelPayloadType: ${type}`);
    }
}