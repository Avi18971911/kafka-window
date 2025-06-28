import {JSONValue} from "../model/MessageDetails.ts";
import JsonView from "@uiw/react-json-view";
import React from "react";
import {darkTheme} from "@uiw/react-json-view/dark";

type JsonViewerProps = {
    jsonData: JSONValue
}

const JsonViewer: React.FC<JsonViewerProps> = (
    { jsonData }
) => {
    if (typeof jsonData === 'object' && jsonData !== null) {
        return (
            <div>
                <JsonView
                    style={darkTheme}
                    value={jsonData}
                    collapsed={2}
                />
            </div>
        );
    }

    return (
        <div>
            <div>
                {jsonData === null ? 'null' : JSON.stringify(jsonData)}
            </div>
        </div>
    );
}

export default JsonViewer;