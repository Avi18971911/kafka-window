import {JSONValue} from "../model/MessageDetails.ts";
import JsonView from "@uiw/react-json-view";
import React from "react";

type JsonViewerProps = {
    jsonData: JSONValue
}

const JsonViewer: React.FC<JsonViewerProps> = (
    { jsonData }
) => {
    console.log("Rendering JsonViewer with data:", jsonData);
    if (typeof jsonData === 'object' && jsonData !== null) {
        return (
            <div>
                <div className="text-xs text-gray-500 mb-2">Key (JSON):</div>
                <JsonView
                    value={jsonData}
                    collapsed={2}
                />
            </div>
        );
    }

    return (
        <div>
            <div className="text-xs text-gray-500 mb-1">Key (JSON):</div>
            <div className="bg-gray-50 p-2 rounded font-mono text-sm">
                {jsonData === null ? 'null' : JSON.stringify(jsonData)}
            </div>
        </div>
    );
}

export default JsonViewer;