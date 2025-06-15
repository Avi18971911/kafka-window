export type MessageDetails = {
    key: string
    keyJsonPayload: JSONValue
    keyPayloadType: PayloadType
    offset: number
    partition: number
    timestamp: string
    topic: string
    value: string
    valueJsonPayload: JSONValue
    valuePayloadType: PayloadType
}

export type PayloadType = 'string' | 'json' | 'consumerOffset'

export type JSONValue = string | number | boolean | null | JSONValue[] | { [key: string]: JSONValue };
